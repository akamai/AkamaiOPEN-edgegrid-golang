package cloudwrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/errs"
)

type (
	// Error is a cloudwrapper error implementation
	// For details on possible error types, refer to: https://techdocs.akamai.com/cloud-wrapper/reference/errors
	Error struct {
		Type        string      `json:"type,omitempty"`
		Title       string      `json:"title,omitempty"`
		Instance    string      `json:"instance"`
		Status      int         `json:"status"`
		Detail      string      `json:"detail"`
		Errors      []ErrorItem `json:"errors"`
		Method      string      `json:"method"`
		ServerIP    string      `json:"serverIp"`
		ClientIP    string      `json:"clientIp"`
		RequestID   string      `json:"requestId"`
		RequestTime string      `json:"requestTime"`
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

const (
	configurationNotFoundType = "/cloud-wrapper/error-types/not-found"
	deletionNotAllowedType    = "/cloud-wrapper/error-types/forbidden"
)

var (
	// ErrConfigurationNotFound is returned when configuration was not found
	ErrConfigurationNotFound = errors.New("configuration not found")
	// ErrDeletionNotAllowed is returned when user has insufficient permissions to delete configuration
	ErrDeletionNotAllowed = errors.New("deletion not allowed")
)

// Error parses an error from the response
func (c *cloudwrapper) Error(r *http.Response) error {
	var result Error
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		result.Status = r.StatusCode
		result.Title = "Failed to read error body"
		result.Detail = err.Error()
		return &result
	}

	if err = json.Unmarshal(body, &result); err != nil {
		c.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		result.Title = "Failed to unmarshal error body. Cloud Wrapper API failed. Check details for more information."
		result.Detail = errs.UnescapeContent(string(body))
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
	if errors.Is(target, ErrConfigurationNotFound) {
		return e.Status == http.StatusNotFound && e.Type == configurationNotFoundType
	}
	if errors.Is(target, ErrDeletionNotAllowed) {
		return e.Status == http.StatusForbidden && e.Type == deletionNotAllowedType
	}

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
