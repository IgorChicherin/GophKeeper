package usecases

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	notes_components "github.com/IgorChicherin/gophkeeper/internal/app/client/ui/components/notes"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/pages/notes"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
	"github.com/nuttech/bell/v2"
)

type NoteUseCase interface {
	Register()
	OnShowNotesList(message bell.Message)
	OnCreateNoteWindow(message bell.Message)
	OnCreateNote(message bell.Message)
	OnChangeCreateNoteType(message bell.Message)
	OnEditNoteWindow(message bell.Message)
	OnEditNote(message bell.Message)
	OnDeleteNote(message bell.Message)
}

func NewNotesUseCase(
	addr string,
	app fyne.App,
	httpClient httpclient.HTTPClientSync,
	certRepo repositories.CertRepository,
	tokenRepo repositories.TokenRepository,
	events *bell.Events,
) NoteUseCase {
	notesRepo := repositories.NewHTTPNoteRepository(addr, httpClient)
	return &noteUseCase{
		app:       app,
		certRepo:  certRepo,
		notesRepo: notesRepo,
		tokenRepo: tokenRepo,
		events:    events,
	}
}

type noteUseCase struct {
	app              fyne.App
	window           fyne.Window
	createNoteWindow fyne.Window
	editNoteWindow   fyne.Window
	certRepo         repositories.CertRepository
	tokenRepo        repositories.TokenRepository
	notesUseCase     usecases.NoteUseCase
	notesRepo        repositories.NotesRepository
	events           *bell.Events
}

func (u *noteUseCase) OnShowNotesList(message bell.Message) {
	pub, err := u.certRepo.Get()

	if err != nil {
		_ = u.events.Ring("error", err.Error())
	}

	enc, err := crypto509.NewEncrypter([]byte(pub))

	if err != nil {
		_ = u.events.Ring("error", err.Error())
	}

	u.notesUseCase = usecases.NewHTTPNoteUseCase(u.tokenRepo, u.notesRepo, enc)
	n, err := u.notesUseCase.GetAllUserNotes()

	if err != nil {
		_ = u.events.Ring("error", err.Error())
	}

	u.window, err = notes.NewNotesListWindow(u.app, n, u.events)

	if err != nil {
		_ = u.events.Ring("error", err.Error())
	}
	u.window.Show()
}

func (u *noteUseCase) OnCreateNoteWindow(message bell.Message) {
	u.createNoteWindow = notes.NewCreateNoteWindow(u.app, u.events)
	u.createNoteWindow.Show()
}

func (u *noteUseCase) OnCreateNote(message bell.Message) {
	n := message.(models.Note)
	_, err := u.notesUseCase.CreateUserNote(n)

	if err != nil {
		_ = u.events.Ring("error", err.Error())
		return
	}
	u.window.Close()
	_ = u.events.Ring("show_notes_list", nil)
	u.createNoteWindow.Close()
}

func (u *noteUseCase) OnEditNoteWindow(message bell.Message) {
	u.window.Close()
	n := message.(models.Note)
	u.editNoteWindow = notes.NewEditNoteWindow(u.app, n, u.events)
	u.editNoteWindow.Show()
}

func (u *noteUseCase) OnChangeCreateNoteType(message bell.Message) {
	noteType := message.(string)

	switch noteType {
	case "TEXT":
		content := notes_components.NewCreateTextNoteForm(u.createNoteWindow, u.events)
		u.createNoteWindow.SetContent(container.New(layout.NewMaxLayout(), content))
	case "CREDIT_CARD":
		content := notes_components.NewCreateCreditCardNoteForm(u.createNoteWindow, u.events)
		u.createNoteWindow.SetContent(container.New(layout.NewMaxLayout(), content))
	case "BINARY":
		content := notes_components.NewCreateBinaryNoteForm(u.createNoteWindow, u.events)
		u.createNoteWindow.SetContent(container.New(layout.NewMaxLayout(), content))
	default:
		content := notes_components.NewCreateAuthNoteForm(u.createNoteWindow, u.events)
		u.createNoteWindow.SetContent(container.New(layout.NewMaxLayout(), content))
	}
}

func (u *noteUseCase) OnDeleteNote(message bell.Message) {
	note := message.(models.Note)
	err := u.notesUseCase.DeleteUserNote(note.ID)

	if err != nil {
		_ = u.events.Ring("error", err.Error())
	}
	u.window.Close()
	_ = u.events.Ring("show_notes_list", nil)
}

func (u *noteUseCase) OnEditNote(message bell.Message) {
	note := message.(models.Note)
	_, err := u.notesUseCase.UpdateUserNote(note)

	if err != nil {
		_ = u.events.Ring("error", err.Error())
	}
	u.editNoteWindow.Close()
	_ = u.events.Ring("show_notes_list", nil)
}

func (u *noteUseCase) Register() {
	u.events.Listen("show_notes_list", u.OnShowNotesList)
	u.events.Listen("create_note_window", u.OnCreateNoteWindow)
	u.events.Listen("create_note", u.OnCreateNote)
	u.events.Listen("create_note_type_change", u.OnChangeCreateNoteType)
	u.events.Listen("edit_note_window", u.OnEditNoteWindow)
	u.events.Listen("edit_note", u.OnEditNote)
	u.events.Listen("delete_note", u.OnDeleteNote)
}
