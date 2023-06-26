package notes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/components/notes"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/nuttech/bell/v2"
)

func NewEditNoteWindow(
	app fyne.App,
	note models.Note,
	events *bell.Events,
) fyne.Window {
	w := app.NewWindow("Create note")
	w.Resize(fyne.NewSize(500, 250))
	w.SetFixedSize(true)

	var notesForm *fyne.Container

	switch note.DataType {
	case "AUTH":
		notesForm = notes.NewEditAuthNoteForm(note, w, events)
	case "TEXT":
		notesForm = notes.NewEditTextNoteForm(note, w, events)
	case "CREDIT_CARD":
		notesForm = notes.NewEditCreditCardNoteForm(note, w, events)
	case "BINARY":
		notesForm = notes.NewEditBinaryNoteForm(note, w, events)
	default:
		notesForm = container.NewMax()
	}

	content := container.New(layout.NewMaxLayout(), notesForm)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}
