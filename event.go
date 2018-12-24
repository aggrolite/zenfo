package zenfo

import (
	"time"
)

// Event represents events for DB store and web API
type Event struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Blurb string    `json:"blurb"`
	Desc  string    `json:"desc"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	URL   string    `json:"url"`
	Venue *Venue    `json:"venue"`
}
