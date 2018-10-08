package apiendpoints

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
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

func (settings *APIPrivacySettings) ToTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Endpoint Private",
		"Resource Path",
		"Resource Methods",
		"Inherits from Endpoint",
		"Resource Private",
	})

	table.Append([]string{
		cast.ToString(settings.Public),
		"",
		"",
		"",
		"",
	})

	for _, resource := range settings.Resources {
		table.Append([]string{
			"",
			cast.ToString(resource.Path),
			cast.ToString(strings.Join(resource.Methods[:], ",")),
			cast.ToString(resource.InheritsFromEndpoint),
			cast.ToString(resource.Public),
		})
	}
	return table
}
