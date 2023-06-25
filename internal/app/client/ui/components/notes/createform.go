package notes

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ui_models "github.com/IgorChicherin/gophkeeper/internal/app/client/ui/models"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/nuttech/bell/v2"
	"os"
	"strconv"
)

var types = []string{"AUTH", "TEXT", "BINARY", "CREDIT_CARD"}

func NewCreateAuthNoteForm(
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	login := widget.NewEntry()
	pwd := widget.NewPasswordEntry()
	note := widget.NewEntry()
	note_type := widget.NewSelect(types, func(value string) {
		_ = events.Ring("create_note_type_change", value)
	})
	note_type.Selected = types[0]

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type", Widget: note_type},
			{Text: "Login", Widget: login},
			{Text: "Password", Widget: pwd},
			{Text: "Tag", Widget: note},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Create", theme.AccountIcon(), func() {
				authData := ui_models.AuthData{Login: login.Text, Password: pwd.Text}
				data, err := json.Marshal(&authData)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				_ = events.Ring("create_note", models.Note{
					DataType: note_type.Selected,
					Metadata: note.Text,
					Data:     data,
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}

func NewCreateTextNoteForm(
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	text := widget.NewMultiLineEntry()
	note := widget.NewEntry()
	note_type := widget.NewSelect(types, func(value string) {
		_ = events.Ring("create_note_type_change", value)
	})
	note_type.Selected = types[1]

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type", Widget: note_type},
			{Text: "Text", Widget: text},
			{Text: "Tag", Widget: note},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Create", theme.AccountIcon(), func() {
				_ = events.Ring("create_note", models.Note{
					DataType: note_type.Selected,
					Metadata: note.Text,
					Data:     []byte(text.Text),
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}

func NewCreateBinaryNoteForm(
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	var fileData []byte
	text := widget.NewButton("Upload", func() {
		openDlg := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer != nil {
				fileData, err = os.ReadFile(closer.URI().Path())

				if err != nil {
					_ = events.Ring("error", err.Error())
				}
			}
			closer.URI()
		}, window)
		window.Resize(fyne.NewSize(640, 480))
		openDlg.Resize(fyne.NewSize(640, 480))
		openDlg.Show()
	})

	note := widget.NewEntry()
	note_type := widget.NewSelect(types, func(value string) {
		_ = events.Ring("create_note_type_change", value)
	})
	note_type.Selected = types[2]

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type", Widget: note_type},
			{Text: "File", Widget: text},
			{Text: "Tag", Widget: note},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Create", theme.AccountIcon(), func() {
				_ = events.Ring("create_note", models.Note{
					DataType: note_type.Selected,
					Metadata: note.Text,
					Data:     fileData,
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}

func NewCreateCreditCardNoteForm(
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	number := widget.NewEntry()
	month := widget.NewEntry()
	year := widget.NewEntry()
	cvv := widget.NewPasswordEntry()
	holder := widget.NewEntry()

	note := widget.NewEntry()
	note_type := widget.NewSelect(types, func(value string) {
		_ = events.Ring("create_note_type_change", value)
	})
	note_type.Selected = types[3]

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type", Widget: note_type},
			{Text: "Number", Widget: number},
			{Text: "Month", Widget: month},
			{Text: "Year", Widget: year},
			{Text: "CVV", Widget: cvv},
			{Text: "Holder", Widget: holder},
			{Text: "Tag", Widget: note},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Create", theme.AccountIcon(), func() {
				n, err := strconv.Atoi(number.Text)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				m, err := strconv.Atoi(month.Text)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				y, err := strconv.Atoi(year.Text)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				c, err := strconv.Atoi(cvv.Text)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				ccData := ui_models.CreditCard{
					Number: uint(n),
					Month:  uint(m),
					Year:   uint(y),
					CVV:    uint(c),
					Holder: holder.Text,
				}
				data, err := json.Marshal(&ccData)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				_ = events.Ring("create_note", models.Note{
					DataType: note_type.Selected,
					Metadata: note.Text,
					Data:     data,
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}
