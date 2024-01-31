package db_sqlite3

import (
	"context"
	"database/sql"
	"testing"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func getTestAppInstance() *app.App[Database_Sqlite3] {
	app := app.App[Database_Sqlite3]{
		Database: Init(),
	}
	return &app
}

func getTestAppInstanceAndDb() (*app.App[Database_Sqlite3], *sql.DB){
	app := getTestAppInstance()
	db := app.Database.GetInstance();
	return app, db
}

func getTestContext() context.Context{
	return context.Background();
}

func TestCreateRepositoryNormal(t *testing.T) {
	app := getTestAppInstance();

	rep,err := CreateRepository(app, false, getTestContext());
	testhelper.ExpectNoError(t, err)

	// context should be nil because normal repositories don't need it
	context := rep.GetContext()
	testhelper.ExpectEqual(t, context, nil)
}

func TestCreateRepositoryTx(t *testing.T) {
	app := getTestAppInstance();

	rep,err := CreateRepository(app, true, getTestContext());
	testhelper.ExpectNoError(t, err)
	
	// context should NOT be nil because repositories that use transactions need it
	context := rep.GetContext()
	testhelper.ExpectEqual(t, context, context)
}

func TestCreateRepositoryTxErrorNoContext(t *testing.T){
	app := getTestAppInstance();

	_,err := CreateRepository(app, true, nil);
	testhelper.ExpectError(t, err)
}