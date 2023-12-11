/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import (
	"fmt"

	"github.com/UPSxACE/go-local-diary/database"
	"github.com/UPSxACE/go-local-diary/server"
)

func main() {
	db := database.Init()
	defer db.Close()

	// Print database object
	fmt.Println("Database:")
	fmt.Println(database.Database)

	server.Init()
}