package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	ui_usecases "github.com/IgorChicherin/gophkeeper/internal/app/client/ui/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
	"github.com/nuttech/bell/v2"
)

func RunApp(
	defaultServerAddr string,
	userUseCase usecases.UserUseCase,
	tokenRepo repositories.TokenRepository,
	certRepo repositories.CertRepository,
	httpClient httpclient.HTTPClientSync,
) {
	fyneApp := app.New()
	events := bell.New()

	events.Listen("error", func(message bell.Message) {
		onAppError(message, fyneApp)
	})

	loginUseCase := ui_usecases.NewLoginUseCase(
		fyneApp, defaultServerAddr, events, userUseCase, tokenRepo, certRepo, httpClient)
	loginUseCase.Register()

	registrationUseCase := ui_usecases.NewRegistrationUseCase(fyneApp, userUseCase, events)
	registrationUseCase.Register()

	notesUseCase := ui_usecases.NewNotesUseCase(defaultServerAddr, fyneApp, httpClient, certRepo, tokenRepo, events)
	notesUseCase.Register()

	_ = events.Ring("show_login_window", nil)

	fyneApp.Run()
}

func onAppError(message bell.Message, fyneApp fyne.App) {
	msg := message.(string)
	errorWindow := fyneApp.NewWindow("Error")
	errorWindow.SetFixedSize(true)
	errorWindow.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			widget.NewLabel(msg),
			widget.NewButton("ok", func() {
				errorWindow.Close()
			})),
	)
	errorWindow.CenterOnScreen()
	errorWindow.Show()
}
