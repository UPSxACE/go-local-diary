package db_sqlite3

import (
	"context"
	"database/sql"
	"errors"
)

type RepositoryNormal struct {
	db           *sql.DB
	cleanupQueue []func() error
}

func (repository *RepositoryNormal) queueCleanup(closeFunction func() error) {
	repository.cleanupQueue = append(repository.cleanupQueue, closeFunction)
}

func (repository *RepositoryNormal) GetContext() context.Context {
	return nil
}

func (repository *RepositoryNormal) Prepare(query string) (*sql.Stmt, error) {
	statement, err := repository.db.Prepare(query)
	if err != nil {
		return statement, err
	}
	repository.queueCleanup(func() error {
		return statement.Close()
	})

	return statement, nil
}

func (repository *RepositoryNormal) Query(statement *sql.Stmt, args ...any) (*sql.Rows, error) {
	rows, err := statement.Query(args)
	if err != nil {
		return rows, err
	}
	repository.queueCleanup(func() error {
		return rows.Close()
	})

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
	result, err := statement.Exec(args)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (repository *RepositoryNormal) Close() []error {
	var errors []error

	for _, closeFunction := range repository.cleanupQueue {
		err := closeFunction()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (repository *RepositoryNormal) Reset() error {
	return errors.New("db_sqlite3.*RepositoryNormal.Reset(): this method should never be called in normal repositories")
}
