package db_sqlite3

import (
	"context"
	"database/sql"
	"errors"
)

type RepositoryTx struct {
	db           *sql.DB
	tx           *sql.Tx
	context      context.Context
	cleanupQueue []func() error
	failed       bool
	done         bool
}

func (repository *RepositoryTx) queueCleanup(closeFunction func() error) {
	repository.cleanupQueue = append(repository.cleanupQueue, closeFunction)
}

func (repository *RepositoryTx) GetContext() context.Context {
	return repository.context
}

func (repository *RepositoryTx) Prepare(query string) (*sql.Stmt, error) {
	statement, err := repository.tx.PrepareContext(repository.context, query)
	if err != nil {
		repository.failed = true
		return statement, err
	}
	repository.queueCleanup(statement.Close)

	return statement, nil
}

func (repository *RepositoryTx) Query(statement *sql.Stmt, args ...any) (*sql.Rows, error) {
	rows, err := statement.QueryContext(repository.context, args...)
	if err != nil {
		repository.failed = true
		return rows, err
	}

	repository.queueCleanup(rows.Close)

	return rows, nil
}
func (repository *RepositoryTx) QueryRow(statement *sql.Stmt, args ...any) *sql.Row {
	if(repository.failed){
		return &sql.Row{};
	}
	if(repository.done){
		return &sql.Row{};
	}
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}
	
	row := statement.QueryRowContext(repository.context, args...)
	return row
}


func (repository *RepositoryTx) Exec(statement *sql.Stmt, args ...any) (sql.Result, error) {
	result, err := statement.ExecContext(repository.context, args...)
	if err != nil {
		repository.failed = true
		return result, err
	}

	return result, nil
}

func (repository *RepositoryTx) Close() []error {
	var errors []error
	if(repository.done){
		return errors;
	}
	
	if repository.failed {
		err := repository.tx.Rollback()
		if err != nil {
			errors = append(errors, err)
		}
	}
	if !repository.failed {
		err := repository.tx.Commit()
		if err != nil {
			errors = append(errors, err)
		}
	}

	repository.failed = false
	repository.done = true

	for _, cleanupFunction := range repository.cleanupQueue {
		err := cleanupFunction()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (repository *RepositoryTx) Reset() error {
	if !repository.done {
		return errors.New("db_sqlite3.*RepositoryTx.Reset(): this method should only be called after Close() is called")
	}
	newTx, err := repository.db.BeginTx(repository.context, nil)
	if err != nil {
		return err
	}
	repository.tx = newTx
	return nil
}
