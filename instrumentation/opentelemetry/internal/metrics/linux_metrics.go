//go:build linux

package metrics

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tklauser/go-sysconf"
)

const procStatArrayLength = 52

var (
	clkTck   = getClockTicks()
	pageSize = float64(os.Getpagesize())
)

type processStats struct {
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	rss    float64
}

type linuxMetrics struct {
	memory          float64
	cpuSecondsTotal float64
}

func newSystemMetrics() (systemMetrics, error) {
	lm := &linuxMetrics{}
	stats, err := lm.processStatsFromPid(os.Getpid())
	if err != nil {
		return nil, err
	}
	lm.memory = stats.rss * pageSize
	lm.cpuSecondsTotal = (stats.stime + stats.utime + stats.cstime + stats.cutime) / clkTck
	return lm, nil
}

func (lm *linuxMetrics) getMemory() float64 {
	return lm.memory
}

func (lm *linuxMetrics) getCPU() float64 {
	return lm.cpuSecondsTotal
}

func (lm *linuxMetrics) processStatsFromPid(pid int) (*processStats, error) {
	procFilepath := filepath.Join("/proc", strconv.Itoa(pid), "stat")
	var err error
	if procStatFileBytes, err := os.ReadFile(filepath.Clean(procFilepath)); err == nil {
		if stat, err := lm.parseProcStatFile(procStatFileBytes, procFilepath); err == nil {
			if err != nil {
				return nil, err
			}
			return stat, nil
		}
		return nil, err
	}
	return nil, err
}

// ref: /proc/pid/stat section of https://man7.org/linux/man-pages/man5/proc.5.html
func (lm *linuxMetrics) parseProcStatFile(bytesArr []byte, procFilepath string) (*processStats, error) {
	infos := strings.Split(string(bytesArr), " ")
	if len(infos) != procStatArrayLength {
		return nil, fmt.Errorf("%s file could not be parsed", procFilepath)
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
