package repositories

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/client/http/repositories"
)

type UserRepository interface {
	Register(login, password string) ([]byte, error)
	Authenticate(login, password string) ([]byte, error)
}

func NewUserRepository(baseURL string) UserRepository {
	return repositories.HTTPUserRepository{BaseURL: baseURL}
}
