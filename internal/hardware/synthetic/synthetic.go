package synthetic

import (
	"fmt"
	"time"

	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/jmhobbs/martha/internal/hardware"
	"github.com/jmhobbs/martha/internal/util/pointer"
)

/*
todo: this doesn't provide an exact timing if reads are not constant
i.e. if it's been a multiple of Interval + Duration with no read,
any read will set the state active.

An alternative implementation would be to start a ticker, and state keep
with a mutex.  Then a read will always be accurate.  I'm uncertain is this is
more desirable than the current implementation.
*/

type Device struct {
	interval time.Duration
	duration time.Duration
	activeInterval time.Duration
	effect configuration.ControlEffect
	idleValue float64
	activeValue float64
	lastActive time.Time
}

func New(parameters configuration.Synthetic) (*Device, error) {
	interval, err := time.ParseDuration(parameters.Interval)
	if err != nil {
		return nil, fmt.Errorf("unable to parse interval for synthetic sensor %q: %w", parameters.ID, err)
	}
	duration, err := time.ParseDuration(parameters.Duration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration for synthetic sensor %q: %w", parameters.ID, err)
	}

	return &Device{
		interval: interval,
		duration: duration,
		activeInterval: interval + duration,
		effect: parameters.Effect,
		idleValue: parameters.Values.Idle,
		activeValue: parameters.Values.Active,
		lastActive: time.Now().Add(duration),
	}, nil
}

func (s *Device) Read() *hardware.Reading {
	return s.readForTime(time.Now())
}

// Calls to readForTime must be in chronological order for state keeping.
func (s *Device) readForTime(now time.Time) *hardware.Reading {
	elapsed := now.Sub(s.lastActive)
	active := elapsed < s.duration
	if elapsed >= s.activeInterval {
		s.lastActive = now
		active = true
	}
	reading := hardware.Reading{}

	value := s.idleValue
	if active {
		value = s.activeValue
	}

	switch s.effect {
	case configuration.CONTROL_EFFECT_TYPE_CO2:
		reading.CO2 = pointer.Float64(value)
	case configuration.CONTROL_EFFECT_TYPE_HUMIDITY:
		reading.Humidity = pointer.Float64(value)
	case configuration.CONTROL_EFFECT_TYPE_TEMPERATURE:
		reading.Temperature = pointer.Float64(value)
	}

	return &reading
}
