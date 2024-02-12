package db_sqlite3

import (
	"testing"

	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func TestGetTables(t *testing.T) {
	dbWrapper := Init(true, ":memory:");
	db := dbWrapper.GetInstance()
	defer db.Close()

	buildQuery := func(tableName string)string{
		return "CREATE TABLE " + tableName + "(id INTEGER PRIMARY KEY, col1 VARCHAR(10), col2 INTEGER)"
	}

	tablesToCreate := []string{"table1", "table2", "table3", "table4"}
	for _, tableName := range tablesToCreate {
		_, err := db.Exec(buildQuery(tableName))
		testhelper.ExpectNoError(t, err)
	}

	allTables := dbWrapper.GetTables();
	testhelper.ExpectEqual(t, allTables, []string{"table1", "table2", "table3", "table4"})
}