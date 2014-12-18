package main

import (
	"flag"
	"log"
	"net/http"
	"os/user"

	"github.com/ap4y/gpgdb/lib"
)

func main() {
	usr, _ := user.Current()
	keyStore := flag.String("keyStore", usr.HomeDir+"/.gpgdb", "path to the database")
	dbStore := flag.String("dbStore", "./db", "path to the database")
	port := flag.String("port", "8080", "port used by daemon")
	flag.Parse()

	es, err := lib.NewEntityStorage(*keyStore)
	if err != nil {
		log.Printf("Enable to read private keys storage: %s", err)
		return
	}

	db, err := lib.NewDB(*dbStore)
	if err != nil {
		log.Printf("Enable to open database: %s", err)
		return
	}

	server := NewServer(es, db)
	log.Printf("Started daemon on port %s", *port)
	http.ListenAndServe(":"+*port, server.Router)
}
