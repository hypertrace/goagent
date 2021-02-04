package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInt32Env(t *testing.T) {
	_, ok := getInt32Env("my_int")
	assert.False(t, ok)

	os.Setenv("my_int", "1")
	intVal, ok := getInt32Env("my_int")
	assert.True(t, ok)
	assert.Equal(t, int32(1), intVal)

	os.Setenv("my_int", "a")
	_, ok = getInt32Env("my_int")
	assert.False(t, ok)
}

func TestGetBoolEnv(t *testing.T) {
	_, ok := getBoolEnv("my_bool")
	assert.False(t, ok)

	os.Setenv("my_bool", "true")
	boolVal, ok := getBoolEnv("my_bool")
	assert.True(t, ok)
	assert.Equal(t, true, boolVal)

	os.Setenv("my_bool", "false")
	boolVal, ok = getBoolEnv("my_bool")
	assert.True(t, ok)
	assert.Equal(t, false, boolVal)

	os.Setenv("my_bool", "x")
	_, ok = getBoolEnv("my_bool")
	assert.False(t, ok)
}

func TestGetStringEnv(t *testing.T) {
	_, ok := getStringEnv("my_string")
	assert.False(t, ok)

	os.Setenv("my_string", "my_value")
	stringVal, ok := getStringEnv("my_string")
	assert.True(t, ok)
	assert.Equal(t, "my_value", stringVal)
}

func TestGetArrayStringEnv(t *testing.T) {
	_, ok := getArrayStringEnv("my_vowels")
	assert.False(t, ok)

	os.Setenv("my_vowels", "a,b")
	vals, ok := getArrayStringEnv("my_vowels")
	assert.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, vals)
}
