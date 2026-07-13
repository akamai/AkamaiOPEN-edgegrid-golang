package datastream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomHTPPSValidation(t *testing.T) {
	baseConnectorStruct := CustomHTTPSConnector{
		DisplayName:        "Test Connector",
		AuthenticationType: AuthenticationTypeNone,
		Endpoint:           "https://example.com",
	}
	baseConnectorStruct.SetDestinationType()

	tests := map[string]struct {
		connectorBuilder      func() CustomHTTPSConnector
		expectValidationError bool
	}{
		"AuthenticationType not in specified set": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.AuthenticationType = "NOTEXISTING"
				return c
			},
			expectValidationError: true,
		},
		"UserName required for auth type BASIC": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.AuthenticationType = AuthenticationTypeBasic
				c.Password = "password"
				return c
			},
			expectValidationError: true,
		},
		"Password required for auth type BASIC": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.AuthenticationType = AuthenticationTypeBasic
				c.UserName = "username"
				return c
			},
			expectValidationError: true,
		},
		"UserName and Password not required for type NONE": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.AuthenticationType = AuthenticationTypeNone
				return c
			},
			expectValidationError: false,
		},
		"UserName and Password required for type BASIC": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.AuthenticationType = AuthenticationTypeBasic
				return c
			},
			expectValidationError: true,
		},
		"CustomHeaderName specified without CustomHeaderValue": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.CustomHeaderName = "Custom_Name"
				return c
			},
			expectValidationError: true,
		},
		"CustomHeaderValue specified without CustomHeaderName": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.CustomHeaderValue = "Custom header value"
				return c
			},
			expectValidationError: true,
		},
		"CustomHeaderValue and CustomHeaderName both specified": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.CustomHeaderName = "Custom_Name"
				c.CustomHeaderValue = "Custom header value"
				return c
			},
			expectValidationError: false,
		},
		"CustomHeaderName contains forbidden characters": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.CustomHeaderName = "azAZ09_-!?>"
				c.CustomHeaderValue = "Custom header value"
				return c
			},
			expectValidationError: true,
		},
		"CustomHeaderName contains only allowed characters": {
			connectorBuilder: func() CustomHTTPSConnector {
				var c = baseConnectorStruct
				c.CustomHeaderName = "azAZ09_-"
				c.CustomHeaderValue = "Custom header value"
				return c
			},
			expectValidationError: false,
		},
		"CustomHeaderValue and CustomHeaderName are optional": {
			connectorBuilder: func() CustomHTTPSConnector {
				return baseConnectorStruct
			},
			expectValidationError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			connectorStructure := test.connectorBuilder()
			err := connectorStructure.Validate()
			fmt.Println(err)
			assert.True(t, (err != nil) == test.expectValidationError, err)
		})
	}
}

func TestNetStorageValidation(t *testing.T) {
	t.Parallel()
	baseConnector := NetStorageConnector{
		DisplayName:     "Test Connector",
		UserName:        "testUser",
		Endpoint:        "test.example.akamaihd.net",
		SecretAccessKey: "testSecretKey",
		Bucket:          "123456",
	}
	baseConnector.SetDestinationType()

	tests := map[string]struct {
		connectorBuilder      func() NetStorageConnector
		expectValidationError bool
	}{
		"valid connector": {
			connectorBuilder: func() NetStorageConnector {
				return baseConnector
			},
			expectValidationError: false,
		},
		"missing UserName": {
			connectorBuilder: func() NetStorageConnector {
				c := baseConnector
				c.UserName = ""
				return c
			},
			expectValidationError: true,
		},
		"missing Endpoint": {
			connectorBuilder: func() NetStorageConnector {
				c := baseConnector
				c.Endpoint = ""
				return c
			},
			expectValidationError: true,
		},
		"missing DisplayName": {
			connectorBuilder: func() NetStorageConnector {
				c := baseConnector
				c.DisplayName = ""
				return c
			},
			expectValidationError: true,
		},
		"missing SecretAccessKey": {
			connectorBuilder: func() NetStorageConnector {
				c := baseConnector
				c.SecretAccessKey = ""
				return c
			},
			expectValidationError: true,
		},
		"missing Bucket (CP code)": {
			connectorBuilder: func() NetStorageConnector {
				c := baseConnector
				c.Bucket = ""
				return c
			},
			expectValidationError: true,
		},
		"Path is optional": {
			connectorBuilder: func() NetStorageConnector {
				c := baseConnector
				c.Path = ""
				return c
			},
			expectValidationError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := test.connectorBuilder()
			err := c.Validate()
			assert.True(t, (err != nil) == test.expectValidationError, err)
		})
	}
}
