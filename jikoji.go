package zenfo

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Jikoji crawls jikoji.org, satisfies Worker interface
type Jikoji struct {
	client   *Client
	venueMap map[string]*Venue
	log      chan string
}

// Name returns human-friendly name for worker logs
func (j *Jikoji) Name() string {
	return "Jikoji (jikoji.org)"
}

// Init sets HTTP client and defines internal venue map
func (j *Jikoji) Init(client *Client, log chan string) error {
	j.client = client
	j.venueMap = make(map[string]*Venue)
	j.log = log

	j.venueMap["jikoji"] = &Venue{
		Name:    "Jikoji Zen Center",
		Addr:    "12100 Skyline Blvd, Los Gatos, CA",
		Phone:   "+1 (408) 741-9562",
		Email:   "info@jikoji.org",
		Lat:     37.2728165,
		Lng:     -122.1466097,
		Website: "https://www.jikoji.org",
	}

	return nil
}

// Desc returns description for website crawled
func (j *Jikoji) Desc() string {
	return "Jikoji (jikoji.org)"
}

// Events hits jikoji events page and returns slice of Event types
// https://www.jikoji.org/jikoji-events
func (j *Jikoji) Events() ([]*Event, error) {

	u := "https://www.jikoji.org/jikoji-events"

	resp, err := j.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var (
		domErr error
		events []*Event
	)

	doc.Find("div.eventlist eventlist--upcoming").EachWithBreak(func(_ int, div *goquery.Selection) bool {

		name := div.Find(".eventlist-title").First().Text()
		desc := div.Find(".eventlist-excerpt").First().Text()

		e := &Event{
			URL:   u,
			Name:  name,
			Desc:  desc,
			Venue: j.venueMap["jikoji"],
		}
		events = append(events, e)

		return true
	})
	j.log <- fmt.Sprintf("Found %d total events", len(events))
	return events, domErr
}
