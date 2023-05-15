package cloudwrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Error is a cloudwrapper error implementation
	// For details on possible error types, refer to: https://techdocs.akamai.com/cloud-wrapper/reference/errors
	Error struct {
		Type     string      `json:"type,omitempty"`
		Title    string      `json:"title,omitempty"`
		Instance string      `json:"instance"`
		Status   int         `json:"status"`
		Detail   string      `json:"detail"`
		Errors   []ErrorItem `json:"errors"`
	}

	// ErrorItem is a cloud wrapper error's item
	ErrorItem struct {
		Type             string `json:"type"`
		Title            string `json:"title"`
		Detail           string `json:"detail"`
		IllegalValue     any    `json:"illegalValue"`
		IllegalParameter string `json:"illegalParameter"`
	}
)

// Error parses an error from the response
func (e *cloudwrapper) Error(r *http.Response) error {
	var result Error
	var body []byte
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		result.Status = r.StatusCode
		result.Title = "Failed to read error body"
		result.Detail = err.Error()
		return &result
	}

	if err := json.Unmarshal(body, &result); err != nil {
		e.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		result.Title = string(body)
		result.Status = r.StatusCode
	}
	return &result
}

func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons
func (e *Error) Is(target error) bool {

	var t *Error
	if !errors.As(target, &t) {
		return false
	}

	if e == t {
		return true
	}

	if e.Status != t.Status {
		return false
	}

	return e.Error() == t.Error()
}
