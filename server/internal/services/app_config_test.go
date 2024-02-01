package services

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/models"
	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

var initialDbSqlFile = "../../sql/initial.sql"
var testValuesSqlFile = "../../sql/initial_test.sql"
var test_app *app.App[db_sqlite3.Database_Sqlite3]
var test_db *sql.DB

func getTestAppInstance() *app.App[db_sqlite3.Database_Sqlite3] {
	app := app.App[db_sqlite3.Database_Sqlite3]{
		Database: db_sqlite3.Init(true),
	}
	return &app
}

func getTestAppInstanceAndDb() (*app.App[db_sqlite3.Database_Sqlite3], *sql.DB) {
	app := getTestAppInstance()
	db := app.Database.GetInstance()
	return app, db
}

func init() {
	app, db := getTestAppInstanceAndDb()
	// load default db
	sqlFileReader, err := db_sqlite3.OpenSqlFile(initialDbSqlFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlFileReader.ExecuteAllFromApp(app)
	if err != nil {
		log.Fatal(err)
	}

	// load test values into the default db
	sqlFileReader, err = db_sqlite3.OpenSqlFile(testValuesSqlFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlFileReader.ExecuteAllFromApp(app)
	if err != nil {
		log.Fatal(err)
	}

	test_app = app
	test_db = db
}

func TestAppConfiguration(t *testing.T) {
	store, err := models.CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err)
	ctx := context.Background()

	// test SetConfiguration method, setting a setting that didn't exist at first
	count, err := store.CountByName("test_test_test")
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, count, 0)

	_, err = AppConfig.SetConfiguration(test_app, ctx, "test_test_test", "random-value")
	testhelper.ExpectNoError(t, err)

	count, err = store.CountByName("test_test_test")
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, count, 1)

	// at the start of the test the configuration of name "configured"
	// does not exist
	count, err = store.CountByName("configured")
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, count, 0)

	// when called, the IsAppConfigured method must set the value to "0"
	// automatically when it doesn't exist
	isConfigured, err := AppConfig.IsAppConfigured(test_app, ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, isConfigured, false)

	// so, now it shall exist
	count, err = store.CountByName("configured")
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, count, 1)

	// if called again, must still return false
	isConfigured, err = AppConfig.IsAppConfigured(test_app, ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, isConfigured, false)

	// set the value to "1" with SetConfiguration
	newval, err := AppConfig.SetConfiguration(test_app, ctx, "configured", "1")
	testhelper.ExpectNoError(t, err)
	// check if its value is properly set
	model, err := store.GetFirstByName("configured")
	testhelper.ExpectNoError(t, err)
	// check if both values are in sync
	testhelper.ExpectEqual(t, newval, model.Value)

	// check if IsAppConfigured now returns true
	isConfigured, err = AppConfig.IsAppConfigured(test_app, ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, isConfigured, true)
}
