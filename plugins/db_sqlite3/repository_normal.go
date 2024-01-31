package db_sqlite3

import (
	"context"
	"database/sql"
	"errors"
)

type RepositoryNormal struct {
	db           *sql.DB
	cleanupQueue []func() error
	openRows     *sql.Rows // prevent leaving rows open
}

func (repository *RepositoryNormal) queueCleanup(closeFunction func() error) {
	repository.cleanupQueue = append(repository.cleanupQueue, closeFunction)
}

func (repository *RepositoryNormal) GetContext() context.Context {
	return nil
}

func (repository *RepositoryNormal) Prepare(query string) (*sql.Stmt, error) {
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}

	statement, err := repository.db.Prepare(query)
	if err != nil {
		return statement, err
	}
	repository.queueCleanup(statement.Close)

	return statement, nil
}

func (repository *RepositoryNormal) Query(statement *sql.Stmt, args ...any) (*sql.Rows, error) {
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}
	
	rows, err := statement.Query(args...)
	if err != nil {
		return rows, err
	}

	// prevent leaving rows open
	repository.openRows = rows

	return rows, nil
}

func (repository *RepositoryNormal) QueryRow(statement *sql.Stmt, args ...any) *sql.Row {
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}

	row := statement.QueryRow(args...)
	return row
}

func (repository *RepositoryNormal) Exec(statement *sql.Stmt, args ...any) (sql.Result, error) {
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}

	result, err := statement.Exec(args...)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (repository *RepositoryNormal) Close() []error {
	var errors []error
	var newQueue []func()error;

	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}

	for _, closeFunction := range repository.cleanupQueue {
		err := closeFunction()
		if err != nil {
			newQueue = append(newQueue, closeFunction)
			errors = append(errors, err)
		}
	}

	repository.cleanupQueue = newQueue;
	return errors
}

func (repository *RepositoryNormal) Reset() error {
	return errors.New("db_sqlite3.*RepositoryNormal.Reset(): this method should never be called in normal repositories (you don't need to reset it to use it again)")
}
