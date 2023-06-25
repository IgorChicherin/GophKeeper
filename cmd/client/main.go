package main

import (
	"net/http"

	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/ui"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
)

func main() {
	baseURL := "http://localhost:3001"
	client := httpclient.NewHTTPClientSync(&http.Client{})

	userRepo := repositories.NewUserRepository(baseURL, client)
	tokenRepo := repositories.NewTokenRepository()
	certRepo := repositories.NewCertRepository()

	userUseCase := usecases.NewHTTPClientUserUseCase(tokenRepo, certRepo, userRepo)

	ui.RunApp(baseURL, userUseCase, tokenRepo, certRepo, client)
}
