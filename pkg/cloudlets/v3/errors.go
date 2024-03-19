package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/errs"
)

// Error is a cloudlets error interface.
type Error struct {
	Type        string          `json:"type,omitempty"`
	Title       string          `json:"title,omitempty"`
	Instance    string          `json:"instance,omitempty"`
	Status      int             `json:"status,omitempty"`
	Errors      json.RawMessage `json:"errors,omitempty"`
	Detail      string          `json:"detail"`
	RequestID   string          `json:"requestId,omitempty"`
	RequestTime string          `json:"requestTime,omitempty"`
	ClientIP    string          `json:"clientIp,omitempty"`
	ServerIP    string          `json:"serverIp,omitempty"`
	Method      string          `json:"method,omitempty"`
}

// ErrPolicyNotFound is returned when policy was not found
var ErrPolicyNotFound = errors.New("policy not found")

// Error parses an error from the response.
func (c *cloudlets) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = "Failed to read error body"
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		c.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Cloudlets API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.Status = r.StatusCode

	return &e
}

func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons.
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
