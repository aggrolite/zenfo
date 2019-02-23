package zenfo

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

// http://www.floatingzendo.org/events/

// Fz crawls floatingzendo.org, satisfies Worker interface
type Fz struct {
	venueMap map[string]*Venue
	client   *Client
}

// Init sets HTTP client and defines internal venue map
func (s *Fz) Init(client *Client) error {
	s.client = client
	s.venueMap = make(map[string]*Venue)

	return nil
}

// Desc returns description for website crawled
func (s *Fz) Desc() string {
	return "Floating Zendo (floatingzendo.org)"
}

// Events hits floating zendo events page and returns slice of Event types
// http://www.floatingzendo.org/events/
func (s *Fz) Events() ([]*Event, error) {
	resp, err := s.client.Get("https://www.aczc.org/schedule/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var (
		//domErr error
		events []*Event
	)
	doc.Find("div.entry-content table tr").EachWithBreak(func(_ int, div *goquery.Selection) bool {
		log.Fatal(div)
		return true
	})
	return events, nil
}
