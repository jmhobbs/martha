package device

import (
	"sync"

	"periph.io/x/conn/v3/gpio"
)

// Relay is a concurrency safe way to access a relay and it's current state.
type Relay struct {
	sync.Mutex
	pin  gpio.PinIO
	isOpen bool
	openLevel gpio.Level
	closedLevel gpio.Level
}

func NewRelay(pin gpio.PinIO, normallyOpen bool) *Relay {
	closedLevel := gpio.High
	openLevel := gpio.Low
	if normallyOpen {
		closedLevel = gpio.Low
		openLevel = gpio.High
	}
	_ = pin.Out(openLevel) // Default OFF
	return &Relay{
		pin: pin,
		isOpen: true,
		openLevel: openLevel,
		closedLevel: closedLevel,
	}
}

func (r *Relay) Off() {
	r.Open()
}

func (r *Relay) Open() {
	r.Lock()
	defer r.Unlock()

	if !r.isOpen {
		_ = r.pin.Out(r.openLevel)
		r.isOpen = true
	}
}

func (r *Relay) On() {
	r.Close()
}

func (r *Relay) Close() {
	r.Lock()
	defer r.Unlock()

	if r.isOpen {
		_ = r.pin.Out(r.closedLevel)
		r.isOpen = false
	}
}

func (r *Relay) Toggle() {
	state := r.IsOpen()

	if state {
		r.Close()
	} else {
		r.Open()
	}
}

func (r *Relay) IsOpen() bool {
	r.Lock()
	defer r.Unlock()

	return r.isOpen
}

func (r *Relay) IsClosed() bool {
	return !r.IsOpen()
}
