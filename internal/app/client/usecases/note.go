package usecases

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases/http"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
)

type NoteUseCase interface {
	CreateUserNote(note models.Note) (models.Note, error)
	GetUserNote(noteID int) (models.Note, error)
	GetAllUserNotes() ([]models.Note, error)
	UpdateUserNote(note models.Note) (models.Note, error)
}

func NewHTTPNoteUseCase(
	tokenRepo repositories.TokenRepository,
	noteRepo repositories.NotesRepository,
	enc crypto509.Encrypter,
) NoteUseCase {
	return http.HTTPNoteUseCase{TokenRepo: tokenRepo, NoteRepo: noteRepo, Encrypter: enc}
}
