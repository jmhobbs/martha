package provider

import "github.com/jmhobbs/martha/internal/configuration"

type Provider interface {
	Load() (*configuration.Config, error)
	Store(*configuration.Config) error
}
