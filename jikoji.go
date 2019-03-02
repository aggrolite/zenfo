package zenfo

import (
	"fmt"
	"strings"
	"time"

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

	doc.Find(".main-content .eventlist--upcoming article.eventlist-event div.eventlist-column-info").EachWithBreak(func(_ int, div *goquery.Selection) bool {

		// Init with name and venue
		e := &Event{
			Name:  div.Find(".eventlist-title").First().Text(),
			Venue: j.venueMap["jikoji"],
		}

		// Fetch profile page
		div.Find(".eventlist-excerpt a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			if clean(s.Text()) == "more information" {
				e.URL, _ = s.Attr("href")
				return false
			}
			return true
		})

		// Error if failed
		if e.URL == "" {
			domErr = fmt.Errorf("Failed to get More Info link for event: %s\n", e.Name)
			return false
		}
		if err := j.getProfileDetails(e); err != nil {
			domErr = err
			return false
		}

		events = append(events, e)

		return true
	})
	j.log <- fmt.Sprintf("Found %d total events", len(events))
	return events, domErr
}

func (j *Jikoji) getProfileDetails(e *Event) error {
	j.log <- fmt.Sprintf("Fetching profile page: %s", e.URL)
	resp, err := j.client.Get(e.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	var (
		domErr error
	)

	// Scrape start-end date
	// This one's a doozey
	doc.Find(".event-time-24hr").EachWithBreak(func(i int, s *goquery.Selection) bool {

		// One day events
		// https://www.jikoji.org/jikoji-events/sunrise-zazenkai-march-2019
		start := s.Find(".event-time-24hr-start").First()

		hour := start.Text()
		if hour == "" {
			hour = s.Text()
			if hour == "" {
				domErr = fmt.Errorf("Failed to extract datetime text! Profile=%s", e.URL)
				return false
			}
		}

		date, ok := start.Attr("datetime")
		if !ok {

			// Multi day events
			// https://www.jikoji.org/jikoji-events/spring-nature-sesshin-march-2019
			date, ok = s.Attr("datetime")
			if !ok {
				domErr = fmt.Errorf("Failed to extract datetime attr! Profile=%s", e.URL)
				return false
			}

			parsed, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s:00-07:00", date, hour))
			if err != nil {
				domErr = err
				return false
			}
			if i == 0 {
				e.Start = parsed
			} else {
				e.End = parsed
			}

		} else {
			startParsed, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s:00-07:00", date, hour))
			if err != nil {
				domErr = err
				return false
			}
			e.Start = startParsed

			// Mistake on the site. 24hr end time listed as 12hr
			endHour := s.Find(".event-time-24hr-end").First().Text()
			if endHour == "" {
				endHour = s.Find(".event-time-12hr-end").Last().Text()
				if endHour == "" {
					domErr = fmt.Errorf("Failed to extract datetime text! Profile=%s", e.URL)
					return false
				}
			}

			endDate, ok := start.Attr("datetime")
			if !ok {
				domErr = fmt.Errorf("Failed to get end date: %s\n", e.URL)
				return false
			}
			j.log <- fmt.Sprintf("end=%s", endDate)

			endParsed, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s:00-07:00", endDate, endHour))
			if err != nil {
				domErr = err
				return false
			}
			e.End = endParsed
		}

		return true
	})

	// Blurb / tagline
	blurb := doc.Find(".sqs-block-html .sqs-block-content h3 strong").First().Text()
	if blurb == "" {
		// non-fatal
		// may not exist
		// https://www.jikoji.org/jikoji-events/practice-period-2019
		j.log <- fmt.Sprintf("Warning, no blurb: %s", e.URL)
	}
	e.Blurb = blurb

	// Event description
	// TODO preserve HTML formatting
	var desc []string
	doc.Find(".sqs-block-html .sqs-block-content").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			return true
		}
		h, err := s.Html()
		if err != nil {
			domErr = err
			return false
		}
		desc = append(desc, h)
		return true
	})
	if err != nil {
		return fmt.Errorf("Failed to get description: %s: %s", e.URL, err)
	}
	e.Desc = strings.Join(desc, "\n")

	//j.log <- fmt.Sprintf("%+v", e)
	return domErr
}
