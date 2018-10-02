package apikeymanager

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Collections []Collection

type Collection struct {
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	KeyCount    int      `json:"keyCount,omitempty"`
	Dirty       bool     `json:"dirty,omitempty"`
	ContractId  string   `json:"contractId,omitempty"`
	GroupId     int      `json:"groupId,omitempty"`
	GrantedACL  []string `json:"grantedACL,omitempty"`
	DirtyACL    []string `json:"dirtyACL,omitempty"`
	Quota       struct {
		Enabled  bool   `json:"enabled,omitempty"`
		Value    int    `json:"value,omitempty"`
		Interval string `json:"interval,omitempty"`
		Headers  struct {
			DenyLimitHeaderShown      bool `json:"denyLimitHeaderShown,omitempty"`
			DenyRemainingHeaderShown  bool `json:"denyRemainingHeaderShown,omitempty"`
			DenyNextHeaderShown       bool `json:"denyNextHeaderShown,omitempty"`
			AllowLimitHeaderShown     bool `json:"allowLimitHeaderShown,omitempty"`
			AllowRemainingHeaderShown bool `json:"allowRemainingHeaderShown,omitempty"`
			AllowResetHeaderShown     bool `json:"allowResetHeaderShown,omitempty"`
		} `json:"headers,omitempty"`
	} `json:"quota,omitempty"`
}

func ListCollections() (*Collections, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		"/apikey-manager-api/v1/collections",
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

	rep := &Collections{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type CreateCollectionOptions struct {
	ContractId  string `json:"contractId,omitempty"`
	GroupId     int    `json:"groupId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func CreateCollection(options *CreateCollectionOptions) (*Collection, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/apikey-manager-api/v1/collections",
		options,
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

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}
