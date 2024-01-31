package db_sqlite3

import (
	"context"
	"testing"
	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func TestCreateStore(t *testing.T) {
	app, db := getTestAppInstanceAndDb()
	defer db.Close()

	// Create normal store
	nstore, err := CreateStore(app, false, nil)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, nstore.TransactionMode(), false)
	testhelper.ExpectDifferent(t, nstore.Repository(), (Repository)(nil))
	// CloseAndResetTransaction on normal store should raise error
	errs := nstore.CloseAndResetTransaction()
	testhelper.ExpectDifferent(t, len(errs), 0)	
	errs = nstore.Close()
	testhelper.ExpectEqual(t, len(errs), 0)


	// Create transaction store
	ctx := context.Background()
	tstore, err := CreateStore(app, true, ctx)
	testhelper.ExpectNoError(t, err)
	testhelper.ExpectEqual(t, tstore.TransactionMode(), true)
	testhelper.ExpectDifferent(t, tstore.Repository(), (Repository)(nil))

	errs = tstore.CloseAndResetTransaction()
	testhelper.ExpectEqual(t, len(errs), 0)	
	errs = tstore.CloseAndResetTransaction()
	testhelper.ExpectEqual(t, len(errs), 0)	
	errs = tstore.CloseAndResetTransaction()
	testhelper.ExpectEqual(t, len(errs), 0)	
	errs = tstore.Close()
	testhelper.ExpectEqual(t, len(errs), 0)

	// the remaining things must be tested in the structs that implement this one
}