package main

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	log "github.com/sirupsen/logrus"
)

func main() {
	userRepo := repositories.NewUserRepository("http://localhost:3001")
	tokenRepo := repositories.NewTokenRepository()
	certRepo := repositories.NewCertRepository()

	userUseCase := usecases.NewHTTPClientUserUseCase(tokenRepo, certRepo, userRepo)

	err := userUseCase.AuthUser(models.RequestUserModel{Login: "string", Password: "string"})

	if err != nil {
		log.Panicln(err)
	}
}
