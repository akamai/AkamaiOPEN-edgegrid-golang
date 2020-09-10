package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"io/ioutil"
	"net/http"
)

type APIError struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	StatusCode int    `json:"-"`
}

var ErrNotFound = errors.New("resource not found")

func NewAPIError(r *http.Response, logger log.Interface) APIError {
	var apiError APIError
	var body []byte
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("reading error response body: %s", err)
		apiError.StatusCode = r.StatusCode
		return apiError
	}
	if err := json.Unmarshal(body, &apiError); err != nil {
		logger.Errorf("could not unmarshal API error: %s", err)
	}
	apiError.StatusCode = r.StatusCode
	return apiError
}

func (e APIError) Error() string {
	return fmt.Sprintf("Title: %s; Type: %s; Details: %s", e.Title, e.Type, e.Detail)
}
