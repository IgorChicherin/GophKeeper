package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/IgorChicherin/gophkeeper/internal/app/client/models"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
)

type HTTPUserRepository struct {
	BaseURL    string
	HTTPClient httpclient.HTTPClientSync
}

func (r HTTPUserRepository) Register(login, password string) ([]byte, error) {
	data, err := json.Marshal(models.RequestUserModel{Login: login, Password: password})
	if err != nil {
		return []byte{}, err
	}
	URL := fmt.Sprintf("%s%s", r.BaseURL, "/api/user/register")
	headers := map[string]string{"Content-Type": "application/json"}
	return r.HTTPClient.Post(URL, headers, data)
}

func (r HTTPUserRepository) Authenticate(login, password string) ([]byte, error) {
	data, err := json.Marshal(models.RequestUserModel{Login: login, Password: password})
	if err != nil {
		return []byte{}, err
	}
	URL := fmt.Sprintf("%s%s", r.BaseURL, "/api/user/login")
	headers := map[string]string{"Content-Type": "application/json"}
	return r.HTTPClient.Post(URL, headers, data)
}
