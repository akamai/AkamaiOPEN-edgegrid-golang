package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
)

// APIError contains detailed data on API error returned to the client
type APIError struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	Detail        string `json:"detail"`
	Instance      string `json:"instance,omitempty"`
	BehaviorName  string `json:"behaviorName,omitempty"`
	ErrorLocation string `json:"errorLocation,omitempty"`
	StatusCode    int    `json:"-"`
}

// NewAPIError returns new APIError object containing data from the error response
// In case of body reading and unmarshaling errors, an APIError with only status code is returned
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
