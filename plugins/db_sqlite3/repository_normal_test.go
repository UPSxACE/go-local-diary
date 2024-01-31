package db_sqlite3

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func TestUseRepositoryNormal(t *testing.T) {
	app, db := getTestAppInstanceAndDb()
	defer db.Close()

	rep, err := CreateRepository(app, false, nil)
	testhelper.ExpectNoError(t, err)

	nrep := testhelper.ExpectType[*RepositoryNormal](t, rep)

	queryCreate := "CREATE TABLE table1(id INTEGER PRIMARY KEY, name VARCHAR(15))"
	stmtCreate, err := nrep.Prepare(queryCreate)
	testhelper.ExpectNoError(t, err)

	_, err = nrep.Exec(stmtCreate)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, len(nrep.cleanupQueue), 1)

	queryCount := "SELECT count(*) FROM table1"
	stmtCount, err := nrep.Prepare(queryCount)
	testhelper.ExpectNoError(t, err)

	querySelectAll := "SELECT * FROM table1 ORDER BY id"
	stmtSelectAll, err := nrep.Prepare(querySelectAll)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, len(nrep.cleanupQueue), 3)

	queryInsert := "INSERT INTO table1(name) VALUES (?)"
	stmtInsert, err := nrep.Prepare(queryInsert)
	testhelper.ExpectNoError(t, err)

	names := []string{"John", "Mario", "Maria", "Peter"}
	for _, name := range names {
		_, err := nrep.Exec(stmtInsert, name)
		testhelper.ExpectNoError(t, err)
	}

	var rowCount string
	row := nrep.QueryRow(stmtCount)
	row.Scan(&rowCount)

	queriedNames := []string{}
	rows, err := nrep.Query(stmtSelectAll)
	testhelper.ExpectNoError(t, err)
	for rows.Next() {
		var id int
		var newName string
		rows.Scan(&id, &newName)
		queriedNames = append(queriedNames, newName)
	}

	testhelper.ExpectEqual(t, rowCount, strconv.Itoa(len(names)))
	testhelper.ExpectEqual(t, names, queriedNames)
	testhelper.ExpectDifferent(t, nrep.openRows, (*sql.Rows)(nil))
	testhelper.ExpectEqual(t, len(nrep.cleanupQueue), 4)

	errs := rep.Close()
	if len(errs) > 1 {
		t.Fatal(errs[0])
	}

	testhelper.ExpectEqual(t, nrep.openRows, (*sql.Rows)(nil))
	testhelper.ExpectEqual(t, len(nrep.cleanupQueue), 0)

	err = rep.Reset()
	testhelper.ExpectError(t, err)
}

func TestRepositoryNormalInvalidQueries(t *testing.T) {
	app, db := getTestAppInstanceAndDb()
	defer db.Close()

	rep, err := CreateRepository(app, false, nil)
	testhelper.ExpectNoError(t, err)

	// valid
	queryCreate := "CREATE TABLE table1(id INTEGER PRIMARY KEY, name VARCHAR(15))"
	stmtCreate, err := rep.Prepare(queryCreate)
	testhelper.ExpectNoError(t, err)
	_, err = rep.Exec(stmtCreate)
	testhelper.ExpectNoError(t, err)

	// valid
	queryInsert := "INSERT INTO table1(name) VALUES (?)"
	stmtInsert, err := rep.Prepare(queryInsert)
	testhelper.ExpectNoError(t, err)

	// valid
	_, err = rep.Exec(stmtInsert, "John")
	testhelper.ExpectNoError(t, err)

	// invalid query
	invalidQuery := "INSERT INTO ()"
	_, err = rep.Prepare(invalidQuery)

	testhelper.ExpectError(t, err)

	// invalid query
	// NOTE: Using the wrong function sometimes won't return error,
	// specially if the statement is valid. It will just do nothing,
	// or possibly even execute it. The behavior seems a bit unpredictable.
	rep.Query(stmtCreate)

	
	// despite everything the database should still work
	validQuery :=  "INSERT INTO table1(name) VALUES (?)"
	validStmt,err := rep.Prepare(validQuery)
	testhelper.ExpectNoError(t, err)
	rep.Exec(validStmt, "Ed")

	var johnName string
	queryFindJohn := "SELECT name FROM table1 WHERE name='John'"
	stmtFindJohn, err := rep.Prepare(queryFindJohn)
	testhelper.ExpectNoError(t, err)
	rep.QueryRow(stmtFindJohn).Scan(&johnName)
	var edName string
	queryFindEd := "SELECT name FROM table1 WHERE name='Ed'"
	stmtFindEd, err := rep.Prepare(queryFindEd)
	testhelper.ExpectNoError(t, err)
	rep.QueryRow(stmtFindEd).Scan(&edName)
	
	testhelper.ExpectEqual(t, johnName, "John")
	testhelper.ExpectEqual(t, edName, "Ed")

	errs := rep.Close()
	if len(errs) > 1 {
		t.Fatal(errs[0])
	}
}