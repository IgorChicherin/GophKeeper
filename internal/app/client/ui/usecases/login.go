package usecases

import (
	"fyne.io/fyne/v2"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui/pages"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
	"github.com/nuttech/bell/v2"
)

type LoginUseCase interface {
	OnLoginSuccess(message bell.Message)
	Register()
}

func NewLoginUseCase(
	app fyne.App,
	addr string,
	events *bell.Events,
	userUseCase usecases.UserUseCase,
	tokenRepo repositories.TokenRepository,
	certRepo repositories.CertRepository,
	client httpclient.HTTPClientSync,
) LoginUseCase {
	window := pages.NewLoginWindow(app, userUseCase, addr, events)
	return &loginUseCase{
		addr:       addr,
		window:     window,
		App:        app,
		tokenRepo:  tokenRepo,
		certRepo:   certRepo,
		httpClient: client,
		events:     events,
	}
}

type loginUseCase struct {
	App        fyne.App
	window     fyne.Window
	addr       string
	events     *bell.Events
	tokenRepo  repositories.TokenRepository
	certRepo   repositories.CertRepository
	httpClient httpclient.HTTPClientSync
}

func (u *loginUseCase) OnLoginSuccess(message bell.Message) {
	u.window.Hide()
	_ = u.events.Ring("show_notes_list", nil)
}

func (u *loginUseCase) Register() {
	u.events.Listen("show_login_window", func(message bell.Message) {
		u.window.Show()
	})

	u.events.Listen("hide_login_window", func(message bell.Message) {
		u.window.Hide()
	})

	u.events.Listen("login_successful", u.OnLoginSuccess)
}
