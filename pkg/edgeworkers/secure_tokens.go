package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CreateSecureTokenRequest represents parameters for CreateSecureToken
	CreateSecureTokenRequest struct {
		ACL        string            `json:"acl,omitempty"`
		Expiry     int               `json:"expiry,omitempty"`
		Hostname   string            `json:"hostname,omitempty"`
		Network    ActivationNetwork `json:"network,omitempty"`
		PropertyID string            `json:"propertyId,omitempty"`
		URL        string            `json:"url,omitempty"`
	}

	// CreateSecureTokenResponse contains response from CreateSecureToken
	CreateSecureTokenResponse struct {
		AkamaiEWTrace string `json:"akamaiEwTrace"`
	}
)

// Validate validates CreateSecureTokenRequest
func (c CreateSecureTokenRequest) Validate() error {
	return validation.Errors{
		"ACL":      validation.Validate(c.ACL, validation.Empty.When(c.URL != "").Error("If you specify an acl don't specify a url.")),
		"Expiry":   validation.Validate(c.Expiry, validation.Min(1), validation.Max(720)),
		"Hostname": validation.Validate(c.Hostname, validation.Required.When(c.PropertyID == "").Error("To create an authentication token, provide either the hostname, or the propertyId")),
		"Network": validation.Validate(c.Network, validation.In(ActivationNetworkStaging, ActivationNetworkProduction).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '' (empty)", c.Network, ActivationNetworkStaging, ActivationNetworkProduction))), // If not specified, the token is created for the network where the last Property version activation occurred.
		"PropertyID": validation.Validate(c.PropertyID, validation.Required.When(c.Hostname == "").Error("To create an authentication token, provide either the hostname, or the propertyId")),
		"URL":        validation.Validate(c.URL, validation.Empty.When(c.ACL != "").Error(" If you specify a url don't specify an acl")),
	}.Filter()
}

var (
	// ErrCreateSecureToken is returned in case an error occurs on CreateSecureToken operation
	ErrCreateSecureToken = errors.New("create secure token")
)

func (e *edgeworkers) CreateSecureToken(ctx context.Context, params CreateSecureTokenRequest) (*CreateSecureTokenResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("CreateSecureToken")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateSecureToken, ErrStructValidation, err)
	}

	uri := "/edgeworkers/v1/secure-token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateSecureToken, err)
	}

	var result CreateSecureTokenResponse
	resp, err := e.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateSecureToken, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateSecureToken, e.Error(resp))
	}

	return &result, nil
}
