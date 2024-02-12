package db_sqlite3

import (
	"testing"

	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

var filePath string = "./file_reader_test.sql"
var fileWithErrorPath string = "./file_reader_test_error.sql"

func TestOpenSqlFile(t *testing.T) {
	fr, err := OpenSqlFile(filePath)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, fr.totalLines, 14)
}

func TestOpenSqlFileErrorUnexistant(t *testing.T) {
	_, err := OpenSqlFile("./file_that_doesnt_exist.sql")
	testhelper.ExpectError(t, err)
}

func TestExecuteAll(t *testing.T) {
	dbWrapper := Init(true, ":memory:");
	db := dbWrapper.GetInstance()
	defer db.Close()

	fr, err := OpenSqlFile(filePath)
	testhelper.ExpectNoError(t, err)
	_, err = fr.ExecuteAll(dbWrapper)
	testhelper.ExpectNoError(t, err)

	th := testhelper.CreateTestHelper[any, any]()
	th.AddTestcase(fr.IgnoredLines(), []int{1, 2, 11})
	th.AddTestcase(fr.LinesParsed(), 14)
	th.AddTestcase(fr.TotalLines(), 14)

	for index, input := range th.GetInputs() {
		testhelper.ExpectEqual(t, input, th.GetOutput(index))
	}

	var count string
	err = db.QueryRow("SELECT count(*) FROM 'table1'").Scan(&count)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, count, "6")
}

func TestExecuteAllError(t *testing.T) {
	dbWrapper := Init(true, ":memory:");
	db := dbWrapper.GetInstance()
	defer db.Close()

	fr, err := OpenSqlFile(fileWithErrorPath)
	testhelper.ExpectNoError(t, err)
	queryThatFailed, err := fr.ExecuteAll(dbWrapper)
	testhelper.ExpectError(t, err)

	th := testhelper.CreateTestHelper[any, any]()
	th.AddTestcase(queryThatFailed, "INSERT INTO 'table1'('name') VALUES ('Maria')INSERT INTO 'table1'('name') VALUES ('Curry');")
	th.AddTestcase(fr.IgnoredLines(), []int{1, 2})
	th.AddTestcase(fr.LinesParsed(), 10)
	th.AddTestcase(fr.TotalLines(), 12)

	for index, input := range th.GetInputs() {
		testhelper.ExpectEqual(t, input, th.GetOutput(index))
	}

	var count string
	err = db.QueryRow("SELECT count(*) FROM 'table1'").Scan(&count)
	testhelper.ExpectNoError(t, err)

	testhelper.ExpectEqual(t, count, "1")
}
