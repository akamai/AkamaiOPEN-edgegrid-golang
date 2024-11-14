package datastream

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
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
)

func (d *ds) GetProperties(ctx context.Context, params GetPropertiesRequest) (*PropertiesDetails, error) {
	logger := d.Log(ctx)
	logger.Debug("GetProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v2/log/groups/%d/properties",
		params.GroupId))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrGetProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProperties, err)
	}

	var rval PropertiesDetails
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProperties, err)
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

	uri, err := url.Parse("/datastream-config-api/v2/log/datasets-fields")
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrGetDatasetFields, err)
	}

	q := uri.Query()
	if params.ProductID != nil {
		q.Add("productId", fmt.Sprintf("%s", *params.ProductID))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetDatasetFields, err)
	}

	var rval DataSets
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetDatasetFields, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetDatasetFields, d.Error(resp))
	}

	return &rval, nil
}
