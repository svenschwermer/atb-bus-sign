package atb_test

import (
	"testing"

	"svenschwermer.de/bus-sign/atb"
)

func TestATB(t *testing.T) {
	maskinagentur := atb.BusStop{NodeID: 16011297}
	d, err := maskinagentur.GetDepartures()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d.GetString("11", "19"))
}
