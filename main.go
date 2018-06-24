package main

import (
	"log"
	"time"

	"svenschwermer.de/bus-sign/atb"
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

	c := make(chan *frame.Frame)
	go display(max, c)
	query(c)
}

func query(c chan<- *frame.Frame) {
	maskinagentur := atb.BusStop{NodeID: 16011297}
	for {
		departures, err := maskinagentur.GetDepartures()
		if err != nil {
			log.Print(err)
		} else {
			str := departures.GetString("11")
			c <- frame.FromText(str, 32)
		}
		time.Sleep(10 * time.Second)
	}
}

func display(max *max7219.Device, c <-chan *frame.Frame) {
	for f := range c {
		for i := 32; i < f.Width(); i++ {
			sub := f.SubFrame(i-32, i)
			max.Frame(sub.ConcatenateLines())
			if i == 32 || i == f.Width()-1 {
				time.Sleep(1 * time.Second)
			} else {
				time.Sleep(75 * time.Millisecond)
			}
		}
	}
}
