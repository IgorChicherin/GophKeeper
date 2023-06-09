package usecases

import "github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"

type NoteUseCase interface {
	CreateUserNote()
	GetUserNote()
	GetAllUserNotes()
	UpdateUserNote()
}

type noteUseCase struct {
	TokenRepo repositories.TokenRepository
	CertRepo  repositories.CertRepository
	NoteRepo  repositories.NotesRepository
}

func NewNoteUseCase(
	tokenRepo repositories.TokenRepository,
	certRepo repositories.CertRepository,
	noteRepo repositories.NotesRepository,
) NoteUseCase {
	return noteUseCase{TokenRepo: tokenRepo, CertRepo: certRepo, NoteRepo: noteRepo}
}

func (n noteUseCase) CreateUserNote() {

}

func (n noteUseCase) GetUserNote() {

}

func (n noteUseCase) GetAllUserNotes() {

}

func (n noteUseCase) UpdateUserNote() {

}
