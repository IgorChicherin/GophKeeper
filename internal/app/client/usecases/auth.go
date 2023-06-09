package usecases

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/usecases/http"
)

type UserUseCase interface {
	AuthUser(user models.RequestUserModel) error
	RegisterUser(user models.RequestUserModel) error
}

func NewHTTPClientUserUseCase(
	tokenRepo repositories.TokenRepository,
	certRepo repositories.CertRepository,
	userRepo repositories.UserRepository,
) UserUseCase {
	return http.UserUseCase{TokenRepo: tokenRepo, CertRepo: certRepo, UserRepo: userRepo}
}
