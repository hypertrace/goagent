package metrics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tklauser/go-sysconf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "hypertrace.goagent.metrics"

type systemMetrics struct {
	memory          float64
	cpuSecondsTotal float64
}

type processStats struct {
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	rss    float64
}

const procStatArrayLength = 52

var (
	clkTck   = getClockTicks()
	pageSize = float64(os.Getpagesize())
)

func InitialiseMetrics() {
	meterProvider := otel.GetMeterProvider()
	meter := meterProvider.Meter(meterName)
	err := setUpMetricRecorder(meter)
	if err != nil {
		fmt.Println("error initialising metrics, failed to setup metric recorder")
	}
}

func processStatsFromPid(pid int) (*systemMetrics, error) {
	sysInfo := &systemMetrics{}
	procFilepath := filepath.Join("/proc", strconv.Itoa(pid), "stat")
	var err error
	if procStatFileBytes, err := os.ReadFile(filepath.Clean(procFilepath)); err == nil {
		if stat, err := parseProcStatFile(procStatFileBytes, procFilepath); err == nil {
			sysInfo.memory = stat.rss * pageSize
			sysInfo.cpuSecondsTotal = (stat.stime + stat.utime + stat.cstime + stat.cutime) / clkTck
			return sysInfo, nil
		}
		return nil, err
	}
	return nil, err
}

// ref: /proc/pid/stat section of https://man7.org/linux/man-pages/man5/proc.5.html
func parseProcStatFile(bytesArr []byte, procFilepath string) (*processStats, error) {
	infos := strings.Split(string(bytesArr), " ")
	if len(infos) != procStatArrayLength {
		return nil, errors.New(fmt.Sprintf("%s file could not be parsed", procFilepath))
	}
	return &processStats{
		utime:  parseFloat(infos[13]),
		stime:  parseFloat(infos[14]),
		cutime: parseFloat(infos[15]),
		cstime: parseFloat(infos[16]),
		rss:    parseFloat(infos[23]),
	}, nil
}
func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

// sysconf for go. claims to work without cgo or external binaries
// https://pkg.go.dev/github.com/tklauser/go-sysconf@v0.3.14#section-readme
func getClockTicks() float64 {
	clktck, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	if err != nil {
		return float64(100)
	}
	return float64(clktck)
}

func setUpMetricRecorder(meter metric.Meter) error {
	if meter == nil {
		return fmt.Errorf("error while setting up metric recorder: meter is nil")
	}
	cpuSeconds, err := meter.Float64ObservableCounter("cpu.seconds.total", metric.WithDescription("Metric to monitor total CPU seconds"))
	if err != nil {
		return fmt.Errorf("error while setting up cpu seconds metric counter: %v", err)
	}
	memory, err := meter.Float64ObservableGauge("memory", metric.WithDescription("Metric to monitor memory usage"))
	if err != nil {
		return fmt.Errorf("error while setting up memory metric counter: %v", err)
	}
	// Register the callback function for both cpu_seconds and memory observable gauges
	_, err = meter.RegisterCallback(
		func(ctx context.Context, result metric.Observer) error {
			systemMetrics, err := processStatsFromPid(os.Getpid())
			result.ObserveFloat64(cpuSeconds, systemMetrics.cpuSecondsTotal)
			result.ObserveFloat64(memory, systemMetrics.memory)
			return err
		},
		cpuSeconds, memory,
	)
	if err != nil {
		log.Fatalf("failed to register callback: %v", err)
		return err
	}
	return nil
}
