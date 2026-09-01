package accountprotection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListProtectedOperationsRequest is used to retrieve list of account protector transactional endpoint for a configuration.
	ListProtectedOperationsRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string
	}

	// GetProtectedOperationByIDRequest is used to retrieve the account protector transactional endpoint by operationID.
	GetProtectedOperationByIDRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string

		// OperationID is the ID of the API operation
		OperationID string
	}

	// CreateProtectedOperationsRequest is used to create a list new account protector protected operation for a specific security configuration.
	CreateProtectedOperationsRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string

		// JsonPayload contains the request payload of the protected operations.
		JsonPayload json.RawMessage
	}

	// UpdateProtectedOperationRequest is used to update details for a account protector protected operation by operationID.
	UpdateProtectedOperationRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string

		// OperationID is the ID of the API operation
		OperationID string

		// JsonPayload contains the payload of the protected operation to be updated
		JsonPayload json.RawMessage
	}

	// RemoveProtectedOperationRequest is used to remove a specific account protector protected operation by operationID.
	RemoveProtectedOperationRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string

		// OperationID is the ID of the API operation
		OperationID string
	}

	// ListProtectedOperationsResponse is the response for GetTransactionalEndpointListRequest.
	ListProtectedOperationsResponse struct {
		// Metadata contains the metadata of the response
		Metadata Metadata `json:"metadata"`

		// Operations contains the list of protected operations
		Operations []map[string]interface{} `json:"operations"`
	}

	// Metadata represents the metadata for the response.
	Metadata struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64 `json:"configId"`

		// Version is the version of the security configuration.
		ConfigVersion int64 `json:"configVersion"`

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string `json:"securityPolicyId"`
	}
)

// Validate a ListProtectedOperationsRequest.
func (v ListProtectedOperationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate a GetProtectedOperationByIDRequest.
func (v GetProtectedOperationByIDRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"OperationID":      validation.Validate(v.OperationID, validation.Required),
	}.Filter()
}

// Validate a CreateProtectedOperationsRequest.
func (v CreateProtectedOperationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate a UpdateProtectedOperationRequest.
func (v UpdateProtectedOperationRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"OperationID":      validation.Validate(v.OperationID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate a RemoveProtectedOperationRequest.
func (v RemoveProtectedOperationRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"OperationID":      validation.Validate(v.OperationID, validation.Required),
	}.Filter()
}

func (ap *accountProtection) GetProtectedOperationByID(ctx context.Context, params GetProtectedOperationByIDRequest) (*ListProtectedOperationsResponse, error) {
	logger := ap.Log(ctx)
	logger.Debug("GetProtectedOperationByID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/account-protection/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.OperationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetProtectedOperationByIDRequest: %w", err)
	}

	var result ListProtectedOperationsResponse
	resp, err := ap.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetProtectedOperationByIDRequest request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return &result, nil
}

func (ap *accountProtection) ListProtectedOperations(ctx context.Context, params ListProtectedOperationsRequest) (*ListProtectedOperationsResponse, error) {
	logger := ap.Log(ctx)
	logger.Debug("ListProtectedOperations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/account-protection",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListProtectedOperationsRequest: %w", err)
	}

	var result ListProtectedOperationsResponse
	resp, err := ap.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListTransactionalEndpointsRequest request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}
	return &result, nil
}

func (ap *accountProtection) CreateProtectedOperations(ctx context.Context, params CreateProtectedOperationsRequest) (*ListProtectedOperationsResponse, error) {
	logger := ap.Log(ctx)
	logger.Debug("CreateProtectedOperations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/account-protection",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateProtectedOperationsRequest: %w", err)
	}

	var result ListProtectedOperationsResponse
	resp, err := ap.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateProtectedOperationsRequest failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, ap.Error(resp)
	}

	return &result, nil
}

func (ap *accountProtection) UpdateProtectedOperation(ctx context.Context, params UpdateProtectedOperationRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("UpdateProtectedOperation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/account-protection/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.OperationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateProtectedOperationRequest: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateProtectedOperationRequest request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}

func (ap *accountProtection) RemoveProtectedOperation(ctx context.Context, params RemoveProtectedOperationRequest) error {
	logger := ap.Log(ctx)
	logger.Debug("RemoveProtectedOperation")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/account-protection/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.OperationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveProtectedOperationRequest: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveProtectedOperationRequest request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return ap.Error(resp)
	}

	return nil
}
