/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/UPSxACE/go-local-diary/app_config"
	"github.com/UPSxACE/go-local-diary/server"
	"github.com/boltdb/bolt"
)

func main() {
	// Open the my.db data file in the current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	appConfig := app_config.AppConfig{Database: db}

	// Print database object
	fmt.Println("Database:")
	fmt.Println(db)

	server.Init(appConfig)
}