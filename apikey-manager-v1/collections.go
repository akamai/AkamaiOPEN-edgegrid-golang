package apikeymanager

import (
	"fmt"

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

func GetCollection(collectionId int) (*Collection, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d", collectionId),
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

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

func CollectionAclAllow(collectionId int, acl []string) (*Collection, error) {
	collection, err := GetCollection(collectionId)
	if err != nil {
		return collection, err
	}

	acl = append(acl, collection.GrantedACL...)

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d/acl", collectionId),
		acl,
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

func CollectionAclDeny(collectionId int, acl []string) (*Collection, error) {
	collection, err := GetCollection(collectionId)
	if err != nil {
		return collection, err
	}

	for cIndex, currentAcl := range collection.GrantedACL {
		for _, newAcl := range acl {
			if newAcl == currentAcl {
				collection.GrantedACL = append(
					collection.GrantedACL[:cIndex],
					collection.GrantedACL[cIndex+1:]...,
				)
			}
		}
	}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d/acl", collectionId),
		collection.GrantedACL,
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

type Key struct {
	Id                  int      `json:"id,omitempty"`
	Value               string   `json:"value,omitempty"`
	Label               string   `json:"label,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	CollectionName      string   `json:"collectionName,omitempty"`
	CollectionId        int      `json:"collectionId,omitempty"`
	Description         string   `json:"description,omitempty"`
	Revoked             bool     `json:"revoked,omitempty"`
	Dirty               bool     `json:"dirty,omitempty"`
	CreatedAt           string   `json:"createdAt,omitempty"`
	RevokedAt           string   `json:"revokedAt,omitempty"`
	TerminationAt       string   `json:"terminationAt,omitempty"`
	QuotaUsage          int      `json:"quotaUsage,omitempty"`
	QuotaUsageTimestamp string   `json:"quotaUsageTimestamp,omitempty"`
	QuotaUpdateState    string   `json:"quotaUpdateState,omitempty"`
}

type CreateKey struct {
	Value        string   `json:"value,omitempty"`
	Label        string   `json:"label,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	CollectionId int      `json:"collectionId,omitempty"`
	Description  string   `json:"description,omitempty"`
	Mode         string   `json:"mode,omitempty"`
}

func CollectionAddKey(collectionId int, name, value string) (*Key, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/apikey-manager-api/v1/keys",
		&CreateKey{
			Label:        name,
			Value:        value,
			CollectionId: collectionId,
			Mode:         "CREATE_ONE",
		},
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

	rep := &Key{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}
