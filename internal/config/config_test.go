package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadUsesDefaults(t *testing.T) {
	appConfig, err := Load(nil)
	require.NoError(t, err)
	require.Equal(t, DefaultConfigPath, appConfig.ConfigPath)
	require.Equal(t, DefaultEventsPath, appConfig.EventsPath)
}

func TestLoadUsesEnvironment(t *testing.T) {
	t.Setenv(EnvConfigPath, "env-config.json")
	t.Setenv(EnvEventsPath, "env-events")

	appConfig, err := Load(nil)
	require.NoError(t, err)
	require.Equal(t, "env-config.json", appConfig.ConfigPath)
	require.Equal(t, "env-events", appConfig.EventsPath)
}

func TestLoadFlagsOverrideEnvironment(t *testing.T) {
	t.Setenv(EnvConfigPath, "env-config.json")
	t.Setenv(EnvEventsPath, "env-events")

	appConfig, err := Load([]string{
		"-config", "flag-config.json",
		"-events", "flag-events",
	})
	require.NoError(t, err)
	require.Equal(t, "flag-config.json", appConfig.ConfigPath)
	require.Equal(t, "flag-events", appConfig.EventsPath)
}

func TestLoadRejectsPositionalArgs(t *testing.T) {
	_, err := Load([]string{"config.json", "events"})
	require.Error(t, err)
}

func TestLoadRejectsUnknownFlag(t *testing.T) {
	_, err := Load([]string{"-unknown", "value"})
	require.Error(t, err)
}
