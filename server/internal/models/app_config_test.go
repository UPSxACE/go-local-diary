package models

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"testing"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

var initialDbSqlFile = "../../sql/initial.sql"
var testValuesSqlFile = "../../sql/initial_test.sql"
var test_app *app.App[db_sqlite3.Database_Sqlite3];
var test_db *sql.DB

func getTestAppInstance() *app.App[db_sqlite3.Database_Sqlite3] {
	app := app.App[db_sqlite3.Database_Sqlite3]{
		Database: db_sqlite3.Init(),
	}
	return &app
}

func getTestAppInstanceAndDb() (*app.App[db_sqlite3.Database_Sqlite3], *sql.DB) {
	app := getTestAppInstance()
	db := app.Database.GetInstance()
	return app, db
}

func TestCreateStoreAppConfig(t *testing.T) {
	app, db := getTestAppInstanceAndDb()
	defer db.Close()

	// Create normal nstore
	nstore, err := CreateStoreAppConfig(app, false, nil)
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, nstore.TransactionMode(), false)
	testhelper.ExpectDifferent(t, nstore.Repository(), (db_sqlite3.Repository)(nil))

	// Create transaction store
	ctx := context.Background()
	tstore, err := CreateStoreAppConfig(app, true, ctx)
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, tstore.TransactionMode(), true)
	testhelper.ExpectDifferent(t, tstore.Repository(), (db_sqlite3.Repository)(nil))
}

func init(){
	app, db := getTestAppInstanceAndDb();
	// load default db
	sqlFileReader, err := db_sqlite3.OpenSqlFile(initialDbSqlFile)
	if(err != nil){
		log.Fatal(err)
	}

	_, err = sqlFileReader.ExecuteAllFromApp(app)
	if(err != nil){
		log.Fatal(err)
	}

	// load test values into the default db
	sqlFileReader, err = db_sqlite3.OpenSqlFile(testValuesSqlFile)
	if(err != nil){
		log.Fatal(err)
	}

	_, err = sqlFileReader.ExecuteAllFromApp(app)
	if(err != nil){
		log.Fatal(err)
	}

	test_app = app
	test_db = db
}

func TestAppConfigCheckById(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	// test when it exists
	exists, model, err := store.CheckById(1)
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectEqual(t, exists, true)
	testhelper.ExpectEqual(t, model, AppConfigModel{1, "config1", "0"})

	// test when it doesn't
	exists, _, err = store.CheckById(999)
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectEqual(t, exists, false)
}

func TestAppConfigCheckByName(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	// test when it exists
	exists, model, err := store.CheckByName("config3")
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectEqual(t, exists, true)
	testhelper.ExpectEqual(t, model, AppConfigModel{3, "config3", ""})

	// test when it doesn't
	exists, _, err = store.CheckByName("zzzzzzzz")
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectEqual(t, exists, false)
}

func TestAppConfigGetFirstById(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	model, err := store.GetFirstById(2)
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectEqual(t, model, AppConfigModel{2, "config2", "this is a sentence"})
}

func TestAppConfigGetFirstByName(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	model, err := store.GetFirstByName("config2")
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectEqual(t, model, AppConfigModel{2, "config2", "this is a sentence"})
}

func TestAppConfigCreate(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	// test creating with auto assigned id
	newModel := AppConfigModel{0, "config4","random thing"}
	model, err := store.Create(newModel)
	testhelper.ExpectNoError(t,err);

	testhelper.ExpectEqual(t, model.Name, newModel.Name)
	testhelper.ExpectEqual(t, model.Value, newModel.Value) 
	testhelper.ExpectDifferent(t, model.Id, 0)


	// test creating with manual assignment of id
	newModel = AppConfigModel{1422, "config5","random thing"}
	model, err = store.Create(newModel)
	testhelper.ExpectNoError(t,err);
	
	testhelper.ExpectEqual(t, model.Name, newModel.Name)
	testhelper.ExpectEqual(t, model.Value, newModel.Value) 
	testhelper.ExpectEqual(t, model.Id, 1422)
}

func TestAppConfigUpdateById(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	model, err := store.GetFirstById(4)
	testhelper.ExpectNoError(t, err);

	oldName := model.Name
	newName := oldName + " updated"
	oldValue := model.Value
	newValue := oldValue + " updated"
	
	model.Name = newName;
	model.Value = newValue;

	model, err = store.UpdateById(4, model)
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectDifferent(t, model.Name, oldName)
	testhelper.ExpectEqual(t, model.Name, newName)
	testhelper.ExpectDifferent(t, model.Value, oldValue)
	testhelper.ExpectEqual(t, model.Value, newValue)
}

func TestAppConfigUpdateByName(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	model, err := store.GetFirstByName("up2")
	testhelper.ExpectNoError(t, err);

	oldName := model.Name
	newName := oldName + " updated"
	oldValue := model.Value
	newValue := oldValue + " updated"
	
	model.Name = newName;
	model.Value = newValue;

	model, err = store.UpdateByName("up2", model)
	testhelper.ExpectNoError(t,err);
	testhelper.ExpectDifferent(t, model.Name, oldName)
	testhelper.ExpectEqual(t, model.Name, newName)
	testhelper.ExpectDifferent(t, model.Value, oldValue)
	testhelper.ExpectEqual(t, model.Value, newValue)
}

func TestAppConfigDeleteById(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	// first the model should exist
	exists, _, err := store.CheckById(6);
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, exists, true);

	// then test deleting it
	deleted, err := store.DeleteById(6)
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, deleted, true)

	// then the model should not exist anymore
	exists, _, err = store.CheckById(6);
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, exists, false);
}

func TestAppConfigDeleteByName(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, false, nil)
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	// first the model should exist
	exists, _, err := store.CheckByName("del2");
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, exists, true);

	// then test deleting it
	deleted, err := store.DeleteByName("del2")
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, deleted, true)

	// then the model should not exist anymore
	exists, _, err = store.CheckByName("del2");
	testhelper.ExpectNoError(t, err);
	testhelper.ExpectEqual(t, exists, false);
}

func TestAppConfigNonExistantValuesError(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, true, context.Background())
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	_, err = store.GetFirstById(1000)
	testhelper.ExpectError(t, err);
	_, err = store.GetFirstByName("zzzzzzzzzzzzzzzz")
	testhelper.ExpectError(t, err);

	model := AppConfigModel{0, "abc", "def"}
	_, err = store.UpdateById(1000, model)
	testhelper.ExpectError(t, err);
	_, err = store.UpdateByName("zzzzzzzzzzzzzzzz", model)
	testhelper.ExpectError(t, err);

	_, err = store.DeleteById(950);
	testhelper.ExpectError(t,err);
	_, err = store.DeleteByName("zzzzzzzzzzzzzzzz");
	testhelper.ExpectError(t,err);
}

func TestAppConfigInvalidCreateOrUpdateError(t *testing.T){
	store,err := CreateStoreAppConfig(test_app, true, context.Background())
	testhelper.ExpectNoError(t, err);
	defer store.Close()

	existingModel,err := store.GetFirstById(1)
	testhelper.ExpectNoError(t, err);

	// id can't be smaller than 0
	_, err = store.Create(AppConfigModel{-1, "validname", "validvalue"})
	testhelper.ExpectError(t, err)
	copy := existingModel;
	copy.Id = -1
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);

	// id can't be updated to 0
	copy = existingModel;
	copy.Id = 0
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);
	
	// id can't exist
	_, err = store.Create(AppConfigModel{1, "validname", "validvalue"})
	testhelper.ExpectError(t, err)
	copy = existingModel;
	copy.Id = 2
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);

	// name can't exist
	_, err = store.Create(AppConfigModel{123, existingModel.Name, "validvalue"})
	testhelper.ExpectError(t, err)
	copy = existingModel;
	copy.Name = "config2"
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);

	// name can't be null
	_, err = store.Create(AppConfigModel{123, "", "validvalue"})
	testhelper.ExpectError(t, err)
	copy = existingModel;
	copy.Name = ""
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);


	// name can't be bigger than 100
	hugeName := strings.Repeat("a", 101)

	_, err = store.Create(AppConfigModel{123, hugeName, "validvalue"})
	testhelper.ExpectError(t, err)
	copy = existingModel;
	copy.Name = hugeName
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);

	// value can't be bigger than 255
	hugeValue := strings.Repeat("a", 256)

	_, err = store.Create(AppConfigModel{123, "validname", hugeValue})
	testhelper.ExpectError(t, err)
	copy = existingModel;
	copy.Value = hugeValue
	_, err = store.UpdateById(1, copy)
	testhelper.ExpectError(t,err);
	_, err = store.UpdateByName("config1", copy)
	testhelper.ExpectError(t,err);
}