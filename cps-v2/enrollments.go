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

type AcknowledgementResponse struct {
	Change string `json:"change"`
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
		Error struct {
			Timestamp	string `json:"timestamp"`
			Code		string `json:"code"`
			Description	string `json:"description"`
		} `json:"error"`
	} `json:"statusInfo"`
	AllowedInput []struct {
		Type              string `json:"type"`
		RequiredToProceed bool   `json:"requiredToProceed"`
		Info              string `json:"info"`
		Update            string `json:"update"`
	} `json:"allowedInput"`
}

type ThirdPartyCSR struct {
	Csr string `json:"csr"`
}

type ThirdPartyCert struct {
	Certificate string `json:"certificate"`
	TrustChain  string `json:"trustchain"`
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

func (enrollment *Enrollment) GetChangeStatus() (*StatusResponse, error) {

	if enrollment.Location == nil {
		return nil, errors.New("No location set")
	}

	enrollment.GetEnrollment()

	currentchangelist := *enrollment.PendingChanges
        if len(currentchangelist) == 0 {
		return nil, nil
	}

        changeurl := currentchangelist[0]


        req, err := client.NewRequest(
                Config,
                "GET",
                changeurl,
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

func (enrollment *Enrollment) AcknowledgeDVChallenges() (*AcknowledgementResponse, error) {

	statusresponse, err := enrollment.GetChangeStatus()
	if err != nil {
		return nil, err
	}

	if statusresponse == nil {
		return nil, nil
	}

	if len(statusresponse.AllowedInput) == 0 {
		return nil, nil
	}

	if statusresponse.AllowedInput[0].Type != "lets-encrypt-challenges" {
		return nil, nil
	}

	var acknowledgement Acknowledgement
	acknowledgement.Acknowledgement = "acknowledge"

        s, err := json.Marshal(acknowledgement)
	if err != nil {
		return nil, err
	}

        req, err := client.NewRequest(
                Config,
                "POST",
                statusresponse.AllowedInput[0].Update,
                bytes.NewReader(s),
        )

	if err != nil {
		return nil, err
	}

        req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
        req.Header.Set("Content-Type", "application/vnd.akamai.cps.acknowledgement.v1+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	if res.StatusCode == 200 {
		var response AcknowledgementResponse
		if err = client.BodyJSON(res, &response); err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil

}

func (enrollment *Enrollment) SubmitThirdPartyCert(thirdpartycert ThirdPartyCert) (*AcknowledgementResponse, error) {

	statusresponse, err := enrollment.GetChangeStatus()
	if err != nil {
		return nil, err
	}

	if statusresponse == nil {
		return nil, nil
	}

	if len(statusresponse.AllowedInput) == 0 {
		return nil, nil
	}

	if statusresponse.AllowedInput[0].Type != "third-party-csr" {
		return nil, nil
	}

        s, err := json.Marshal(thirdpartycert)
	if err != nil {
		return nil, err
	}

        req, err := client.NewRequest(
                Config,
                "POST",
                statusresponse.AllowedInput[0].Update,
                bytes.NewReader(s),
        )

	if err != nil {
		return nil, err
	}

        req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
        req.Header.Set("Content-Type", "application/vnd.akamai.cps.certificate-and-trust-chain.v1+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	if res.StatusCode == 200 {
		var response AcknowledgementResponse
		if err = client.BodyJSON(res, &response); err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (enrollment *Enrollment) GetThirdPartyCSR() (*ThirdPartyCSR, error) {

	statusresponse, err := enrollment.GetChangeStatus()

	if err != nil {
		return nil, err
	}

	if statusresponse == nil {
		return nil, nil
	}

	if len(statusresponse.AllowedInput) == 0 {
		return nil, nil
	}

	if statusresponse.AllowedInput[0].Type != "third-party-csr" {
		return nil, nil
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

        req.Header.Set("Accept", "application/vnd.akamai.cps.csr.v1+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

        var response ThirdPartyCSR
        if err = client.BodyJSON(res, &response); err != nil {
                return nil, err
        }

	return &response, nil
}

func (enrollment *Enrollment) GetDVChallenges() (*DomainValidations, error) {

	statusresponse, err := enrollment.GetChangeStatus()
	if err != nil {
		return nil, err
	}

	if statusresponse == nil {
		return nil, nil
	}

	if len(statusresponse.AllowedInput) == 0 {
		return nil, nil
	}

	if statusresponse.AllowedInput[0].Type != "lets-encrypt-challenges" {
		return nil, nil
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

func (enrollment *Enrollment) Delete() (*CreateEnrollmentResponse, error) {

        req, err := client.NewRequest(
                Config,
                "DELETE",
		fmt.Sprintf("%s?allow-cancel-pending-changes=true", *enrollment.Location),
                nil,
        )

	if err != nil {
		return nil, err
	}

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
