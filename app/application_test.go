package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApplicationShouldPanicIfDBTypeNotSet(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Should have panic")
		}
	}()

	NewApplication("")

}

func TestGetEnvVarShouldReturnEnvironmentVariableValue(t *testing.T) {
	t.Parallel()

	envKeyName := "TEST_VARIABLE"
	envValue := "testValue"

	os.Setenv(envKeyName, envValue)
	value := getEnvVar(envKeyName, "")

	assert.True(t, value == envValue, "Environment variable value sould match")
}

func TestGetEnvVarShouldReturnDefaultEnvironmentVariableValue(t *testing.T) {
	t.Parallel()

	envKeyName := "TEST_GET_NON_EXISTING_ENV_VAR"
	defaultValue := "THIS_IS_A_DEFAULT_VALUE"

	os.Unsetenv(envKeyName)
	value := getEnvVar(envKeyName, defaultValue)

	assert.True(t, value == defaultValue, "Should return default value")
}
