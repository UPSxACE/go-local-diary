package app_config

import (
	"github.com/boltdb/bolt"
)

type AppConfig struct {
	Database *bolt.DB;
	DevMode bool;
	Plugins PluginsData;
}

type PluginsData = map[string]interface{};