/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server"
)

func main() {
	devFlag := flag.Bool("dev", false, "Run server on developer mode")
	flag.Parse()

	// Init server with Sqlite
	db := db_sqlite3.Init(*devFlag)

	server := server.NewServer(db, *devFlag)

	fmt.Println("Initializing app...")
	server.Init()
}
