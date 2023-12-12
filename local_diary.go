/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/UPSxACE/go-local-diary/app_config"
	"github.com/UPSxACE/go-local-diary/server"
	"github.com/UPSxACE/go-local-diary/server/dev_component_parser"
	"github.com/boltdb/bolt"
)

func main() {
	devFlag := flag.Bool("dev", false, "Run server on developer mode")
	flag.Parse()


	// Open the my.db data file in the current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	appConfig := app_config.AppConfig{Database: db, DevMode: *devFlag}

	// Plugins
	appConfig.Plugins = map[string]interface{}{}
	dev_component_parser.LoadPlugin(&appConfig)

	// Print server config
	fmt.Println("App Config:")
	fmt.Println(appConfig)

	server.Init(&appConfig)
}