package database

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var Database *bolt.DB = nil;

func Init() *bolt.DB{
	// Open the my.db data file in the current directory.
	// It will be created if it doesn't exist.
	var err error;
	Database, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return Database;
}