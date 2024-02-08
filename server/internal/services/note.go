package services

import (
	"context"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/models"
)

var Note = noteService{}

type noteService struct{}

func (service *noteService) CreateNote(app *app.App[db_sqlite3.Database_Sqlite3], context context.Context, title string, content string) (valid bool, validationErrorMessage string, err error) {
	// FIXME validation

	// AFTER VALIDATION

	newNote := models.NoteModel{
		Title:   title,
		Content: content,
	}

	store, err := models.CreateStoreNote(app, true, context)
	if err != nil {
		return false, "", err
	}
	defer store.Close()

	_, err = store.Create(newNote)
	if err != nil {
		return false, "", err
	}

	return true, "", nil
}

func (service *noteService) GetNotesOrderByCreateDateDesc(app *app.App[db_sqlite3.Database_Sqlite3], searchFilter string, includePreParsed bool) (notes []models.NoteModel, err error){
	store, err := models.CreateStoreNote(app, false, nil);
	if err != nil {
		return nil, err;
	}
	defer store.Close()

	models, err := store.GetAllOrderByCreateDateDesc(searchFilter, includePreParsed)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (service *noteService) GetNote(app *app.App[db_sqlite3.Database_Sqlite3], id int) (note models.NoteModel, err error){
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

/*
Updates user name configuration. 
*/
func (service *noteService) UpdateNote(app *app.App[db_sqlite3.Database_Sqlite3], context context.Context, id int, title string, content string) (valid bool, validationErrorMessage string, err error) {
	// FIXME validation

	// After validation

	updatedNote := models.NoteModel{
		Title:   title,
		Content: content,
	}

	store, err := models.CreateStoreNote(app, true, context)
	if err != nil {
		return false, "",  err
	}
	defer store.Close()

	_, err = store.UpdateById(id, updatedNote)
	if err != nil {
		return false, "", err
	}

	if(err != nil){
		return false, "", err
	}

	return true, "", nil;
}