package repositories

import "github.com/IgorChicherin/gophkeeper/internal/app/server/http/models"

type HTTPNotesRepository struct {
	BaseURL string
}

func (r HTTPNotesRepository) GetNote(token string, noteID int) ([]byte, error) {
	return []byte{}, nil
}

func (r HTTPNotesRepository) GetNotes(token string) ([]byte, error) {
	return []byte{}, nil
}

func (r HTTPNotesRepository) CreateNote(token string, note models.CreateNoteRequest) ([]byte, error) {
	return []byte{}, nil
}

func (r HTTPNotesRepository) UpdateNote(token string, note models.CreateNoteRequest) ([]byte, error) {
	return []byte{}, nil
}
