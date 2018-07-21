package zenfo

// All events here, but too much stuff:
// https://www.aczc.org/events-calendar/?view=calendar&month=July-2018
// https://www.aczc.org/api/open/GetItemsByMonth?month=July-2018&collectionId=594a9d2920099e63d87a096e

// Special events listed here:
// https://www.aczc.org/schedule/

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

// Aczc crawls aczc.org, satisfies Worker interface
type Aczc struct {
	venueMap map[string]*Venue
	client   *Client
}

// Init sets HTTP client and defines internal venue map
func (s *Aczc) Init(client *Client) error {
	s.client = client
	s.venueMap = make(map[string]*Venue)

	s.venueMap["Angel City"] = &Venue{
		Name:    "Angel City Zen Center",
		Addr:    "1407 W 2nd St Los Angeles, CA 90026",
		Phone:   "+1 (323) 426-6269",
		Email:   "angelcityzencenter@gmail.com",
		Lat:     34.060979,
		Lng:     -118.260530,
		Website: "https://aczc.org",
	}
	s.venueMap["Mount Baldy"] = &Venue{
		Name:    "Mount Baldy Zen Center",
		Addr:    "7901 Mount Baldy Road, Mount Baldy, CA 91759",
		Phone:   "+1 (909) 985-6410",
		Email:   "office@mbzc.org",
		Lat:     34.264256,
		Lng:     -117.632916,
		Website: "http://mbzc.org",
	}

	return nil
}

// Desc returns description for website crawled
func (s *Aczc) Desc() string {
	return "Angel City Zen Center (aczc.org)"
}

// Events hits aczc events page and returns slice of Event types
// https://www.aczc.org/schedule/
func (s *Aczc) Events() ([]*Event, error) {

	resp, err := s.client.Get("https://www.aczc.org/schedule/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	errs := make(chan error)
	go func() {
		doc.Find("h1").Each(func(_ int, h1 *goquery.Selection) {

			if h1.Text() != "Upcoming Special Events" {
				return
			}
			h1.Siblings().Each(func(_ int, p *goquery.Selection) {
				href, ok := p.Find("a").First().Attr("href")
				if !ok {
					log.Printf("Yikes! Event did not have a tag! event=%s\n", h1.Text())
				}

				log.Printf("Fetching aczc.org event: %s\n", href)
				resp, err := s.client.Get(href)
				if err != nil {
					errs <- err
					return
				}

				eventDoc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					errs <- err
					return
				}

				title := eventDoc.Find("h1.eventitem-title").Text()
				log.Printf("title=%s\n", title)

				// a tag text = name
				// href = website
				// some are broken, have multiple a tags, but same link

				// em = date
				// easier to parse from url, probably

			})
		})
		close(errs)
	}()
	if err := <-errs; err != nil {
		return nil, err
	}

	// TODO construct and return events
	return nil, nil
}
