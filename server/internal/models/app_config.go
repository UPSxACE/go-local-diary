package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
)

type AppConfigModel struct {
	Id    int
	Name  string // (unique)
	Value string
}

type AppConfigStore struct {
	repository      db_sqlite3.Repository
	transactionMode bool
}

func CreateStoreAppConfig(appInstance *app.App[db_sqlite3.Database_Sqlite3], useTransactions bool, context context.Context) (AppConfigStore, error) {
	rep, err := db_sqlite3.CreateRepository(appInstance, useTransactions, context)
	if err != nil {
		return AppConfigStore{}, err
	}

	return AppConfigStore{repository: rep, transactionMode: useTransactions}, nil
}

func (store *AppConfigStore) Close() []error {
	errs := store.repository.Close()
	return errs
}

func (store *AppConfigStore) CloseAndResetTransaction() []error {
	if !store.transactionMode {
		return []error{errors.New("db_sqlite3.*AppConfigStore.CloseAndResetTransaction(): This function should never be called when the model was created with useTransactions argument set to false")}
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

func (store *AppConfigStore) CheckById(id int) (itExists bool, model AppConfigModel, err error) {
	modelFound, err := store.GetFirstById(id)
	if err != nil {
		return false, AppConfigModel{}, err
	}
	if modelFound.Id == 0 {
		return false, AppConfigModel{}, err
	}
	return true, modelFound, nil
}

func (store *AppConfigStore) CheckByName(name string) (itExists bool, model AppConfigModel, err error) {
	modelFound, err := store.GetFirstByName(name)
	if err != nil {
		return false, AppConfigModel{}, err
	}
	if modelFound.Id == 0 {
		return false, AppConfigModel{}, err
	}
	return true, modelFound, nil
}


func (store *AppConfigStore) GetFirstById(id int) (AppConfigModel, error) {
	query := `SELECT * FROM app_config WHERE id = ?`

	statement, err := store.repository.Prepare(query)
	if err != nil {
		return AppConfigModel{}, nil
	}

	rows, err := store.repository.Query(statement, id)
	if err != nil {
		return AppConfigModel{}, nil
	}

	var result AppConfigModel
	if rows.Next() {
		rows.Scan(&result.Id, &result.Name, &result.Value)
	}

	return result, nil
}

func (store *AppConfigStore) GetFirstByName(name string) (AppConfigModel, error) {
	query := `SELECT * FROM app_config WHERE name = ?`

	statement, err := store.repository.Prepare(query)
	if err != nil {
		return AppConfigModel{}, nil
	}

	rows, err := store.repository.Query(statement, name)
	if err != nil {
		return AppConfigModel{}, nil
	}

	var result AppConfigModel
	if rows.Next() {
		rows.Scan(&result.Id, &result.Name, &result.Value)
	}

	return result, nil
}

func (store *AppConfigStore) Create(model AppConfigModel) (AppConfigModel, error) {
	var query string
	if model.Id != 0 {
		query = `INSERT INTO app_config(id, name, value) VALUES (?, ?, ?)`
	} else {
		query = `INSERT INTO app_config(name, value) VALUES (?, ?)`
	}

	statement, err := store.repository.Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	var res sql.Result
	if model.Id != 0 {
		res, err = store.repository.Exec(statement, model.Id, model.Name, model.Value)
	} else {
		res, err = store.repository.Exec(statement, model.Name, model.Value)
	}
	if err != nil {
		return AppConfigModel{}, err
	}

	insertedId, err := res.LastInsertId()
	if err != nil {
		return AppConfigModel{}, err
	}

	inserted, err := store.GetFirstById(int(insertedId))
	if err != nil {
		return AppConfigModel{}, err
	}

	return inserted, nil
}

func (store *AppConfigStore) UpdateById(id int, model AppConfigModel) (AppConfigModel, error) {
	query := `UPDATE app_config SET
	id = ?,
	name = ?,
	value = ?
	WHERE
	id = ?
	`

	statement, err := store.repository.Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	_, err = store.repository.Exec(statement, model.Id, model.Name, model.Value, id)
	if err != nil {
		return AppConfigModel{}, err
	}

	updated, err := store.GetFirstById(id)
	if err != nil {
		return AppConfigModel{}, err
	}

	return updated, nil
}

func (store *AppConfigStore) UpdateByName(name string, model AppConfigModel) (AppConfigModel, error) {
	query := `UPDATE app_config SET
	id = ?,
	name = ?,
	value = ?
	WHERE
	name = ?
	`

	statement, err := store.repository.Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	_, err = store.repository.Exec(statement, model.Id, model.Name, model.Value, name)
	if err != nil {
		return AppConfigModel{}, err
	}

	updated, err := store.GetFirstByName(name)
	if err != nil {
		return AppConfigModel{}, err
	}

	return updated, nil
}
