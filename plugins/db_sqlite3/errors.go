package db_sqlite3;

type EmptyQueryResults struct{}
func (e *EmptyQueryResults) Error() string {
	return "query returned 0 results";
}

type InvalidModelAction struct{}
func (e *InvalidModelAction) Error() string {
	return "the model data, or the attempted action with it is invalid";
}