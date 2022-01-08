package app

import (
	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/jmhobbs/martha/internal/configuration/provider"
	"github.com/jmhobbs/martha/internal/hardware"
	"github.com/jmhobbs/martha/internal/sampler"

	"github.com/rs/zerolog"
)

func Run(debug bool, configPath string) {
	logger := makeLogger(debug)

	config, configProvider, err := loadConfigFromFile(logger, configPath)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("Starting Hardware Manager")
	manager, err := hardware.NewManager(logger, config)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	defer manager.Close()

	logger.Info().Msg("Starting Sampling System")
	sampler, err := sampler.New(logger, config, manager)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	sampler.Run()

	defer func (config *configuration.Config, configProvider provider.Provider, logger zerolog.Logger) {
		logger.Info().Msg("Shutting down...")
		if err := configProvider.Store(config); err != nil {
			logger.Error().Err(err).Msg("Error writing config file at shutdown.")
		}
	}(config, configProvider, logger)
}
