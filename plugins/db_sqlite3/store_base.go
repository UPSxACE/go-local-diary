package db_sqlite3;

import (
	"context"
	"errors"

	"github.com/UPSxACE/go-local-diary/app"
)

type StoreBase struct {
	repository      Repository
	transactionMode bool
}

func CreateStore(appInstance *app.App[Database_Sqlite3], useTransactions bool, context context.Context) (StoreBase, error){
	rep, err := CreateRepository(appInstance, useTransactions, context)
	if err != nil {
		return StoreBase{}, err
	}

	return StoreBase{repository: rep, transactionMode: useTransactions}, nil
}

func (store *StoreBase) Repository() Repository {
	return store.repository
}

func (store *StoreBase) TransactionMode() bool {
	return store.transactionMode
}


func (store *StoreBase) Close() []error {
	errs := store.repository.Close()
	return errs
}

func (store *StoreBase) CloseAndResetTransaction() []error {
	if !store.transactionMode {
		return []error{errors.New("(store).CloseAndResetTransaction(): This function should never be called when the model was created with useTransactions argument set to false")}
	}
	errs := store.repository.Close()
	if len(errs) != 0 {
		return errs
	}
	err := store.repository.Reset()
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}