package main

import (
	"flag"
	"log"

	"github.com/aggrolite/zenfo"
)

var (
	dbName string
	dbUser string
)

func init() {
	flag.StringVar(&dbName, "dbname", "", "Postgres DB name")
	flag.StringVar(&dbUser, "dbuser", "", "Postgres DB user")
	flag.Parse()
}

func main() {
	m, err := zenfo.NewManager(dbName, dbUser)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Run(); err != nil {
		log.Fatal(err)
	}
}
