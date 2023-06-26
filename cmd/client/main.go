package main

import (
	"net/http"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/config"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
)

func main() {
	cfg, err := config.GetClientConfig()

	if err != nil {
		fyneApp := app.New()
		errorWindow := fyneApp.NewWindow("Error")
		errorWindow.SetFixedSize(true)
		errorWindow.CenterOnScreen()
		errorWindow.SetContent(
			container.New(
				layout.NewVBoxLayout(),
				widget.NewLabel(err.Error()),
				widget.NewButton("ok", func() {
					fyneApp.Quit()
				})),
		)
		errorWindow.Show()
		fyneApp.Run()
		return
	}

	client := httpclient.NewHTTPClientSync(&http.Client{})

	userRepo := repositories.NewUserRepository(cfg.Address, client)
	tokenRepo := repositories.NewTokenRepository()
	certRepo := repositories.NewCertRepository()

	userUseCase := usecases.NewHTTPClientUserUseCase(tokenRepo, certRepo, userRepo)

	ui.RunApp(cfg.Address, userUseCase, tokenRepo, certRepo, client)
}
