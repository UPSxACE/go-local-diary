package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/utils"
)

type NoteModel struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Views      int `json:"-"` // automatic
	LastreadAt string `json:"-"`// automatic
	CreatedAt  string `json:"createdAt"` // automatic
	UpdatedAt  string `json:"updatedAt"` // automatic
	DeletedAt string `json:"-"`// automatic
	Deleted bool `json:"-"`// automatic
}

type NoteStore struct {
	db_sqlite3.StoreBase
}

func (store *NoteStore) validateModelRules(model NoteModel) (valid bool) {
	return true
}

func (store *NoteStore) validateModelCreate(model NoteModel) (valid bool, err error) {
	return true, nil
}

func (store *NoteStore) validateModelUpdate(oldModel NoteStore, newModel NoteStore) (valid bool, err error) {
	return true, nil
}

func (store *NoteStore) validateModelDelete(model NoteModel) (valid bool, err error) {
	return true, nil
}

func CreateStoreNote(appInstance *app.App[db_sqlite3.Database_Sqlite3], useTransactions bool, context context.Context) (NoteStore, error) {
	sb, err := db_sqlite3.CreateStore(appInstance, useTransactions, context)
	return NoteStore{StoreBase: sb}, err
}

func (store *NoteStore) GetFirstById(id int) (NoteModel, error) {
	query := `SELECT * FROM note WHERE id = ? ORDER BY id`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return NoteModel{}, err
	}

	row := store.Repository().QueryRow(statement, id)
	var result NoteModel
	var deletedInt int;
	row.Scan(&result.Id, &result.Title, &result.Content, &result.Views, &result.LastreadAt, &result.CreatedAt, &result.UpdatedAt, &result.DeletedAt, &deletedInt)
	result.Deleted = utils.IntToBool(deletedInt)

	if result.Id == 0 {
		return result, &db_sqlite3.EmptyQueryResults{}
	}

	return result, nil
}

func (store *NoteStore) Create(model NoteModel) (NoteModel, error) {
	dateNow := time.Now().Format("20060102")

	model.Views = 0;
	model.LastreadAt = "";
	model.CreatedAt = dateNow;
	model.UpdatedAt = dateNow;
	model.DeletedAt = "";
	model.Deleted = false;
	
	valid, err := store.validateModelCreate(model)
	if err != nil {
		return NoteModel{}, err
	}
	if !valid {
		return NoteModel{}, &db_sqlite3.InvalidModelAction{}
	}

	var query string
	if model.Id != 0 {
		query = `INSERT INTO note(id, title, content, views, lastread_at, created_at, updated_at, deleted_at, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	} else {
		query = `INSERT INTO note(title, content, views, lastread_at, created_at, updated_at, deleted_at, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	}

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return NoteModel{}, err
	}

	var res sql.Result
	if model.Id != 0 {
		res, err = store.Repository().Exec(statement, model.Id, model.Title, model.Content, model.Views, model.LastreadAt, model.CreatedAt, model.UpdatedAt, model.DeletedAt, utils.BoolToInt(model.Deleted))
	} else {
		res, err = store.Repository().Exec(statement, model.Title, model.Content, model.Views, model.LastreadAt, model.CreatedAt, model.UpdatedAt, model.DeletedAt, utils.BoolToInt(model.Deleted))
	}
	if err != nil {
		return NoteModel{}, err
	}

	insertedId, err := res.LastInsertId()
	if err != nil {
		return NoteModel{}, err
	}

	inserted, err := store.GetFirstById(int(insertedId))
	if err != nil {
		return NoteModel{}, err
	}

	return inserted, nil
}