package models

import (
	"context"
	"database/sql"
	"reflect"
	"unicode/utf8"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
)

type AppConfigModel struct {
	Id    int
	Name  string // (unique)
	Value string
}

type AppConfigStore struct {
	db_sqlite3.StoreBase
}

func (store *AppConfigStore) validateModelRules(model AppConfigModel) (valid bool) {
	if model.Id < 0 {
		return false
	}
	if  model.Name == "" {
		return false;
	}
	if utf8.RuneCountInString(model.Name) > 100 {
		return false
	}
	if utf8.RuneCountInString(model.Value) > 255 {
		return false
	}
	return true
}

func (store *AppConfigStore) validateModelCreate(model AppConfigModel) (valid bool, err error) {
	rules := store.validateModelRules(model)
	if !rules {
		return false, nil
	}
	// id can't exist already
	count, err := store.CountById(model.Id)
	if count > 0 || err != nil {
		return false, err
	}
	// name can't exist already
	count, err = store.CountByName(model.Name)
	if count > 0 || err != nil {
		return false, err
	}
	return true, nil
}

func (store *AppConfigStore) validateModelUpdate(oldModel AppConfigModel, newModel AppConfigModel) (valid bool, err error) {
	if newModel.Id == 0 {
		return false, nil
	}
	rules := store.validateModelRules(newModel)
	if !rules {
		return false, nil
	}
	// prevent duplicated Ids
	count, err := store.CountById(newModel.Id)
	if err != nil {
		return false, err;
	}
	if count == 1 {
		modelToCompare, err := store.GetFirstById(newModel.Id)
		if err != nil {
			return false, err
		}
		if !reflect.DeepEqual(modelToCompare, oldModel) {
			return false, err
		}
	}

	// prevent duplicated Names
	count, err = store.CountByName(newModel.Name)
	if err != nil {
		return false, err;
	}
	if count == 1 {
		modelToCompare, err := store.GetFirstByName(newModel.Name)
		if err != nil {
			return false, err
		}
		if !reflect.DeepEqual(modelToCompare, oldModel) {
			return false, err
		}
	}

	return true, nil

}

func (store *AppConfigStore) validateModelDelete(model AppConfigModel) (valid bool, err error) {
	return true, nil
}

func CreateStoreAppConfig(appInstance *app.App[db_sqlite3.Database_Sqlite3], useTransactions bool, context context.Context) (AppConfigStore, error) {
	sb, err := db_sqlite3.CreateStore(appInstance, useTransactions, context)
	return AppConfigStore{StoreBase: sb}, err
}

func (store *AppConfigStore) CheckById(id int) (itExists bool, model AppConfigModel, err error) {
	modelFound, err := store.GetFirstById(id)
	if err != nil {
		_, noResultsError := err.(*db_sqlite3.EmptyQueryResults)
		if noResultsError {
			return false, AppConfigModel{}, nil
		}

		return false, AppConfigModel{}, err
	}
	return true, modelFound, nil
}

func (store *AppConfigStore) CheckByName(name string) (itExists bool, model AppConfigModel, err error) {
	modelFound, err := store.GetFirstByName(name)
	if err != nil {
		_, noResultsError := err.(*db_sqlite3.EmptyQueryResults)
		if noResultsError {
			return false, AppConfigModel{}, nil
		}

		return false, AppConfigModel{}, err
	}
	return true, modelFound, nil
}

func (store *AppConfigStore) GetFirstById(id int) (AppConfigModel, error) {
	query := `SELECT * FROM app_config WHERE id = ? ORDER BY id`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	row := store.Repository().QueryRow(statement, id)
	var result AppConfigModel
	row.Scan(&result.Id, &result.Name, &result.Value)

	if result.Id == 0 {
		return result, &db_sqlite3.EmptyQueryResults{}
	}

	return result, nil
}

func (store *AppConfigStore) GetFirstByName(name string) (AppConfigModel, error) {
	query := `SELECT * FROM app_config WHERE name = ? ORDER BY id`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	row := store.Repository().QueryRow(statement, name)
	var result AppConfigModel
	row.Scan(&result.Id, &result.Name, &result.Value)

	if result.Id == 0 {
		return result, &db_sqlite3.EmptyQueryResults{}
	}

	return result, nil
}

func (store *AppConfigStore) Create(model AppConfigModel) (AppConfigModel, error) {
	valid, err := store.validateModelCreate(model)
	if err != nil {
		return AppConfigModel{}, err
	}
	if !valid {
		return AppConfigModel{}, &db_sqlite3.InvalidModelAction{}
	}

	var query string
	if model.Id != 0 {
		query = `INSERT INTO app_config(id, name, value) VALUES (?, ?, ?)`
	} else {
		query = `INSERT INTO app_config(name, value) VALUES (?, ?)`
	}

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	var res sql.Result
	if model.Id != 0 {
		res, err = store.Repository().Exec(statement, model.Id, model.Name, model.Value)
	} else {
		res, err = store.Repository().Exec(statement, model.Name, model.Value)
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
	oldModel, err := store.GetFirstById(id)
	if err != nil {
		return AppConfigModel{}, err
	}
	valid, err := store.validateModelUpdate(oldModel, model)
	if err != nil {
		return AppConfigModel{}, err
	}
	if !valid {
		return AppConfigModel{}, &db_sqlite3.InvalidModelAction{}
	}

	query := `UPDATE app_config SET
	id = ?,
	name = ?,
	value = ?
	WHERE
	id = ?
	`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	_, err = store.Repository().Exec(statement, model.Id, model.Name, model.Value, id)
	if err != nil {
		return AppConfigModel{}, err
	}

	updated, err := store.GetFirstById(model.Id)
	if err != nil {
		return AppConfigModel{}, err
	}

	return updated, nil
}

func (store *AppConfigStore) UpdateByName(name string, model AppConfigModel) (AppConfigModel, error) {
	oldModel, err := store.GetFirstByName(name)
	if err != nil {
		return AppConfigModel{}, err
	}
	valid, err := store.validateModelUpdate(oldModel, model)
	if err != nil {
		return AppConfigModel{}, err
	}
	if !valid {
		return AppConfigModel{}, &db_sqlite3.InvalidModelAction{}
	}

	query := `UPDATE app_config SET
	id = ?,
	name = ?,
	value = ?
	WHERE
	name = ?
	`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return AppConfigModel{}, err
	}

	_, err = store.Repository().Exec(statement, model.Id, model.Name, model.Value, name)
	if err != nil {
		return AppConfigModel{}, err
	}

	updated, err := store.GetFirstByName(model.Name)
	if err != nil {
		return AppConfigModel{}, err
	}

	return updated, nil
}

func (store *AppConfigStore) DeleteById(id int) (deleted bool, err error){
	modelToDelete, err := store.GetFirstById(id)
	if err != nil {
		return false, err
	}
	valid, err := store.validateModelDelete(modelToDelete)
	if err != nil {
		return false, err
	}
	if !valid {
		return false,  &db_sqlite3.InvalidModelAction{}
	}

	query := "DELETE FROM app_config WHERE id = ?";
	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return false, err
	}

	_, err = store.Repository().Exec(statement, id)
	if err != nil {
		return false, err
	}

	count, err := store.CountById(modelToDelete.Id)
	if err != nil {
		return false, err
	}

	deleted = count == 0;
	return deleted, nil
}

func (store *AppConfigStore) DeleteByName(name string) (deleted bool, err error){
	modelToDelete, err := store.GetFirstByName(name)
	if err != nil {
		return false, err
	}
	valid, err := store.validateModelDelete(modelToDelete)
	if err != nil {
		return false, err
	}
	if !valid {
		return false,  &db_sqlite3.InvalidModelAction{}
	}

	query := "DELETE FROM app_config WHERE name = ?";
	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return false, err
	}

	_, err = store.Repository().Exec(statement, name)
	if err != nil {
		return false, err
	}

	count, err := store.CountByName(modelToDelete.Name)
	if err != nil {
		return false, err
	}

	deleted = count == 0;
	return deleted, nil
}

func (store *AppConfigStore) CountById(id int) (int, error) {
	query := `SELECT count(*) FROM app_config WHERE id = ? ORDER BY id`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return 0, err
	}

	var count int
	store.Repository().QueryRow(statement, id).Scan(&count)

	return count, nil
}

func (store *AppConfigStore) CountByName(name string) (int, error) {
	query := `SELECT count(*) FROM app_config WHERE name = ? ORDER BY id`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return 0, err
	}

	var count int
	store.Repository().QueryRow(statement, name).Scan(&count)

	return count, nil
}
