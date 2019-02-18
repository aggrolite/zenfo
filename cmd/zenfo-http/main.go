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
	temp   bool
)

func init() {
	flag.IntVar(&port, "port", 8081, "HTTP port to listen on")
	flag.BoolVar(&temp, "temp", false, "Show temporary 'coming soon' page")
	flag.StringVar(&dbName, "dbname", "zenfo", "Postgres DB name")
	flag.StringVar(&dbUser, "dbuser", "postgres", "Postgres DB user")
	flag.Parse()

	if len(dbName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if len(dbUser) == 0 {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	api, err := zenfo.NewAPI(dbUser, dbName, port, temp)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
