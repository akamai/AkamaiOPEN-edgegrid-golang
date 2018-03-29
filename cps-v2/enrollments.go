package cps

import (
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Enrollments represents a CPS enrollments list
type Enrollments struct {
  Enrollments []Enrollment `json:"enrollments"`
}

// NewEnrollments creates an empty Enrollments
func NewEnrollments() *Enrollments {
  e := &Enrollments{}
  return e
}

// GetEnrollments populates Enrollments with the list of avaiable enrollments
//
// API Docs: https://developer.akamai.com/api/luna/cps/resources.html#getenrollments
func GetEnrollments() (*Enrollments, error) {
  e := NewEnrollments()
  if req, err := client.NewRequest(
    Config.EdgeGridConfig,
    "GET",
    "/cps/v2/enrollments",
    nil,
    ); err == nil {
    for k, v := range Config.Headers {
      req.Header.Add(k, v)
    }

    if res, err := client.Do(Config.EdgeGridConfig, req); err == nil {
      if client.IsError(res) {
    		return nil, client.NewAPIError(res)
    	}
      if err := client.BodyJSON(res, e); err != nil {
    		return nil, err
    	}
    } else {
      return nil, err
    }
  } else {
    return nil, err
  }

  return e, nil
}

// Enrollment represents a CPS enrollment resource
type Enrollment struct {
  Location string `json:"location"`
  RootAuth string `json:"ra"`
  ValidationType string `json:"validationType"`
  CertificateType string `json:"certificateType"`
  NetworkConfiguration NetworkConfiguration `json:"networkConfiguration"`
  CSR CSR `json:"csr"`
  SignatureAlgorithm string `json:"signatureAlgorithm"`
  ChangeManagement bool `json:"changeManagement"`
  Org Org `json:"org"`
  AdminContact Contact `json:"adminContact"`
  TechContact Contact `json:"techContact"`
  ThirdParty ThirdParty `json:"thirdParty"`
  EnableMultiStackedCertificates bool `json:"enableMultiStackedCertificates"`
  PendingChanges []string `json:"pendingChanges"`
}

// NewEnrollment creates an empty Enrollment
func NewEnrollment() *Enrollment {
  e := &Enrollment{}
  return e
}

// GetEnrollment populates an Enrollment
//
// API Docs: https://developer.akamai.com/api/luna/cps/resources.html#getasingleenrollment
func GetEnrollment(id string) (*Enrollment, error) {
  e := NewEnrollment()
  if req, err := client.NewRequest(
    Config.EdgeGridConfig,
    "GET",
    "/cps/v2/enrollments/" + id,
    nil,
    ); err == nil {
    for k, v := range Config.Headers {
      req.Header.Add(k, v)
    }

    if res, err := client.Do(Config.EdgeGridConfig, req); err == nil {
      if client.IsError(res) {
    		return nil, client.NewAPIError(res)
    	}
      if err := client.BodyJSON(res, e); err != nil {
    		return nil, err
    	}
    } else {
      return nil, err
    }
  } else {
    return nil, err
  }

  return e, nil
}

// ThirdParty represents an enrollment thirdParty
type ThirdParty struct {
  ExcludeSans bool `json:"excludeSans"`
}

// Org represents an enrollment org
type Org struct {
  Name string `json:"name"`
  AddressLineOne string `json:"addressLineOne"`
  AddressLineTwo string `json:"addressLineTwo"`
  City string `json:"city"`
  Region string `json:"region"`
  PostalCode string `json:"postalCode"`
  Country string `json:"country"`
  Phone string `json:"phone"`
}

// Contact represents an enrollment contact (adminContact|techContact)
type Contact struct {
  FirstName string `json:"firstName"`
  LastName string `json:"lastName"`
  Phone string `json:"phone"`
  Email string `json:"email"`
  AddressLineOne string `json:"addressLineOne"`
  AddressLineTwo string `json:"addressLineTwo"`
  City string `json:"city"`
  Country string `json:"country"`
  OrganizationName string `json:"organizationName"`
  PostalCode string `json:"postalCode"`
  Region string `json:"region"`
  Title string `json:"title"`
}

// NetworkConfiguration represents an enrollment networkConfiguration
type NetworkConfiguration struct {
  Geography string `json:"geography"`
  SecureNetwork string `json:"secureNetwork"`
  MustHaveCiphers string `json:"mustHaveCiphers"`
  PreferredCiphers string `json:"preferredCiphers"`
  DisallowedTlsVersions []string `json:"disallowedTlsVersions"`
  SNI SNI `json:"sni"`
}

// SNI represents a networkConfiguration sni
type SNI struct {
  CloneDnsNames bool `json:"cloneDnsNames"`
  DnsNames []string `json:"dnsNames"`
}

// CSR represents an enrollment csr
type CSR struct {
  CommonName string `json:"cn"`
  Country string `json:"c"`
  State string `json:"st"`
  Location string `json:"l"`
  Org string `json:"o"`
  OrgUnit string `json:"ou"`
  SANS []string `json:"sans"`
}
