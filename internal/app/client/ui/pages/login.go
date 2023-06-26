package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/components"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/nuttech/bell/v2"
)

func NewLoginWindow(
	app fyne.App,
	userUseCase usecases.UserUseCase,
	defaultSrvAddr string,
	events *bell.Events,
) fyne.Window {
	w := app.NewWindow("Login")
	w.Resize(fyne.NewSize(500, 200))
	w.SetFixedSize(true)

	loginForm := components.NewLoginForm(app, userUseCase, events)
	content := container.New(layout.NewMaxLayout(), loginForm)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}
