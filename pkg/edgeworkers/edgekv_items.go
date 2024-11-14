package edgeworkers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListItemsRequest represents the request params used to list items
	ListItemsRequest struct {
		ItemsRequestParams
	}

	// GetItemRequest represents the request params used to get a single item
	GetItemRequest struct {
		ItemID string
		ItemsRequestParams
	}

	// UpsertItemRequest represents the request params and body used to create or update an item
	UpsertItemRequest struct {
		ItemID   string
		ItemData Item
		ItemsRequestParams
	}

	// DeleteItemRequest represents the request params used to delete an item
	DeleteItemRequest struct {
		ItemID string
		ItemsRequestParams
	}

	// ItemsRequestParams represents the params used to list items
	ItemsRequestParams struct {
		Network     ItemNetwork
		NamespaceID string
		GroupID     string
	}

	// Item represents a single item
	Item string

	// ItemNetwork represents available item network types
	ItemNetwork string

	// ListItemsResponse represents the response from list items
	ListItemsResponse []string
)

const (
	// ItemStagingNetwork is the staging network
	ItemStagingNetwork ItemNetwork = "staging"

	// ItemProductionNetwork is the staging network
	ItemProductionNetwork ItemNetwork = "production"
)

// Validate validates ItemsRequestParams
func (r ItemsRequestParams) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(ItemStagingNetwork, ItemProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, ItemStagingNetwork, ItemProductionNetwork))),
		"NamespaceID": validation.Validate(r.NamespaceID, validation.Required),
		"GroupID":     validation.Validate(r.GroupID, validation.Required),
	}.Filter()
}

// Validate validates ListItemsRequest
func (r ListItemsRequest) Validate() error {
	return validation.Errors{
		"ItemsRequestParams": validation.Validate(r.ItemsRequestParams, validation.Required),
	}.Filter()
}

// Validate validates GetItemRequest
func (r GetItemRequest) Validate() error {
	return validation.Errors{
		"ItemID":             validation.Validate(r.ItemID, validation.Required),
		"ItemsRequestParams": validation.Validate(r.ItemsRequestParams, validation.Required),
	}.Filter()
}

// Validate validates UpsertItemRequest
func (r UpsertItemRequest) Validate() error {
	return validation.Errors{
		"ItemID":             validation.Validate(r.ItemID, validation.Required),
		"ItemData":           validation.Validate(r.ItemData, validation.Required),
		"ItemsRequestParams": validation.Validate(r.ItemsRequestParams, validation.Required),
	}.Filter()
}

// Validate validates DeleteItemRequest
func (r DeleteItemRequest) Validate() error {
	return validation.Errors{
		"ItemID":             validation.Validate(r.ItemID, validation.Required),
		"ItemsRequestParams": validation.Validate(r.ItemsRequestParams, validation.Required),
	}.Filter()
}

// IsJSON returns true if the given string is a valid json
func IsJSON(str Item) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

var (
	// ErrListItems is returned in case an error occurs on ListItems operation
	ErrListItems = errors.New("list items")
	// ErrGetItem is returned in case an error occurs on GetItem operation
	ErrGetItem = errors.New("get item")
	// ErrUpsertItem is returned in case an error occurs on UpsertItem operation
	ErrUpsertItem = errors.New("create or update item")
	// ErrDeleteItem is returned in case an error occurs on DeleteItem operation
	ErrDeleteItem = errors.New("delete item")
)

func (e *edgeworkers) ListItems(ctx context.Context, params ListItemsRequest) (*ListItemsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListItems")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListItems, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s/groups/%s", params.Network, params.NamespaceID, params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListItems, err)
	}

	var result ListItemsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListItems, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListItems, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetItem(ctx context.Context, params GetItemRequest) (*Item, error) {
	logger := e.Log(ctx)
	logger.Debug("GetItem")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetItem, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s/groups/%s/items/%s", params.Network,
		params.NamespaceID, params.GroupID, params.ItemID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetItem, err)
	}

	resp, err := e.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetItem, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetItem, e.Error(resp))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data: %s", ErrGetItem, err)
	}

	result := Item(data)
	return &result, nil
}

func (e *edgeworkers) UpsertItem(ctx context.Context, params UpsertItemRequest) (*string, error) {
	logger := e.Log(ctx)
	logger.Debug("UpsertItem")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpsertItem, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s/groups/%s/items/%s", params.Network,
		params.NamespaceID, params.GroupID, params.ItemID)

	data := []byte(params.ItemData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, ioutil.NopCloser(bytes.NewBuffer(data)))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpsertItem, err)
	}
	req.ContentLength = int64(len(data))

	if !IsJSON(params.ItemData) {
		req.Header.Set("Content-Type", "text/plain")
	}

	resp, err := e.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpsertItem, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpsertItem, e.Error(resp))
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data: %s", ErrUpsertItem, err)
	}

	result := string(data)
	return &result, nil
}

func (e *edgeworkers) DeleteItem(ctx context.Context, params DeleteItemRequest) (*string, error) {
	logger := e.Log(ctx)
	logger.Debug("DeleteItem")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteItem, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s/groups/%s/items/%s", params.Network,
		params.NamespaceID, params.GroupID, params.ItemID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteItem, err)
	}

	resp, err := e.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteItem, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeleteItem, e.Error(resp))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data: %s", ErrDeleteItem, err)
	}

	var result = string(data)
	return &result, nil
}
