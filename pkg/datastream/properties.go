package datastream

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Properties is an interface for listing various DS API properties
	Properties interface {
		// GetProperties returns properties that are active on the production and staging network for a specific product type that are available within a group
		//
		// See: https://developer.akamai.com/api/core_features/datastream2_config/v1.html#getproperties
		GetProperties(context.Context, GetPropertiesRequest) ([]Property, error)

		// GetPropertiesByGroup returns properties that are active on the production and staging network and available within a specific group
		//
		// See: https://developer.akamai.com/api/core_features/datastream2_config/v1.html#getpropertiesbygroup
		GetPropertiesByGroup(context.Context, GetPropertiesByGroupRequest) ([]Property, error)

		// GetDatasetFields returns groups of data set fields available in the template.
		//
		// See: https://developer.akamai.com/api/core_features/datastream2_config/v1.html#gettemplatename
		GetDatasetFields(context.Context, GetDatasetFieldsRequest) ([]DataSets, error)
	}

	// GetPropertiesRequest contains parameters necessary to send a GetProperties request
	GetPropertiesRequest struct {
		GroupId   int
		ProductId string
	}

	// GetPropertiesByGroupRequest contains parameters necessary to send a GetPropertiesByGroup request
	GetPropertiesByGroupRequest struct {
		GroupId int
	}

	// GetDatasetFieldsRequest contains parameters necessary to send a GetDatasetFields request
	GetDatasetFieldsRequest struct {
		TemplateName TemplateName
	}
)

// Validate performs validation on GetPropertiesRequest
func (r GetPropertiesRequest) Validate() error {
	return validation.Errors{
		"GroupId":   validation.Validate(r.GroupId, validation.Required),
		"ProductId": validation.Validate(r.ProductId, validation.Required),
	}.Filter()
}

// Validate performs validation on GetPropertiesRequest
func (r GetPropertiesByGroupRequest) Validate() error {
	return validation.Errors{
		"GroupId": validation.Validate(r.GroupId, validation.Required),
	}.Filter()
}

// Validate performs validation on GetDatasetFieldsRequest
func (r GetDatasetFieldsRequest) Validate() error {
	return validation.Errors{
		"TemplateName": validation.Validate(r.TemplateName, validation.Required, validation.In(TemplateNameEdgeLogs)),
	}.Filter()
}

var (
	// ErrGetProperties is returned when GetProperties fails
	ErrGetProperties = errors.New("list properties")
	// ErrGetPropertiesByGroup is returned when GetPropertiesByGroup fails
	ErrGetPropertiesByGroup = errors.New("list properties by group")
	// ErrGetDatasetFields is returned when GetDatasetFields fails
	ErrGetDatasetFields = errors.New("list data set fields")
)

func (d *ds) GetProperties(ctx context.Context, params GetPropertiesRequest) ([]Property, error) {
	logger := d.Log(ctx)
	logger.Debug("GetProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/properties/product/%s/group/%d",
		params.ProductId, params.GroupId))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrGetProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProperties, err)
	}

	var rval []Property
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProperties, d.Error(resp))
	}

	return rval, nil
}

func (d *ds) GetPropertiesByGroup(ctx context.Context, params GetPropertiesByGroupRequest) ([]Property, error) {
	logger := d.Log(ctx)
	logger.Debug("GetPropertiesByGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertiesByGroup, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/properties/group/%d",
		params.GroupId))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrGetPropertiesByGroup, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertiesByGroup, err)
	}

	var rval []Property
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertiesByGroup, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertiesByGroup, d.Error(resp))
	}

	return rval, nil
}

func (d *ds) GetDatasetFields(ctx context.Context, params GetDatasetFieldsRequest) ([]DataSets, error) {
	logger := d.Log(ctx)
	logger.Debug("GetDatasetFields")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetDatasetFields, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/datasets/template/%s",
		params.TemplateName))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrGetDatasetFields, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetDatasetFields, err)
	}

	var rval []DataSets
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetDatasetFields, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetDatasetFields, d.Error(resp))
	}

	return rval, nil
}
