package db_sqlite3

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/UPSxACE/go-local-diary/app"
)

type SqlFileReader struct {
	data []string;
	IgnoredLines []int
	LinesParsed int;
	TotalLines int;
}

// utilizing go routines is worth?
func (fr *SqlFileReader) ExecuteAll(dbInstance *sql.DB) (err error, queryThatFailed string){
	var nextQuery string = "";

	for i, line := range fr.data {
		fr.LinesParsed++;

		// ignore comments
		if(strings.HasPrefix(line, "--")){
			fr.IgnoredLines = append(fr.IgnoredLines, i+1)
			continue;
		}

		nextQuery+=line;

		// check if its the end of the next query
		if(strings.HasSuffix(line, ";")){
			_, err := dbInstance.Exec(nextQuery)
			if err != nil{
				return err, nextQuery
			}
			fmt.Println("Query Executed: ", nextQuery)
			nextQuery = "";
		}
		
	}

	fmt.Printf("Parsed: %v/%v lines\n", fr.LinesParsed, fr.TotalLines)
	return nil, "";
}

func (fr *SqlFileReader) ExecuteAllFromApp(appInstance *app.App[Database_Sqlite3]) (err error, queryThatFailed string){
	db := appInstance.Database.GetInstance()
	return fr.ExecuteAll(db)
}

func OpenSqlFile(path string) (*SqlFileReader, error) {
	data, err := os.ReadFile(path);
	if(err != nil){
		return &SqlFileReader{}, err;
	}
	
	dataString := string(data)
	dataStringUniversal := strings.ReplaceAll(dataString,"\r\n", "\n")
	dataStringInLines  := strings.Split(dataStringUniversal, "\n")
	
	return &SqlFileReader{data:dataStringInLines, TotalLines: len(dataStringInLines)}, nil;
}