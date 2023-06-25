package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/nuttech/bell/v2"
)

func NewLoginForm(
	app fyne.App,
	userUseCase usecases.UserUseCase,
	defaultSrvAddr string,
	events *bell.Events,
) *fyne.Container {
	addr := widget.NewEntry()
	addr.SetText(defaultSrvAddr)

	login := widget.NewEntry()
	pwd := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Server address", Widget: addr},
			{Text: "Login", Widget: login},
			{Text: "Password", Widget: pwd},
		},
	}

	return container.NewVBox(
		form,
		layout.NewSpacer(),
		container.New(
			layout.NewGridLayout(3),
			widget.NewButtonWithIcon("Sign up", theme.ConfirmIcon(), func() {
				_ = events.Ring("sign_up_click", nil)
			}),
			widget.NewButtonWithIcon("Sign in", theme.AccountIcon(), func() {
				err := userUseCase.AuthUser(models.RequestUserModel{Login: login.Text, Password: pwd.Text})

				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}
				err = events.Ring("login_successful", nil)
				if err != nil {
					_ = events.Ring("error", err.Error())
					return
				}
			}),
			widget.NewButtonWithIcon("Close", theme.CancelIcon(), app.Quit),
		),
		layout.NewSpacer(),
	)
}
