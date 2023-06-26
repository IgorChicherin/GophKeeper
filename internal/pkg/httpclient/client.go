package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/IgorChicherin/gophkeeper/internal/shared/models"
	log "github.com/sirupsen/logrus"
)

type HTTPClientSync interface {
	Get(URL string, headers, params map[string]string) ([]byte, error)
	Post(URL string, headers map[string]string, data []byte) ([]byte, error)
	Put(URL string, headers map[string]string, data []byte) ([]byte, error)
}

type httpClientSync struct {
	client *http.Client
}

func NewHTTPClientSync(client *http.Client) HTTPClientSync {
	return httpClientSync{client: client}
}

func (c httpClientSync) Get(URL string, headers, params map[string]string) ([]byte, error) {
	return c.sendData(http.MethodGet, URL, headers, params, nil)
}

func (c httpClientSync) Post(URL string, headers map[string]string, data []byte) ([]byte, error) {
	return c.sendData(http.MethodPost, URL, headers, nil, data)
}

func (c httpClientSync) Put(URL string, headers map[string]string, data []byte) ([]byte, error) {
	return c.sendData(http.MethodPut, URL, headers, nil, data)
}

func (c httpClientSync) makeRequest(
	method string,
	URL string,
	headers, query map[string]string,
	body []byte,
) (*http.Request, error) {
	var req *http.Request
	var err error

	switch method {
	case http.MethodPost, http.MethodPut:
		req, err = http.NewRequest(string(method), URL, nil)

		if len(body) > 0 {
			req, err = http.NewRequest(string(method), URL, bytes.NewBuffer(body))
		}

		if err != nil {
			return nil, err
		}
	case http.MethodGet, http.MethodDelete:
		req, err = http.NewRequest(string(method), URL, nil)

		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("method not not implemented")
	}

	if len(query) > 0 {
		c.addQuery(req, query)
	}

	if len(headers) > 0 {
		c.addHeaders(req, headers)
	}

	return req, nil
}

func (c httpClientSync) addQuery(req *http.Request, query map[string]string) {
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (c httpClientSync) addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (c httpClientSync) sendData(
	method string,
	URL string,
	headers, query map[string]string,
	data []byte,
) ([]byte, error) {
	req, err := c.makeRequest(method, URL, headers, query, data)
	if err != nil {
		return []byte{}, err
	}

	resp, err := c.client.Do(req)

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
		var e models.DefaultErrorResponse
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
