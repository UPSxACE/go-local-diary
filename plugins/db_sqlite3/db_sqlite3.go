/*
The package db_sqlite is a plugin that shall be used to plug a sqlite3 database driver
into the app. It comes also with the Repository interface that must be used to create models,
and the SqlFileReader class that can be used to read sql files and execute their instructions.
*/
package db_sqlite3

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database_Sqlite3 struct {
	instance *sql.DB
	version  string
}

func Init(devMode bool, path string) *Database_Sqlite3 {
	// If devmode, reset db
	if devMode && path != ":memory:" {
		os.Remove(path)
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", path)
	if devMode {
		db.SetMaxOpenConns(1) // NOTE: Necessary when using ":memory:" connection
	}

	if err != nil {
		log.Fatal(err)
	}

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sqlite3 initialized!")
	fmt.Println("Sqlite3 version: " + version)

	return &Database_Sqlite3{instance: db, version: version}
}

func (db *Database_Sqlite3) GetInstance() *sql.DB {
	return db.instance
}

func (db *Database_Sqlite3) GetTables() []string {
	query := `SELECT name FROM 
			(SELECT * FROM sqlite_schema UNION ALL
		 	SELECT * FROM sqlite_temp_schema)
 			WHERE type='table'
 			ORDER BY name`

	statement, err := db.instance.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, name)
	}

	return tables
}
