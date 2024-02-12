package db_sqlite3

import (
	"context"
	"testing"

	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func getTestContext() context.Context{
	return context.Background();
}

func TestCreateRepositoryNormal(t *testing.T) {
	dbWrapper := Init(true, ":memory:");
	db := dbWrapper.GetInstance()
	defer db.Close()

	rep,err := CreateRepository(dbWrapper, false, getTestContext());
	testhelper.ExpectNoError(t, err)

	// context should be nil because normal repositories don't need it
	context := rep.GetContext()
	testhelper.ExpectEqual(t, context, nil)
}

func TestCreateRepositoryTx(t *testing.T) {
	dbWrapper := Init(true, ":memory:");
	db := dbWrapper.GetInstance()
	defer db.Close()

	rep,err := CreateRepository(dbWrapper, true, getTestContext());
	testhelper.ExpectNoError(t, err)
	
	// context should NOT be nil because repositories that use transactions need it
	context := rep.GetContext()
	testhelper.ExpectEqual(t, context, context)
}

func TestCreateRepositoryTxErrorNoContext(t *testing.T){
	dbWrapper := Init(true, ":memory:");
	db := dbWrapper.GetInstance()
	defer db.Close()

	_,err := CreateRepository(dbWrapper, true, nil);
	testhelper.ExpectError(t, err)
}