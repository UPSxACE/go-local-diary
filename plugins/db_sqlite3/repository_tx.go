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
	openRows     *sql.Rows // prevent leaving rows open
}

// TODO: Add a wrapper to the statements returned by these functions
//       that limits the user to only use those in the interface Repository
//       and so it detects when one of them is used after the connection was 
//       closed already

type emptyResult struct {
}
func (er *emptyResult) LastInsertId() (int64, error){
	panic("This is an empty sql.Result, originated by an error.")
}
func (er *emptyResult) RowsAffected() (int64, error){
	panic("This is an empty sql.Result, originated by an error.")
}

func (repository *RepositoryTx) queueCleanup(closeFunction func() error) {
	repository.cleanupQueue = append(repository.cleanupQueue, closeFunction)
}

func (repository *RepositoryTx) GetContext() context.Context {
	return repository.context
}

func (repository *RepositoryTx) Prepare(query string) (*sql.Stmt, error) {
	if(repository.failed){
		return &sql.Stmt{},errors.New("the transaction has failed. Please use Close() to rollback the transaction, and then Reset() to start a new one");
	}
	if(repository.done){
		return &sql.Stmt{},errors.New("this transaction was closed over already. Please use Reset() to start a new one");
	}

	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}

	statement, err := repository.tx.PrepareContext(repository.context, query)
	if err != nil {
		repository.failed = true
		return statement, err
	}
	repository.queueCleanup(statement.Close)

	return statement, nil
}

func (repository *RepositoryTx) Query(statement *sql.Stmt, args ...any) (*sql.Rows, error) {
	if(repository.failed){
		return &sql.Rows{},errors.New("the transaction has failed. Please use Close() to rollback the transaction, and then Reset() to start a new one");
	}
	if(repository.done){
		return &sql.Rows{},errors.New("this transaction was closed over already. Please use Reset() to start a new one");
	}
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}
	
	rows, err := statement.QueryContext(repository.context, args...)
	if err != nil {
		repository.failed = true
		return rows, err
	}

	// prevent leaving rows open
	repository.openRows = rows

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
	if(repository.failed){
		return &emptyResult{},errors.New("the transaction has failed. Please use Close() to rollback the transaction, and then Reset() to start a new one");
	}
	if(repository.done){
		return &emptyResult{},errors.New("this transaction was closed over already. Please use Reset() to start a new one");
	}
	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}
	
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

	// prevent leaving rows open
	if(repository.openRows != nil){
		repository.openRows.Close()
		repository.openRows = nil;
	}

	var newQueue []func()error;

	for _, cleanupFunction := range repository.cleanupQueue {
		err := cleanupFunction()
		if err != nil {
			newQueue = append(newQueue, cleanupFunction)
			errors = append(errors, err)
		}
	}

	repository.cleanupQueue = newQueue;
	
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
	repository.done = false;
	repository.tx = newTx
	return nil
}
