package atb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const departuresURL = "https://atbapi.tar.io/api/v1/departures"

type myTime struct{ time.Time }

func (t *myTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	loc, _ := time.LoadLocation("Europe/Oslo")
	t.Time, err = time.ParseInLocation("2006-01-02T15:04:05.000", s, loc)
	return err
}

type Departure struct {
	Line                    string
	RegisteredDepartureTime myTime
	ScheduledDepartureTime  myTime
	Destination             string
	IsRealtimeData          bool
}

type Departures struct {
	URL                   string
	IsGoingTowardsCentrum bool
	Departures            []Departure
}

type BusStop struct {
	NodeID int
}

func (b *BusStop) GetDepartures() (*Departures, error) {
	url := departuresURL + "/" + strconv.Itoa(b.NodeID)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if r.StatusCode != 200 {
		return nil, fmt.Errorf("Server returned error code: %s", r.Status)
	}
	defer r.Body.Close()
	d := &Departures{}
	return d, json.NewDecoder(r.Body).Decode(d)
}
