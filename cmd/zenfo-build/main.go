package main

import (
	"flag"
	"log"
	"os"

	"github.com/aggrolite/zenfo"
)

var (
	dbName string
	dbUser string
)

func init() {
	flag.StringVar(&dbName, "dbname", os.Getenv("DBNAME"), "Postgres DB name")
	flag.StringVar(&dbUser, "dbuser", os.Getenv("DBUSER"), "Postgres DB user")
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
