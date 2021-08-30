package cloudlets

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"net/url"
)

type (
	// LoadBalancer is a cloudlets LoadBalancer API interface
	LoadBalancer interface {
		// ListOrigins lists all origins of specified type for the current account
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getloadbalancingconfigs
		ListOrigins(context.Context, ListOriginsRequest) (Origins, error)

		// GetOrigin gets specific origin by originID
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getorigin
		GetOrigin(context.Context, string) (*Origin, error)
	}

	// Origin is a response returned by GetOrigin
	Origin struct {
		OriginID    string     `json:"originId"`
		Hostname    string     `json:"hostname"`
		Type        OriginType `json:"type"`
		Checksum    string     `json:"checksum"`
		Description string     `json:"description"`
		Akamaized   bool       `json:"akamaized"`
	}

	// Origins is a response returned by ListOrigins
	Origins []Origin

	// OriginType is a type for Origin Type
	OriginType string

	// ListOriginsRequest describes the parameters of the ListOrigins request
	ListOriginsRequest struct {
		Type OriginType
	}
)

const (
	// OriginTypeAll is a value to use when you want ListOrigins to return origins of all types
	OriginTypeAll OriginType = ""
	// OriginTypeCustomer is a value to use when you want ListOrigins to return only origins of CUSTOMER type
	OriginTypeCustomer OriginType = "CUSTOMER"
	// OriginTypeApplicationLoadBalancer is a value to use when you want ListOrigins to return only origins of APPLICATION_LOAD_BALANCER type
	OriginTypeApplicationLoadBalancer OriginType = "APPLICATION_LOAD_BALANCER"
	// OriginTypeNetStorage is a value to use when you want ListOrigins to return only origins of NETSTORAGE type
	OriginTypeNetStorage OriginType = "NETSTORAGE"
)

var (
	// ErrListOrigins is returned when ListOrigins fails
	ErrListOrigins = errors.New("list origins")
	// ErrGetOrigin is returned when GetOrigin fails
	ErrGetOrigin = errors.New("get origin")
)

// Validate validates ListOriginsRequest
func (v ListOriginsRequest) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(v.Type, validation.In(OriginTypeCustomer, OriginTypeApplicationLoadBalancer, OriginTypeNetStorage, OriginTypeAll)),
	}.Filter()
}

func (c *cloudlets) ListOrigins(ctx context.Context, params ListOriginsRequest) (Origins, error) {
	logger := c.Log(ctx)
	logger.Debug("ListOrigins")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListOrigins, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cloudlets/api/v2/origins")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListOrigins, err)
	}
	if params.Type != OriginTypeAll {
		q := uri.Query()
		q.Add("type", string(params.Type))
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListOrigins, err)
	}

	var result Origins
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListOrigins, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListOrigins, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) GetOrigin(ctx context.Context, originID string) (*Origin, error) {
	logger := c.Log(ctx)
	logger.Debug("GetOrigin")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s", originID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetOrigin, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetOrigin, err)
	}

	var result Origin
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetOrigin, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetOrigin, c.Error(resp))
	}

	return &result, nil
}
