//go:build !linux

package metrics

type noopMetrics struct{}

func newSystemMetrics() systemMetrics {
	return &noopMetrics{}
}

func (nm *noopMetrics) getMemory() (float64, error) {
	return 0, nil
}

func (nm *noopMetrics) getCPU() (float64, error) {
	return 0, nil
}

func (nm *noopMetrics) getCurrentMetrics() error {
	return nil
}
