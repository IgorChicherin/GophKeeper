package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	shared_models "github.com/IgorChicherin/gophkeeper/internal/shared/models"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type HTTPUserRepository struct {
	BaseURL string
}

func (r HTTPUserRepository) Register(login, password string) ([]byte, error) {
	data, err := json.Marshal(models.RequestUserModel{Login: login, Password: password})
	if err != nil {
		return []byte{}, err
	}
	URL := fmt.Sprintf("%s%s", r.BaseURL, "/api/user/register")
	return r.senData(URL, data)
}

func (r HTTPUserRepository) Authenticate(login, password string) ([]byte, error) {
	data, err := json.Marshal(models.RequestUserModel{Login: login, Password: password})
	if err != nil {
		return []byte{}, err
	}
	URL := fmt.Sprintf("%s%s", r.BaseURL, "/api/user/login")
	return r.senData(URL, data)
}

func (r HTTPUserRepository) senData(URL string, data []byte) ([]byte, error) {
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errorf("Send Data error: %s", err.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var e shared_models.DefaultErrorResponse
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &e)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(e.Error)
	}

	return io.ReadAll(resp.Body)
}
