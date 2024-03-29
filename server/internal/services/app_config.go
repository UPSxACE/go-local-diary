package services

import (
	"context"
	"fmt"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/models"
)

// NOTE: Do not call service methods inside each other, to avoid multiple instances of stores at the same time

var AppConfig = AppConfigService{}

type AppConfigService struct{}

/*
Returns true if the configuration of name "configured" is set to "1" in the app_config table.
Returns false if it is set to "0".

If the value is not set in the app_config table yet, the record will be created with the value "0",
and then false will be returned.
*/
func (service *AppConfigService) IsAppConfigured(app *app.App[db_sqlite3.Database_Sqlite3], context context.Context) (appConfigured bool, err error) {
	appConfigStore, err := models.CreateStoreAppConfig(app, true, context)
	if err != nil {
		return false, err
	}
	defer appConfigStore.Close()

	// check if exists
	exists, model, err := appConfigStore.CheckByName("configured")
	if err != nil {
		return false, err
	}

	// does not exist (create and then return false)
	if !exists {
		_, err = appConfigStore.Create(models.AppConfigModel{Id: 0, Name: "configured", Value: "0"})
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// exists (just return its value)
	if model.Value == "1" {
		return true, nil
	} else if model.Value == "0" {
		return false, nil
	} else {
		fmt.Println("Invalid configuration value in config: 'configured'", err)
		return false, nil
	}
}

/*
Sets a configuration in app_config table.
If the configuration already exists it will be updated, if it doesn't it will  be created.
*/
func (service *AppConfigService) SetConfiguration(app *app.App[db_sqlite3.Database_Sqlite3], context context.Context, configName string, configValue string) (newValue string, err error) {
	store, err := models.CreateStoreAppConfig(app, true, context)
	if err != nil {
		return "", err
	}
	defer store.Close()

	// check if exists
	exists, model, err := store.CheckByName(configName)
	if err != nil {
		return "", err
	}

	// does not exist
	if !exists {
		model = models.AppConfigModel{Name: configName, Value: configValue}
		model, err = store.Create(model)
		if err != nil {
			return "", err
		}
	}
	// exists 
	if exists {
		model.Value = "1"
		model, err = store.UpdateByName(model.Name, model)
		if err != nil {
			return "", err
		}
	}

	return model.Value, nil;
}
