package zenfo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // postgres driver
)

// Manager handles running all workers
type Manager struct {
	workers map[string]Worker
	db      *sql.DB
}

// NewManager builds internal worker map and returns new Mananger object
func NewManager(dbName, dbUser string) (*Manager, error) {
	m := new(Manager)

	// TODO It'd be nice to keep worker code inside a subdir, e.g. workers/
	m.workers = map[string]Worker{
		"aczc": &Aczc{},
		"sfzc": &Sfzc{},
	}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName))
	if err != nil {
		return nil, err
	}
	m.db = db

	return m, nil
}

// Run iterates all workers for execution
func (m *Manager) Run() error {

	client := NewClient()

	for name, w := range m.workers {

		if err := w.Init(client); err != nil {
			return err
		}

		log.Printf("Running worker: %s - %s\n", name, w.Desc())

		events, err := w.Events()
		if err != nil {
			return err
		}

		venueSelect, err := m.db.Prepare(`SELECT id FROM venues WHERE name=$1`)
		if err != nil {
			return err
		}
		venueInsert, err := m.db.Prepare(`INSERT INTO venues (name, addr, geo, website, phone, email) VALUES ($1, $2, point($3, $4), $5, $6, $7) RETURNING id`)
		if err != nil {
			return err
		}
		eventStmt, err := m.db.Prepare("INSERT INTO events (venue_id, name, blurb, description, start_date, end_date, url) VALUES ($1, $2, $3, $4, $5, $6, $7)")
		if err != nil {
			return err
		}

		for _, e := range events {
			log.Printf("event=%+v\n", e)

			venue := e.Venue

			var venueID int
			if err := venueSelect.QueryRow(venue.Name).Scan(&venueID); err != nil {
				if err == sql.ErrNoRows {
					if err := venueInsert.QueryRow(venue.Name, venue.Addr, venue.Lat, venue.Lng, venue.Website, venue.Phone, venue.Email).Scan(&venueID); err != nil {
						return err
					}

				} else {
					return err
				}
			}

			_, err := eventStmt.Exec(venueID, e.Name, e.Blurb, e.Desc, e.Start, e.End, e.URL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
