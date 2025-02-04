package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PatchPropertyHostnameBucketRequest contains path, query and body params used for adding or removing hostnames to and from property's hostname bucket
	PatchPropertyHostnameBucketRequest struct {
		PropertyID string
		ContractID string
		GroupID    string
		Body       PatchPropertyHostnameBucketBody
	}
	// PatchPropertyHostnameBucketResponse contains PATCH response returned when patching property hostname bucket
	PatchPropertyHostnameBucketResponse struct {
		ActivationLink string              `json:"activationLink"`
		ActivationID   string              `json:"activationId"`
		Hostnames      []PatchHostnameItem `json:"hostnames"`
	}
	// PatchPropertyHostnameBucketBody contains body params for PatchPropertyHostnameBucket
	PatchPropertyHostnameBucketBody struct {
		Add          []PatchPropertyHostnameBucketAdd `json:"add,omitempty"`
		Remove       []string                         `json:"remove,omitempty"`
		Network      string                           `json:"network"`
		NotifyEmails []string                         `json:"notifyEmails,omitempty"`
		Note         string                           `json:"note,omitempty"`
	}
	// PatchPropertyHostnameBucketAdd contains params for adding property hostname to a bucket
	PatchPropertyHostnameBucketAdd struct {
		EdgeHostnameID       string            `json:"edgeHostnameId"`
		CertProvisioningType CertType          `json:"certProvisioningType"`
		CnameType            HostnameCnameType `json:"cnameType"`
		CnameFrom            string            `json:"cnameFrom"`
	}
	// PatchHostnameItem contains hostname details returned after PATCH operation
	PatchHostnameItem struct {
		CertProvisioningType CertType          `json:"certProvisioningType"`
		CnameFrom            string            `json:"cnameFrom"`
		CnameTo              string            `json:"cnameTo"`
		CnameType            HostnameCnameType `json:"cnameType"`
		EdgeHostnameID       string            `json:"edgeHostnameId"`
		CertStatus           CertStatusItem    `json:"certStatus"`
		Action               string            `json:"action"`
	}
)

// Validate validates PatchPropertyHostnameBucketRequest
func (r PatchPropertyHostnameBucketRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Body":       validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates PatchPropertyHostnameBucketBody
func (b PatchPropertyHostnameBucketBody) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(b.Network, validation.Required),
		"Add":     validation.Validate(b.Add),
		"": validation.Validate(nil,
			validation.By(func(interface{}) error {
				if len(b.Add) == 0 && len(b.Remove) == 0 {
					return fmt.Errorf("at least one hostname is required in add or remove list")
				}
				return nil
			})),
	}.Filter()
}

// Validate validates PatchPropertyHostnameBucketBody
func (b PatchPropertyHostnameBucketAdd) Validate() error {
	return validation.Errors{
		"EdgeHostnameID":       validation.Validate(b.EdgeHostnameID, validation.Required),
		"CertProvisioningType": validation.Validate(b.CertProvisioningType, validation.Required),
		"CnameType":            validation.Validate(b.CnameType, validation.Required),
		"CnameFrom":            validation.Validate(b.CnameType, validation.Required),
	}.Filter()
}

var (
	// ErrPatchPropertyHostnameBucket represents error when patching property hostname bucket fails
	ErrPatchPropertyHostnameBucket = errors.New("patching property hostname bucket")
)

func (p *papi) PatchPropertyHostnameBucket(ctx context.Context, params PatchPropertyHostnameBucketRequest) (*PatchPropertyHostnameBucketResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("PatchPropertyHostnameBucket")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrPatchPropertyHostnameBucket, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/papi/v1/properties/%s/hostnames",
		params.PropertyID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrPatchPropertyHostnameBucket, err)
	}
	q := uri.Query()
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrPatchPropertyHostnameBucket, err)
	}

	var result PatchPropertyHostnameBucketResponse

	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrPatchPropertyHostnameBucket, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrPatchPropertyHostnameBucket, p.Error(resp))
	}
	return &result, nil
}
