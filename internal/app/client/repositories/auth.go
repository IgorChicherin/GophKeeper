package repositories

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/client/http/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
)

type UserRepository interface {
	Register(login, password string) ([]byte, error)
	Authenticate(login, password string) ([]byte, error)
}

func NewUserRepository(baseURL string, client httpclient.HTTPClientSync) UserRepository {
	return repositories.HTTPUserRepository{BaseURL: baseURL, HTTPClient: client}
}
