package hardware

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmhobbs/martha/internal/configuration"
)

type Measurement struct {
	Timestamp time.Time
	Type configuration.DeviceType
	ID string
	Reading *Reading
	Error error
}

type Reading struct {
	CO2 *float64
	Humidity *float64
	Open *bool
	Pressure *float64
	Temperature *float64
}

type Sensor interface {
	Read() *Reading
}

func (r Reading) String() string {
	parts := []string{}

	if r.CO2 != nil {
		parts = append(parts, fmt.Sprintf("co2=%0.1f", *r.CO2))
	}
	if r.Humidity != nil {
		parts = append(parts, fmt.Sprintf("humidity=%0.1f%%", *r.Humidity))
	}
	if r.Open != nil {
		parts = append(parts, fmt.Sprintf("open=%t", *r.Open))
	}
	if r.Pressure != nil {
		parts = append(parts, fmt.Sprintf("pressure=%0.1f", *r.Pressure))
	}
	if r.Temperature != nil {
		parts = append(parts, fmt.Sprintf("temperature=%0.1fF", *r.Temperature))
	}

	return strings.Join(parts, " ")
}

func (m Measurement) String() string {
	parts := []string {
		"[" + m.Timestamp.Format(time.RFC3339) + " : " + string(m.Type) + "]",
	}

	var prefix string
	switch m.Type{
	case configuration.DeviceType(configuration.SENSOR_TYPE_EXTERNAL):
		prefix = "external"
	case configuration.DeviceType(configuration.SENSOR_TYPE_BME280):
		prefix = "i2c"
	case configuration.DeviceType(configuration.CONTROL_TYPE_NORMALLY_OPEN_RELAY):
		fallthrough
	case configuration.DeviceType(configuration.CONTROL_TYPE_NORMALLY_CLOSED_RELAY):
		prefix = "pin"
	case configuration.DeviceType(configuration.SENSOR_TYPE_SYNTHETIC):
		prefix = "synth"
	case configuration.DeviceType(configuration.SENSOR_TYPE_MAX6675):
		prefix = "spi"
	}

	parts = append(parts, "(" + prefix + "=" + m.ID + ")")
	if m.Reading != nil {
		parts = append(parts, m.Reading.String())
	}
	if m.Error != nil {
		parts = append(parts, fmt.Sprintf("error=%q", m.Error.Error()))
	}

	return strings.Join(parts, " ")
}