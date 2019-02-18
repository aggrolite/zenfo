//go:generate esc -o static.go -pkg zenfo -prefix static static

package zenfo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// API provides HTTP endpoints for serving events
type API struct {
	Port   int
	dbUser string
	dbName string
}

// NewAPI returns new API object
func NewAPI(dbUser, dbName string, port int, temp bool) (*API, error) {
	return &API{
		Port:   port,
		dbUser: dbUser,
		dbName: dbName,
	}, nil
}

// Run starts web server to listen on configured port
func (api *API) Run() error {
	http.HandleFunc("/api/events", api.getEvents)
	log.Printf("HTTP API listening on port %d\n", api.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", api.Port), nil)
}

func (api *API) getEvents(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", api.dbUser, api.dbName))
	if err != nil {
		log.Printf("url=%s err=%s\n", r.URL, err)
		http.Error(w, "Oops! Something went wrong!", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	q := "SELECT id, venue_id, name, blurb, description, start_date, end_date, url FROM events"

	keys, ok := r.URL.Query()["id"]
	if ok {
		q = fmt.Sprintf("%s WHERE id=%s", q, keys[0])
	}

	q = fmt.Sprintf("%s ORDER BY start_date", q)

	// Fetch events
	rows, err := db.Query(q)
	if err != nil {
		log.Printf("url=%s err=%s q=%s\n", r.URL, err, q)
		http.Error(w, "Oops! Something went wrong!", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var (
			venueID int
			event   Event
		)
		if err := rows.Scan(&event.ID, &venueID, &event.Name, &event.Blurb, &event.Desc, &event.Start, &event.End, &event.URL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch venue tied to venue_id
		venue := new(Venue)
		venueRow := db.QueryRow(`SELECT name, addr, geo[0], geo[1], website, phone, email FROM venues WHERE id=$1`, venueID)
		err := venueRow.Scan(&venue.Name, &venue.Addr, &venue.Lat, &venue.Lng, &venue.Website, &venue.Phone, &venue.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.Venue = venue
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}
