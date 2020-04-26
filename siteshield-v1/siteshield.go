package siteshield

import (
	"fmt"
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type SiteShieldMapResponse struct {
	AcknowledgeRequiredBy int64    `json:"acknowledgeRequiredBy"`
	Acknowledged          bool     `json:"acknowledged"`
	AcknowledgedBy        string   `json:"acknowledgedBy"`
	AcknowledgedOn        int64    `json:"acknowledgedOn"`
	Contacts              []string `json:"contacts"`
	CurrentCidrs          []string `json:"currentCidrs"`
	ID                    int      `json:"id"`
	LatestTicketID        int      `json:"latestTicketId"`
	MapAlias              string   `json:"mapAlias"`
	McmMapRuleID          int      `json:"mcmMapRuleId"`
	ProposedCidrs         []string `json:"proposedCidrs"`
	RuleName              string   `json:"ruleName"`
	Service               string   `json:"service"`
	Shared                bool     `json:"shared"`
	SureRouteName         string   `json:"sureRouteName"`
	Type                  string   `json:"type"`
}

func GetMap(id string) (*SiteShieldMapResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/siteshield/v1/maps/%s", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response SiteShieldMapResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func Acknowledge(id string) (*SiteShieldMapResponse, error) {
	req, err := client.NewRequest(
		Config,
		"POST",
		fmt.Sprintf("/siteshield/v1/maps/%s/acknowledge", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response SiteShieldMapResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
