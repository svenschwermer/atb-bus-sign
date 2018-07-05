// +build sim

package max7219

import (
	"fmt"
)

type Device struct{}

func Open(dev string, cascade int) (*Device, error) {
	return &Device{}, nil
}

func (d *Device) Close() error {
	return nil
}

func (d *Device) Init() error {
	return nil
}

// Frame draws the passed frame in the terminal after clearing it.
func (d *Device) Frame(frame fmt.Stringer) error {
	fmt.Print("\033[H\033[2J", frame.String())
	return nil
}
