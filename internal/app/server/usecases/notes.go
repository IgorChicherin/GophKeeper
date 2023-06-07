package usecases

import (
	"encoding/base64"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
)

type NotesUseCase interface {
	CreateUserNote(user models.User, req models.CreateNoteRequest) (models.Note, error)
	GetNote(user models.User, noteID int) (models.CreateNoteRequest, error)
	GetUserNotes(user models.User) ([]models.Note, error)
}

type notesUseCase struct {
	NotesRepo repositories.NotesRepository
	Decrypter crypto509.Decrypter
}

func NewNotesUseCase(notesRepo repositories.NotesRepository, dec crypto509.Decrypter) NotesUseCase {
	return notesUseCase{NotesRepo: notesRepo, Decrypter: dec}
}

func (n notesUseCase) CreateUserNote(user models.User, req models.CreateNoteRequest) (models.Note, error) {
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return models.Note{}, err
	}
	return n.NotesRepo.CreateNote(models.Note{
		UserID:   user.UserID,
		Data:     data,
		Metadata: req.Metadata,
		DataType: req.DataType,
	})
}

func (n notesUseCase) GetNote(user models.User, noteID int) (models.CreateNoteRequest, error) {
	note, err := n.NotesRepo.GetNote(user.UserID, noteID)
	if err != nil {
		return models.CreateNoteRequest{}, err
	}
	data := base64.StdEncoding.EncodeToString(note.Data)
	return models.CreateNoteRequest{
		Metadata: note.Metadata,
		DataType: note.DataType,
		Data:     data,
	}, nil
}

func (n notesUseCase) GetUserNotes(user models.User) ([]models.Note, error) {
	notes, err := n.NotesRepo.GetUserNotesList(user.UserID)
	if err != nil {
		return nil, err
	}

	if len(notes) > 0 {
		return notes, nil
	}

	return []models.Note{}, err
}
