package main

import (
	"testing"

	"svenschwermer.de/bus-sign/atb"
	"svenschwermer.de/bus-sign/frame"
)

func TestIntegration(t *testing.T) {
	maskinagentur := atb.BusStop{NodeID: 16011297}
	departures, err := maskinagentur.GetDepartures()
	if err != nil {
		t.Error(err)
	}
	str := departures.GetString("11", "19") + " | "

	f := frame.FromText(str, 32)
	t.Logf("\"%s\":\n%s", str, f.CompactString())
}
