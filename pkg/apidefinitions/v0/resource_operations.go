// Package v0 provides access to the Akamai APIDefinitions V0 API
package v0

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type (
	// ResourceOperations represents the resource operations belonging to API endpoint ID

	// GetResourceOperationRequest contains parameters for GetResourceOperation method
	GetResourceOperationRequest struct {
		VersionNumber int64
		APIID         int64
	}

	// UpdateResourceOperationRequest contains parameters for UpdateResourceOperation method
	UpdateResourceOperationRequest struct {
		VersionNumber int64
		APIID         int64
		Body          ResourceOperationsRequestBody
	}

	// DeleteResourceOperationRequest contains parameters for DeleteResourceOperation method
	DeleteResourceOperationRequest struct {
		VersionNumber int64
		APIID         int64
	}

	// DeleteResourceOperationResponse contains parameters for DeleteResourceOperation method
	DeleteResourceOperationResponse struct {
		APIID         int64  `json:"apiEndpointId,omitempty"`
		VersionNumber int64  `json:"versionNumber,omitempty"`
		Status        int64  `json:"status,omitempty"`
		Detail        string `json:"detail,omitempty"`
	}

	// GetResourceOperationResponse holds response parameter for GetResourceOperation method
	GetResourceOperationResponse ResourceOperationResponse

	// ResourceOperationsRequestBody holds request body parameter for UpdateResourceOperationRequest method
	ResourceOperationsRequestBody ResourceOperationResponse

	// UpdateResourceOperationResponse holds response parameter for UpdateResourceOperation method
	UpdateResourceOperationResponse ResourceOperationResponse

	// ResourceOperationResponse holds the response for all Operation methods
	ResourceOperationResponse struct {
		ResourceOperations *orderedmap.OrderedMap[string, *orderedmap.OrderedMap[string, Operation]] `json:"operations"`
	}

	// Operation represents single resource operation
	Operation struct {
		Method                *string                                            `json:"method"`
		Purpose               *string                                            `json:"purpose"`
		MultistepGroupName    *string                                            `json:"multistepGroupName,omitempty"`
		Parameters            *orderedmap.OrderedMap[string, OperationParameter] `json:"parameters,omitempty"`
		Conditions            []ParameterPathCondition                           `json:"conditions,omitempty"`
		FailureConditions     []OperationCondition                               `json:"failureConditions,omitempty"`
		SuccessConditions     []OperationCondition                               `json:"successConditions,omitempty"`
		StepSuccessConditions []OperationCondition                               `json:"stepSuccessConditions,omitempty"`
		OriginUserIDCondition *OperationCondition                                `json:"originUserIdCondition,omitempty"`
	}

	// OperationCondition represents condition data for a resource operation
	OperationCondition struct {
		HeaderName                 *string  `json:"headerName,omitempty"`
		PositiveMatch              *bool    `json:"positiveMatch,omitempty"`
		SuppressFromClientResponse *bool    `json:"suppressFromClientResponse,omitempty"`
		Type                       *string  `json:"type,omitempty"`
		ValueCase                  *bool    `json:"valueCase,omitempty"`
		ValueWildcard              *bool    `json:"valueWildcard,omitempty"`
		Path                       *string  `json:"xPath,omitempty"`
		Values                     []string `json:"values,omitempty"`
	}

	// ParameterPathCondition represents condition for parameter path
	ParameterPathCondition struct {
		Path          []string `json:"path,omitempty"`
		Location      *string  `json:"location,omitempty"`
		PositiveMatch *bool    `json:"positiveMatch,omitempty"`
		Values        []string `json:"values,omitempty"`
	}

	// OperationParameter parameter details
	OperationParameter struct {
		Path                       []string `json:"path,omitempty"`
		Location                   *string  `json:"location,omitempty"`
		ReferenceParameterLocation *string  `json:"referenceParameterLocation,omitempty"`
		UsedForLogin               *bool    `json:"usedForLogin,omitempty"`
	}
)

const (
	// API URL
	operationsURI = "/api-definitions/v0/endpoints/%d/versions/%d/operations"
)

var (
	// ErrGetResourceOperation is returned when GetResourceOperation fails
	ErrGetResourceOperation = errors.New("get resource operations")
	// ErrUpdateResourceOperation is returned when UpdateResourceOperation fails
	ErrUpdateResourceOperation = errors.New("update resource operations")
	// ErrDeleteResourceOperation is returned when DeleteResourceOperation fails
	ErrDeleteResourceOperation = errors.New("delete resource operations")
)

// Validate validates GetResourceOperationRequest
func (r GetResourceOperationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIID":         validation.Validate(r.APIID, validation.Required),
		"VersionNumber": validation.Validate(r.VersionNumber, validation.Required),
	})
}

// Validate validates DeleteResourceOperationRequest
func (r DeleteResourceOperationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIID":         validation.Validate(r.APIID, validation.Required),
		"VersionNumber": validation.Validate(r.VersionNumber, validation.Required),
	})
}

// Validate validates UpdateResourceOperationRequest
func (u UpdateResourceOperationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIID":         validation.Validate(u.APIID, validation.Required),
		"VersionNumber": validation.Validate(u.VersionNumber, validation.Required),
		"Body":          validation.Validate(u.Body, validation.Required),
	})
}

func (a *apidefinitions) GetResourceOperation(ctx context.Context, params GetResourceOperationRequest) (*GetResourceOperationResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("GetResourceOperation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetResourceOperation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf(operationsURI, params.APIID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetResourceOperation, err)
	}

	var result GetResourceOperationResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetResourceOperation, err)
	}

	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetResourceOperation, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) UpdateResourceOperation(ctx context.Context, params UpdateResourceOperationRequest) (*UpdateResourceOperationResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("UpdateResourceOperation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateResourceOperation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf(operationsURI, params.APIID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateResourceOperation, err)
	}

	var result UpdateResourceOperationResponse
	resp, err := a.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateResourceOperation, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateResourceOperation, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) DeleteResourceOperation(ctx context.Context, params DeleteResourceOperationRequest) (*DeleteResourceOperationResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("DeleteResourceOperation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteResourceOperation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf(operationsURI, params.APIID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: Error while creating new http context: %s", ErrDeleteResourceOperation, err)
	}

	requestBody := UpdateResourceOperationResponse{
		ResourceOperations: orderedmap.New[string, *orderedmap.OrderedMap[string, Operation]](),
	}

	var result UpdateResourceOperationResponse

	resp, err := a.Exec(req, &result, requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteResourceOperation, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeleteResourceOperation, a.Error(resp))
	}

	response := DeleteResourceOperationResponse{}

	response.APIID = params.APIID
	response.VersionNumber = params.VersionNumber
	response.Status = http.StatusOK
	response.Detail = "Api resource operations for Endpoint is Deleted"

	return &response, nil
}
