package models

import (
	"context"
	"errors"
	"time"

	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type NoteDifModel struct {
	Id       int    `json:"id"`
	NoteId   int    `json:"noteId"`   // automatic
	Content  string `json:"content"`  // automatic
	EditedAt string `json:"editedAt"` // automatic
}

type NoteDifStore struct {
	db_sqlite3.StoreBase
}

func CreateStoreNoteDif(database *db_sqlite3.Database_Sqlite3, useTransactions bool, context context.Context) (NoteDifStore, error) {
	sb, err := db_sqlite3.CreateStore(database, useTransactions, context)
	return NoteDifStore{StoreBase: sb}, err
}

//FIXME func GetLatestById

func (store *NoteDifStore) CheckLatestByNoteId(note_id int) (itExists bool, model NoteDifModel, err error) {
	modelFound, err := store.GetLatestByNoteId(note_id)
	if err != nil {
		_, noResultsError := err.(*db_sqlite3.EmptyQueryResults)
		if noResultsError {
			return false, NoteDifModel{}, nil
		}

		return false, NoteDifModel{}, err
	}
	return true, modelFound, nil
}

func (store *NoteDifStore) GetLatestByNoteId(noteId int) (NoteDifModel, error) {
	query := `SELECT * FROM note_dif WHERE note_id = ? ORDER BY edited_at DESC`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return NoteDifModel{}, err
	}

	row := store.Repository().QueryRow(statement, noteId)
	var result NoteDifModel
	row.Scan(&result.Id, &result.NoteId, &result.Content, &result.EditedAt)

	if result.Id == 0 {
		return result, &db_sqlite3.EmptyQueryResults{}
	}

	return result, nil
}

func (store *NoteDifStore) GetAllByNoteIdOrderByOldestFirst(noteId int) ([]NoteDifModel, error) {
	query := `SELECT * FROM note_dif WHERE note_id == ? ORDER BY edited_at ASC`

	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := store.Repository().Query(statement, noteId)
	if err != nil {
		return nil, err
	}

	var result []NoteDifModel
	for rows.Next() {
		model := NoteDifModel{}
		rows.Scan(&model.Id, &model.NoteId, &model.Content, &model.EditedAt)
		result = append(result, model)
	}

	return result, nil
}

func (store *NoteDifStore) RegisterChange(oldModel *NoteModel, validNoteModelAboutToUpdateContent *NoteModel) (updateDate string, difModel NoteDifModel, err error) {
	if oldModel.Id != validNoteModelAboutToUpdateContent.Id {
		return "", NoteDifModel{}, errors.New("ids do not match")
	}

	noteId := oldModel.Id

	noteDifs, err := store.GetAllByNoteIdOrderByOldestFirst(noteId)
	if err != nil {
		return "", NoteDifModel{}, err
	}
	first := len(noteDifs) == 0

	dmp := diffmatchpatch.New()

	var finalContent string
	if first {
		finalContent = oldModel.Content
	}
	if !first {
		for i, noteDif := range noteDifs {
			if i == 0 {
				finalContent = noteDif.Content
			}
			if i != 0 {
				// Rebuild the string
				// (if this process becomes too heavy, simply cache the newest note edits
				// instead of rebuilding it all from scratch each edit, 
			    // or make this a secondary thing that happens asynchronously)
				patches, err := dmp.PatchFromText(noteDif.Content)
				if err != nil {
					return "", NoteDifModel{}, err
				}
				finalContent_, patchesBool := dmp.PatchApply(patches, finalContent)
				for _, patchBool := range patchesBool {
					if !patchBool {
						return "", NoteDifModel{}, err
					}
				}
				finalContent = finalContent_
			}
		}

		difs := dmp.DiffMain(finalContent, validNoteModelAboutToUpdateContent.Content, false)
		patches := dmp.PatchMake(difs)
		patchesText := dmp.PatchToText(patches)
		finalContent = patchesText
	}

	model := NoteDifModel{}
	// Automatic fields
	model.NoteId = noteId
	dateNow := time.Now().Format("20060102")
	model.EditedAt = dateNow
	model.Content = finalContent

	query := `INSERT INTO note_dif(note_id, content, edited_at) VALUES (?, ?, ?)`
	statement, err := store.Repository().Prepare(query)
	if err != nil {
		return "", NoteDifModel{}, err
	}

	res, err := store.Repository().Exec(statement, model.NoteId, model.Content, model.EditedAt)
	if err != nil {
		return "", NoteDifModel{}, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return "", NoteDifModel{}, err
	}

	inserted, err := store.GetLatestByNoteId(int(noteId))
	if err != nil {
		return "", NoteDifModel{}, err
	}

	return inserted.EditedAt, inserted, nil
}
