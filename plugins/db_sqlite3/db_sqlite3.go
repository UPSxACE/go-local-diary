// Web, Api, App
package db_sqlite3

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database_Sqlite3 struct {
	instance *sql.DB;
	version string;
}

func Init() *Database_Sqlite3 {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sqlite3 version: " + version)

	return &Database_Sqlite3{instance: db,version: version}
}

func (db *Database_Sqlite3) Create() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Sqlite3) ReadOne() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Sqlite3) ReadAll() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Sqlite3) UpdateOne() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Sqlite3) UpdateAll() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Sqlite3) DeleteOne() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Sqlite3) DeleteAll() any {
	panic("not implemented") // TODO: Implement
}
