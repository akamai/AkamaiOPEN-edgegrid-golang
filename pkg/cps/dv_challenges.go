package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type (
	// DVChallenges is a CPS DV challenges API interface
	DVChallenges interface {
		// GetChangeLetsEncryptChallenges gets detailed information about Domain Validation challenges
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeLetsEncryptChallenges(context.Context, GetChangeRequest) (*DVArray, error)

		// AcknowledgeDVChallenges sends acknowledgement request to CPS informing that the validation is completed
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		AcknowledgeDVChallenges(context.Context, AcknowledgementRequest) error
	}

	// DVArray is an array of DV objects
	DVArray struct {
		DV []DV `json:"dv"`
	}

	// DV is a Domain Validation entity
	DV struct {
		Challenges         []Challenge `json:"challenges"`
		Domain             string      `json:"domain"`
		Error              string      `json:"error"`
		Expires            string      `json:"expires"`
		RequestTimestamp   string      `json:"requestTimestamp"`
		Status             string      `json:"status"`
		ValidatedTimestamp string      `json:"validatedTimestamp"`
		ValidationStatus   string      `json:"validationStatus"`
	}

	// Challenge contains domain information of a specific domain to be validated
	Challenge struct {
		Error             string             `json:"error"`
		FullPath          string             `json:"fullPath"`
		RedirectFullPath  string             `json:"redirectFullPath"`
		ResponseBody      string             `json:"responseBody"`
		Status            string             `json:"status"`
		Token             string             `json:"token"`
		Type              string             `json:"type"`
		ValidationRecords []ValidationRecord `json:"validationRecords"`
	}

	// ValidationRecord represents validation attempt
	ValidationRecord struct {
		Authorities []string `json:"authorities"`
		Hostname    string   `json:"hostname"`
		Port        string   `json:"port"`
		ResolvedIP  []string `json:"resolvedIp"`
		TriedIP     string   `json:"triedIp"`
		URL         string   `json:"url"`
		UsedIP      string   `json:"usedIp"`
	}
)

var (
	// ErrGetChangeLetsEncryptChallenges is returned when GetChangeLetsEncryptChallenges fails
	ErrGetChangeLetsEncryptChallenges = errors.New("fetching change for lets-encrypt-challenges")
	// ErrAcknowledgeLetsEncryptChallenges when AcknowledgeDVChallenges fails
	ErrAcknowledgeLetsEncryptChallenges = errors.New("acknowledging lets-encrypt-challenges")
)

func (c *cps) GetChangeLetsEncryptChallenges(ctx context.Context, params GetChangeRequest) (*DVArray, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeLetsEncryptChallenges, ErrStructValidation, err)
	}

	var rval DVArray

	logger := c.Log(ctx)
	logger.Debug("GetChangeLetsEncryptChallenges")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d/input/info/lets-encrypt-challenges",
		params.EnrollmentID,
		params.ChangeID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetChangeLetsEncryptChallenges, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeLetsEncryptChallenges, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.dv-challenges.v2+json")

	resp, err := c.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeLetsEncryptChallenges, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeLetsEncryptChallenges, c.Error(resp))
	}

	return &rval, nil
}

func (c *cps) AcknowledgeDVChallenges(ctx context.Context, params AcknowledgementRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrAcknowledgeLetsEncryptChallenges, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("AcknowledgeDVVhallenges")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d/input/update/lets-encrypt-challenges-completed",
		params.EnrollmentID, params.ChangeID))
	if err != nil {
		return fmt.Errorf("%w: parsing URL: %s", ErrAcknowledgeLetsEncryptChallenges, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrAcknowledgeLetsEncryptChallenges, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.acknowledgement.v1+json; charset=utf-8")

	resp, err := c.Exec(req, nil, params.Acknowledgement)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrAcknowledgeLetsEncryptChallenges, err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", ErrAcknowledgeLetsEncryptChallenges, c.Error(resp))
	}

	return nil
}
