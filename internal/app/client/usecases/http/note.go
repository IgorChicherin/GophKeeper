package http

import (
	"encoding/base64"
	"encoding/json"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	shared_models "github.com/IgorChicherin/gophkeeper/internal/shared/models"
)

type HTTPNoteUseCase struct {
	TokenRepo repositories.TokenRepository
	NoteRepo  repositories.NotesRepository
	Encrypter crypto509.Encrypter
}

func (n HTTPNoteUseCase) CreateUserNote(note models.Note) (models.Note, error) {
	token, err := n.TokenRepo.Get()
	if err != nil {
		return models.Note{}, err
	}

	noteData, err := n.Encrypter.EncryptData(note.Data)
	noteData = base64.StdEncoding.EncodeToString([]byte(noteData))

	noteDataRequest := shared_models.CreateNoteRequest{Metadata: note.Metadata, Data: noteData, DataType: note.DataType}
	body, err := n.NoteRepo.CreateNote(token, noteDataRequest)
	var noteResponse models.Note
	err = json.Unmarshal(body, &noteResponse)

	if err != nil {
		return models.Note{}, err
	}

	return noteResponse, nil
}

func (n HTTPNoteUseCase) GetUserNote(noteID int) (models.Note, error) {
	token, err := n.TokenRepo.Get()
	if err != nil {
		return models.Note{}, err
	}
	body, err := n.NoteRepo.GetNote(token, noteID)
	if err != nil {
		return models.Note{}, err
	}

	var noteResponse models.Note
	err = json.Unmarshal(body, &noteResponse)
	if err != nil {
		return models.Note{}, err
	}

	return noteResponse, nil
}

func (n HTTPNoteUseCase) GetAllUserNotes() ([]models.Note, error) {
	token, err := n.TokenRepo.Get()
	if err != nil {
		return []models.Note{}, err
	}

	body, err := n.NoteRepo.GetNotes(token)
	if err != nil {
		return []models.Note{}, err
	}

	var noteResponse []models.Note
	err = json.Unmarshal(body, &noteResponse)
	if err != nil {
		return []models.Note{}, err
	}

	return noteResponse, nil
}

func (n HTTPNoteUseCase) UpdateUserNote(note models.Note) (models.Note, error) {
	token, err := n.TokenRepo.Get()
	if err != nil {
		return models.Note{}, err
	}

	noteData, err := n.Encrypter.EncryptData(note.Data)
	noteData = base64.StdEncoding.EncodeToString([]byte(noteData))

	noteDataRequest := shared_models.CreateNoteRequest{Metadata: note.Metadata, Data: noteData, DataType: note.DataType}
	body, err := n.NoteRepo.UpdateNote(token, noteDataRequest)

	var noteResponse models.Note
	err = json.Unmarshal(body, &noteResponse)

	if err != nil {
		return models.Note{}, err
	}
	return noteResponse, nil
}
