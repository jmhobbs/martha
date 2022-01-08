package device

import (
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
)

type MAX6675 struct {
	conn spi.Conn
}

func NewMAX6675(port spi.Port) (*MAX6675, error) {
	conn, err := port.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	return &MAX6675{conn}, nil
}

func (m *MAX6675) Read() (float64, error) {
	var (
		write []byte = []byte{0x00, 0x00}
		read  []byte = make([]byte, len(write))
		value uint16
	)

	if err := m.conn.Tx(write, read); err != nil {
		return 0.0, err
	}
	value = (uint16(read[0])<<5 | uint16(read[1])>>3)
	return float64(value) * 0.25, nil
}
