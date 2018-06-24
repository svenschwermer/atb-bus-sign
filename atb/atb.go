package atb

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const departuresURL = "https://atbapi.tar.io/api/v1/departures"

type myTime struct{ time.Time }

func (t *myTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	loc, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return err
	}
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

func (d *Departures) GetString(lines ...string) (out string) {
	sort.Slice(d.Departures, func(i, j int) bool {
		return d.Departures[i].RegisteredDepartureTime.Before(d.Departures[j].RegisteredDepartureTime.Time)
	})

	for i, line := range lines {
		if i > 0 {
			out += " | "
		}
		out += line + "->"
		found := 0
		for _, x := range d.Departures {
			if x.Line == line {
				d := time.Until(x.RegisteredDepartureTime.Time)
				if found == 0 {
					out += fmt.Sprintf("%s %.0fm", x.Destination, math.Round(d.Minutes()))
				} else {
					out += fmt.Sprintf(" (%.0fm)", math.Round(d.Minutes()))
				}
				found++
			}
			if found >= 2 {
				break
			}
		}
		if found == 0 {
			out += "?"
		}
	}
	return
}
