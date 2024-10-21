package metrics

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "goagent.hypertrace.org/metrics"

type systemMetrics interface {
	getMemory() float64
	getCPU() float64
}

func InitializeSystemMetrics() {
	meterProvider := otel.GetMeterProvider()
	meter := meterProvider.Meter(meterName)
	err := setUpMetricRecorder(meter)
	if err != nil {
		log.Printf("error initializing metrics, failed to setup metric recorder: %v\n", err)
	}
}

func setUpMetricRecorder(meter metric.Meter) error {
	if meter == nil {
		return fmt.Errorf("error while setting up metric recorder: meter is nil")
	}
	cpuSeconds, err := meter.Float64ObservableCounter("hypertrace.agent.cpu.seconds.total", metric.WithDescription("Metric to monitor total CPU seconds"))
	if err != nil {
		return fmt.Errorf("error while setting up cpu seconds metric counter: %v", err)
	}
	memory, err := meter.Float64ObservableGauge("hypertrace.agent.memory", metric.WithDescription("Metric to monitor memory usage"))
	if err != nil {
		return fmt.Errorf("error while setting up memory metric counter: %v", err)
	}
	// Register the callback function for both cpu_seconds and memory observable gauges
	_, err = meter.RegisterCallback(
		func(ctx context.Context, result metric.Observer) error {
			sysMetrics, err := newSystemMetrics()
			if err != nil {
				return err
			}
			result.ObserveFloat64(cpuSeconds, sysMetrics.getCPU())
			result.ObserveFloat64(memory, sysMetrics.getMemory())
			return nil
		},
		cpuSeconds, memory,
	)
	if err != nil {
		log.Fatalf("failed to register callback: %v", err)
		return err
	}
	return nil
}
