package zenfo

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// http://www.floatingzendo.org/events/

// Fz crawls floatingzendo.org, satisfies Worker interface
type Fz struct {
	venueMap map[string]*Venue
	client   *Client
	log      chan string
}

// Name returns human-friendly name for worker logs
func (f *Fz) Name() string {
	return "Floating Zendo (floatingzendo.org)"
}

// Init sets HTTP client and defines internal venue map
func (f *Fz) Init(client *Client, log chan string) error {
	f.client = client
	f.venueMap = make(map[string]*Venue)
	f.log = log

	f.venueMap["1041 morse st., san jose"] = &Venue{
		Name:    "Floating Zendo - San Jose Friends Meeting House",
		Addr:    "1041 Morse St, San Jose, CA 95126",
		Email:   "secretary@floatingzendo.org",
		Lat:     37.341372,
		Lng:     121.928258,
		Website: "http://www.floatingzendo.org",
	}
	f.venueMap["jikoji"] = &Venue{
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
func (f *Fz) Desc() string {
	return "Floating Zendo (floatingzendo.org)"
}

// Events hits floating zendo events page and returns slice of Event types
// http://www.floatingzendo.org/events/
func (f *Fz) Events() ([]*Event, error) {

	u := "http://www.floatingzendo.org/events/"

	resp, err := f.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	year := "2019" // doc.Find("div.entry-content blockquote em strong").First().Text()
	if year == "" {
		return nil, errors.New("Failed to extract year")
	}

	var (
		domErr error
		events []*Event
		re     = regexp.MustCompile(`\([^)]+\)`)
	)

	doc.Find("div.entry-content table tr").EachWithBreak(func(_ int, div *goquery.Selection) bool {
		parts := div.Find("td")
		if parts.Length() != 3 {
			domErr = fmt.Errorf("Expected 3 parts to event, found %d", parts.Length())
			return false
		}

		partSlice := parts.Map(func(_ int, s *goquery.Selection) string {
			return s.Text()
		})

		date := clean(re.ReplaceAllString(partSlice[0], ""))

		name := partSlice[1]
		venue := partSlice[2]
		desc := strings.Join([]string{name, venue}, "\n")

		d := strings.FieldsFunc(date, func(r rune) bool {
			return r == '-' || r == '–'
		})
		if len(d) != 2 {
			domErr = fmt.Errorf("Expected 2 date items, got %d: %s", len(d), date)
			return false
		}

		for _, dt := range d {
			dt = clean(dt)
			t, err := time.Parse("april 1 9pm 2019", fmt.Sprintf("%s %s", dt, year))
			if err != nil {
				domErr = fmt.Errorf("Failed to parse date: %s: %s", dt, err)
				return false
			}
			f.log <- fmt.Sprintf("t=%s\n", t)
		}

		clean := clean(venue)
		v, ok := f.venueMap[clean]
		if !ok {

			if strings.Contains(clean, "jikoji") {
				v = f.venueMap["jikoji"]
			} else {

				domErr = fmt.Errorf("Failed to match venue for: %s", venue)
				return false
			}
		}

		e := &Event{
			URL:   u,
			Name:  name,
			Desc:  desc,
			Venue: v,
		}
		events = append(events, e)

		f.log <- fmt.Sprintf("Found event: %s | %s | %s\n", date, name, venue)

		return true
	})
	f.log <- fmt.Sprintf("Found %d total events", len(events))
	return events, domErr
}
