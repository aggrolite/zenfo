package zenfo

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq" // postgres driver
)

// Manager handles running all workers
type Manager struct {
	workers []Worker
	db      *sql.DB
	log     chan string
	DryRun  bool
}

// NewManager builds internal worker map and returns new Mananger object
func NewManager(dbName, dbUser string) (*Manager, error) {
	m := new(Manager)

	// TODO It'd be nice to keep worker code inside a subdir, e.g. workers/
	// Currently it creates circular imports
	m.workers = []Worker{
		&Aczc{},
		&Sfzc{},
		&Jikoji{},
	}

	if !m.DryRun {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName))
		if err != nil {
			return nil, err
		}
		m.db = db
	}

	return m, nil
}

// Run iterates all workers for execution
func (m *Manager) Run() error {

	// http client
	client := NewClient()

	var (
		venueSelect *sql.Stmt
		venueInsert *sql.Stmt
		eventInsert *sql.Stmt
	)

	if !m.DryRun {

		var err error

		venueSelect, err = m.db.Prepare(`SELECT id FROM venues WHERE name=$1`)
		if err != nil {
			return err
		}
		venueInsert, err = m.db.Prepare(`INSERT INTO venues (name, addr, geo, website, phone, email) VALUES ($1, $2, point($3, $4), $5, $6, $7) RETURNING id`)
		if err != nil {
			return err
		}
		eventInsert, err = m.db.Prepare("INSERT INTO events (venue_id, name, blurb, description, start_date, end_date, url) VALUES ($1, $2, $3, $4, $5, $6, $7)")
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	for _, worker := range m.workers {
		wg.Add(2)

		w := worker
		out := make(chan string)
		errs := make(chan error)

		// Handle errors
		go func() {
			if err := <-errs; err != nil {
				log.Fatalf("[%s] -> FATAL: %s\n", w.Name(), err)
			}

			// Wait for any error before marking as done
			wg.Done()
		}()

		// Print worker log to screen
		go func() {
			for msg := range out {
				log.Printf("[%s] -> %s\n", w.Name(), msg)
			}

			// Also wait for log output before marking done
			wg.Done()
		}()

		go func() {
			defer close(out)
			defer close(errs)

			if err := w.Init(client, out); err != nil {
				errs <- err
			}

			log.Printf("Running worker: %s - %s\n", w.Name(), w.Desc())

			events, err := w.Events()
			if err != nil {
				errs <- err
			}

			for _, e := range events {
				//log.Printf("event=%+v\n", e)

				venue := e.Venue

				var venueID int

				// Create any new venues found from events
				if !m.DryRun {
					err := venueSelect.QueryRow(venue.Name).Scan(&venueID)
					if err == sql.ErrNoRows {
						if err := venueInsert.QueryRow(venue.Name, venue.Addr, venue.Lat, venue.Lng, venue.Website, venue.Phone, venue.Email).Scan(&venueID); err != nil {
							errs <- err
						}
					} else if err != nil {
						errs <- err
					}

					// Store event
					if _, err := eventInsert.Exec(venueID, e.Name, e.Blurb, e.Desc, e.Start, e.End, e.URL); err != nil {
						errs <- err
					}
				}
			}
		}()
	}
	wg.Wait()

	return nil
}
