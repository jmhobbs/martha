package hardware

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jmhobbs/martha/internal/configuration"
	"github.com/jmhobbs/martha/internal/util/pointer"
)

func Test_MeasurementToString(t *testing.T) {
	t.Run("normal read", func (t *testing.T) {
		m := Measurement{
			Timestamp: time.Unix(1234567890, 0).UTC(),
			Type: configuration.DeviceType(configuration.SENSOR_TYPE_BME280),
			ID: "I2C0",
			Reading: &Reading{
				Temperature: pointer.Float64(32.0),
				Humidity: pointer.Float64(84.539),
			},
		}
	
		expected := "[2009-02-13T23:31:30Z : BME280] (i2c=I2C0) humidity=84.5% temperature=32.0F"
		actual := fmt.Sprint(m)
	
		if expected != actual {
			t.Errorf("Measurement did not stringify correctly.\nexpected: %q\n  actual: %q", expected, actual)
		}
	})

	t.Run("error", func(t *testing.T) {
		m := Measurement{
			Timestamp: time.Unix(1234567890, 0).UTC(),
			Type: configuration.DeviceType(configuration.SENSOR_TYPE_MAX6675),
			ID: "SPI1",
			Error: errors.New("an error occurred when reading the device"),
		}

		expected := "[2009-02-13T23:31:30Z : MAX6675] (spi=SPI1) error=\"an error occurred when reading the device\""
		actual := fmt.Sprint(m)
	
		if expected != actual {
			t.Errorf("Measurement did not stringify correctly.\nexpected: %q\n  actual: %q", expected, actual)
		}
	})

	t.Run("error with measurement", func (t *testing.T) {
		m := Measurement{
			Timestamp: time.Unix(1234567890, 0).UTC(),
			Type: configuration.DeviceType(configuration.CONTROL_TYPE_NORMALLY_OPEN_RELAY),
			ID: "22",
			Reading: &Reading{
				Open: pointer.Bool(true),
			},
			Error: errors.New("short error"),
		}

		expected := "[2009-02-13T23:31:30Z : NORMALLY_OPEN_RELAY] (pin=22) open=true error=\"short error\""
		actual := fmt.Sprint(m)
	
		if expected != actual {
			t.Errorf("Measurement did not stringify correctly.\nexpected: %q\n  actual: %q", expected, actual)
		}
	})
}

func Test_ReadingToString(t *testing.T) {
	r := Reading {
		CO2: pointer.Float64(567.89),
		Humidity: pointer.Float64(76.99),
		Open: pointer.Bool(false),
		Pressure: pointer.Float64(11.3),
		Temperature: pointer.Float64(65.478),
	}

	expected := "co2=567.9 humidity=77.0% open=false pressure=11.3 temperature=65.5F"
	actual := fmt.Sprint(r)

	if expected != actual {
		t.Errorf("Reading did not stringify correctly.\nexpected: %q\n  actual: %q", expected, actual)
	}
}