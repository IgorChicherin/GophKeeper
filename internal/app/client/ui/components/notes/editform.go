package notes

import (
	"encoding/json"
	"fmt"
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

func NewEditAuthNoteForm(
	note models.Note,
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	var authData ui_models.AuthData

	err := json.Unmarshal(note.Data, &authData)

	if err != nil {
		_ = events.Ring("error", err.Error())
		return container.NewMax()
	}

	login := widget.NewEntry()
	login.SetText(authData.Login)

	pwd := widget.NewPasswordEntry()
	pwd.SetText(authData.Password)

	tag := widget.NewEntry()
	tag.SetText(note.Metadata)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Login", Widget: login},
			{Text: "Password", Widget: pwd},
			{Text: "Tag", Widget: tag},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Save", theme.AccountIcon(), func() {
				authData := ui_models.AuthData{Login: login.Text, Password: pwd.Text}
				data, err := json.Marshal(&authData)

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}

				_ = events.Ring("edit_note", models.Note{
					ID:       note.ID,
					UserID:   note.UserID,
					DataType: note.DataType,
					Metadata: tag.Text,
					Data:     data,
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}

func NewEditTextNoteForm(
	note models.Note,
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	text := widget.NewMultiLineEntry()
	text.SetText(string(note.Data))

	tag := widget.NewEntry()
	tag.SetText(note.Metadata)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Text", Widget: text},
			{Text: "Tag", Widget: tag},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Save", theme.AccountIcon(), func() {
				_ = events.Ring("edit_note", models.Note{
					ID:       note.ID,
					UserID:   note.UserID,
					DataType: note.DataType,
					Metadata: tag.Text,
					Data:     []byte(text.Text),
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}

func NewEditBinaryNoteForm(
	note models.Note,
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

	tag := widget.NewEntry()
	tag.SetText(note.Metadata)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "File", Widget: text},
			{Text: "Tag", Widget: tag},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Save", theme.AccountIcon(), func() {
				_ = events.Ring("edit_note", models.Note{
					ID:       note.ID,
					UserID:   note.UserID,
					DataType: note.DataType,
					Metadata: tag.Text,
					Data:     fileData,
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}

func NewEditCreditCardNoteForm(
	note models.Note,
	window fyne.Window,
	events *bell.Events,
) *fyne.Container {
	var ccData ui_models.CreditCard

	err := json.Unmarshal(note.Data, &ccData)

	if err != nil {
		_ = events.Ring("error", err.Error())
		return container.NewMax()
	}

	number := widget.NewEntry()
	number.SetText(fmt.Sprintf("%d", ccData.Number))

	month := widget.NewEntry()
	month.SetText(fmt.Sprintf("%d", ccData.Number))
	year := widget.NewEntry()
	year.SetText(fmt.Sprintf("%d", ccData.Number))

	cvv := widget.NewPasswordEntry()
	cvv.SetText(fmt.Sprintf("%d", ccData.Number))

	holder := widget.NewEntry()
	holder.SetText(ccData.Holder)

	tag := widget.NewEntry()
	tag.SetText(note.Metadata)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Number", Widget: number},
			{Text: "Month", Widget: month},
			{Text: "Year", Widget: year},
			{Text: "CVV", Widget: cvv},
			{Text: "Holder", Widget: holder},
			{Text: "Tag", Widget: tag},
		},
	}
	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Save", theme.AccountIcon(), func() {
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

				_ = events.Ring("edit_note", models.Note{
					DataType: note.DataType,
					Metadata: tag.Text,
					Data:     data,
				})
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), window.Close),
		),
		layout.NewSpacer(),
	)
}
