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

	f := frame.New(8, 32)
	mod := [][]byte{{0x40, 0x00}, {0x20, 0x00}, {0x10, 0x00}, {0x08, 0x00}, {0x04, 0x00}, {0x02, 0x00}, {0x01, 0x00}, {0x00, 0x80}}
	for {
		f[0][0] = 0x80
		f[1][0] = 0x40
		f[2][0] = 0x20
		f[3][0] = 0x10
		f[4][0] = 0x08
		f[5][0] = 0x04
		f[6][0] = 0x02
		f[7][0] = 0x01
		for i := range f {
			f[i][3] = 0x00
		}
		max.Frame(f.ConcatenateLines())
		time.Sleep(25 * time.Millisecond)

		for i := 0; i < 24; i++ {
			f.Modify(i, i+9, mod)
			max.Frame(f.ConcatenateLines())
			time.Sleep(25 * time.Millisecond)
		}
	}
}
