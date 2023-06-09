package repositories

import (
	"errors"
	"sync"
)

type CertRepository interface {
	Get() (string, error)
	Set(cert string)
}

type certRepository struct {
	data sync.Map
}

func NewCertRepository() TokenRepository {
	return &certRepository{}
}

func (r *certRepository) Get() (string, error) {
	v, ok := r.data.Load("cert")
	if !ok {
		err := errors.New("cert isn't set")
		return "", err
	}
	val := v.(string)
	return val, nil
}

func (r *certRepository) Set(cert string) {
	r.data.Store("cert", cert)
}
