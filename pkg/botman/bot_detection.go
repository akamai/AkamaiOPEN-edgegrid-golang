package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// The BotDetection interface supports retrieving bot detections
	BotDetection interface {
		// GetBotDetectionList todo: add link
		GetBotDetectionList(ctx context.Context, params GetBotDetectionListRequest) (*GetBotDetectionListResponse, error)
	}

	// GetBotDetectionListRequest is used to retrieve the akamai bot category actions for a policy.
	GetBotDetectionListRequest struct {
		DetectionName string
	}

	// GetBotDetectionListResponse is returned from a call to GetBotDetectionList.
	GetBotDetectionListResponse struct {
		Detections []map[string]interface{} `json:"detections"`
	}
)

func (b *botman) GetBotDetectionList(ctx context.Context, params GetBotDetectionListRequest) (*GetBotDetectionListResponse, error) {

	logger := b.Log(ctx)
	logger.Debug("GetBotDetectionList")

	uri := "/appsec/v1/bot-detections"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotDetectionList request: %w", err)
	}

	var result GetBotDetectionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotDetectionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetBotDetectionListResponse
	if params.DetectionName != "" {
		for _, val := range result.Detections {
			if val["detectionName"].(string) == params.DetectionName {
				filteredResult.Detections = append(filteredResult.Detections, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}
