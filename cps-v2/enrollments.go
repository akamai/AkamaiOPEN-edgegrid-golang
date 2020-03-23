package cps

import (
	"encoding/json"
	"fmt"
	"errors"
	"bytes"
	"time"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Enrollments represents an enrollment
//
// API Docs: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#enrollments
type Enrollment struct {
	client.Resource
	AdminContact              *Contact              `json:"adminContact"`
	CertificateChainType      *string               `json:"certificateChainType"`
	CertificateType           CertificateType       `json:"certificateType"`
	CertificateSigningRequest *CSR                  `json:"csr"`
	ChangeManagement          bool                  `json:"changeManagement"`
	EnableMultiStacked        bool                  `json:"enableMultiStackedCertificates"`
	Location                  *string               `json:"location"`
	MaxAllowedSans            *int                  `json:"maxAllowedSanNames"`
	MaxAllowedWildcardSans    *int                  `json:"maxAllowedWildcardSanNames"`
	NetworkConfiguration      *NetworkConfiguration `json:"networkConfiguration"`
	Organization              *Organization         `json:"org"`
	PendingChanges            *[]string             `json:"pendingChanges"`
	RegistrationAuthority     RegistrationAuthority `json:"ra"`
	SignatureAuthority        *SHA                  `json:"signatureAlgorithm"`
	TechContact               *Contact              `json:"techContact"`
	ThirdParty                *ThirdParty           `json:"thirdParty"`
	ValidationType            ValidationType        `json:"validationType"`
}

type CreateEnrollmentQueryParams struct {
	ContractID      string
	DeployNotAfter  *string
	DeployNotBefore *string
}

type Acknowledgement struct {
	Acknowledgement string `json:"acknowledgement"`
}

type ListEnrollmentsQueryParams struct {
	ContractID string
}

type CreateEnrollmentResponse struct {
	Location string   `json:"enrollment"`
	Changes  []string `json:"changes"`
}

type StatusResponse struct {
	StatusInfo struct {
		Status             string `json:"status"`
		State              string `json:"state"`
		Description        string `json:"description"`
		DeploymentSchedule struct {
			NotBefore interface{} `json:"notBefore"`
			NotAfter  interface{} `json:"notAfter"`
		} `json:"deploymentSchedule"`
		Error interface{} `json:"error"`
	} `json:"statusInfo"`
	AllowedInput []struct {
		Type              string `json:"type"`
		RequiredToProceed bool   `json:"requiredToProceed"`
		Info              string `json:"info"`
		Update            string `json:"update"`
	} `json:"allowedInput"`
}

type DomainValidations struct {
	Dv []struct {
		Domain             string      `json:"domain"`
		Status             string      `json:"status"`
		Error              interface{} `json:"error"`
		ValidationStatus   string      `json:"validationStatus"`
		RequestTimestamp   time.Time   `json:"requestTimestamp"`
		ValidatedTimestamp interface{} `json:"validatedTimestamp"`
		Expires            time.Time   `json:"expires"`
		Challenges         []struct {
			Type              string        `json:"type"`
			Status            string        `json:"status"`
			Error             interface{}   `json:"error"`
			Token             string        `json:"token"`
			ResponseBody      string        `json:"responseBody"`
			FullPath          string        `json:"fullPath"`
			RedirectFullPath  string        `json:"redirectFullPath"`
			ValidationRecords []interface{} `json:"validationRecords"`
		} `json:"challenges"`
	} `json:"dv"`
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func GetLocation(enrollmentid string) (*string) {
	location := fmt.Sprintf("/cps/v2/enrollments/%s", enrollmentid)
	return &location
}

// Create an Enrollment on CPS
//
//
// API Docs: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#5aaa335c
// Endpoint: POST /cps/v2/enrollments{?contractId,deploy-not-after,deploy-not-before}
func (enrollment *Enrollment) Create(params CreateEnrollmentQueryParams) (*CreateEnrollmentResponse, error) {
	var request = fmt.Sprintf(
		"/cps/v2/enrollments?contractId=%s",
		params.ContractID,
	)

	if params.DeployNotAfter != nil {
		request = fmt.Sprintf(
			"%s&deploy-not-after=%s",
			request,
			*params.DeployNotAfter,
		)
	}

	if params.DeployNotBefore != nil {
		request = fmt.Sprintf(
			"%s&deploy-not-before=%s",
			request,
			*params.DeployNotBefore,
		)
	}

	req, err := newRequest(
		"POST",
		request,
		enrollment,
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

	var response CreateEnrollmentResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (enrollment *Enrollment) GetEnrollment() error {

	req, err := client.NewRequest(
		Config,
		"GET",
		*enrollment.Location,
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment.v7+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, &enrollment); err != nil {
		return err
	}

	return nil
}

// Get an enrollment by location
//
//
// API Docs: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#getasingleenrollment
// Endpoint: POST /cps/v2/enrollments/{enrollmentId}
func GetEnrollment(location string) (*Enrollment, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		location,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment.v7+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Enrollment
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (enrollment *Enrollment) GetStatus(enrollmentresponse CreateEnrollmentResponse) (*StatusResponse, error) {

	if len(enrollmentresponse.Changes) == 0 {
		return nil, errors.New("No change id detected")
	}

        req, err := client.NewRequest(
                Config,
                "GET",
                enrollmentresponse.Changes[0],
                nil,
        )

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.akamai.cps.change.v2+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

        var response StatusResponse
        if err = client.BodyJSON(res, &response); err != nil {
                return nil, err
        }

	return &response, nil
}

func getAcknowledgement() (Acknowledgement) {
	var acknowledgement Acknowledgement
	acknowledgement.Acknowledgement = "acknowledge"
	return acknowledgement
}

func (enrollment *Enrollment) AcknowledgeDVChallenges(statusresponse StatusResponse) error {

	if statusresponse.AllowedInput[0].Type != "lets-encrypt-challenges" {
		return errors.New("No Let's Encrypt challenge detected")
	}

	acknowledgement := getAcknowledgement()

        s, err := json.Marshal(acknowledgement)
	if err != nil {
		return err
	}

        req, err := client.NewRequest(
                Config,
                "POST",
                statusresponse.AllowedInput[0].Update,
                bytes.NewReader(s),
        )

	if err != nil {
		return err
	}

        req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
        req.Header.Set("Content-Type", "application/vnd.akamai.cps.acknowledgement.v1+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	return nil

}

func (enrollment *Enrollment) GetDVChallenges(statusresponse StatusResponse) (*DomainValidations, error) {

	if statusresponse.AllowedInput[0].Type != "lets-encrypt-challenges" {
		return nil, errors.New("No Let's Encrypt challenge detected")
	}

        req, err := client.NewRequest(
                Config,
                "GET",
                statusresponse.AllowedInput[0].Info,
                nil,
        )

	if err != nil {
		return nil, err
	}

        req.Header.Set("Accept", "application/vnd.akamai.cps.dv-challenges.v2+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

        var response DomainValidations
        if err = client.BodyJSON(res, &response); err != nil {
                return nil, err
        }

	return &response, nil

}

func (enrollment *Enrollment) Update() (*CreateEnrollmentResponse, error) {

        s, err := json.Marshal(enrollment)
	if err != nil {
		return nil, err
	}

        req, err := client.NewRequest(
                Config,
                "PUT",
                *enrollment.Location,
                bytes.NewReader(s),
        )

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/vnd.akamai.cps.enrollment.v7+json")
	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment-status.v1+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

        var response CreateEnrollmentResponse
        if err = client.BodyJSON(res, &response); err != nil {
                return nil, err
        }

	return &response, nil
}

func ListEnrollments(params ListEnrollmentsQueryParams) ([]Enrollment, error) {
	var enrollments []Enrollment

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/cps/v2/enrollments?contractId={%s}",
			params.ContractID,
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

	if err = client.BodyJSON(res, enrollments); err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (enrollment *Enrollment) Exists(enrollments []Enrollment) bool {
	for _, e := range enrollments {
		if e.CertificateSigningRequest.CommonName == enrollment.CertificateSigningRequest.CommonName {
			return true
		}
	}

	return false
}

// CreateEnrollment wraps enrollment.Create to accept json
func CreateEnrollment(data []byte, params CreateEnrollmentQueryParams) (*CreateEnrollmentResponse, error) {
	var enrollment Enrollment
	if err := json.Unmarshal(data, &enrollment); err != nil {
		return nil, err
	}

	return enrollment.Create(params)
}
