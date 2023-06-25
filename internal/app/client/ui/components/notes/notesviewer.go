package notes

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	ui_models "github.com/IgorChicherin/gophkeeper/internal/app/client/ui/models"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/nuttech/bell/v2"
	"reflect"
)

func NewNotesViewer(
	note models.Note,
	events *bell.Events,
) *fyne.Container {

	if !reflect.DeepEqual(note, models.Note{}) {
		switch note.DataType {
		case "AUTH":
			return NewAuthNoteViewer(note, events)
		case "TEXT":
			return NewTextNoteViewer(note)
		case "BINARY":
			return NewBinaryNoteViewer(note)
		case "CREDIT_CARD":
			return NewCreditCardNoteViewer(note, events)

		}
	}
	return container.New(
		layout.NewFormLayout(),
	)
}

func NewAuthNoteViewer(note models.Note, events *bell.Events) *fyne.Container {
	var data ui_models.AuthData

	err := json.Unmarshal(note.Data, &data)

	if err != nil {
		_ = events.Ring("error", err.Error())
	}

	return container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Tag"), widget.NewLabel(note.Metadata),
		widget.NewLabel("Login"), widget.NewLabel(data.Login),
		widget.NewLabel("Password"), widget.NewLabel(data.Password),
	)
}

func NewTextNoteViewer(note models.Note) *fyne.Container {
	return container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Tag"), widget.NewLabel(note.Metadata),
		widget.NewLabel("Text"), widget.NewLabel(string(note.Data)),
	)
}

func NewBinaryNoteViewer(note models.Note) *fyne.Container {
	return container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Tag"), widget.NewLabel(note.Metadata),
		widget.NewLabel("Data"), widget.NewLabel(string(note.Data)),
	)
}

func NewCreditCardNoteViewer(note models.Note, events *bell.Events) *fyne.Container {
	var data ui_models.CreditCard

	err := json.Unmarshal(note.Data, &data)

	if err != nil {
		_ = events.Ring("error", err.Error())
	}
	validTo := fmt.Sprintf("%d/%d", data.Month, data.Year)
	return container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Valid to"), widget.NewLabel(validTo),
		widget.NewLabel("CVV"), widget.NewLabel(fmt.Sprintf("%d", data.CVV)),
		widget.NewLabel("Card holder"), widget.NewLabel(data.Holder),
	)
}
