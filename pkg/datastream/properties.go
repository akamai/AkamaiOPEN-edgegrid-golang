package datastream

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPropertiesRequest contains parameters necessary to send a GetProperties request
	GetPropertiesRequest struct {
		GroupId int
	}

	// GetDatasetFieldsRequest contains parameters necessary to send a GetDatasetFields request
	GetDatasetFieldsRequest struct {
		ProductID *string
	}

	// GetAppSecConfigsRequest contains parameters necessary to send a GetAppSecConfigs request
	GetAppSecConfigsRequest struct {
		GroupID    int
		ContractID string
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
	return validation.Errors{
		"GroupId": validation.Validate(r.GroupId, validation.Required),
	}.Filter()
}

var (
	// ErrGetProperties is returned when GetProperties fails
	ErrGetProperties = errors.New("list properties")
	// ErrGetDatasetFields is returned when GetDatasetFields fails
	ErrGetDatasetFields = errors.New("list data set fields")
	// ErrGetAppSecConfigs is returned when GetAppSecConfigs fails
	ErrGetAppSecConfigs = errors.New("list appsec configs")
)

// Validate performs validation on GetAppSecConfigsRequest
func (r GetAppSecConfigsRequest) Validate() error {
	return validation.Errors{
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
		"ContractID": validation.Validate(r.ContractID, validation.Required),
	}.Filter()
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

	uri, err := url.Parse("/datastream-config-api/v3/log/cdn/datasets-fields")
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %w", ErrGetDatasetFields, err)
	}

	q := uri.Query()
	if params.ProductID != nil {
		q.Add("productId", *params.ProductID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
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
		return nil, fmt.Errorf("%s: %w", ErrGetDatasetFields, d.Error(resp))
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
