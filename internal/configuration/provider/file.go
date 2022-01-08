package provider

import (
	"os"

	"github.com/jmhobbs/martha/internal/configuration"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type fileProvider struct {
	logger zerolog.Logger
	path string
}

func NewFileProvider(parentLogger zerolog.Logger, path string) Provider {
	return &fileProvider{
		logger: parentLogger.With().Str("component", "configuration.FileProvider").Logger(),
		path: path,
	}
}

func (p *fileProvider) Load() (*configuration.Config, error) {
	p.logger.Info().Str("path", p.path).Msg("Loading configuration file.")
	f, err := os.Open(p.path)
	if err != nil {
		p.logger.Debug().Err(err)
		if os.IsNotExist(err) {
			p.logger.Info().Msg("Configuration file does not exist, returning default configuration.")
			return configuration.DefaultConfig(), nil
		}
		return nil, err
	}
	defer f.Close()

	var config configuration.Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		p.logger.Debug().Err(err)
		return nil, err
	}

	config.HydrateDefaults()

	return &config, config.Validate()
}

func (p *fileProvider) Store(config *configuration.Config) error {
	p.logger.Info().Str("path", p.path).Msg("Storing configuration file.")
	f, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewEncoder(f).Encode(config)
}