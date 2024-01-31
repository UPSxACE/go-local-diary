package db_sqlite3

// GENERIC REPOSITORY PATTERN

import (
	"context"
	"database/sql"
	"errors"

	"github.com/UPSxACE/go-local-diary/app"
)

type Repository interface {
	GetContext() context.Context;
	Prepare(query string) (*sql.Stmt, error)
	Query(statement *sql.Stmt, args ...any) (*sql.Rows, error)
	QueryRow(statement *sql.Stmt, args ...any) *sql.Row
	Exec(statement *sql.Stmt, args ...any) (sql.Result, error)
	Close() []error
	Reset() error
}

func CreateRepository(appInstance *app.App[Database_Sqlite3], transactionMode bool, context context.Context) (Repository, error) {
	if !transactionMode {
		return &RepositoryNormal{db: appInstance.Database.GetInstance()}, nil
	}

	if context == nil {
		return &RepositoryNormal{}, errors.New("db_sqlite3.CreateRepository(): Can't call this function with transactionMode set to true, and yet have a nil context")
	}
	db := appInstance.Database.GetInstance();
	tx, err := db.BeginTx(context, nil)
	if err != nil {
		return &RepositoryNormal{}, err
	}
	return &RepositoryTx{db: db, tx: tx, context: context}, nil
}