package repositories

import "github.com/IgorChicherin/gophkeeper/internal/app/server/http/models"

type NotesRepository interface {
	GetNote(token string, noteID int) ([]byte, error)
	GetNotes(token string) ([]byte, error)
	CreateNote(token string, note models.CreateNoteRequest) ([]byte, error)
	UpdateNote(token string, note models.CreateNoteRequest) ([]byte, error)
}
