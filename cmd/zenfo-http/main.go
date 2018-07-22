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
	port   int
)

func init() {
	flag.IntVar(&port, "port", 8080, "HTTP port to listen on")
	flag.StringVar(&dbName, "dbname", os.Getenv("DBNAME"), "Postgres DB name")
	flag.StringVar(&dbUser, "dbuser", os.Getenv("DBUSER"), "Postgres DB user")
	flag.Parse()
}

func main() {
	api, err := zenfo.NewAPI(dbName, dbUser, port)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
