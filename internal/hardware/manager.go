package hardware

import (
	"fmt"
	"time"

	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/jmhobbs/martha/internal/util/pointer"
	"github.com/jmhobbs/martha/pkg/device"

	"github.com/rs/zerolog"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

// todo: wrap device.* in something to conform to the Sensor interface

type i2cDevice struct {
	Bus i2c.BusCloser
	BME280 *device.BME280
}

type spiDevice struct {
	Port spi.PortCloser
	MAX6675 *device.MAX6675
}

type Manager struct {
	logger zerolog.Logger
	i2cDevices map[string]i2cDevice
	spiDevices map[string]spiDevice
	relays map[string]*device.Relay
}

func NewManager(parentLogger zerolog.Logger, config *configuration.Config) (*Manager, error) {
	_, err := host.Init()
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize host: %w", err)
	}

	var m *Manager = new(Manager)
	m.logger = parentLogger.With().Str("component", "hardware.Manager").Logger()
	return m, m.configure(config)
}

// todo: error capture/handling
func (m *Manager) Close() error {
	for _, dev := range m.i2cDevices {
		dev.Bus.Close()
	}
	for _, dev := range m.spiDevices {
		dev.Port.Close()
	}
	return nil	
}

// todo: error handling
func (m *Manager) configure(config *configuration.Config) error {
	for _, fc := range config.FruitingChambers {
		for _, sensor := range fc.Sensors {
			switch sensor.Type {
			case configuration.SENSOR_TYPE_BME280:
				bus, err := i2creg.Open(sensor.I2C)
				if err != nil {
					return err
				}
				m.i2cDevices[sensor.I2C] = i2cDevice{
					bus,
					device.NewBME280(bus),
				}
			case configuration.SENSOR_TYPE_MAX6675:
				port, err := spireg.Open(sensor.SPIPort)
				if err != nil {
					return err
				}
				dev, err := device.NewMAX6675(port)
				if err != nil {
					defer port.Close()
					return err
				}
				m.spiDevices[sensor.SPIPort] = spiDevice{
					port,
					dev,
				}
			// todo: do these belong here?
			case configuration.SENSOR_TYPE_SYNTHETIC:
			case configuration.SENSOR_TYPE_EXTERNAL:
			}
		}

		for _, control := range fc.Controls {
			m.relays[control.Pin] = device.NewRelay(
				gpioreg.ByName(control.Pin),
				control.Type == configuration.CONTROL_TYPE_NORMALLY_OPEN_RELAY,
			)
		}
	}

	return nil
}

func (m *Manager) ReadAll() []Measurement {
	var measurement Measurement

	measurements := []Measurement{}

	for bus, device := range m.i2cDevices {
		temperature, pressure, humidity, err := device.BME280.Read()

		measurement = Measurement{
			Timestamp: time.Now().UTC(),
			ID: bus,
		}

		if err != nil {
			measurement.Error = err
		} else {
			measurement.Reading = &Reading{
				Temperature: pointer.Float64(temperature),
				Pressure: pointer.Float64(pressure),
				Humidity: pointer.Float64(humidity),
			}
		}
		measurements = append(measurements, measurement)
	}

	for port, device := range m.spiDevices {
		temperature, err := device.MAX6675.Read()

		measurement = Measurement{
			Timestamp: time.Now().UTC(),
			ID: port,
		}

		if err != nil {
			measurement.Error = err
		} else {
			measurement.Reading = &Reading{
				Temperature: pointer.Float64(temperature),
			}
		}
		measurements = append(measurements, measurement)
	}

	for pin, device := range m.relays {
		measurements = append(measurements, Measurement{
			Timestamp: time.Now().UTC(),
			ID: pin,
			Reading: &Reading{
				Open: pointer.Bool(device.IsOpen()),
			},
		})
	}

	return measurements
}
