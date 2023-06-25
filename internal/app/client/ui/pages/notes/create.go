package notes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/components/notes"
	"github.com/nuttech/bell/v2"
)

func NewCreateNoteWindow(
	app fyne.App,
	events *bell.Events,
) fyne.Window {
	w := app.NewWindow("Create note")
	w.Resize(fyne.NewSize(500, 250))
	w.SetFixedSize(true)

	notesForm := notes.NewCreateAuthNoteForm(w, events)
	content := container.New(layout.NewMaxLayout(), notesForm)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}
