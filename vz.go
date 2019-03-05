package zenfo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Vz crawls villagezendo.org, satisfies Worker interface
type Vz struct {
	venueMap map[string]*Venue
	client   *Client
	log      chan string
}

type villageEventJSON struct {
	Name  string    `json:"title"`
	Desc  string    `json:"description"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	URL   string    `json:"url"`
	Venue string    `json:"venue_slug"`
}

// Name return human-friendly name for worker for logs
func (v *Vz) Name() string {
	return "Village Zendo (villagezendo.org)"
}

// Init sets HTTP client and defines internal venue map
func (v *Vz) Init(client *Client, log chan string) error {

	v.client = client
	v.venueMap = make(map[string]*Venue)
	v.log = log

	v.venueMap["village-zendo"] = &Venue{
		Name:    "Village Zendo",
		Addr:    "588 Broadway, Suite 1108, New York, NY 10012",
		Phone:   "+1 (212) 340-4656",
		Email:   "info@villagezendo.org",
		Lat:     40.724682,
		Lng:     -73.997087,
		Website: "https://villagezendo.org",
	}
	v.log <- "Initialized!"

	return nil
}

// Desc returns description for website crawled
func (v *Vz) Desc() string {
	return "Village Zendo (villagezendo.org)"
}

// Events hits JSON API and returns slice of Event types
func (v *Vz) Events() ([]*Event, error) {
	u := "https://villagezendo.org/wp-admin/admin-ajax.php?action=eventorganiser-fullcal&start=2019-02-24&timeformat=g%3Ai+a"
	resp, err := v.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var events []villageEventJSON
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, errors.New("No events found")
	}

	var final []*Event
	for _, e := range events {

		if e.End.Before(e.Start) {
			e.End = e.Start
		}

		venue, ok := v.venueMap[e.Venue]
		if !ok {
			return nil, fmt.Errorf("Failed to match venue for '%s' - event=%+v", e.Venue, e)
		}

		finalEvent := &Event{
			Name:  e.Name,
			Desc:  e.Desc,
			Start: e.Start,
			End:   e.End,
			URL:   u,
			Venue: venue,
		}
		v.log <- fmt.Sprintf("Found event: %s", e.Name)

		final = append(final, finalEvent)
	}
	v.log <- fmt.Sprintf("Found %d total events", len(final))

	return final, nil
}
