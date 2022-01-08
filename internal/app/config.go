package app

import (
	"fmt"
	"path/filepath"

	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/jmhobbs/martha/internal/configuration/provider"

	"github.com/rs/zerolog"
)


func loadConfigFromFile(logger zerolog.Logger, path string) (*configuration.Config, provider.Provider, error) {
	absoluteConfigPath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to determine configuration file path: %w", err)
	}

	configProvider := provider.NewFileProvider(logger, absoluteConfigPath)
	config, err := configProvider.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't load configuration file: %w", err)
	}

	return config, configProvider, nil
}