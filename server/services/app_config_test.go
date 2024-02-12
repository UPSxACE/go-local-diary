package services

import (
	"context"
	"log"
	"testing"

	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/models"
	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

var initialDbSqlFile = "../../server/sql/initial.sql"
var testValuesSqlFile = "../../server/sql/initial_test.sql"
var test_db *db_sqlite3.Database_Sqlite3

func init() {
	test_db = db_sqlite3.Init(true, ":memory:")

	// load default db
	sqlFileReader, err := db_sqlite3.OpenSqlFile(initialDbSqlFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlFileReader.ExecuteAll(test_db)
	if err != nil {
		log.Fatal(err)
	}

	// load test values into the default db
	sqlFileReader, err = db_sqlite3.OpenSqlFile(testValuesSqlFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlFileReader.ExecuteAll(test_db)
	if err != nil {
		log.Fatal(err)
	}
}

func TestAppConfiguration(t *testing.T) {
	services := NewServices(test_db)

	store, err := models.CreateStoreAppConfig(test_db, false, nil)
	testhelper.ExpectNoError(t, err)
	ctx := context.Background()

	// test SetConfiguration method, setting a setting that didn't exist at first
	count, err := store.CountByName("test_test_test")
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, count, 0)

	_, err = services.SetConfiguration(ctx, "test_test_test", "random-value")
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
	isConfigured, err := services.IsAppConfigured(ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, isConfigured, false)

	// so, now it shall exist
	count, err = store.CountByName("configured")
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, count, 1)

	// if called again, must still return false
	isConfigured, err = services.IsAppConfigured(ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, isConfigured, false)

	// set the value to "1" with SetConfiguration
	newval, err := services.SetConfiguration(ctx, "configured", "1")
	testhelper.ExpectNoError(t, err)
	// check if its value is properly set
	model, err := store.GetFirstByName("configured")
	testhelper.ExpectNoError(t, err)
	// check if both values are in sync
	testhelper.ExpectEqual(t, newval, model.Value)

	// check if IsAppConfigured now returns true
	isConfigured, err = services.IsAppConfigured(ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, isConfigured, true)
}
