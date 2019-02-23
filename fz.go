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
	log      chan string
}

// Name return human-friendly name for worker for logs
func (f *Fz) Name() string {
	return "Floating Zendo (floatingzendo.org)"
}

// Init sets HTTP client and defines internal venue map
func (f *Fz) Init(client *Client, log chan string) error {
	f.client = client
	f.venueMap = make(map[string]*Venue)

	return nil
}

// Desc returns description for website crawled
func (f *Fz) Desc() string {
	return "Floating Zendo (floatingzendo.org)"
}

// Events hits floating zendo events page and returns slice of Event types
// http://www.floatingzendo.org/events/
func (f *Fz) Events() ([]*Event, error) {
	resp, err := f.client.Get("https://www.aczc.org/schedule/")
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
