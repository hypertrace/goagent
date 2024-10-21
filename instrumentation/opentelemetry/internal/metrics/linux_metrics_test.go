//go:build linux

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// Mock data simulating /proc/[pid]/stat content
	mockData = "85 (md) I 2 0 0 0 -1 698880 0 0 0 0 100 200 300 400 0 -20 1 0 223 0 550 18437615 0 0 0 0 0 0 0 2647 0 1 0 0 17 1 0 0 0 0 0 0 0 0 0 0 0 0 0"
)

func TestParseProcStatFile(t *testing.T) {
	lm := &linuxMetrics{}
	procFilepath := "/proc/123/stat"

	stat, err := lm.parseProcStatFile([]byte(mockData), procFilepath)
	assert.NoError(t, err, "unexpected error while parsing proc stat file")

	assert.Equal(t, 100.0, stat.utime, "utime does not match expected value")
	assert.Equal(t, 200.0, stat.stime, "stime does not match expected value")
	assert.Equal(t, 300.0, stat.cutime, "cutime does not match expected value")
	assert.Equal(t, 400.0, stat.cstime, "cstime does not match expected value")
	assert.Equal(t, 550.0, stat.rss, "rss does not match expected value")
}
