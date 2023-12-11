package app_config

import (
	"github.com/boltdb/bolt"
)

type AppConfig struct {
	Database *bolt.DB;
}