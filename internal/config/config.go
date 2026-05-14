package config

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	DefaultConfigPath = "config.json"
	DefaultEventsPath = "events"

	EnvConfigPath = "DUNGEON_CONFIG_PATH"
	EnvEventsPath = "DUNGEON_EVENTS_PATH"
)

type AppConfig struct {
	ConfigPath string
	EventsPath string
}

func Default() AppConfig {
	return AppConfig{
		ConfigPath: DefaultConfigPath,
		EventsPath: DefaultEventsPath,
	}
}

func Load(args []string) (AppConfig, error) {
	return FromFlags(args, FromEnv(Default()))
}

func FromEnv(base AppConfig) AppConfig {
	appConfig := base
	if value := os.Getenv(EnvConfigPath); value != "" {
		appConfig.ConfigPath = value
	}

	if value := os.Getenv(EnvEventsPath); value != "" {
		appConfig.EventsPath = value
	}

	return appConfig
}

func FromFlags(args []string, base AppConfig) (AppConfig, error) {
	appConfig := base
	flags := flag.NewFlagSet("dungeon-event-processor", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&appConfig.ConfigPath, "config", appConfig.ConfigPath, "path to dungeon config file")
	flags.StringVar(&appConfig.EventsPath, "events", appConfig.EventsPath, "path to events file")

	if err := flags.Parse(args); err != nil {
		return AppConfig{}, err
	}

	if flags.NArg() > 0 {
		return AppConfig{}, fmt.Errorf("usage: app [-config path] [-events path]")
	}

	return appConfig, nil
}
