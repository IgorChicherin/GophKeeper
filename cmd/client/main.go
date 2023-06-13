package main

import (
	"fmt"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	db_models "github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	baseURL := "http://localhost:3001"
	client := httpclient.NewHTTPClientSync(&http.Client{})

	userRepo := repositories.NewUserRepository(baseURL, client)
	tokenRepo := repositories.NewTokenRepository()
	certRepo := repositories.NewCertRepository()

	userUseCase := usecases.NewHTTPClientUserUseCase(tokenRepo, certRepo, userRepo)

	err := userUseCase.AuthUser(models.RequestUserModel{Login: "string", Password: "string"})

	if err != nil {
		log.Panicln(err)
	}

	pubKey, err := certRepo.Get()
	if err != nil {
		log.Panicln(err)
	}

	enc, err := crypto509.NewEncrypter([]byte(pubKey))

	if err != nil {
		log.Panicln(err)
	}

	noteRepo := repositories.NewHTTPNoteRepository(baseURL, client)
	notesUseCase := usecases.NewHTTPNoteUseCase(tokenRepo, noteRepo, enc)
	n := db_models.Note{
		Data:     []byte("Hello!"),
		Metadata: "test",
		DataType: "BINARY",
	}

	createdNote, err := notesUseCase.CreateUserNote(n)

	fmt.Printf("%v", createdNote)
	data, err := notesUseCase.GetUserNote(3)

	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("%v\n", data)
	fmt.Printf("%s\n", string(data.Data))
}
