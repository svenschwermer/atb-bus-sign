package main

import (
	"log"
	"time"

	"svenschwermer.de/bus-sign/atb"
	"svenschwermer.de/bus-sign/frame"
	"svenschwermer.de/bus-sign/max7219"
)

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
		var str string
		if err != nil {
			log.Print(err)
			str = "ERROR: " + err.Error()
		} else {
			str = departures.GetString("11", "19")
		}
		str += " | "
		c <- frame.FromText(str, 32)
		time.Sleep(10 * time.Second)
	}
}

func display(max *max7219.Device, c <-chan *frame.Frame) {
	f := frame.FromText("waiting...  ", 32)
	for {
		select {
		case newFrame, ok := <-c:
			if !ok {
				return
			}
			f = newFrame
		default:
		}
		for i := 0; i < f.Width(); i++ {
			sub := f.SubFrame(i, i+32)
			max.Frame(sub.ConcatenateLines())
			time.Sleep(50 * time.Millisecond)
		}
	}
}
