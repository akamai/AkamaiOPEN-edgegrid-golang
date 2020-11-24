package iam

import "context"

type (
	// Support is a list of iam supported object methods
	Support interface {
		SupportedCountries(context.Context) ([]string, error)
		SupportedContactTypes(context.Context) ([]string, error)
		SupportedLanguages(context.Context) ([]string, error)
		ListProducts(context.Context) ([]string, error)
		ListTimeoutPolicies(context.Context) ([]TimeoutPolicy, error)
		ListStates(context.Context, ListStatesRequest) ([]string, error)
	}

	// TimeoutPolicy specifies session timeout policy options that can be assigned to each user
	TimeoutPolicy struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}

	// ListStatesRequest specifies the country for the requested states
	ListStatesRequest struct {
		Country string `json:"country"`
	}
)
