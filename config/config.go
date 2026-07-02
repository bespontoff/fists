package config

import "log/slog"

type Config struct {
	LogLevel slog.Level
}

func New() *Config {
	c := &Config{
		LogLevel: slog.LevelDebug,
	}
	return c
}
