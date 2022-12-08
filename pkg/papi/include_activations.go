package papi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// IncludeActivations contains operations available on IncludeVersion resource
	IncludeActivations interface {
		// ActivateInclude creates a new include activation, which deactivates any current activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-include-activation
		ActivateInclude(context.Context, ActivateIncludeRequest) (*ActivationIncludeResponse, error)

		// DeactivateInclude deactivates the include activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-include-activation
		DeactivateInclude(context.Context, DeactivateIncludeRequest) (*DeactivationIncludeResponse, error)

		// CancelIncludeActivation cancels specified include activation, if it is still in `PENDING` state
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/delete-include-activation
		CancelIncludeActivation(context.Context, CancelIncludeActivationRequest) (*CancelIncludeActivationResponse, error)

		// GetIncludeActivation gets details about an activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-activation
		GetIncludeActivation(context.Context, GetIncludeActivationRequest) (*GetIncludeActivationResponse, error)

		// ListIncludeActivations lists all activations for all versions of the include, on both production and staging networks
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-activations
		ListIncludeActivations(context.Context, ListIncludeActivationsRequest) (*ListIncludeActivationsResponse, error)
	}

	// ActivateIncludeRequest contains parameters used to activate include
	ActivateIncludeRequest ActivateOrDeactivateIncludeRequest

	// DeactivateIncludeRequest contains parameters used to deactivate include
	DeactivateIncludeRequest ActivateOrDeactivateIncludeRequest

	// ActivateOrDeactivateIncludeRequest contains parameters used to activate or deactivate include
	ActivateOrDeactivateIncludeRequest struct {
		IncludeID              string            `json:"-"`
		Version                int               `json:"includeVersion"`
		Network                ActivationNetwork `json:"network"`
		Note                   string            `json:"note"`
		NotifyEmails           []string          `json:"notifyEmails"`
		AcknowledgeWarnings    []string          `json:"acknowledgeWarnings,omitempty"`
		AcknowledgeAllWarnings bool              `json:"acknowledgeAllWarnings"`
		IgnoreHTTPErrors       *bool             `json:"ignoreHttpErrors,omitempty"`
		ComplianceRecord       complianceRecord  `json:"complianceRecord,omitempty"`
	}

	// CancelIncludeActivationRequest contains parameters used to cancel pending activation of include
	CancelIncludeActivationRequest struct {
		ContractID   string
		GroupID      string
		IncludeID    string
		ActivationID string
	}

	// CancelIncludeActivationResponse represents a response object returned by CancelIncludeActivation operation
	CancelIncludeActivationResponse ListIncludeActivationsResponse

	// ActivationIncludeResponse represents a response object returned by ActivateInclude operation
	ActivationIncludeResponse struct {
		ActivationID   string `json:"-"`
		ActivationLink string `json:"activationLink"`
	}

	// DeactivationIncludeResponse represents a response object returned by DeactivateInclude operation
	DeactivationIncludeResponse struct {
		ActivationID   string `json:"-"`
		ActivationLink string `json:"activationLink"`
	}

	// GetIncludeActivationRequest contains parameters used to get the include activation
	GetIncludeActivationRequest struct {
		IncludeID    string
		ActivationID string
	}

	// GetIncludeActivationResponse represents a response object returned by GetIncludeActivation
	GetIncludeActivationResponse struct {
		AccountID   string                `json:"accountId"`
		ContractID  string                `json:"contractId"`
		GroupID     string                `json:"groupId"`
		Activations IncludeActivationsRes `json:"activations"`
		Validations *Validations          `json:"validations,omitempty"`
		Activation  IncludeActivation     `json:"-"`
	}

	// Validations represent include activation validation object
	Validations struct {
		ValidationSummary          ValidationSummary  `json:"validationSummary"`
		ValidationProgressItemList ValidationProgress `json:"validationProgressItemList"`
		Network                    ActivationNetwork  `json:"network"`
	}

	// ValidationSummary represent include activation validation summary object
	ValidationSummary struct {
		CompletePercent      float64 `json:"completePercent"`
		HasValidationError   bool    `json:"hasValidationError"`
		HasValidationWarning bool    `json:"hasValidationWarning"`
		HasSystemError       bool    `json:"hasSystemError"`
		HasClientError       bool    `json:"hasClientError"`
		MessageState         string  `json:"messageState"`
	}

	// ValidationProgress represents include activation validation progress object
	ValidationProgress struct {
		ErrorItems []ErrorItem `json:"errorItemsList"`
	}

	// ErrorItem represents validation progress error item object
	ErrorItem struct {
		VersionID             int    `json:"versionId"`
		PropertyName          string `json:"propertyName"`
		VersionNumber         int    `json:"versionNumber"`
		HasValidationError    bool   `json:"hasValidationError"`
		HasValidationWarning  bool   `json:"hasValidationWarning"`
		ValidationResultsLink string `json:"validationResultsLink"`
	}

	// ListIncludeActivationsRequest contains parameters used to list the include activations
	ListIncludeActivationsRequest struct {
		IncludeID  string
		ContractID string
		GroupID    string
	}

	// ListIncludeActivationsResponse represents a response object returned by ListIncludeActivations
	ListIncludeActivationsResponse struct {
		AccountID   string                `json:"accountId"`
		ContractID  string                `json:"contractId"`
		GroupID     string                `json:"groupId"`
		Activations IncludeActivationsRes `json:"activations"`
	}

	// IncludeActivationsRes represents Activations object
	IncludeActivationsRes struct {
		Items []IncludeActivation `json:"items"`
	}

	// IncludeActivation represents an include activation object
	IncludeActivation struct {
		ActivationID        string                  `json:"activationId"`
		Network             ActivationNetwork       `json:"network"`
		ActivationType      ActivationType          `json:"activationType"`
		Status              ActivationStatus        `json:"status"`
		SubmitDate          string                  `json:"submitDate"`
		UpdateDate          string                  `json:"updateDate"`
		Note                string                  `json:"note"`
		NotifyEmails        []string                `json:"notifyEmails"`
		FMAActivationState  string                  `json:"fmaActivationState"`
		FallbackInfo        *ActivationFallbackInfo `json:"fallbackInfo"`
		IncludeID           string                  `json:"includeId"`
		IncludeName         string                  `json:"includeName"`
		IncludeType         IncludeType             `json:"includeType"`
		IncludeVersion      int                     `json:"includeVersion"`
		IncludeActivationID string                  `json:"includeActivationId"`
	}

	// complianceRecord is an interface for ComplianceRecord data type
	complianceRecord interface {
		noncomplianceReasonType() string
	}

	// ComplianceRecordNone holds data relevant for ComplianceRecord with noncomplianceReason 'None'
	ComplianceRecordNone struct {
		CustomerEmail  string `json:"customerEmail"`
		PeerReviewedBy string `json:"peerReviewedBy"`
		UnitTested     bool   `json:"unitTested"`
		TicketID       string `json:"ticketId,omitempty"`
	}

	// ComplianceRecordOther holds data relevant for ComplianceRecord with noncomplianceReason 'Other'
	ComplianceRecordOther struct {
		OtherNoncomplianceReason string `json:"otherNoncomplianceReason"`
		TicketID                 string `json:"ticketId,omitempty"`
	}

	// ComplianceRecordNoProductionTraffic holds data relevant for ComplianceRecord with noncomplianceReason 'NoProductionTraffic'
	ComplianceRecordNoProductionTraffic struct {
		TicketID string `json:"ticketId,omitempty"`
	}

	// ComplianceRecordEmergency holds data relevant for ComplianceRecord with noncomplianceReason 'Emergency'
	ComplianceRecordEmergency struct {
		TicketID string `json:"ticketId,omitempty"`
	}
)

const (
	// NoncomplianceReasonNoProductionTraffic is noncompliance reason type for compliance record
	NoncomplianceReasonNoProductionTraffic = "NO_PRODUCTION_TRAFFIC"
	// NoncomplianceReasonOther is noncompliance reason type for compliance record
	NoncomplianceReasonOther = "OTHER"
	// NoncomplianceReasonEmergency is noncompliance reason type for compliance record
	NoncomplianceReasonEmergency = "EMERGENCY"
	// NoncomplianceReasonNone is noncompliance reason type for compliance record
	NoncomplianceReasonNone = "NONE"
)

// Validate validates ActivateIncludeRequest
func (i ActivateIncludeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"Version":      validation.Validate(i.Version, validation.Required),
		"Network":      validation.Validate(i.Network, validation.Required),
		"NotifyEmails": validation.Validate(i.NotifyEmails, validation.Required),
		"ComplianceRecord": validation.Validate(i.ComplianceRecord,
			validation.Required.When(i.Network == ActivationNetworkProduction).
				Error("ComplianceRecord is required for production network"),
			validation.When(i.Network == ActivationNetworkProduction, validation.By(unitTestedFieldValidationRule))),
	})
}

func unitTestedFieldValidationRule(value interface{}) error {
	switch value.(type) {
	case *ComplianceRecordNone:
		if value.(*ComplianceRecordNone).UnitTested == false {
			return errors.New("for PRODUCTION activation network and nonComplianceRecord, UnitTested value has to be set to true, otherwise API will not work correctly")
		}
	}
	return nil
}

func (c *ComplianceRecordNone) noncomplianceReasonType() string {
	return NoncomplianceReasonNone
}

// Validate validates ComplianceRecordNone
func (c *ComplianceRecordNone) Validate() error {
	return validation.Errors{
		"CustomerEmail":  validation.Validate(c.CustomerEmail, validation.Required),
		"PeerReviewedBy": validation.Validate(c.PeerReviewedBy, validation.Required),
	}.Filter()
}

// MarshalJSON is a custom marshaller for ComplianceRecordNone struct
func (c ComplianceRecordNone) MarshalJSON() ([]byte, error) {
	type ComplianceRecord ComplianceRecordNone
	v := struct {
		ComplianceRecord
		NoncomplianceReason string `json:"noncomplianceReason"`
	}{
		ComplianceRecord(c),
		c.noncomplianceReasonType(),
	}
	return json.Marshal(v)
}

func (c *ComplianceRecordOther) noncomplianceReasonType() string {
	return NoncomplianceReasonOther
}

// Validate validates ComplianceRecordOther
func (c *ComplianceRecordOther) Validate() error {
	return validation.Errors{
		"OtherNoncomplianceReason": validation.Validate(c.OtherNoncomplianceReason, validation.Required),
	}.Filter()
}

// MarshalJSON is a custom marshaller for ComplianceRecordOther struct
func (c ComplianceRecordOther) MarshalJSON() ([]byte, error) {
	type ComplianceRecord ComplianceRecordOther
	v := struct {
		ComplianceRecord
		NoncomplianceReason string `json:"noncomplianceReason"`
	}{
		ComplianceRecord(c),
		c.noncomplianceReasonType(),
	}
	return json.Marshal(v)
}

func (c *ComplianceRecordNoProductionTraffic) noncomplianceReasonType() string {
	return NoncomplianceReasonNoProductionTraffic
}

// MarshalJSON is a custom marshaller for ComplianceRecordNoProductionTraffic struct
func (c ComplianceRecordNoProductionTraffic) MarshalJSON() ([]byte, error) {
	type ComplianceRecord ComplianceRecordNoProductionTraffic
	v := struct {
		ComplianceRecord
		NoncomplianceReason string `json:"noncomplianceReason"`
	}{
		ComplianceRecord(c),
		c.noncomplianceReasonType(),
	}
	return json.Marshal(v)
}

func (c *ComplianceRecordEmergency) noncomplianceReasonType() string {
	return NoncomplianceReasonEmergency
}

// MarshalJSON is a custom marshaller for ComplianceRecordEmergency struct
func (c ComplianceRecordEmergency) MarshalJSON() ([]byte, error) {
	type ComplianceRecord ComplianceRecordEmergency
	v := struct {
		ComplianceRecord
		NoncomplianceReason string `json:"noncomplianceReason"`
	}{
		ComplianceRecord(c),
		c.noncomplianceReasonType(),
	}
	return json.Marshal(v)
}

// Validate validates DeactivateIncludeRequest
func (i DeactivateIncludeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"Version":      validation.Validate(i.Version, validation.Required),
		"Network":      validation.Validate(i.Network, validation.Required),
		"NotifyEmails": validation.Validate(i.NotifyEmails, validation.Required),
	})
}

// Validate validates GetIncludeActivationRequest
func (i GetIncludeActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"ActivationID": validation.Validate(i.ActivationID, validation.Required),
	})
}

// Validate validates CancelIncludeActivationRequest
func (i CancelIncludeActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractID":   validation.Validate(i.ContractID, validation.Required),
		"GroupID":      validation.Validate(i.GroupID, validation.Required),
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"ActivationID": validation.Validate(i.ActivationID, validation.Required),
	})
}

// Validate validates ListIncludeActivationsRequest
func (i ListIncludeActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":  validation.Validate(i.IncludeID, validation.Required),
		"ContractID": validation.Validate(i.ContractID, validation.Required),
		"GroupID":    validation.Validate(i.GroupID, validation.Required),
	})
}

var (
	// ErrActivateInclude is returned in case an error occurs on ActivateInclude operation
	ErrActivateInclude = errors.New("activate include")
	// ErrDeactivateInclude is returned in case an error occurs on DeactivateInclude operation
	ErrDeactivateInclude = errors.New("deactivate include")
	// ErrCancelIncludeActivation is returned in case an error occurs on CancelIncludeActivation operation
	ErrCancelIncludeActivation = errors.New("cancel include activation")
	// ErrGetIncludeActivation is returned in case an error occurs on GetIncludeActivation operation
	ErrGetIncludeActivation = errors.New("get include activation")
	// ErrListIncludeActivations is returned in case an error occurs on ListIncludeActivations operation
	ErrListIncludeActivations = errors.New("list include activations")
)

func (p *papi) ActivateInclude(ctx context.Context, params ActivateIncludeRequest) (*ActivationIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ActivateInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateInclude, ErrStructValidation, err)
	}

	if params.IgnoreHTTPErrors == nil {
		params.IgnoreHTTPErrors = tools.BoolPtr(true)
	}

	requestBody := struct {
		ActivateIncludeRequest
		ActivationType ActivationType `json:"activationType"`
	}{
		params,
		ActivationTypeActivate,
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/activations", params.IncludeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateInclude, err)
	}

	var result ActivationIncludeResponse
	resp, err := p.Exec(req, &result, requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateInclude, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrActivateInclude, p.Error(resp))
	}

	id, err := ResponseLinkParse(result.ActivationLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateInclude, ErrInvalidResponseLink, err)
	}
	result.ActivationID = id

	return &result, nil
}

func (p *papi) DeactivateInclude(ctx context.Context, params DeactivateIncludeRequest) (*DeactivationIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeactivateInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateInclude, ErrStructValidation, err)
	}

	if params.IgnoreHTTPErrors == nil {
		params.IgnoreHTTPErrors = tools.BoolPtr(true)
	}

	requestBody := struct {
		DeactivateIncludeRequest
		ActivationType ActivationType `json:"activationType"`
	}{
		params,
		ActivationTypeDeactivate,
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/activations", params.IncludeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateInclude, err)
	}

	var result DeactivationIncludeResponse
	resp, err := p.Exec(req, &result, requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateInclude, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateInclude, p.Error(resp))
	}

	id, err := ResponseLinkParse(result.ActivationLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateInclude, ErrInvalidResponseLink, err)
	}
	result.ActivationID = id

	return &result, nil
}

func (p *papi) CancelIncludeActivation(ctx context.Context, params CancelIncludeActivationRequest) (*CancelIncludeActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CancelIncludeActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCancelIncludeActivation, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/activations/%s", params.IncludeID, params.ActivationID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCancelIncludeActivation, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCancelIncludeActivation, err)
	}

	var result CancelIncludeActivationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCancelIncludeActivation, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCancelIncludeActivation, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) GetIncludeActivation(ctx context.Context, params GetIncludeActivationRequest) (*GetIncludeActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetIncludeActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetIncludeActivation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/activations/%s", params.IncludeID, params.ActivationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetIncludeActivation, err)
	}

	var result GetIncludeActivationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetIncludeActivation, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetIncludeActivation, p.Error(resp))
	}

	if len(result.Activations.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: ActivationID: %s", ErrGetIncludeActivation, ErrNotFound, params.ActivationID)
	}
	result.Activation = result.Activations.Items[0]

	return &result, nil
}

func (p *papi) ListIncludeActivations(ctx context.Context, params ListIncludeActivationsRequest) (*ListIncludeActivationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludeActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludeActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/activations", params.IncludeID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListIncludeActivations, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludeActivations, err)
	}

	var result ListIncludeActivationsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListIncludeActivations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludeActivations, p.Error(resp))
	}

	return &result, nil
}
