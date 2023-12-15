package db_bolt

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type Database_Bolt struct {
	instance *bolt.DB
}

func Init() *Database_Bolt {
	// Open the my.db data file in the current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return &Database_Bolt{instance: db}
}

func (db *Database_Bolt) Create() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Bolt) ReadOne() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Bolt) ReadAll() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Bolt) UpdateOne() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Bolt) UpdateAll() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Bolt) DeleteOne() any {
	panic("not implemented") // TODO: Implement
}
func (db *Database_Bolt) DeleteAll() any {
	panic("not implemented") // TODO: Implement
}
