package siteshield

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// SiteShieldMap represents a collection of Site Shield
//
// See: SiteShieldMap.GetSiteShieldMaps()
// API Docs: // site_shield v1
//
// https://developer.akamai.com/api/cloud_security/site_shield/v1.html

type (
	// SiteShieldMap contains operations available on SiteShieldMap resource
	// See: // site_shield v1
	//
	// https://developer.akamai.com/api/cloud_security/site_shield/v1.html#getamap
	SiteShieldMap interface {
		GetSiteShieldMaps(ctx context.Context) (*GetSiteShieldMapsResponse, error)
		GetSiteShieldMap(ctx context.Context, params SiteShieldMapRequest) (*SiteShieldMapResponse, error)
		AckSiteShieldMap(ctx context.Context, params SiteShieldMapRequest) (*SiteShieldMapResponse, error)
	}

	SiteShieldMapRequest struct {
		UniqueID int
	}

	GetSiteShieldMapsResponse struct {
		SiteShieldMaps []SiteShieldMapResponse `json:"siteShieldMaps"`
	}

	SiteShieldMapResponse struct {
		Acknowledged             bool     `json:"acknowledged"`
		Contacts                 []string `json:"contacts"`
		CurrentCidrs             []string `json:"currentCidrs"`
		ProposedCidrs            []string `json:"proposedCidrs"`
		RuleName                 string   `json:"ruleName"`
		Type                     string   `json:"type"`
		Service                  string   `json:"service"`
		Shared                   bool     `json:"shared"`
		AcknowledgeRequiredBy    int64    `json:"acknowledgeRequiredBy"`
		PreviouslyAcknowledgedOn int64    `json:"previouslyAcknowledgedOn"`
		ID                       int      `json:"id,omitempty"`
		LatestTicketID           int      `json:"latestTicketId,omitempty"`
		MapAlias                 string   `json:"mapAlias,omitempty"`
		McmMapRuleID             int      `json:"mcmMapRuleId,omitempty"`
	}
)

// Validate validates SiteShieldMapRequest
func (v SiteShieldMapRequest) Validate() error {
	return validation.Errors{
		"UniqueID": validation.Validate(v.UniqueID, validation.Required),
	}.Filter()
}

// GetSiteShieldMaps will get a list of SiteShieldMap.
//
// API Docs: // site_shield v1
//
// https://developer.akamai.com/api/cloud_security/site_shield/v1.html#listmaps

func (s *siteshieldmap) GetSiteShieldMaps(ctx context.Context) (*GetSiteShieldMapsResponse, error) {
	logger := s.Log(ctx)
	logger.Debug("GetSiteShieldMaps")

	var rval GetSiteShieldMapsResponse

	uri := "/siteshield/v1/maps"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getSiteShieldMaps request: %s", err.Error())
	}

	resp, err := s.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getsiteshieldmaps request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.Error(resp)
	}

	return &rval, nil

}

// GetSiteShieldMap will get a SiteShieldMap by unique ID.
//
// API Docs: // site_shield v1
//
// https://developer.akamai.com/api/cloud_security/site_shield/v1.html#getamap

func (s *siteshieldmap) GetSiteShieldMap(ctx context.Context, params SiteShieldMapRequest) (*SiteShieldMapResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := s.Log(ctx)
	logger.Debug("GetSiteShieldMap")

	var rval SiteShieldMapResponse

	uri := fmt.Sprintf("/siteshield/v1/maps/%d", params.UniqueID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getSiteShieldMap request: %s", err.Error())
	}

	resp, err := s.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getSiteShieldMap request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.Error(resp)
	}

	return &rval, nil
}

// AckSiteShieldMap will acknowledge changes to a SiteShieldMap.
//
// API Docs: // site_shield v1
//
// https://developer.akamai.com/api/cloud_security/site_shield/v1.html#acknowledgeamap

func (s *siteshieldmap) AckSiteShieldMap(ctx context.Context, params SiteShieldMapRequest) (*SiteShieldMapResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := s.Log(ctx)
	logger.Debug("AckSiteShieldMap")

	postURL := fmt.Sprintf("/siteshield/v1/maps/%d/acknowledge", params.UniqueID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create AckSiteShieldMap: %s", err.Error())
	}

	var rval SiteShieldMapResponse
	resp, err := s.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("AckSiteShieldMap request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.Error(resp)
	}

	return &rval, nil
}
