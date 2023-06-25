package usecases

import (
	"encoding/base64"

	"github.com/IgorChicherin/gophkeeper/internal/app/server/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	shared_models "github.com/IgorChicherin/gophkeeper/internal/shared/models"
)

type NotesUseCase interface {
	CreateUserNote(user models.User, req shared_models.CreateNoteRequest) (shared_models.DecodedNoteResponse, error)
	GetNote(user models.User, noteID int) (shared_models.DecodedNoteResponse, error)
	GetUserNotes(user models.User) ([]shared_models.DecodedNoteResponse, error)
	UpdateUserNote(user models.User, noteID int, req shared_models.CreateNoteRequest) (shared_models.DecodedNoteResponse, error)
	DeleteUserNote(user models.User, noteID int) error
}

type notesUseCase struct {
	NotesRepo repositories.NotesRepository
	Decrypter crypto509.Decrypter
}

func NewNotesUseCase(notesRepo repositories.NotesRepository, dec crypto509.Decrypter) NotesUseCase {
	return notesUseCase{NotesRepo: notesRepo, Decrypter: dec}
}

func (n notesUseCase) CreateUserNote(user models.User, req shared_models.CreateNoteRequest) (shared_models.DecodedNoteResponse, error) {
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}

	encodedNote, err := n.NotesRepo.CreateNote(models.Note{
		UserID:   user.UserID,
		Data:     data,
		Metadata: req.Metadata,
		DataType: req.DataType,
	})

	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}

	decodedBytes, err := n.Decrypter.DecryptData(encodedNote.Data)
	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}

	return shared_models.DecodedNoteResponse{
		ID:        encodedNote.ID,
		UserID:    encodedNote.UserID,
		Data:      base64.StdEncoding.EncodeToString(decodedBytes),
		Metadata:  encodedNote.Metadata,
		DataType:  encodedNote.DataType,
		UpdatedAt: encodedNote.UpdatedAt,
		CreatedAt: encodedNote.CreatedAt,
	}, nil
}

func (n notesUseCase) GetNote(user models.User, noteID int) (shared_models.DecodedNoteResponse, error) {
	note, err := n.NotesRepo.GetNote(user.UserID, noteID)
	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}
	data, err := n.Decrypter.DecryptData(note.Data)

	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}
	return shared_models.DecodedNoteResponse{
		ID:        note.ID,
		UserID:    note.UserID,
		Data:      base64.StdEncoding.EncodeToString(data),
		Metadata:  note.Metadata,
		DataType:  note.DataType,
		UpdatedAt: note.UpdatedAt,
		CreatedAt: note.CreatedAt,
	}, nil
}

func (n notesUseCase) GetUserNotes(user models.User) ([]shared_models.DecodedNoteResponse, error) {
	notes, err := n.NotesRepo.GetUserNotesList(user.UserID)
	if err != nil {
		return nil, err
	}

	if len(notes) > 0 {
		var decodedNotes []shared_models.DecodedNoteResponse
		for _, note := range notes {
			data, err := n.Decrypter.DecryptData(note.Data)

			if err != nil {
				return []shared_models.DecodedNoteResponse{}, err
			}

			decodedNotes = append(decodedNotes, shared_models.DecodedNoteResponse{
				ID:        note.ID,
				UserID:    note.UserID,
				Data:      base64.StdEncoding.EncodeToString(data),
				Metadata:  note.Metadata,
				DataType:  note.DataType,
				UpdatedAt: note.UpdatedAt,
				CreatedAt: note.CreatedAt,
			})
		}
		return decodedNotes, nil
	}

	return []shared_models.DecodedNoteResponse{}, err
}

func (n notesUseCase) DeleteUserNote(user models.User, noteID int) error {
	return n.NotesRepo.DeleteNote(user.UserID, noteID)
}

func (n notesUseCase) UpdateUserNote(user models.User, noteID int, req shared_models.CreateNoteRequest) (shared_models.DecodedNoteResponse, error) {
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}

	encodedNote, err := n.NotesRepo.UpdateNote(models.Note{
		ID:       noteID,
		UserID:   user.UserID,
		Data:     data,
		Metadata: req.Metadata,
		DataType: req.DataType,
	})

	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}
	decodedBytes, err := n.Decrypter.DecryptData(encodedNote.Data)
	if err != nil {
		return shared_models.DecodedNoteResponse{}, err
	}

	return shared_models.DecodedNoteResponse{
		ID:        encodedNote.ID,
		UserID:    encodedNote.UserID,
		Data:      base64.StdEncoding.EncodeToString(decodedBytes),
		Metadata:  encodedNote.Metadata,
		DataType:  encodedNote.DataType,
		UpdatedAt: encodedNote.UpdatedAt,
		CreatedAt: encodedNote.CreatedAt,
	}, nil
}
