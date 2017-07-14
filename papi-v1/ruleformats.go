package papi

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
)

// RuleFormats is a collection of available rule formats
type RuleFormats struct {
	client.Resource
	RuleFormats struct {
		Items []string `json:"items"`
	} `json:"ruleFormats"`
}

// NewRuleFormats creates a new RuleFormats
func NewRuleFormats() *RuleFormats {
	ruleFormats := &RuleFormats{}
	ruleFormats.Init()

	return ruleFormats
}

// GetRuleFormats populates RuleFormats
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listruleformats
// Endpoint: GET /papi/v0/rule-formats
func (ruleFormats *RuleFormats) GetRuleFormats() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/papi/v0/rule-formats",
		nil,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err := client.BodyJSON(res, ruleFormats); err != nil {
		return err
	}

	return nil
}

// GetSchema fetches the schema for a given product and rule format
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruleformatsschema
// Endpoint: /papi/v0/schemas/products/{productId}/{ruleFormat}
func (ruleFormats *RuleFormats) GetSchema(product string, ruleFormat string) (*gojsonschema.Schema, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v0/schemas/products/%s/%s",
			product,
			ruleFormat,
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

	schemaBytes, _ := ioutil.ReadAll(res.Body)
	schemaBody := string(schemaBytes)
	loader := gojsonschema.NewStringLoader(schemaBody)
	schema, err := gojsonschema.NewSchema(loader)

	return schema, err
}
