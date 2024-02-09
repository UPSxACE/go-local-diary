package services

import (
	"context"

	"github.com/UPSxACE/go-local-diary/server/models"
)

func (services *Services) CreateNote(context context.Context, title string, content string) (valid bool, validationErrorMessage string, err error) {
	// FIXME validation

	// AFTER VALIDATION

	newNote := models.NoteModel{
		Title:   title,
		Content: content,
	}

	store, err := models.CreateStoreNote(services.database, true, context)
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

func (services *Services) GetNotesOrderByCreateDateDesc(searchFilter string, includePreParsed bool) (notes []models.NoteModel, err error){
	store, err := models.CreateStoreNote(services.database, false, nil);
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

func (service *Services) GetNote(id int) (note models.NoteModel, err error){
	store, err := models.CreateStoreNote(service.database, false, nil);
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
func (services *Services) UpdateNote(context context.Context, id int, title string, content string) (valid bool, validationErrorMessage string, err error) {
	// FIXME validation

	// After validation

	updatedNote := models.NoteModel{
		Title:   title,
		Content: content,
	}

	store, err := models.CreateStoreNote(services.database, true, context)
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