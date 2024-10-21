//go:build !linux

package metrics

type noopMetrics struct{}

func newSystemMetrics() (systemMetrics, error) {
	return &noopMetrics{}, nil
}

func (nm *noopMetrics) getMemory() float64 {
	return 0
}

func (nm *noopMetrics) getCPU() float64 {
	return 0
}
