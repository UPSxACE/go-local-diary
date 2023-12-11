/*
The main package is the entry point for the application,
and will be responsible for initializing the server,
the database, and other dependencies.
*/
package main

import "github.com/UPSxACE/go-local-diary/server"

func main() {
	println("Local Diary executed")
	server.Init()
}