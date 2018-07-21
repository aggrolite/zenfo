package main

import (
	"flag"
	"log"

	"github.com/aggrolite/zenfo"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 8080, "HTTP port to listen on")
	flag.Parse()
}

func main() {
	api, err := zenfo.NewAPI(port)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
