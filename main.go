package main

import (
	"log"

	"svenschwermer.de/bus-sign/frame"
	"svenschwermer.de/bus-sign/max7219"
)

/*
func msgRequestCode(n uint32) uintptr {
	return uintptr(0x80006B00 + (n * 0x200000))
	                ---
}
*/

func main() {
	max, err := max7219.Open("/dev/spidev0.1", 4)
	if err != nil {
		log.Fatal(err)
	}
	defer max.Close()

	if err := max.Init(); err != nil {
		log.Fatal("Failed to initialize MAX7219: ", err)
	}

	f := frame.New(8, 32)
	f.Text(0, "11 ~5")
	max.Frame(f.ConcatenateLines())
}
