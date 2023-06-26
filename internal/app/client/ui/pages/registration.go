package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/components"
	"github.com/nuttech/bell/v2"
)

func NewRegistrationWindow(
	app fyne.App,
	events *bell.Events,
) fyne.Window {
	w := app.NewWindow("Registration")
	w.Resize(fyne.NewSize(500, 200))
	w.SetFixedSize(true)

	regForm := components.NewRegistrationForm(w, events)

	w.SetContent(
		container.New(layout.NewMaxLayout(), regForm))
	w.CenterOnScreen()

	return w
}
