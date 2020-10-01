package dns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

type (
	// Authoritiess contains operations available on Authorities data sources
	// See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html#getauthoritativenameserverdata
	Authorities interface {
		// GetAuthorities provides a list of structured read-only list of name serveers
		// See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html#getauthoritativenameserverdata
		GetAuthorities(context.Context, string) (*AuthorityResponse, error)
		// GetNameServerRecordList provides a list of name server records
		// See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html#getauthoritativenameserverdata
		GetNameServerRecordList(context.Context, string) ([]string, error)
		//
		NewAuthorityResponse(context.Context, string) *AuthorityResponse
	}

	AuthorityResponse struct {
		Contracts []struct {
			ContractID  string   `json:"contractId"`
			Authorities []string `json:"authorities"`
		} `json:"contracts"`
	}
)

func (p *dns) NewAuthorityResponse(ctx context.Context, contract string) *AuthorityResponse {

	logger := p.Log(ctx)
	logger.Debug("NewAuthorityResponse")

	authorities := &AuthorityResponse{}
	return authorities
}

func (p *dns) GetAuthorities(ctx context.Context, contractId string) (*AuthorityResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetAuthorities")

	if contractId == "" {
		return nil, fmt.Errorf("GetAuthorities reqs valid contractId")
	}

	getURL := fmt.Sprintf("/config-dns/v2/data/authorities?contractIds=%s", contractId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getauthorities request: %w", err)
	}

	var authNames AuthorityResponse
	resp, err := p.Exec(req, &authNames)
	if err != nil {
		return nil, fmt.Errorf("getauthorities request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	} else {
		return &authNames, nil
	}
}

func (p *dns) GetNameServerRecordList(ctx context.Context, contractId string) ([]string, error) {

	logger := p.Log(ctx)
	logger.Debug("GetNameServerRecordList")

	if contractId == "" {
		return nil, fmt.Errorf("GetAuthorities reqs valid contractId")
	}

	NSrecords, err := p.GetAuthorities(ctx, contractId)

	if err != nil {
		return nil, err
	}

	var arrLength int
	for _, c := range NSrecords.Contracts {
		arrLength = len(c.Authorities)
	}

	ns := make([]string, 0, arrLength)

	for _, r := range NSrecords.Contracts {
		for _, n := range r.Authorities {
			ns = append(ns, n)
		}
	}
	return ns, nil
}
