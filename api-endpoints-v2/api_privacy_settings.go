package apiendpoints

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type APIPrivacySettings struct {
	Resources map[int]APIPrivacyResource `json:"resources"`
	Public    bool                       `json:"public"`
}

type APIPrivacyResource struct {
	ResourceSettings
	Notes  string `json:"notes"`
	Public bool   `json:"public"`
}

func GetAPIPrivacySettings(endpointId, version int) (*APIPrivacySettings, error) {
	if version == 0 {
		versions, err := ListVersions(&ListVersionsOptions{EndpointId: endpointId})
		if err != nil {
			return nil, err
		}

		loc := len(versions.APIVersions) - 1
		v := versions.APIVersions[loc]
		version = v.VersionNumber
	}

	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/settings/api-privacy",
			endpointId,
			version,
		),
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

	rep := &APIPrivacySettings{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

func UpdateAPIPrivacySettings(endpointId, version int, settings *APIPrivacySettings) (*APIPrivacySettings, error) {
	if version == 0 {
		versions, err := ListVersions(&ListVersionsOptions{EndpointId: endpointId})
		if err != nil {
			return nil, err
		}

		loc := len(versions.APIVersions) - 1
		v := versions.APIVersions[loc]
		version = v.VersionNumber
	}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/settings/api-privacy",
			endpointId,
			version,
		),
		settings,
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

	return settings, nil
}
