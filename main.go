package main

import (
	"log"

	"golang.org/x/exp/io/spi"
)

var spiDeviceConfig = spi.Devfs{
	Dev:      "/dev/spidev1.0",
	Mode:     spi.Mode0,
	MaxSpeed: 10000000, // 10 MHz
}

func main() {
	spiDevice, err := spi.Open(&spiDeviceConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer spiDevice.Close()
}
