package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Includes contains operations available on Include resource
	Includes interface {
		// ListIncludes lists Includes available for the current contract and group
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-includes
		ListIncludes(context.Context, ListIncludesRequest) (*ListIncludesResponse, error)

		// ListIncludeParents lists parents of a specific Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-parents
		ListIncludeParents(context.Context, ListIncludeParentsRequest) (*ListIncludeParentsResponse, error)

		// GetInclude gets information about a specific Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include
		GetInclude(context.Context, GetIncludeRequest) (*GetIncludeResponse, error)

		// CreateInclude creates a new Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-includes
		CreateInclude(context.Context, CreateIncludeRequest) (*CreateIncludeResponse, error)

		// DeleteInclude deletes an Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/delete-include
		DeleteInclude(context.Context, DeleteIncludeRequest) (*DeleteIncludeResponse, error)
	}

	// ListIncludesRequest contains parameters used to list includes
	ListIncludesRequest struct {
		ContractID string
		GroupID    string
	}

	// ListIncludesResponse represents a response object returned by ListIncludes
	ListIncludesResponse struct {
		Includes IncludeItems `json:"includes"`
	}

	// ListIncludeParentsRequest contains parameters used to list parents of an include
	ListIncludeParentsRequest struct {
		ContractID string
		GroupID    string
		IncludeID  string
	}

	// ListIncludeParentsResponse represents a response object returned by ListIncludeParents
	ListIncludeParentsResponse struct {
		Properties ParentPropertyItems `json:"properties"`
	}

	// GetIncludeRequest contains parameters used to fetch an include
	GetIncludeRequest struct {
		ContractID string
		GroupID    string
		IncludeID  string
	}

	// GetIncludeResponse represents a response object returned by GetInclude
	GetIncludeResponse struct {
		Includes IncludeItems `json:"includes"`
		Include  Include      `json:"-"`
	}

	// CreateIncludeRequest contains parameters used to create an include
	CreateIncludeRequest struct {
		ContractID       string            `json:"-"`
		GroupID          string            `json:"-"`
		IncludeName      string            `json:"includeName"`
		IncludeType      IncludeType       `json:"includeType"`
		ProductID        string            `json:"productId"`
		RuleFormat       string            `json:"ruleFormat,omitempty"`
		CloneIncludeFrom *CloneIncludeFrom `json:"cloneFrom,omitempty"`
	}

	// CreateIncludeResponse represents a response object returned by CreateInclude
	CreateIncludeResponse struct {
		IncludeID       string `json:"-"`
		IncludeLink     string `json:"includeLink"`
		ResponseHeaders CreateIncludeResponseHeaders
	}

	// DeleteIncludeRequest contains parameters used to delete an include
	DeleteIncludeRequest struct {
		ContractID string
		GroupID    string
		IncludeID  string
	}

	// DeleteIncludeResponse represents a response object returned by DeleteInclude
	DeleteIncludeResponse struct {
		Message string `json:"message"`
	}

	// CloneIncludeFrom optionally identifies another include instance to clone when making a request to create new include
	CloneIncludeFrom struct {
		CloneFromVersionEtag string `json:"cloneFromVersionEtag,omitempty"`
		IncludeID            string `json:"includeId"`
		Version              int    `json:"version"`
	}

	// CreateIncludeResponseHeaders contains information received in response headers when making a request to create new include
	CreateIncludeResponseHeaders struct {
		IncludesLimitTotal     string
		IncludesLimitRemaining string
	}

	// Include represents an Include object
	Include struct {
		AccountID         string      `json:"accountId"`
		AssetID           string      `json:"assetId"`
		ContractID        string      `json:"contractId"`
		GroupID           string      `json:"groupId"`
		IncludeID         string      `json:"includeId"`
		IncludeName       string      `json:"includeName"`
		IncludeType       IncludeType `json:"includeType"`
		LatestVersion     int         `json:"latestVersion"`
		ProductionVersion *int        `json:"productionVersion"`
		PropertyType      *string     `json:"propertyType"`
		StagingVersion    *int        `json:"stagingVersion"`
	}

	// IncludeItems represents a list of Include objects
	IncludeItems struct {
		Items []Include `json:"items"`
	}

	// ParentProperty represents an include parent object
	ParentProperty struct {
		AccountID         string `json:"accountId"`
		AssetID           string `json:"assetId"`
		ContractID        string `json:"contractId"`
		GroupID           string `json:"groupId"`
		ProductionVersion *int   `json:"productionVersion,omitempty"`
		PropertyID        string `json:"propertyId"`
		PropertyName      string `json:"propertyName"`
		StagingVersion    *int   `json:"stagingVersion,omitempty"`
	}

	// ParentPropertyItems represents a list of ParentProperty objects
	ParentPropertyItems struct {
		Items []ParentProperty `json:"items"`
	}

	// IncludeType is type of include
	IncludeType string
)

const (
	// IncludeTypeMicroServices is used for creating a new microservices include
	IncludeTypeMicroServices IncludeType = "MICROSERVICES"

	// IncludeTypeCommonSettings is used for creating a new common_settings include
	IncludeTypeCommonSettings IncludeType = "COMMON_SETTINGS"
)

// Validate validates ListIncludesRequest
func (i ListIncludesRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(i.ContractID, validation.Required),
	}.Filter()
}

// Validate validates ListIncludeParentsRequest
func (i ListIncludeParentsRequest) Validate() error {
	return validation.Errors{
		"IncludeID": validation.Validate(i.IncludeID, validation.Required),
	}.Filter()
}

// Validate validates GetIncludeRequest
func (i GetIncludeRequest) Validate() error {
	errs := validation.Errors{
		"ContractID": validation.Validate(i.ContractID, validation.Required),
		"GroupID":    validation.Validate(i.GroupID, validation.Required),
		"IncludeID":  validation.Validate(i.IncludeID, validation.Required),
	}

	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates CreateIncludeRequest
func (i CreateIncludeRequest) Validate() error {
	errs := validation.Errors{
		"ContractID":                 validation.Validate(i.ContractID, validation.Required),
		"GroupID":                    validation.Validate(i.GroupID, validation.Required),
		"IncludeName":                validation.Validate(i.IncludeName, validation.Required),
		"IncludeType":                validation.Validate(i.IncludeType, validation.Required, validation.In(IncludeTypeMicroServices, IncludeTypeCommonSettings)),
		"ProductID":                  validation.Validate(i.ProductID, validation.Required),
		"CloneIncludeFrom.IncludeID": validation.Validate(i.CloneIncludeFrom, validation.When(i.CloneIncludeFrom != nil, validation.By(validateCloneIncludeID))),
		"CloneIncludeFrom.Version":   validation.Validate(i.CloneIncludeFrom, validation.When(i.CloneIncludeFrom != nil, validation.By(validateCloneVersion))),
	}

	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates DeleteIncludeRequest
func (i DeleteIncludeRequest) Validate() error {
	return validation.Errors{
		"IncludeID": validation.Validate(i.IncludeID, validation.Required),
	}.Filter()
}

// validateCloneIncludeID validates IncludeID under CloneIncludeFrom
func validateCloneIncludeID(value interface{}) error {
	v, ok := value.(*CloneIncludeFrom)
	if !ok {
		return fmt.Errorf("type %T is invalid. Must be *CloneIncludeFrom", value)
	}

	if v.IncludeID == "" {
		return fmt.Errorf("cannot be blank")
	}

	return nil
}

// validateCloneVersion validates Version under CloneIncludeFrom
func validateCloneVersion(value interface{}) error {
	v, ok := value.(*CloneIncludeFrom)
	if !ok {
		return fmt.Errorf("type %T is invalid. Must be *CloneIncludeFrom", value)
	}

	if v.Version == 0 {
		return fmt.Errorf("cannot be blank")
	}

	return nil
}

var (
	// ErrListIncludes is returned in case an error occurs on ListIncludes operation
	ErrListIncludes = errors.New("list Includes")
	// ErrListIncludeParents is returned in case an error occurs on ListIncludeParents operation
	ErrListIncludeParents = errors.New("list Include Parents")
	// ErrGetInclude is returned in case an error occurs on GetInclude operation
	ErrGetInclude = errors.New("get an Include")
	// ErrCreateInclude is returned in case an error occurs on CreateInclude operation
	ErrCreateInclude = errors.New("create an Include")
	// ErrDeleteInclude is returned in case an error occurs on DeleteInclude operation
	ErrDeleteInclude = errors.New("delete an Include")
)

func (p *papi) ListIncludes(ctx context.Context, params ListIncludesRequest) (*ListIncludesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludes, ErrStructValidation, err)
	}

	uri, err := url.Parse("/papi/v1/includes")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListIncludes, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludes, err)
	}

	var result ListIncludesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", ErrListIncludes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludes, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) ListIncludeParents(ctx context.Context, params ListIncludeParentsRequest) (*ListIncludeParentsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludeParents")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludeParents, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/parents", params.IncludeID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListIncludeParents, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludeParents, err)
	}

	var result ListIncludeParentsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListIncludeParents, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludeParents, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) GetInclude(ctx context.Context, params GetIncludeRequest) (*GetIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetInclude, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s", params.IncludeID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetInclude, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetInclude, err)
	}

	var result GetIncludeResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetInclude, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetInclude, p.Error(resp))
	}

	if len(result.Includes.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: IncludeID: %s", ErrGetInclude, ErrNotFound, params.IncludeID)
	}
	result.Include = result.Includes.Items[0]

	return &result, nil
}

func (p *papi) CreateInclude(ctx context.Context, params CreateIncludeRequest) (*CreateIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateInclude, ErrStructValidation, err)
	}

	uri, err := url.Parse("/papi/v1/includes")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateInclude, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateInclude, err)
	}

	var result CreateIncludeResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateInclude, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateInclude, p.Error(resp))
	}

	result.ResponseHeaders.IncludesLimitTotal = resp.Header.Get("x-limit-includes-per-contract-limit")
	result.ResponseHeaders.IncludesLimitRemaining = resp.Header.Get("x-limit-includes-per-contract-remaining")

	id, err := ResponseLinkParse(result.IncludeLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateInclude, ErrInvalidResponseLink, err)
	}
	result.IncludeID = id

	return &result, nil
}

func (p *papi) DeleteInclude(ctx context.Context, params DeleteIncludeRequest) (*DeleteIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeleteInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteInclude, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s", params.IncludeID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrDeleteInclude, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteInclude, err)
	}

	var result DeleteIncludeResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteInclude, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeleteInclude, p.Error(resp))
	}

	return &result, nil
}
