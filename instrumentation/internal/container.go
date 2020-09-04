package internal

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

const (
	dockerPrefix   = "/docker/"
	kubepodsPrefix = "/kubepods/"
)

// ErrNotInContainerEnv is returned when the GetContainerID function is
// called in a non container environment
var ErrNotInContainerEnv = errors.New("not in a container environment")

func getContainerIDFromReader(f io.Reader) (string, error) {
	s := bufio.NewScanner(f)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}

		group := strings.SplitN(s.Text(), ":", 3)[2]
		if strings.HasPrefix(group, dockerPrefix) {
			return group[len(dockerPrefix):], nil
		} else if strings.HasPrefix(group, kubepodsPrefix) {
			return group[len(kubepodsPrefix):], nil
		}
	}

	return "", ErrNotInContainerEnv
}

// GetContainerID returns the container ID when in a containerized environment.
func GetContainerID() *string {
	f, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if containerID, err := getContainerIDFromReader(f); err == nil {
		return &containerID
	}

	return null
}
