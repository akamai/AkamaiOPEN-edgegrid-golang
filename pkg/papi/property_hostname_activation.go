package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPropertyHostnameActivationRequest contains parameters required to get a property hostname activation
	GetPropertyHostnameActivationRequest struct {
		PropertyID           string
		HostnameActivationID string
		ContractID           string
		GroupID              string
		IncludeHostnames     bool
	}
	// GetPropertyHostnameActivationResponse contains GET response returned while fetching property hostname activation
	GetPropertyHostnameActivationResponse struct {
		AccountID          string                    `json:"accountId"`
		ContractID         string                    `json:"contractId"`
		GroupID            string                    `json:"groupId"`
		HostnameActivation HostnameActivationGetItem `json:"hostnameActivation"`
	}
	// ListPropertyHostnameActivationsRequest contains path and query params used for listing property hostname activations
	ListPropertyHostnameActivationsRequest struct {
		PropertyID string
		Offset     int
		Limit      int
		ContractID string
		GroupID    string
	}
	// ListPropertyHostnameActivationsResponse contains GET response returned while fetching property hostname activations
	ListPropertyHostnameActivationsResponse struct {
		AccountID           string                  `json:"accountId"`
		ContractID          string                  `json:"contractId"`
		GroupID             string                  `json:"groupId"`
		HostnameActivations HostnameActivationsList `json:"hostnameActivations"`
	}
	// CancelPropertyHostnameActivationRequest contains path and query params used for canceling hostname activation
	CancelPropertyHostnameActivationRequest struct {
		PropertyID           string
		HostnameActivationID string
		ContractID           string
		GroupID              string
	}
	// CancelPropertyHostnameActivationResponse contains DELETE response returned when canceling property hostname activation
	CancelPropertyHostnameActivationResponse struct {
		AccountID          string                       `json:"accountId"`
		ContractID         string                       `json:"contractId"`
		GroupID            string                       `json:"groupId"`
		HostnameActivation HostnameActivationCancelItem `json:"hostnameActivation"`
	}
	// HostnameActivationsList contains returned hostname activations
	HostnameActivationsList struct {
		Items            []HostnameActivationListItem `json:"items"`
		TotalItems       *int                         `json:"totalItems"`
		CurrentItemCount *int                         `json:"currentItemCount"`
		NextLink         *string                      `json:"nextLink"`
		PreviousLink     *string                      `json:"previousLink"`
	}
	// HostnameActivationGetItem contains property hostname activation details
	HostnameActivationGetItem struct {
		ActivationType       string                 `json:"activationType"`
		HostnameActivationID string                 `json:"hostnameActivationId"`
		PropertyName         string                 `json:"propertyName"`
		PropertyID           string                 `json:"propertyId"`
		Network              string                 `json:"network"`
		Status               string                 `json:"status"`
		SubmitDate           time.Time              `json:"submitDate,omitempty"`
		UpdateDate           time.Time              `json:"updateDate,omitempty"`
		Note                 string                 `json:"note,omitempty"`
		NotifyEmails         []string               `json:"notifyEmails,omitempty"`
		Hostnames            []PropertyHostnameItem `json:"hostnames,omitempty"`
	}
	// HostnameActivationListItem contains property hostname activation details
	HostnameActivationListItem struct {
		ActivationType       string    `json:"activationType"`
		HostnameActivationID string    `json:"hostnameActivationId"`
		PropertyName         string    `json:"propertyName"`
		PropertyID           string    `json:"propertyId"`
		Network              string    `json:"network"`
		Status               string    `json:"status"`
		SubmitDate           time.Time `json:"submitDate,omitempty"`
		UpdateDate           time.Time `json:"updateDate,omitempty"`
		Note                 string    `json:"note,omitempty"`
		NotifyEmails         []string  `json:"notifyEmails,omitempty"`
	}
	// HostnameActivationCancelItem contains property hostname activation details
	HostnameActivationCancelItem struct {
		ActivationType       string    `json:"activationType"`
		HostnameActivationID string    `json:"hostnameActivationId"`
		PropertyName         string    `json:"propertyName"`
		PropertyID           string    `json:"propertyId"`
		Network              string    `json:"network"`
		Status               string    `json:"status"`
		SubmitDate           time.Time `json:"submitDate,omitempty"`
		UpdateDate           time.Time `json:"updateDate,omitempty"`
		Note                 string    `json:"note,omitempty"`
		NotifyEmails         []string  `json:"notifyEmails,omitempty"`
		PropertyVersion      int       `json:"propertyVersion,omitempty"`
	}
	// propertyHostnamesHelper is used to unwrap property hostnames returned by API
	propertyHostnamesHelper struct {
		HostnameActivationGetItem
		Hostnames propertyHostnameItems `json:"hostnames,omitempty"`
	}
	propertyHostnameItems struct {
		Items []PropertyHostnameItem `json:"items"`
	}
	// PropertyHostnameItem contains hostname details returned after GET operation
	PropertyHostnameItem struct {
		CertProvisioningType string         `json:"certProvisioningType"`
		CnameFrom            string         `json:"cnameFrom"`
		CnameTo              string         `json:"cnameTo"`
		EdgeHostnameID       string         `json:"edgeHostnameId"`
		CertStatus           CertStatusItem `json:"certStatus"`
		Action               string         `json:"action"`
	}
	// hostnameActivationGetHelper is used to unwrap single element list of hostname activations returned by API
	hostnameActivationGetHelper struct {
		GetPropertyHostnameActivationResponse
		HostnameActivations hostnameActivationGetItems `json:"hostnameActivations"`
	}
	hostnameActivationGetItems struct {
		Items []propertyHostnamesHelper `json:"items"`
	}
	//hostnameActivationCancelHelper is used to unwrap single element list of hostname activations returned by API
	hostnameActivationCancelHelper struct {
		CancelPropertyHostnameActivationResponse
		HostnameActivations hostnameActivationsCancelItems `json:"hostnameActivations"`
	}
	hostnameActivationsCancelItems struct {
		Items []HostnameActivationCancelItem `json:"items"`
	}
)

// unwrapSingleElement extracts hostname activation from single element list
func (h hostnameActivationGetHelper) unwrapSingleElement() *GetPropertyHostnameActivationResponse {
	response := h.GetPropertyHostnameActivationResponse
	if len(h.HostnameActivations.Items) > 0 {
		activation := h.HostnameActivations.Items[0]
		response.HostnameActivation = activation.HostnameActivationGetItem
		response.HostnameActivation.Hostnames = activation.Hostnames.Items
	}
	return &response
}

// unwrapSingleElement extracts hostname activation from single element list
func (h hostnameActivationCancelHelper) unwrapSingleElement() *CancelPropertyHostnameActivationResponse {
	response := h.CancelPropertyHostnameActivationResponse
	if len(h.HostnameActivations.Items) > 0 {
		response.HostnameActivation = h.HostnameActivations.Items[0]
	}
	return &response
}

// Validate validates GetPropertyHostnameActivationRequest
func (r GetPropertyHostnameActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID":           validation.Validate(r.PropertyID, validation.Required),
		"ContractID":           validation.Validate(r.ContractID, validation.Required),
		"GroupID":              validation.Validate(r.GroupID, validation.Required),
		"HostnameActivationID": validation.Validate(r.HostnameActivationID, validation.Required),
	})
}

// Validate validates ListPropertyHostnameActivationsRequest
func (r ListPropertyHostnameActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"ContractID": validation.Validate(r.ContractID, validation.Required.When(r.GroupID != "").Error("cannot be blank when GroupID is provided")),
		"GroupID":    validation.Validate(r.GroupID, validation.Required.When(r.ContractID != "").Error("cannot be blank when ContractID is provided")),
		"Offset":     validation.Validate(r.Offset, validation.Min(0)),
		"Limit":      validation.Validate(r.Limit, validation.Min(1)),
	})
}

// Validate validates CancelPropertyHostnameActivationRequest
func (r CancelPropertyHostnameActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID":           validation.Validate(r.PropertyID, validation.Required),
		"ContractID":           validation.Validate(r.ContractID, validation.Required),
		"GroupID":              validation.Validate(r.GroupID, validation.Required),
		"HostnameActivationID": validation.Validate(r.PropertyID, validation.Required),
	})
}

var (
	// ErrGetPropertyHostnameActivation represents error when fetching hostname activation fails
	ErrGetPropertyHostnameActivation = errors.New("fetching hostname activation")
	// ErrListPropertyHostnameActivations represents error when fetching hostname activations fails
	ErrListPropertyHostnameActivations = errors.New("fetching hostname activations")
	//ErrCancelPropertyHostnameActivation represents error when canceling hostname activation fails
	ErrCancelPropertyHostnameActivation = errors.New("canceling hostname activation")
	// ErrCancelPropertyHostnameActivationAlreadyAborted represents error when canceling aborted hostname activation
	ErrCancelPropertyHostnameActivationAlreadyAborted = errors.New("activation already aborted")
)

func (p *papi) GetPropertyHostnameActivation(ctx context.Context, params GetPropertyHostnameActivationRequest) (*GetPropertyHostnameActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetPropertyHostnameActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyHostnameActivation, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/papi/v1/properties/%s/hostname-activations/%s",
		params.PropertyID,
		params.HostnameActivationID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPropertyHostnameActivation, err)
	}
	q := uri.Query()
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.IncludeHostnames {
		q.Add("includeHostnames", "true")
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyHostnameActivation, err)
	}

	result := hostnameActivationGetHelper{}
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyHostnameActivation, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyHostnameActivation, p.Error(resp))
	}
	return result.unwrapSingleElement(), nil
}

func (p *papi) ListPropertyHostnameActivations(ctx context.Context, params ListPropertyHostnameActivationsRequest) (*ListPropertyHostnameActivationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListPropertyHostnameActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyHostnameActivation, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/papi/v1/properties/%s/hostname-activations",
		params.PropertyID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPropertyHostnameActivation, err)
	}
	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.Offset != 0 {
		q.Add("offset", strconv.Itoa(params.Offset))
	}
	if params.Limit != 0 {
		q.Add("limit", strconv.Itoa(params.Limit))
	}
	uri.RawQuery = q.Encode()

	var result ListPropertyHostnameActivationsResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyHostnameActivation, err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPropertyHostnameActivations, err)
	}

	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPropertyHostnameActivations, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) CancelPropertyHostnameActivation(ctx context.Context, params CancelPropertyHostnameActivationRequest) (*CancelPropertyHostnameActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeletePropertyHostnameActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCancelPropertyHostnameActivation, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/papi/v1/properties/%s/hostname-activations/%s",
		params.PropertyID,
		params.HostnameActivationID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCancelPropertyHostnameActivation, err)
	}
	q := uri.Query()
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCancelPropertyHostnameActivation, err)
	}

	result := hostnameActivationCancelHelper{}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCancelPropertyHostnameActivation, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("%s: %s", ErrCancelPropertyHostnameActivation, ErrCancelPropertyHostnameActivationAlreadyAborted)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCancelPropertyHostnameActivation, p.Error(resp))
	}

	return result.unwrapSingleElement(), nil
}
