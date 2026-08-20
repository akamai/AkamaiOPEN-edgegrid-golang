package datastream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/errs"
)

type (
	// Error is a ds error interface
	Error struct {
		Type       string          `json:"type"`
		Title      string          `json:"title"`
		Detail     string          `json:"detail"`
		Instance   string          `json:"instance"`
		StatusCode int             `json:"statusCode"`
		Errors     []RequestErrors `json:"errors"`
	}

	// RequestErrors is an optional errors array that lists potentially more than one problem detected in the request
	RequestErrors struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Instance string `json:"instance,omitempty"`
		Detail   string `json:"detail"`
	}
)

var (
	// ErrGetProperties is returned when GetProperties fails
	ErrGetProperties = errors.New("list properties")

	// ErrGetDatasetFields is returned when GetDatasetFields fails
	ErrGetDatasetFields = errors.New("list data set fields")

	// ErrGetAppSecConfigs is returned when GetAppSecConfigs fails
	ErrGetAppSecConfigs = errors.New("list appsec configs")

	// ErrListAnswerXServiceIDs is returned when ListAnswerXServiceIDs fails
	ErrListAnswerXServiceIDs = errors.New("list answerx service ids")

	// ErrCreateStream represents an error when stream creation fails
	ErrCreateStream = errors.New("creating stream")

	// ErrGetStream represents error when fetching stream fails
	ErrGetStream = errors.New("fetching stream information")

	// ErrUpdateStream represents error when updating stream fails
	ErrUpdateStream = errors.New("updating stream")

	// ErrDeleteStream represents error when deleting stream fails
	ErrDeleteStream = errors.New("deleting stream")

	// ErrListStreams represents error when listing streams fails
	ErrListStreams = errors.New("listing streams")

	// ErrActivateStream is returned when ActivateStream fails
	ErrActivateStream = errors.New("activate stream")

	// ErrDeactivateStream is returned when DeactivateStream fails
	ErrDeactivateStream = errors.New("deactivate stream")

	// ErrGetActivationHistory is returned when GetActivationHistory fails
	ErrGetActivationHistory = errors.New("view activation history")
)

// Error parses an error from the response
func (d *ds) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.StatusCode = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		d.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. DataStream2 API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.StatusCode = r.StatusCode

	return &e
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

	if e.StatusCode != t.StatusCode {
		return false
	}

	return e.Error() == t.Error()
}
