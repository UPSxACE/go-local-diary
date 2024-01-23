// Web, Api, App
package db_sqlite3

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database_Sqlite3 struct {
	instance *sql.DB
	version  string
}

func Init() *Database_Sqlite3 {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sqlite3 version: " + version)

	return &Database_Sqlite3{instance: db, version: version}
}

func (db *Database_Sqlite3) GetInstance() *sql.DB{
	return db.instance;
}

func (db *Database_Sqlite3) GetTables() []string{
	query := `SELECT name FROM 
			(SELECT * FROM sqlite_schema UNION ALL
		 	SELECT * FROM sqlite_temp_schema)
 			WHERE type='table'
 			ORDER BY name`

	statement,err := db.instance.Prepare(query)
	if(err != nil){
		log.Fatal(err)
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string;
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, name)
	}
	
	return tables;
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
