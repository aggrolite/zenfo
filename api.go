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
	Port int
	dbh  *sql.DB
	temp bool
	cert string
	key  string
}

// NewAPI returns new API object
func NewAPI(dbUser, dbName, cert, key string, port int, temp bool) (*API, error) {
	a := &API{
		Port: port,
		temp: temp,
		cert: cert,
		key:  key,
	}
	if !a.temp {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName))
		if err != nil {
			return nil, err
		}
		a.dbh = db
	}
	return a, nil
}

// Run starts web server to listen on configured port
func (api *API) Run() error {
	http.HandleFunc("/__health", api.getHealth)
	http.HandleFunc("/api/events", api.getEvents)
	http.HandleFunc("/api/venues", api.getVenues)
	log.Printf("HTTP API listening on port %d\n", api.Port)

	return http.ListenAndServeTLS(fmt.Sprintf(":%d", api.Port), api.cert, api.key, nil)
}

// Close closes DB handler
func (api *API) Close() error {
	if api.temp {
		return api.dbh.Close()
	}
	return nil
}

func (api *API) getVenues(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	q := "SELECT id, name, geo[0], geo[1], website, phone, email FROM venues"

	rows, err := api.dbh.Query(q)
	if err != nil {
		log.Printf("url=%s err=%s q=%s\n", r.URL, err, q)
		http.Error(w, "Oops! Something went wrong!", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var venues []*Venue
	for rows.Next() {
		var v Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Lat, &v.Lng, &v.Website, &v.Phone, &v.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		venues = append(venues, &v)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(venues)
}

func (api *API) getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func (api *API) getEvents(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	q := "SELECT id, venue_id, name, blurb, description, start_date, end_date, url FROM events"

	keys, ok := r.URL.Query()["id"]
	if ok {
		q = fmt.Sprintf("%s WHERE id=%s", q, keys[0])
	}

	q = fmt.Sprintf("%s ORDER BY start_date", q)

	// Fetch events
	rows, err := api.dbh.Query(q)
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
		venueRow := api.dbh.QueryRow(`SELECT id, name, addr, geo[0], geo[1], website, phone, email FROM venues WHERE id=$1`, venueID)
		err := venueRow.Scan(&venue.ID, &venue.Name, &venue.Addr, &venue.Lat, &venue.Lng, &venue.Website, &venue.Phone, &venue.Email)
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
