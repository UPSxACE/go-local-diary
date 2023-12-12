package app_config

import (
	"github.com/UPSxACE/go-local-diary/server/plugin"
	"github.com/boltdb/bolt"
)

type AppConfig struct {
	Database *bolt.DB;
	DevMode bool;
	Plugins plugin.PluginsData;
}