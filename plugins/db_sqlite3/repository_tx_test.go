package db_sqlite3

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func TestUseRepositoryTx(t *testing.T) {
	app, db := getTestAppInstanceAndDb()
	defer db.Close()

	ctx := context.Background()

	rep, err := CreateRepository(app, true, ctx)
	testhelper.ExpectNoError(t, err)

	trep := testhelper.ExpectType[*RepositoryTx](t, rep)

	queryCreate := "CREATE TABLE table1(id INTEGER PRIMARY KEY, name VARCHAR(15))"
	stmtCreate, err := trep.Prepare(queryCreate)
	testhelper.ExpectNoError(t, err)

	_, err = trep.Exec(stmtCreate)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, len(trep.cleanupQueue), 1)

	queryCount := "SELECT count(*) FROM table1"
	stmtCount, err := trep.Prepare(queryCount)
	testhelper.ExpectNoError(t, err)

	querySelectAll := "SELECT * FROM table1 ORDER BY id"
	stmtSelectAll, err := trep.Prepare(querySelectAll)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, len(trep.cleanupQueue), 3)

	queryInsert := "INSERT INTO table1(name) VALUES (?)"
	stmtInsert, err := trep.Prepare(queryInsert)
	testhelper.ExpectNoError(t, err)

	names := []string{"John", "Mario", "Maria", "Peter"}
	for _, name := range names {
		_, err := trep.Exec(stmtInsert, name)
		testhelper.ExpectNoError(t, err)
	}

	var rowCount string
	row := trep.QueryRow(stmtCount)
	row.Scan(&rowCount)

	queriedNames := []string{}
	rows, err := trep.Query(stmtSelectAll)
	testhelper.ExpectNoError(t, err)
	for rows.Next() {
		var id int
		var newName string
		rows.Scan(&id, &newName)
		queriedNames = append(queriedNames, newName)
	}

	testhelper.ExpectEqual(t, trep.failed, false)
	testhelper.ExpectEqual(t, trep.done, false)
	testhelper.ExpectEqual(t, rowCount, strconv.Itoa(len(names)))
	testhelper.ExpectEqual(t, names, queriedNames)
	testhelper.ExpectDifferent(t, trep.openRows, (*sql.Rows)(nil))
	testhelper.ExpectEqual(t, len(trep.cleanupQueue), 4)

	errs := rep.Close()
	if len(errs) > 1 {
		t.Fatal(errs[0])
	}

	testhelper.ExpectEqual(t, trep.failed, false)
	testhelper.ExpectEqual(t, trep.done, true)
	testhelper.ExpectEqual(t, trep.openRows, (*sql.Rows)(nil))
	testhelper.ExpectEqual(t, len(trep.cleanupQueue), 0)

	err = rep.Reset()

	testhelper.ExpectEqual(t, trep.failed, false)
	testhelper.ExpectEqual(t, trep.done, false)
	testhelper.ExpectEqual(t, err, nil)

	// Check if things were commited
	queryFind := "SELECT name FROM table1 WHERE name='Peter'"
	stmtFind, err := trep.Prepare(queryFind)
	testhelper.ExpectNoError(t, err)
	var peterName string;
	trep.QueryRow(stmtFind).Scan(&peterName)

	testhelper.ExpectEqual(t, peterName, "Peter")
}

func TestRepositoryTxInvalidQueries(t *testing.T) {
	app, db := getTestAppInstanceAndDb()
	defer db.Close()

	ctx := context.Background()

	rep, err := CreateRepository(app, true, ctx)
	testhelper.ExpectNoError(t, err)

	// valid
	queryCreate := "CREATE TABLE table1(id INTEGER PRIMARY KEY, name VARCHAR(15))"
	stmtCreate, err := rep.Prepare(queryCreate)
	testhelper.ExpectNoError(t, err)
	_, err = rep.Exec(stmtCreate)
	testhelper.ExpectNoError(t, err)

	// Reset shall raise error when called before Close
	err = rep.Reset()
	testhelper.ExpectError(t, err);

	errs := rep.Close()
	if len(errs) > 0 {
		t.Fatal(errs[0])
	}

	err = rep.Reset()
	testhelper.ExpectNoError(t, err)

	// Table should be saved and repository ready to be used again

	// valid
	queryInsert := "INSERT INTO table1(name) VALUES (?)"
	stmtInsert, err := rep.Prepare(queryInsert)
	testhelper.ExpectNoError(t, err)

	// valid
	_, err = rep.Exec(stmtInsert, "Roger")
	testhelper.ExpectNoError(t, err)

	// valid (Roger should be in database)
	queryFind := "SELECT name FROM table1 WHERE name='Roger'"
	stmtFind, err := rep.Prepare(queryFind)
	testhelper.ExpectNoError(t, err)
	var rogerName string;
	rep.QueryRow(stmtFind).Scan(&rogerName)
	testhelper.ExpectEqual(t, rogerName, "Roger")

	// invalid query
	invalidQuery := "INSERT INTO ()"
	_, err = rep.Prepare(invalidQuery)
	testhelper.ExpectError(t, err)

	trep := testhelper.ExpectType[*RepositoryTx](t, rep)
	testhelper.ExpectEqual(t, trep.failed, true)

	// After a failed query, the methods should raise an error when called before Close()
	_, err = rep.Exec(stmtFind)
	testhelper.ExpectError(t, err);
	_, err = rep.Prepare(queryFind)
	testhelper.ExpectError(t, err);
	_, err = rep.Query(stmtFind)
	testhelper.ExpectError(t, err);
	// Query Row doesn't raise errors, but it should return something empty instead
	returnedRow := rep.QueryRow(stmtFind)
	testhelper.ExpectEqual(t, returnedRow, &sql.Row{}); 

	// now that it failed, on close it should rollback everything
	errs = rep.Close()
	if len(errs) > 0 {
		t.Fatal(errs[0])
	}
	err = rep.Reset()
	testhelper.ExpectNoError(t, err)

	// since it was rolled back, roger should not be found
	rogerName = ""
	stmtFind, err = rep.Prepare("SELECT name FROM table1 WHERE name='Roger'")
	testhelper.ExpectNoError(t, err)
	rep.QueryRow(stmtFind).Scan(&rogerName)
	testhelper.ExpectEqual(t, rogerName, "")

	// despite everything the database should still work and be usable
	validQuery := "INSERT INTO table1(name) VALUES (?)"
	validStmt, err := rep.Prepare(validQuery)
	testhelper.ExpectNoError(t, err)
	rep.Exec(validStmt, "Ed")
	errs = rep.Close()
	if len(errs) > 0{
		t.Fatal(errs[0])
	}
	err = rep.Reset()
	testhelper.ExpectNoError(t, err)

	var edName string
	queryFindEd := "SELECT name FROM table1 WHERE name='Ed'"
	stmtFindEd, err := rep.Prepare(queryFindEd)
	testhelper.ExpectNoError(t, err)
	rep.QueryRow(stmtFindEd).Scan(&edName)

	testhelper.ExpectEqual(t, edName, "Ed")

	errs = rep.Close()
	if len(errs) > 1 {
		t.Fatal(errs[0])
	}

	// after closed, trying to attempt any types of query should also raise error
	_, err = rep.Exec(stmtFind)
	testhelper.ExpectError(t, err);
	_, err = rep.Prepare(queryFind)
	testhelper.ExpectError(t, err);
	_, err = rep.Query(stmtFind)
	testhelper.ExpectError(t, err);
	// Query Row doesn't raise errors, but it should return something empty instead
	returnedRow = rep.QueryRow(stmtFind)
	testhelper.ExpectEqual(t, returnedRow, &sql.Row{}); 
}