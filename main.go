package main

import (
	"log"
	"time"

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

	f := frame.New(8, 64)
	f.Text(0, "abcdefghij")
	for {
		for i := 0; i <= 32; i++ {
			sub := f.SubFrame(i, i+31)
			max.Frame(sub.ConcatenateLines())
			time.Sleep(100 * time.Millisecond)
		}
	}
}
