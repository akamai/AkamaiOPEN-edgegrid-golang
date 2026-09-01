package datastream

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPropertiesRequest contains parameters necessary to send a GetProperties request
	GetPropertiesRequest struct {
		GroupId int
	}

	// GetDatasetFieldsRequest contains parameters necessary to send a GetDatasetFields request
	GetDatasetFieldsRequest struct {
		//LogType identifies the log type for which to retrieve dataset fields. Valid values are "CDN" and "ANSWERX".
		LogType LogType

		// ProductID is an optional parameter that identifies the product for which to retrieve dataset fields for CDN logtype. If provided, only dataset fields associated with the specified product will be returned.
		ProductID string
	}

	// GetAppSecConfigsRequest contains parameters necessary to send a GetAppSecConfigs request
	GetAppSecConfigsRequest struct {
		GroupID    int
		ContractID string
	}

	// ListAnswerXServiceIDsRequest contains parameters necessary to send a ListAnswerXServiceIDs request.
	ListAnswerXServiceIDsRequest struct {
		// ContractID identifies the contract for which to list available AnswerX service IDs.
		ContractID string

		// PageSize is the number of service IDs to return in a single response page.
		PageSize int64

		// Page identifies the response page to retrieve.
		Page int64
	}

	// AppSecConfigDetails contains information about an AppSec configuration available for streaming
	AppSecConfigDetails struct {
		FileType          string `json:"fileType"`
		ID                int    `json:"id"`
		LatestVersion     int    `json:"latestVersion"`
		Name              string `json:"name"`
		ProductionVersion int    `json:"productionVersion"`
		TargetProduct     string `json:"targetProduct"`
	}

	// ListAnswerXServiceIDsResponse contains the service IDs belonging to the given contract.
	ListAnswerXServiceIDsResponse struct {
		// Metadata contains pagination information for the current response page.
		Metadata *PaginationMetadata `json:"metadata"`

		// ContractID identifies the contract for which the service IDs are listed.
		ContractID string `json:"contractId"`

		// AnswerXServiceIDs contains the AnswerX service IDs associated with the contract.
		AnswerXServiceIDs []AnswerXServiceDetail `json:"serviceSubletterIds"`
	}

	// PaginationMetadata contains pagination information returned by the API.
	PaginationMetadata struct {
		// LastPage is the index of the final available response page.
		LastPage int64 `json:"lastPage"`

		// PageSize is the maximum number of items returned in the current page.
		PageSize int64 `json:"pageSize"`

		// Page is the index of the current response page.
		Page int64 `json:"page"`

		// TotalElements is the total number of available service IDs across all pages.
		TotalElements int64 `json:"totalElements"`
	}

	// PropertiesDetails identifies the properties belong to the given group.
	PropertiesDetails struct {
		Properties []PropertyDetails `json:"properties"`
		GroupID    int               `json:"groupId"`
	}

	// PropertyDetails identifies detailed info about the properties monitored in the stream.
	PropertyDetails struct {
		Hostnames    []string `json:"hostnames"`
		ProductID    string   `json:"productId"`
		ProductName  string   `json:"productName"`
		PropertyID   int      `json:"propertyId"`
		PropertyName string   `json:"propertyName"`
		ContractID   string   `json:"contractId"`
	}
)

// Validate performs validation on GetPropertiesRequest
func (r GetPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"GroupId": validation.Validate(r.GroupId, validation.Required),
	})
}

// Validate performs validation on GetAppSecConfigsRequest
func (r GetAppSecConfigsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
		"ContractID": validation.Validate(r.ContractID, validation.Required),
	})
}

// Validate performs validation on ListAnswerXServiceIDsRequest.
func (r ListAnswerXServiceIDsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractID": validation.Validate(r.ContractID, validation.Required),
		"PageSize": validation.Validate(r.PageSize,
			validation.Min(int64(0)),
			validation.When(r.Page != 0, validation.Required.Error("page and pageSize must be provided together")),
		),
		"Page": validation.Validate(r.Page,
			validation.Min(int64(0)),
			validation.When(r.PageSize != 0, validation.Required.Error("page and pageSize must be provided together")),
		),
	})
}

// Validate performs validation on GetDatasetFieldsRequest
func (r GetDatasetFieldsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LogType":   validation.Validate(r.LogType, validation.Required, validation.In(LogTypeCDN, LogTypeAnswerX)),
		"ProductID": validation.Validate(r.ProductID, validation.When(r.LogType != LogTypeCDN, validation.Empty)),
	})
}

func (d *ds) GetProperties(ctx context.Context, params GetPropertiesRequest) (*PropertiesDetails, error) {
	logger := d.Log(ctx)
	logger.Debug("GetProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v3/log/cdn/groups/%d/properties",
		params.GroupId))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %w", ErrGetProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetProperties, err)
	}

	var rval PropertiesDetails
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProperties, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) GetDatasetFields(ctx context.Context, params GetDatasetFieldsRequest) (*DataSets, error) {
	logger := d.Log(ctx)
	logger.Debug("GetDatasetFields")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetDatasetFields, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/datastream-config-api/v3/log/%s/datasets-fields", params.LogType.ToPathValue()).
		AddQueryParamIf("productId", params.ProductID, params.ProductID != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetDatasetFields, err)
	}

	var rval DataSets
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetDatasetFields, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetDatasetFields, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) GetAppSecConfigs(ctx context.Context, params GetAppSecConfigsRequest) ([]AppSecConfigDetails, error) {
	logger := d.Log(ctx)
	logger.Debug("GetAppSecConfigs")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetAppSecConfigs, ErrStructValidation, err)
	}

	req, err := request.NewGet(
		ctx,
		"/datastream-config-api/v3/log/appsec/groups/%d/contracts/%s/configs",
		params.GroupID, params.ContractID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetAppSecConfigs, err)
	}

	var rval []AppSecConfigDetails
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetAppSecConfigs, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAppSecConfigs, d.Error(resp))
	}

	return rval, nil
}

func (d *ds) ListAnswerXServiceIDs(ctx context.Context, params ListAnswerXServiceIDsRequest) (*ListAnswerXServiceIDsResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ListAnswerXServiceIDs")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListAnswerXServiceIDs, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/datastream-config-api/v3/log/answerx/contracts/%s/answerxSSIDs", params.ContractID).
		AddQueryParamIf("pageSize", strconv.FormatInt(params.PageSize, 10), params.PageSize > 0).
		AddQueryParamIf("page", strconv.FormatInt(params.Page, 10), params.Page > 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListAnswerXServiceIDs, err)
	}

	var rval ListAnswerXServiceIDsResponse
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListAnswerXServiceIDs, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListAnswerXServiceIDs, d.Error(resp))
	}

	return &rval, nil
}
