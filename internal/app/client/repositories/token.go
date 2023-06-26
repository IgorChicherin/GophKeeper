package repositories

import (
	"errors"
	"sync"
)

type TokenRepository interface {
	Get() (string, error)
	Set(token string)
}

type tokenRepository struct {
	data sync.Map
}

func NewTokenRepository() TokenRepository {
	return &tokenRepository{}
}

func (r *tokenRepository) Get() (string, error) {
	v, ok := r.data.Load("token")
	if !ok {
		err := errors.New("token isn't set")
		return "", err
	}
	val := v.(string)
	return val, nil
}

func (r *tokenRepository) Set(token string) {
	r.data.Store("token", token)
}
