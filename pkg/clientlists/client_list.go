package clientlists

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ClientListType represents client list type
	ClientListType string

	// GetClientListsRequest contains request parameters for GetClientLists method
	GetClientListsRequest struct {
		Type               []ClientListType
		Name               string
		Search             string
		IncludeItems       bool
		IncludeDeprecated  bool
		IncludeNetworkList bool
		Page               *int
		PageSize           *int
		Sort               []string
	}

	// GetClientListsResponse contains response parameters from GetClientLists method
	GetClientListsResponse struct {
		Content []ClientList
	}

	// ClientList contains list content and items
	ClientList struct {
		ListContent
		Items []ListItemContent
	}

	// ListContent contains list content
	ListContent struct {
		Name                       string         `json:"name"`
		Type                       ClientListType `json:"type"`
		Notes                      string         `json:"notes"`
		Tags                       []string       `json:"tags"`
		ListID                     string         `json:"listId"`
		Version                    int64          `json:"version"`
		ItemsCount                 int64          `json:"itemsCount"`
		CreateDate                 string         `json:"createDate"`
		CreatedBy                  string         `json:"createdBy"`
		UpdateDate                 string         `json:"updateDate"`
		UpdatedBy                  string         `json:"updatedBy"`
		ProductionActivationStatus string         `json:"productionActivationStatus"`
		StagingActivationStatus    string         `json:"stagingActivationStatus"`
		ProductionActiveVersion    int64          `json:"productionActiveVersion"`
		StagingActiveVersion       int64          `json:"stagingActiveVersion"`
		ListType                   string         `json:"listType"`
		Shared                     bool           `json:"shared"`
		ReadOnly                   bool           `json:"readOnly"`
		Deprecated                 bool           `json:"deprecated"`
	}

	// ListItemContent contains client list item information
	ListItemContent struct {
		Value            string         `json:"value"`
		Tags             []string       `json:"tags"`
		Description      string         `json:"description"`
		ExpirationDate   string         `json:"expirationDate"`
		CreateDate       string         `json:"createDate"`
		CreatedBy        string         `json:"createdBy"`
		CreatedVersion   int64          `json:"createdVersion"`
		ProductionStatus string         `json:"productionStatus"`
		StagingStatus    string         `json:"stagingStatus"`
		Type             ClientListType `json:"type"`
		UpdateDate       string         `json:"updateDate"`
		UpdatedBy        string         `json:"updatedBy"`
	}

	// ListItemPayload contains item's editable fields to use as update/create/delete payload
	ListItemPayload struct {
		Value          string   `json:"value"`
		Tags           []string `json:"tags"`
		Description    string   `json:"description"`
		ExpirationDate string   `json:"expirationDate"`
	}

	// GetClientListRequest contains request params for GetClientList method
	GetClientListRequest struct {
		ListID       string
		IncludeItems bool
	}

	// GetClientListResponse contains response from GetClientList method
	GetClientListResponse struct {
		ListContent
		ContractID string            `json:"contractId"`
		GroupID    int64             `json:"groupId"`
		GroupName  string            `json:"groupName"`
		Items      []ListItemContent `json:"items"`
	}

	// CreateClientListRequest contains request params for CreateClientList method
	CreateClientListRequest struct {
		ContractID string            `json:"contractId"`
		GroupID    int64             `json:"groupId"`
		Name       string            `json:"name"`
		Type       ClientListType    `json:"type"`
		Notes      string            `json:"notes"`
		Tags       []string          `json:"tags"`
		Items      []ListItemPayload `json:"items"`
	}

	// CreateClientListResponse contains response from CreateClientList method
	CreateClientListResponse GetClientListResponse

	// UpdateClientListRequest contains request params for UpdateClientList method
	UpdateClientListRequest struct {
		UpdateClientList
		ListID string
	}

	// UpdateClientList contains the body of client list update request
	UpdateClientList struct {
		Name  string   `json:"name"`
		Notes string   `json:"notes"`
		Tags  []string `json:"tags"`
	}

	// UpdateClientListResponse contains response from UpdateClientList method
	UpdateClientListResponse struct {
		ListContent
		ContractID string `json:"contractId"`
		GroupName  string `json:"groupName"`
		GroupID    int64  `json:"groupId"`
	}

	// UpdateClientListItemsRequest contains request params for UpdateClientListItems method
	UpdateClientListItemsRequest struct {
		UpdateClientListItems
		ListID string
	}

	// UpdateClientListItems contains the body of client list items update request
	UpdateClientListItems struct {
		Append []ListItemPayload `json:"append"`
		Update []ListItemPayload `json:"update"`
		Delete []ListItemPayload `json:"delete"`
	}

	// UpdateClientListItemsResponse contains response from UpdateClientListItems method
	UpdateClientListItemsResponse struct {
		Appended []ListItemContent `json:"appended"`
		Updated  []ListItemContent `json:"updated"`
		Deleted  []ListItemContent `json:"deleted"`
	}

	// DeleteClientListRequest contains request params for DeleteClientList method
	DeleteClientListRequest struct {
		ListID string
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
	if params.Search != "" {
		q.Add("search", params.Search)
	}
	if params.IncludeItems {
		q.Add("includeItems", strconv.FormatBool(params.IncludeItems))
	}
	if params.IncludeDeprecated {
		q.Add("includeDeprecated", strconv.FormatBool(params.IncludeDeprecated))
	}
	if params.IncludeNetworkList {
		q.Add("includeNetworkList", strconv.FormatBool(params.IncludeNetworkList))
	}
	if params.Page != nil {
		q.Add("page", fmt.Sprintf("%d", *params.Page))
	}
	if params.PageSize != nil {
		q.Add("pageSize", fmt.Sprintf("%d", *params.PageSize))
	}
	if params.Sort != nil {
		for _, v := range params.Sort {
			q.Add("sort", string(v))
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) GetClientList(ctx context.Context, params GetClientListRequest) (*GetClientListResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientList")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/client-list/v1/lists/%s", params.ListID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	q := uri.Query()
	if params.IncludeItems {
		q.Add("includeItems", strconv.FormatBool(params.IncludeItems))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getClientList request: %s", err.Error())
	}

	var rval GetClientListResponse
	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getClientList request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) UpdateClientList(ctx context.Context, params UpdateClientListRequest) (*UpdateClientListResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientList")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/client-list/v1/lists/%s", params.ListID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create updateClientList request: %s", err.Error())
	}

	var rval UpdateClientListResponse
	resp, err := p.Exec(req, &rval, &params.UpdateClientList)
	if err != nil {
		return nil, fmt.Errorf("updateClientList request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) UpdateClientListItems(ctx context.Context, params UpdateClientListItemsRequest) (*UpdateClientListItemsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientListItems")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/client-list/v1/lists/%s/items", params.ListID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateClientListItems request: %s", err.Error())
	}

	var rval UpdateClientListItemsResponse
	resp, err := p.Exec(req, &rval, &params.UpdateClientListItems)
	if err != nil {
		return nil, fmt.Errorf("UpdateClientListItems request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) CreateClientList(ctx context.Context, params CreateClientListRequest) (*CreateClientListResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateClientList")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/client-list/v1/lists", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create createClientList request: %s", err.Error())
	}

	var rval CreateClientListResponse
	resp, err := p.Exec(req, &rval, &params)
	if err != nil {
		return nil, fmt.Errorf("createClientList request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) DeleteClientList(ctx context.Context, params DeleteClientListRequest) error {
	logger := p.Log(ctx)
	logger.Debug("DeleteClientList")

	if err := params.validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("%s/%s", "/client-list/v1/lists", params.ListID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create deleteClientList request: %s", err.Error())
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("deleteClientList request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

func (v GetClientListRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID": validation.Validate(v.ListID, validation.Required),
	})
}

func (v UpdateClientListRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID": validation.Validate(v.ListID, validation.Required),
		"Name":   validation.Validate(v.ListID, validation.Required),
	})
}

func (v UpdateClientListItemsRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID": validation.Validate(v.ListID, validation.Required),
	})
}

func (v CreateClientListRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Name": validation.Validate(v.Name, validation.Required),
		"Type": validation.Validate(v.Type, validation.Required),
	})
}

func (v DeleteClientListRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID": validation.Validate(v.ListID, validation.Required),
	})
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
