/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/plugins/dev_component_parser"
	"github.com/UPSxACE/go-local-diary/server"
)

func main() {
	devFlag := flag.Bool("dev", false, "Run server on developer mode")
	flag.Parse()

	// Init server with Sqlite
	app := app.App[db_sqlite3.Database_Sqlite3]{Database: db_sqlite3.Init(), DevMode: *devFlag, Plugins: map[string]interface{}{}}

	// Load Plugins
	dev_component_parser.LoadPlugin(&app)

	// Print server config
	fmt.Println("App Config:")
	fmt.Println(app)

	server.Init(&app)
}