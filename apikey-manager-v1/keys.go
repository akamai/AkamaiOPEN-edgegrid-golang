package apikeymanager

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

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
