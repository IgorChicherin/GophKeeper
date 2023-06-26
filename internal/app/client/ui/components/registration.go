package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/nuttech/bell/v2"
)

func NewRegistrationForm(
	w fyne.Window,
	events *bell.Events,
) *fyne.Container {
	login := widget.NewEntry()
	pwd := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Login", Widget: login},
			{Text: "Password", Widget: pwd},
		},
		OnSubmit: func() {
			_ = events.Ring("submit_registration_from", models.RequestUserModel{Login: login.Text, Password: pwd.Text})
			w.Close()
			_ = events.Ring("show_login_window", nil)
		},
		OnCancel: func() {
			_ = events.Ring("cancel_registration_form", nil)
		},
	}

	return container.NewVBox(form)
}
