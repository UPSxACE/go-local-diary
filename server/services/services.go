package services

import (
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
)

type Services struct {
	database *db_sqlite3.Database_Sqlite3;
}

func NewServices(database *db_sqlite3.Database_Sqlite3) (*Services) {
	return &Services{database: database}
}