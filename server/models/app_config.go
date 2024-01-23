package models

import (
	"log"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
)

type AppConfig struct {
}

func PrepareAppConfigTable(appInstance *app.App[db_sqlite3.Database_Sqlite3]) {
	db := appInstance.Database.GetInstance();

	query := `CREATE TABLE IF NOT EXISTS app_config (
	id INTEGER PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE,
	value VARCHAR(255)
	)`

	_, err := db.Exec(query)

	if(err != nil){
		log.Fatal(err)
	}
}