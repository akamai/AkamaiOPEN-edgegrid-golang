package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ResponseAction interface supports retrieving response actions of a configuration
	ResponseAction interface {
		// GetResponseActionList https://techdocs.akamai.com/bot-manager/reference/get-response-actions
		GetResponseActionList(ctx context.Context, params GetResponseActionListRequest) (*GetResponseActionListResponse, error)
	}

	// GetResponseActionListRequest is used to retrieve the akamai bot category actions for a policy.
	GetResponseActionListRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// GetResponseActionListResponse is returned from a call to GetResponseActionList.
	GetResponseActionListResponse struct {
		ResponseActions []map[string]interface{} `json:"responseActions"`
	}
)

// Validate validates a GetAkamaiBotCategoryActionListRequest.
func (v GetResponseActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (b *botman) GetResponseActionList(ctx context.Context, params GetResponseActionListRequest) (*GetResponseActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetResponseActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetResponseActionList request: %w", err)
	}

	var result GetResponseActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetResponseActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetResponseActionListResponse
	if params.ActionID != "" {
		for _, val := range result.ResponseActions {
			if val["actionId"].(string) == params.ActionID {
				filteredResult.ResponseActions = append(filteredResult.ResponseActions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}
