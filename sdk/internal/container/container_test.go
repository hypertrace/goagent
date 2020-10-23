package container

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerContainerIDCanBeObtainedFromCGroups(t *testing.T) {
	testCgroups := []string{
		`2:cpu:/docker/ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f
	1:cpuset:/docker/ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f`,
		`2:cpu:/kubepods/ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f
	1:cpuset:/kubepods/ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f`,
	}

	for _, cgroup := range testCgroups {
		containerID, err := getContainerIDFromReader(bytes.NewBufferString(cgroup))
		assert.Nil(t, err)
		assert.Equal(t, "ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f", containerID)
	}
}

func TestDockerContainerIDCannotBeObtainerInNonContainerisedEnvs(t *testing.T) {
	_, err := getContainerIDFromReader(bytes.NewBufferString(`2:cpu:/xyz/ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f
	1:cpuset:/xyz/ba4024e95abb12affe2b0f56ff86536d0abad7e95b09b591b03e6670dd0b5e5f`))
	assert.Equal(t, ErrNotInContainerEnv, err)
}
