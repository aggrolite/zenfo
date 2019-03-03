//go:generate esc -o static.go -pkg zenfo -prefix dist -ignore .*\.sw.* dist

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
	dev             bool
	dbh             *sql.DB
	cert            string
	key             string
	eventsAllStmt   *sql.Stmt
	eventsOneStmt   *sql.Stmt
	eventsAfterStmt *sql.Stmt
}

// NewAPI returns new API object
func NewAPI(dbUser, dbName, cert, key string, dev bool) (*API, error) {
	a := &API{
		dev:  dev,
		cert: cert,
		key:  key,
	}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName))
	if err != nil {
		return nil, err
	}
	a.dbh = db

	q := "SELECT id, venue_id, name, blurb, description, start_date, end_date, url FROM events"

	var stmtErr error
	a.eventsAllStmt, stmtErr = a.dbh.Prepare(fmt.Sprintf("%s WHERE start_date >= now()::date - interval '1 day' ORDER BY start_date LIMIT 20", q))
	if stmtErr != nil {
		return nil, stmtErr
	}
	a.eventsOneStmt, stmtErr = a.dbh.Prepare(fmt.Sprintf("%s WHERE id = $1", q))
	if stmtErr != nil {
		return nil, stmtErr
	}
	a.eventsAfterStmt, stmtErr = a.dbh.Prepare(fmt.Sprintf("%s WHERE start_date > $1 ORDER BY start_date LIMIT 20", q))
	if stmtErr != nil {
		return nil, stmtErr
	}

	return a, nil
}

// Run starts web server to listen on configured port
func (api *API) Run() error {
	http.Handle("/", http.FileServer(FS(false)))
	http.HandleFunc("/__health", api.getHealth)
	http.HandleFunc("/api/events", api.getEvents)
	http.HandleFunc("/api/venues", api.getVenues)

	// Dev mode only runs on http
	if api.dev {
		return api.runDev()
	}

	go func() {
		p := 8081
		log.Printf("HTTP->HTTPS listening on port %d\n", p)
		http.ListenAndServe(fmt.Sprintf(":%d", p), http.HandlerFunc(api.redirect))
	}()

	p := 8082
	log.Printf("HTTPS API listening on port %d\n", p)

	// Firefox is strict
	// https://godoc.org/net/http#ListenAndServeTLS
	// If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", p), api.cert, api.key, nil)
}

func (api *API) runDev() error {
	p := 8081
	return http.ListenAndServe(fmt.Sprintf(":%d", p), nil)
}

func (api *API) setHeaders(h http.Header) {
	h.Set("Content-Type", "application/json")
	if api.dev {
		h.Set("Access-Control-Allow-Origin", "*")
	} else {
		h.Set("Access-Control-Allow-Origin", "https://zenfo.info")
	}
}

func (api *API) redirect(w http.ResponseWriter, req *http.Request) {
	redir := fmt.Sprintf("https://%s%s", req.Host, req.URL)
	http.Redirect(w, req, redir, http.StatusMovedPermanently)
}

// Close closes DB handler
func (api *API) Close() error {
	return api.dbh.Close()
}

func (api *API) getVenues(w http.ResponseWriter, r *http.Request) {
	api.setHeaders(w.Header())

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
	api.setHeaders(w.Header())

	q := r.URL.Query()

	var (
		events []*Event
		err    error
	)
	for param, v := range q {
		if param == "id" {
			events, err = api.events(api.eventsOneStmt, v[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			break
		}

		if param == "after" {
			events, err = api.events(api.eventsAfterStmt, v[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			break

		}
	}

	if len(q) == 0 {
		events, err = api.events(api.eventsAllStmt)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

func (api *API) events(stmt *sql.Stmt, args ...interface{}) ([]*Event, error) {
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var (
			venueID int
			event   Event
		)
		if err := rows.Scan(&event.ID, &venueID, &event.Name, &event.Blurb, &event.Desc, &event.Start, &event.End, &event.URL); err != nil {
			return nil, err
		}

		// Fetch venue tied to venue_id
		venue := new(Venue)
		venueRow := api.dbh.QueryRow(`SELECT id, name, addr, geo[0], geo[1], website, phone, email FROM venues WHERE id=$1`, venueID)
		err := venueRow.Scan(&venue.ID, &venue.Name, &venue.Addr, &venue.Lat, &venue.Lng, &venue.Website, &venue.Phone, &venue.Email)
		if err != nil {
			return nil, err
		}
		event.Venue = venue
		events = append(events, &event)
	}
	return events, rows.Err()
}
