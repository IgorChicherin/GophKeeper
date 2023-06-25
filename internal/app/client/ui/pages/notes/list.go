package notes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	notes_components "github.com/IgorChicherin/gophkeeper/internal/app/client/ui/components/notes"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/nuttech/bell/v2"
)

func NewNotesListWindow(
	app fyne.App,
	notes []models.Note,
	events *bell.Events,
) (fyne.Window, error) {
	w := app.NewWindow("Notes")
	w.Resize(fyne.NewSize(640, 480))
	var notesViewer *fyne.Container
	var n models.Note

	notesList := notes_components.NewNotesList(notes, events)
	notesViewer = notes_components.NewNotesViewer(models.Note{}, events)

	events.Listen("show_note", func(message bell.Message) {
		n = message.(models.Note)
		notesViewer.RemoveAll()
		notesViewer.Add(
			notes_components.NewNotesViewer(n, events),
		)
	})

	if len(notes) > 0 {
		notesViewer = notes_components.NewNotesViewer(notes[0], events)
	}

	notesList.Resize(fyne.NewSize(320, 300))
	notesViewer.Resize(fyne.NewSize(320, 300))

	content := container.NewAdaptiveGrid(
		1,
		container.NewGridWithColumns(2, notesList, notesViewer),
		layout.NewSpacer(),
		container.NewGridWithColumns(
			4,
			container.NewVBox(
				layout.NewSpacer(),
				widget.NewButton("Create", func() {
					_ = events.Ring("create_note_window", nil)
				})),

			container.NewVBox(
				layout.NewSpacer(),
				widget.NewButton("Edit", func() {
					_ = events.Ring("edit_note_window", n)
				})),

			container.NewVBox(
				layout.NewSpacer(),
				widget.NewButton("Delete", func() {
					_ = events.Ring("delete_note", n)
				})),

			container.NewVBox(
				layout.NewSpacer(),
				widget.NewButton("Close", func() {
					w.Close()
					_ = events.Ring("show_login_window", nil)
				})),
		),
	)
	w.SetContent(content)
	w.CenterOnScreen()
	return w, nil
}
