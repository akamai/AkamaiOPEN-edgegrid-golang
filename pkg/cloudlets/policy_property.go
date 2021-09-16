package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type (
	// PolicyProperty interface is a cloudlets API interface for policy associated properties
	PolicyProperty interface {
		// GetPolicyProperties gets all the associated properties by the policyID
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicyproperties
		GetPolicyProperties(context.Context, int64) (GetPolicyPropertiesResponse, error)
	}

	// GetPolicyPropertiesResponse contains response data for GetPolicyProperties
	GetPolicyPropertiesResponse map[string]AssociateProperty

	// AssociateProperty contains the response data for a single property
	AssociateProperty struct {
		GroupID       int64         `json:"groupId"`
		ID            int64         `json:"id"`
		Name          string        `json:"name"`
		NewestVersion NetworkStatus `json:"newestVersion"`
		Production    NetworkStatus `json:"production"`
		Staging       NetworkStatus `json:"staging"`
	}

	// NetworkStatus is the type for NetworkStatus of any activation
	NetworkStatus struct {
		ActivatedBy        string                     `json:"activatedBy"`
		ActivationDate     string                     `json:"activationDate"`
		Version            int64                      `json:"version"`
		CloudletsOrigins   map[string]CloudletsOrigin `json:"cloudletsOrigins"`
		ReferencedPolicies []string                   `json:"referencedPolicies"`
	}

	// CloudletsOrigin is the type for CloudletsOrigins in NetworkStatus
	CloudletsOrigin struct {
		OriginID    string     `json:"id"`
		Hostname    string     `json:"hostname"`
		Type        OriginType `json:"type"`
		Checksum    string     `json:"checksum"`
		Description string     `json:"description"`
	}
)

var (
	// ErrGetPolicyProperties is returned when GetPolicyProperties fails
	ErrGetPolicyProperties = errors.New("get policy properties")
)

// GetPolicyProperties gets all the associated properties by the policyID
func (c *cloudlets) GetPolicyProperties(ctx context.Context, policyID int64) (GetPolicyPropertiesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicyProperties")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/properties", policyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPolicyProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyProperties, err)
	}

	var result GetPolicyPropertiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyProperties, c.Error(resp))
	}

	return result, nil
}
