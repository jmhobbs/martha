package sampler

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/jmhobbs/martha/internal/hardware"
	"github.com/jmhobbs/martha/internal/hardware/synthetic"
	"github.com/jmhobbs/martha/internal/util"
	"github.com/rs/zerolog"
)

type SubscriptionRequest struct {
	All bool
	I2CBuses []string
	SPIPorts []string
	Pins []string
}

type subscription struct {
	matcher SubscriptionRequest
	channel chan hardware.Measurement
}

type Sampler struct {
	sync.RWMutex
	logger zerolog.Logger
	manager *hardware.Manager
	interval time.Duration
	subscriptions map[string]subscription
	synthetics map[string]*synthetic.Device
}

func New(parentLogger zerolog.Logger, config *configuration.Config, manager *hardware.Manager) (*Sampler, error) {
	interval, err := time.ParseDuration(config.SamplingInterval)
	if err != nil {
		return nil, fmt.Errorf("Could not parse sampling_interval in configuration: %w", err)
	}

	synthetics := make(map[string]*synthetic.Device)
	for _, chamber := range config.FruitingChambers {
		for _, sensor := range chamber.Sensors {
			if sensor.Type == configuration.SENSOR_TYPE_SYNTHETIC {
				dev, err := synthetic.New(sensor.Synthetic)
				if err != nil {
					return nil, err
				}
				synthetics[sensor.Synthetic.ID] = dev
			}
		}
	}

	return &Sampler{
		logger: parentLogger.With().Str("component", "sampler.Sampler").Logger(),
		manager: manager,
		interval: interval,
		subscriptions: make(map[string]subscription),
		synthetics: synthetics,
	}, nil
}

func (s *Sampler) Stop() error {
	return errors.New("not implemented")
}

func (s *Sampler) Run() {
	for {
		<- time.After(s.interval)
		for id, device := range s.synthetics {
			reading := device.Read()
			s.publish(hardware.Measurement{
				Timestamp: time.Now(),
				Type: configuration.DeviceType(configuration.SENSOR_TYPE_SYNTHETIC),
				ID: id,
				Reading: reading,
			})
		}

		// todo: externals

		for _, measurement := range s.manager.ReadAll() {
			s.publish(measurement)
		}
	}
}

func (s *Sampler) Subscribe(request SubscriptionRequest) (string, chan hardware.Measurement) {
	s.Lock()
	defer s.Unlock()

	id := util.MakeIdFull(16, 4)
	s.subscriptions[id] = subscription{
		request,
		make(chan hardware.Measurement),
	}
	return id, s.subscriptions[id].channel
}

func (s *Sampler) Unsubscribe(id string) {
	s.Lock()
	defer s.Unlock()

	close(s.subscriptions[id].channel)
	delete(s.subscriptions, id)
}

func (s *Sampler) publish(m hardware.Measurement) {
	// todo: map to subscription channels
	s.logger.Info().Msg(m.String())
}