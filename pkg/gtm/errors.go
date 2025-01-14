package gtm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/errs"
)

var (
	// ErrNotFound used when status code is 404 Not Found
	ErrNotFound = errors.New("404 Not Found")

	// ErrNoDatacenterAssignedToMap occurs when no datacenter is assigned to the map target during the creation of a geographic property.
	ErrNoDatacenterAssignedToMap = errors.New("no datacenter is assigned to map target (all others)")
)

type (
	// Error is a gtm error interface
	Error struct {
		Type          string  `json:"type"`
		Title         string  `json:"title"`
		Detail        string  `json:"detail"`
		Instance      string  `json:"instance,omitempty"`
		BehaviorName  string  `json:"behaviorName,omitempty"`
		ErrorLocation string  `json:"errorLocation,omitempty"`
		StatusCode    int     `json:"-"`
		Errors        []Error `json:"errors"`
	}
)

// Error parses an error from the response
func (g *gtm) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.StatusCode = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		g.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. GTM API failed. Check details for more information."
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

	if errors.Is(target, ErrNotFound) && e.StatusCode == http.StatusNotFound {
		return true
	}

	if errors.Is(target, ErrNoDatacenterAssignedToMap) && strings.Contains(e.Detail, "no datacenter is assigned to map target (all others)") {
		return true
	}

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
