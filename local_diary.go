/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/UPSxACE/go-local-diary/app_config"
	"github.com/UPSxACE/go-local-diary/plugins/db_bolt"
	"github.com/UPSxACE/go-local-diary/plugins/dev_component_parser"
	"github.com/UPSxACE/go-local-diary/server"
)

func main() {
	devFlag := flag.Bool("dev", false, "Run server on developer mode")
	flag.Parse()

	// Init server with BoltDB
	appConfig := app_config.AppConfig[db_bolt.Database_Bolt]{Database: db_bolt.Init(), DevMode: *devFlag, Plugins: map[string]interface{}{}}

	// Load Plugins
	dev_component_parser.LoadPlugin(&appConfig)

	// Print server config
	fmt.Println("App Config:")
	fmt.Println(appConfig)

	server.Init(&appConfig)
}