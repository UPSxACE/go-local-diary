package db_sqlite3

import (
	"database/sql"
	"os"
	"strings"

	"github.com/UPSxACE/go-local-diary/app"
)

/*
This struct represents an entity that is being used for reading an sql file,
and that can be used to execute the commands on it.

IMPORTANT: Make sure every command in the file ends with ';'
*/
type SqlFileReader struct {
	data []string;
	ignoredLines []int
	linesParsed int;
	totalLines int;
}

func (fr *SqlFileReader) IgnoredLines() []int{
	return fr.ignoredLines;
}
func (fr *SqlFileReader) LinesParsed() int{
	return fr.linesParsed;
}
func (fr *SqlFileReader) TotalLines() int{
	return fr.totalLines;
}

/*
Reads an sql file in the given path and returns a SqlFileReader instance,
ready to execute the instructions inside it.
*/
func OpenSqlFile(path string) (*SqlFileReader, error) {
	data, err := os.ReadFile(path);
	if(err != nil){
		return &SqlFileReader{}, err;
	}

	// NOTE: Does not support converting old mac \r line breaks
	dataString := string(data)
	dataStringUniversal := strings.ReplaceAll(dataString,"\r\n", "\n") 
	dataStringInLines  := strings.Split(dataStringUniversal, "\n")
	
	return &SqlFileReader{data:dataStringInLines, totalLines: len(dataStringInLines)}, nil;
}

/*
Execute all commands inside the sql file that was read.
*/
func (fr *SqlFileReader) ExecuteAll(dbInstance *sql.DB) (queryThatFailed string, err error){
	var nextQuery string = "";

	for i, line := range fr.data {
		fr.linesParsed++;

		// ignore comments
		if(strings.HasPrefix(line, "--")){
			fr.ignoredLines = append(fr.ignoredLines, i+1)
			continue;
		}

		nextQuery+=line;

		// check if its the end of the next query
		if(strings.HasSuffix(line, ";")){
			_, err := dbInstance.Exec(nextQuery)
			if err != nil{
				return nextQuery, err
			}
			// fmt.Println("Query Executed: ", nextQuery)
			nextQuery = "";
		}
		
	}

	// fmt.Printf("Parsed: %v/%v lines\n", fr.linesParsed, fr.totalLines)
	return "", nil
}

/*
Execute all commands inside the sql file that was read.
Extract the database instance from the given app instance.
*/
func (fr *SqlFileReader) ExecuteAllFromApp(appInstance *app.App[Database_Sqlite3]) (queryThatFailed string, err error ){
	db := appInstance.Database.GetInstance()
	return fr.ExecuteAll(db)
}

