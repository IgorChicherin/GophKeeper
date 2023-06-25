package usecases

import (
	"fyne.io/fyne/v2"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/pages"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/nuttech/bell/v2"
)

type RegistrationUseCase interface {
	OnSignUpBtnClick(message bell.Message)
	OnSubmit(message bell.Message)
	Register()
}

func NewRegistrationUseCase(
	app fyne.App,
	userUseCase usecases.UserUseCase,
	events *bell.Events,
) RegistrationUseCase {
	return &registrationUseCase{
		app:         app,
		events:      events,
		userUseCase: userUseCase,
	}
}

type registrationUseCase struct {
	app         fyne.App
	window      fyne.Window
	events      *bell.Events
	userUseCase usecases.UserUseCase
}

func (u *registrationUseCase) OnSignUpBtnClick(message bell.Message) {
	u.window = pages.NewRegistrationWindow(u.app, u.events)
	u.window.Show()
}

func (u *registrationUseCase) OnCancelBtnClick(message bell.Message) {
	u.window.Close()
	_ = u.events.Ring("show_login_window", nil)
}

func (u *registrationUseCase) OnSubmit(message bell.Message) {
	data := message.(models.RequestUserModel)
	err := u.userUseCase.RegisterUser(data)

	if err != nil {
		_ = u.events.Ring("error", err.Error())
		return
	}
}

func (u registrationUseCase) Register() {
	u.events.Listen("registration_successful", u.OnCancelBtnClick)
	u.events.Listen("cancel_registration_form", u.OnCancelBtnClick)

	u.events.Listen("submit_registration_from", u.OnSubmit)
	u.events.Listen("sign_up_click", u.OnSignUpBtnClick)
}
