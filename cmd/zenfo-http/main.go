package main

import (
	"flag"
	"log"
	"os"

	"github.com/aggrolite/zenfo"
)

var (
	dbName   string
	dbUser   string
	dev      bool
	certPath string
	keyPath  string
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Toggle development mode")
	flag.StringVar(&dbName, "dbname", "zenfo", "Postgres DB name")
	flag.StringVar(&dbUser, "dbuser", "postgres", "Postgres DB user")
	flag.StringVar(&certPath, "cert", "", "Path to SSL cert file")
	flag.StringVar(&keyPath, "key", "", "Path to SSL key file")
	flag.Parse()

	if len(dbName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if len(dbUser) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	if !dev {
		if len(certPath) == 0 {
			flag.Usage()
			os.Exit(2)
		}
		if len(keyPath) == 0 {
			flag.Usage()
			os.Exit(2)
		}
	}
}

func main() {
	api, err := zenfo.NewAPI(dbUser, dbName, certPath, keyPath, dev)
	if err != nil {
		log.Fatal(err)
	}
	defer api.Close()

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
