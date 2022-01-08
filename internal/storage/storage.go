package storage

import (
	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/rs/zerolog"
)

type StorageEngine struct {}

func New(logger zerolog.Logger, config *configuration.Config) *StorageEngine {
	return &StorageEngine{}
}