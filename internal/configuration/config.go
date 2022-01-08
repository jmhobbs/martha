package configuration

import "github.com/jmhobbs/martha/internal/util"

type Sensor struct {
	Name string `yaml:"name,omitempty"`
	Type SensorType `yaml:"type"`
	// Required for I2C devices (BME280)
	I2C string `yaml:"i2c_bus,omitempty"`
	// Required for SPI devices ()
	SPIPort string `yaml:"spi_port,omitempty"`
	// Unique key used for updating EXTERNAL value via API
	ExternalKey string `yaml:"external_key,omitempty"`
	// Extra fields for configuring a fully synthetic 
	Synthetic Synthetic `yaml:"synthetic,omitempty"`
}

type SyntheticValues struct {
	// Value when not triggered
	Idle float64 `yaml:"idle"`
	// Value when triggered
	Active float64 `yaml:"active"`
}

type Synthetic struct {
	ID string `yaml:"id"`
	// How often to trigger the synthetic state
	Interval string `yaml:"interval"`
	// How long to trigger the synthetic state
	Duration string `yaml:"duration"`
	// What parameter to affect
	Effect ControlEffect `yaml:"effect"`
	// Values for parameter
	Values SyntheticValues `yaml:"values"`
}

type Control struct {
	Name string `yaml:"name,omitempty"`
	Type ControlType `yaml:"type"`
	Effect ControlEffect `yaml:"effect"`
	Pin string `yaml:"pin"`
}
type Profile struct {
	ID string `yaml:"id"`
	Name string `yaml:"name"`
	// todo
}

type FruitingChamber struct {
	ID string `yaml:"id"`
	Name string `yaml:"name"`
	Sensors []Sensor `yaml:"sensors"`
	Controls []Control `yaml:"controls,omitempty"`
	Profile string `yaml:"profile"`
}

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Profiles []Profile `yaml:"profiles"`
	SamplingInterval string `yaml:"sampling_interval"`
	FruitingChambers []FruitingChamber `yaml:"fruiting_chambers"`
}

func (c *Config) Validate() error {
	// todo
	return nil
}

func (c *Config) HydrateDefaults() {
	if c.Server.Address == "" {
		c.Server.Address = "127.0.0.1:8080"
	}

	if c.SamplingInterval == "" {
		c.SamplingInterval = "30s"
	}

	for i, profile := range c.Profiles {
		if profile.ID == "" {
			c.Profiles[i].ID = util.MakeId()
		}
	}

	for i, chamber := range c.FruitingChambers {
		if chamber.ID == "" {
			c.FruitingChambers[i].ID = util.MakeId()
		}
		for j, sensor := range chamber.Sensors {
			if sensor.Type == SENSOR_TYPE_SYNTHETIC && sensor.Synthetic.ID == "" {
				c.FruitingChambers[i].Sensors[j].Synthetic.ID = util.MakeId()
			}
		}
	}
}
