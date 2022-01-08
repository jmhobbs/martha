package device

import (
	"time"

	"periph.io/x/conn/v3/i2c"
)

type BME280 struct {
	dev                 *i2c.Dev
	oversample_temp     uint8
	oversample_pressure uint8
	oversample_humidity uint8
	sampling_mode       uint8
}

func NewBME280(bus i2c.Bus) *BME280 {
	return &BME280{
		&i2c.Dev{Addr: 0x76, Bus: bus},
		2,
		2,
		2,
		1,
	}
}

func (b *BME280) Read() (temperature, pressure, humidity float64, err error) {
	// 1. Set humidity mode 0xF2 and oversample
	if _, err = b.dev.Write([]byte{0xF2, 0x02}); err != nil {
		return
	}

	// 2. Set to forced mode 0xF4
	var control uint8 = b.oversample_temp<<5 | b.oversample_pressure<<2 | b.sampling_mode
	if _, err = b.dev.Write([]byte{0xF4, byte(control)}); err != nil {
		return
	}

	// 3. Get compensation values
	var compensations *bm280compensations

	compensations, err = b.readCompensations()
	if err != nil {
		return
	}

	// 4. Wait for measurements
	wait_time := 1.25 + (2.3 * float64(b.oversample_temp)) + ((2.3 * float64(b.oversample_pressure)) + 0.575) + ((2.3 * float64(b.oversample_humidity)) + 0.575)
	time.Sleep(time.Duration(wait_time) * time.Millisecond)

	// 5. Read registers
	data := make([]byte, 8)
	if err = b.dev.Tx([]byte{0xF7}, data); err != nil {
		return
	}

	var (
		pressureRaw    uint32 = (uint32(data[0]) << 12) | (uint32(data[1]) << 4) | (uint32(data[2]) >> 4)
		temperatureRaw uint32 = (uint32(data[3]) << 12) | (uint32(data[4]) << 4) | (uint32(data[5]) >> 4)
		humidityRaw    uint16 = (uint16(data[6]) << 8) | uint16(data[7])
	)

	// 6. Compensate and return
	temperature, pressure, humidity = compensations.Compensate(temperatureRaw, pressureRaw, humidityRaw)
	return
}

type bm280compensations struct {
	// Humidity
	dig_T1 uint16
	dig_T2 int16
	dig_T3 int16
	// Pressure
	dig_P1 uint16
	dig_P2 int16
	dig_P3 int16
	dig_P4 int16
	dig_P5 int16
	dig_P6 int16
	dig_P7 int16
	dig_P8 int16
	dig_P9 int16
	// Humidity
	dig_H1 uint8
	dig_H2 int16
	dig_H3 uint8
	dig_H4 int16
	dig_H5 int16
	dig_H6 int8
}

func (b *BME280) readCompensations() (*bm280compensations, error) {
	c0 := make([]byte, 24)
	if err := b.dev.Tx([]byte{0x88}, c0); err != nil {
		return nil, err
	}
	c1 := make([]byte, 1)
	if err := b.dev.Tx([]byte{0xA1}, c1); err != nil {
		return nil, err
	}
	c2 := make([]byte, 7)
	if err := b.dev.Tx([]byte{0xE1}, c2); err != nil {
		return nil, err
	}

	compensation := bm280compensations{
		dig_T1: uint16(c0[1])<<8 | uint16(c0[0]),
		dig_T2: int16(c0[3])<<8 | int16(c0[2]),
		dig_T3: int16(c0[5])<<8 | int16(c0[4]),
		dig_P1: uint16(c0[7])<<8 | uint16(c0[6]),
		dig_P2: int16(c0[9])<<8 | int16(c0[8]),
		dig_P3: int16(c0[11])<<8 | int16(c0[10]),
		dig_P4: int16(c0[13])<<8 | int16(c0[12]),
		dig_P5: int16(c0[15])<<8 | int16(c0[14]),
		dig_P6: int16(c0[17])<<8 | int16(c0[16]),
		dig_P7: int16(c0[19])<<8 | int16(c0[18]),
		dig_P8: int16(c0[21])<<8 | int16(c0[20]),
		dig_P9: int16(c0[23])<<8 | int16(c0[22]),
		dig_H1: uint8(c1[0]),
		dig_H2: int16(c2[1])<<8 | int16(c2[0]),
		dig_H3: uint8(c2[2]),
		dig_H4: (int16(c2[3]) << 4) | (int16(c2[4]) & 0x0F),
		dig_H5: (int16(c2[5]) << 4) | (int16(c2[4]) >> 4 & 0x0F),
		dig_H6: int8(c2[6]),
	}

	return &compensation, nil
}

func (c *bm280compensations) Compensate(temperatureRaw, pressureRaw uint32, humidityRaw uint16) (float64, float64, float64) {
	// Temperature compensation
	tvar1 := (((temperatureRaw >> 3) - uint32(c.dig_T1<<1)) * uint32(c.dig_T2)) >> 11
	tvar2 := (((((temperatureRaw >> 4) - uint32(c.dig_T1)) * ((temperatureRaw >> 4) - uint32(c.dig_T1))) >> 12) * uint32(c.dig_T3)) >> 14
	t_fine := tvar1 + tvar2
	temperature := float64(((t_fine * 5) + 128) >> 8)

	// Pressure compensation
	pressure := 0.0

	pvar1 := float64(t_fine)/2.0 - 64000.0
	pvar2 := pvar1 * pvar1 * float64(c.dig_P6) / 32768.0
	pvar2 = pvar2 + pvar1*float64(c.dig_P5)*2.0
	pvar2 = pvar2/4.0 + float64(c.dig_P4)*65536.0
	pvar1 = (float64(c.dig_P3)*pvar1*pvar1/524288.0 + float64(c.dig_P2)*pvar1) / 524288.0
	pvar1 = (1.0 + pvar1/32768.0) * float64(c.dig_P1)
	if pvar1 != 0 {
		pressure = 1048576.0 - float64(pressureRaw)
		pressure = ((pressure - pvar2/4096.0) * 6250.0) / pvar1
		pvar1 = float64(c.dig_P9) * pressure * pressure / 2147483648.0
		pvar2 = pressure * float64(c.dig_P8) / 32768.0
		pressure = pressure + (pvar1+pvar2+float64(c.dig_P7))/16.0
	}

	// Humidity compensation
	humidity := float64(t_fine) - 76800.0
	humidity = (float64(humidityRaw) - (float64(c.dig_H4)*64.0 + float64(c.dig_H5)/16384.0*humidity)) * (float64(c.dig_H2) / 65536.0 * (1.0 + float64(c.dig_H6)/67108864.0*humidity*(1.0+float64(c.dig_H3)/67108864.0*humidity)))
	humidity = float64(humidity) * (1.0 - float64(c.dig_H1)*float64(humidity)/524288.0)

	return temperature / 100.0, pressure / 100.0, humidity
}
