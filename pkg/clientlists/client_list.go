package clientlists

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Lists interface to support creating, retrieving, updating and removing client lists.
	Lists interface {
		// GetClientLists lists all client lists accessible for an authenticated user
		//
		// See: https://techdocs.akamai.com/client-lists/reference/get-lists
		GetClientLists(ctx context.Context, params GetClientListsRequest) (*GetClientListsResponse, error)
	}

	// ClientListType represents client list type
	ClientListType string

	// GetClientListsRequest contains request parameters for GetClientLists method
	GetClientListsRequest struct {
		Type []ClientListType
		Name string
	}

	// GetClientListsResponse contains response parameters from GetClientLists method
	GetClientListsResponse struct {
		Content []ListContent
	}

	// ListContent contains list content
	ListContent struct {
		Name                       string         `json:"name"`
		Type                       ClientListType `json:"type"`
		Notes                      string         `json:"notes"`
		Tags                       []string       `json:"tags"`
		ListID                     string         `json:"listId"`
		Version                    int            `json:"version"`
		ItemsCount                 int            `json:"itemsCount"`
		CreateDate                 string         `json:"createDate"`
		CreatedBy                  string         `json:"createdBy"`
		UpdateDate                 string         `json:"updateDate"`
		UpdatedBy                  string         `json:"updatedBy"`
		ProductionActivationStatus string         `json:"productionActivationStatus"`
		StagingActivationStatus    string         `json:"stagingActivationStatus"`
		ListType                   string         `json:"listType"`
		Shared                     bool           `json:"shared"`
		ReadOnly                   bool           `json:"readOnly"`
		Deprecated                 bool           `json:"deprecated"`
	}
)

func (p *clientlists) GetClientLists(ctx context.Context, params GetClientListsRequest) (*GetClientListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientLists")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri, err := url.Parse("/client-list/v1/lists")
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL: %s", err.Error())
	}

	q := uri.Query()
	if params.Name != "" {
		q.Add("name", params.Name)
	}
	if params.Type != nil {
		for _, v := range params.Type {
			q.Add("type", string(v))
		}
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getClientLists request: %s", err.Error())
	}

	var rval GetClientListsResponse

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getClientLists request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (v GetClientListsRequest) validate() error {
	listTypes := getValidListTypesAsInterface()
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Type": validation.Validate(v.Type, validation.Each(validation.In(listTypes...).Error(
			fmt.Sprintf("Invalid 'type' value(s) provided. Valid values are: %s", listTypes)))),
	})
}

const (
	// IP for ip type list type
	IP ClientListType = "IP"
	// GEO for geo/countries list type
	GEO ClientListType = "GEO"
	// ASN for AS Number list type
	ASN ClientListType = "ASN"
	// TLSFingerprint for TLS Fingerprint list type
	TLSFingerprint ClientListType = "TLS_FINGERPRINT"
	// FileHash for file hash type list
	FileHash ClientListType = "FILE_HASH"
)

func getValidListTypesAsInterface() []interface{} {
	return []interface{}{
		IP,
		GEO,
		ASN,
		TLSFingerprint,
		FileHash,
	}
}
