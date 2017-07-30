package papi

import "errors"

// Error constants
const (
	ErrInvalidPath = iota
	ErrCriteriaNotFound
	ErrBehaviorNotFound
	ErrRuleNotFound
	ErrInvalidRules
)

var (
	ErrorMap = map[int]error{
		ErrInvalidPath:      errors.New("Invalid Path"),
		ErrCriteriaNotFound: errors.New("Criteria not found"),
		ErrBehaviorNotFound: errors.New("Behavior not found"),
		ErrRuleNotFound:     errors.New("Rule not found"),
		ErrInvalidRules:     errors.New("Rule validation failed. See papi.Rules.Errors for details"),
	}
)
