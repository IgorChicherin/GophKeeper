package http

import (
	"encoding/json"

	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/repositories"
)

type UserUseCase struct {
	TokenRepo repositories.TokenRepository
	CertRepo  repositories.CertRepository
	UserRepo  repositories.UserRepository
}

func (u UserUseCase) AuthUser(user models.RequestUserModel) error {
	body, err := u.UserRepo.Authenticate(user.Login, user.Password)
	if err != nil {
		return err
	}

	var resp models.ResponseUserLogin
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	u.TokenRepo.Set(resp.Token)
	u.CertRepo.Set(resp.Cert)

	return nil
}

func (u UserUseCase) RegisterUser(user models.RequestUserModel) error {
	_, err := u.UserRepo.Register(user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}
