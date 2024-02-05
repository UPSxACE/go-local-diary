package services

import (
	"context"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/models"
)

var Note = NoteService{}

type NoteService struct{}

func (service *NoteService) CreateNote(app *app.App[db_sqlite3.Database_Sqlite3], context context.Context, newNote models.NoteModel) (valid bool, createdNote models.NoteModel, validationErrorMessage string, err error) {
	// FIXME validation

	store, err := models.CreateStoreNote(app, true, context)
	if err != nil {
		return false, models.NoteModel{}, "", err
	}
	defer store.Close()

	// check if exists
	model, err := store.Create(newNote)
	if err != nil {
		return false, models.NoteModel{}, "", err
	}

	return true, model, "", nil
}

func (service *NoteService) GetNotesOrderByCreateDateDesc(app *app.App[db_sqlite3.Database_Sqlite3]) (notes []models.NoteModel, err error){
	store, err := models.CreateStoreNote(app, false, nil);
	if err != nil {
		return nil, err;
	}
	defer store.Close()

	models, err := store.GetAllOrderByCreateDateDesc()
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (service *NoteService) GetNote(app *app.App[db_sqlite3.Database_Sqlite3], id int) (note models.NoteModel, err error){
	store, err := models.CreateStoreNote(app, false, nil);
	if err != nil {
		return models.NoteModel{}, err;
	}
	defer store.Close()

	model, err := store.GetFirstById(id)
	if err != nil {
		return models.NoteModel{}, err
	}

	return model, nil
}