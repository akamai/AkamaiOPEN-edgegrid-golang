package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// ThirdPartyCSR is a CPS API enabling management of third-party certificates
	ThirdPartyCSR interface {
		// GetChangeThirdPartyCSR gets certificate signing request
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeThirdPartyCSR(ctx context.Context, params GetChangeRequest) (*ThirdPartyCSRResponse, error)
	}

	// ThirdPartyCSRResponse is a response object containing list of csrs
	ThirdPartyCSRResponse struct {
		CSRs []CertSigningRequest `json:"csrs"`
	}

	// CertSigningRequest holds CSR
	CertSigningRequest struct {
		CSR          string `json:"csr"`
		KeyAlgorithm string `json:"keyAlgorithm"`
	}
)

var (
	// ErrGetChangeThirdPartyCSR is returned when GetChangeThirdPartyCSR fails
	ErrGetChangeThirdPartyCSR = errors.New("get change third-party csr")
)

func (c *cps) GetChangeThirdPartyCSR(ctx context.Context, params GetChangeRequest) (*ThirdPartyCSRResponse, error) {
	c.Log(ctx).Debug("GetChangeThirdPartyCSR")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeThirdPartyCSR, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/info/third-party-csr",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeThirdPartyCSR, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.csr.v2+json")

	var result ThirdPartyCSRResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeThirdPartyCSR, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeThirdPartyCSR, c.Error(resp))
	}

	return &result, nil
}
