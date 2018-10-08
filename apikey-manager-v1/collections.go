package apikeymanager

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
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
	Quota       Quota    `json:"quota,omitempty"`
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

func (c *Collections) ToTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID",
		"Name",
		"Granted ACL",
		"Quota Enabled",
		"Quota Value",
		"Quota Interval",
	})

	for _, collection := range *c {
		table.Append([]string{
			cast.ToString(collection.Id),
			cast.ToString(collection.Name),
			cast.ToString(strings.Join(collection.GrantedACL[:], ",")),
			cast.ToString(collection.Quota.Enabled),
			cast.ToString(collection.Quota.Value),
			cast.ToString(collection.Quota.Interval),
		})
	}
	return table
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

func (collection *Collection) ToTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID",
		"Name",
		"Granted ACL",
		"Quota Enabled",
		"Quota Value",
		"Quota Interval",
	})

	table.Append([]string{
		cast.ToString(collection.Id),
		cast.ToString(collection.Name),
		cast.ToString(strings.Join(collection.GrantedACL[:], ",")),
		cast.ToString(collection.Quota.Enabled),
		cast.ToString(collection.Quota.Value),
		cast.ToString(collection.Quota.Interval),
	})

	return table
}

func GetCollectionMulti(collection string) (*Collections, error) {
	id, err := strconv.Atoi(collection)
	if err == nil {
		req, err := client.NewJSONRequest(
			Config,
			"GET",
			fmt.Sprintf("/apikey-manager-api/v1/collections/%d", id),
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

		return &Collections{*rep}, nil
	} else {
		collections, err := ListCollections()
		if err != nil {
			return nil, err
		}

		if collection[len(collection):] == "*" {
			collection = collection[:len(collection)-1] + ".*"
		}

		ret := Collections{}
		for _, r := range *collections {
			matched, err := regexp.MatchString(collection, r.Name)
			if err != nil {
				return nil, err
			}

			if matched {
				ret = append(ret, r)
			}
		}

		if len(ret) > 0 {
			return &ret, nil
		}

		return nil, errors.New("Collection not found.")
	}
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

	for i := len(collection.GrantedACL) - 1; i >= 0; i-- {
		for _, newAcl := range acl {
			if newAcl == collection.GrantedACL[i] {
				collection.GrantedACL = append(
					collection.GrantedACL[:i],
					collection.GrantedACL[i+1:]...,
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

type Quota struct {
	Enabled  bool   `json:"enabled"`
	Value    int    `json:"value"`
	Interval string `json:"interval"`
	Headers  struct {
		DenyLimitHeaderShown      bool `json:"denyLimitHeaderShown"`
		DenyRemainingHeaderShown  bool `json:"denyRemainingHeaderShown"`
		DenyNextHeaderShown       bool `json:"denyNextHeaderShown"`
		AllowLimitHeaderShown     bool `json:"allowLimitHeaderShown"`
		AllowRemainingHeaderShown bool `json:"allowRemainingHeaderShown"`
		AllowResetHeaderShown     bool `json:"allowResetHeaderShown"`
	} `json:"headers,omitempty"`
}

func CollectionSetQuota(collection string, limit int, interval string) (*Collections, error) {
	collections, err := GetCollectionMulti(collection)
	if err != nil {
		return nil, err
	}

	rcollections := Collections{}
	for _, collection := range *collections {
		collection.Quota.Value = limit
		collection.Quota.Interval = interval
		collection.Quota.Enabled = true
		collection.Quota.Headers.DenyLimitHeaderShown = true
		collection.Quota.Headers.DenyRemainingHeaderShown = true
		collection.Quota.Headers.DenyNextHeaderShown = true
		collection.Quota.Headers.AllowLimitHeaderShown = true
		collection.Quota.Headers.AllowRemainingHeaderShown = true
		collection.Quota.Headers.AllowResetHeaderShown = true
		req, err := client.NewJSONRequest(
			Config,
			"PUT",
			fmt.Sprintf("/apikey-manager-api/v1/collections/%d/quota", collection.Id),
			collection.Quota,
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

		rcollections = append(rcollections, *rep)
	}

	return &rcollections, nil
}
