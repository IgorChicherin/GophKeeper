package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/httpclient"
	"github.com/IgorChicherin/gophkeeper/internal/shared/models"
)

type HTTPNotesRepository struct {
	BaseURL    string
	HTTPClient httpclient.HTTPClientSync
}

func (r HTTPNotesRepository) GetNote(token string, noteID int) ([]byte, error) {
	URL := fmt.Sprintf("%s/api/notes/%d", r.BaseURL, noteID)
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": token,
	}
	return r.HTTPClient.Get(URL, headers, nil)
}

func (r HTTPNotesRepository) GetNotes(token string) ([]byte, error) {
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": token,
	}
	URL := fmt.Sprintf("%s/api/%s", r.BaseURL, "notes")
	return r.HTTPClient.Get(URL, headers, nil)
}

func (r HTTPNotesRepository) CreateNote(token string, note models.CreateNoteRequest) ([]byte, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": token,
	}
	data, err := json.Marshal(&note)
	if err != nil {
		return nil, err
	}
	URL := fmt.Sprintf("%s/api/%s", r.BaseURL, "notes/create")
	return r.HTTPClient.Post(URL, headers, data)
}

func (r HTTPNotesRepository) UpdateNote(token string, note models.CreateNoteRequest) ([]byte, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": token,
	}
	data, err := json.Marshal(&note)
	if err != nil {
		return nil, err
	}
	URL := fmt.Sprintf("%s/api/%s", r.BaseURL, "notes")
	return r.HTTPClient.Put(URL, headers, data)
}
