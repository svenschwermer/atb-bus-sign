package max7219

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/io/spi"
)

const (
	AddressNoOp        = byte(0x0)
	AddressDigit0      = byte(0x1)
	AddressDigit1      = byte(0x2)
	AddressDigit2      = byte(0x3)
	AddressDigit3      = byte(0x4)
	AddressDigit4      = byte(0x5)
	AddressDigit5      = byte(0x6)
	AddressDigit6      = byte(0x7)
	AddressDigit7      = byte(0x8)
	AddressDecodeMode  = byte(0x9)
	AddressIntensity   = byte(0xA)
	AddressScanLimit   = byte(0xB)
	AddressShutdown    = byte(0xC)
	AddressDisplayTest = byte(0xF)
)

type Device struct {
	spiDevice *spi.Device
	cascade   int
}

func Open(dev string, cascade int) (*Device, error) {
	if cascade < 1 {
		return nil, fmt.Errorf("Minimum one MAX7219 must be cascaded")
	}
	spiDeviceConfig := spi.Devfs{
		Dev:      dev,
		Mode:     spi.Mode0,
		MaxSpeed: 1000000, // 1 MHz
	}
	spiDevice, err := spi.Open(&spiDeviceConfig)
	if err != nil {
		return nil, err
	}
	if err := spiDevice.SetBitsPerWord(8); err != nil {
		return nil, err
	}
	if err := spiDevice.SetBitOrder(spi.MSBFirst); err != nil {
		return nil, err
	}
	return &Device{spiDevice, cascade}, nil
}

func (d *Device) Close() error {
	return d.spiDevice.Close()
}

func (d *Device) WriteToAll(address, data byte) error {
	w := bytes.Repeat([]byte{address, data}, d.cascade)
	return d.spiDevice.Tx(w, nil)
}

func (d *Device) Init() error {
	sequence := []struct{ a, d byte }{
		{AddressDecodeMode, 0x00},  // No decoding
		{AddressIntensity, 0x01},   // Minimal intensity
		{AddressScanLimit, 0x07},   // Scan all digits (0..7)
		{AddressDisplayTest, 0x00}, // Disable display test
		{AddressShutdown, 0x01},    // Normal operation
	}
	for _, x := range sequence {
		if err := d.WriteToAll(x.a, x.d); err != nil {
			return err
		}
	}
	return nil
}

func (d *Device) Line(line int, patterns ...byte) error {
	if line < 0 || line > 7 {
		return fmt.Errorf("Line out of valid range 0..7")
	}
	if len(patterns) != d.cascade {
		return fmt.Errorf("Number of patterns must match cascade length")
	}
	w := make([]byte, 0, 2*d.cascade)
	for i := 0; i < d.cascade; i++ {
		w = append(w, byte(line+1), patterns[i])
	}
	return d.spiDevice.Tx(w, nil)
}

// Frame draws the passed concatenated lines (MSB first)
func (d *Device) Frame(data []byte) error {
	if len(data) != 8*d.cascade {
		return fmt.Errorf("Frame data doesn't fit cascade dimensions")
	}
	for line := 0; line < 8; line++ {
		start := line * d.cascade
		lineData := data[start : start+d.cascade]
		if err := d.Line(line, lineData...); err != nil {
			return err
		}
	}
	return nil
}
