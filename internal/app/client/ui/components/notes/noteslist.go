package notes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/nuttech/bell/v2"
)

func NewNotesList(
	notes []models.Note,
	events *bell.Events,
) *fyne.Container {

	list := widget.NewList(
		func() int {
			return len(notes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("tmp")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(notes[i].Metadata)
		})
	list.OnSelected = func(i widget.ListItemID) {
		_ = events.Ring("show_note", notes[i])
	}
	return container.NewMax(list)
}
