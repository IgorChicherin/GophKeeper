package repositories

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/client/http/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
	"github.com/IgorChicherin/gophkeeper/internal/shared/models"
)

type NotesRepository interface {
	GetNote(token string, noteID int) ([]byte, error)
	GetNotes(token string) ([]byte, error)
	CreateNote(token string, note models.CreateNoteRequest) ([]byte, error)
	UpdateNote(token string, noteID int, note models.CreateNoteRequest) ([]byte, error)
	DeleteNote(token string, noteID int) error
}

func NewHTTPNoteRepository(baseUrl string, client httpclient.HTTPClientSync) NotesRepository {
	return repositories.HTTPNotesRepository{baseUrl, client}
}
