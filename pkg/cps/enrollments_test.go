package cps

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEnrollments(t *testing.T) {
	tests := map[string]struct {
		params           ListEnrollmentsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedHeaders  map[string]string
		expectedResponse *ListEnrollmentsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         ListEnrollmentsRequest{ContractID: "Contract-123"},
			responseStatus: http.StatusOK,
			responseBody: ` 
{"enrollments":[ {
  "location" : "/cps-api/enrollments/1",
  "ra" : "third-party",
  "validationType" : "third-party",
  "certificateType" : "third-party",
  "certificateChainType" : "default",
  "networkConfiguration" : {
    "geography" : "core",
    "secureNetwork" : "standard-tls",
    "mustHaveCiphers" : "ak-akamai-2020q1",
    "preferredCiphers" : "ak-akamai-2020q1",
    "disallowedTlsVersions" : [ "TLSv1", "TLSv1_1" ],
    "sniOnly" : true,
    "quicEnabled" : false,
    "dnsNameSettings" : {
      "cloneDnsNames" : true,
      "dnsNames" : [ "res-sqa2-3pss-4111-10-1-3CV382-ui.com", "san1.res-sqa2-3pss-4111-10-1-3CV382-ui.com" ]
    },
    "ocspStapling" : "on",
    "clientMutualAuthentication" : null
  },
  "signatureAlgorithm" : null,
  "changeManagement" : false,
  "csr" : {
    "cn" : "res-sqa2-3pss-4111-10-1-3CV382-ui.com",
    "c" : "IN",
    "st" : "KA",
    "l" : "BLR",
    "o" : "Akamai",
    "ou" : "ETG",
    "sans" : [ "san1.res-sqa2-3pss-4111-10-1-3CV382-ui.com", "res-sqa2-3pss-4111-10-1-3CV382-ui.com" ],
    "preferredTrustChain" : null
  },
  "org" : {
    "name" : "Akamai",
    "addressLineOne" : "EGL",
    "addressLineTwo" : "",
    "city" : "BLR",
    "region" : "KA",
    "postalCode" : "71",
    "country" : "IN",
    "phone" : "12"
  },
  "adminContact" : {
    "firstName" : "R1",
    "lastName" : "D1",
    "phone" : "4356",
    "email" : "rd1@akamai.com",
    "addressLineOne" : "EGL",
    "addressLineTwo" : "",
    "city" : "BLR",
    "country" : "IN",
    "organizationName" : "Akamai",
    "postalCode" : "71",
    "region" : "KA",
    "title" : null
  },
  "techContact" : {
    "firstName" : "R2",
    "lastName" : "D2",
    "phone" : "6456",
    "email" : "rd2@akamai.com",
    "addressLineOne" : "150 Broadway",
    "addressLineTwo" : "",
    "city" : "Cambridge",
    "country" : "US",
    "organizationName" : "Akamai Technologies",
    "postalCode" : "02142",
    "region" : "Massachusetts",
    "title" : null
  },
  "thirdParty" : {
    "excludeSans" : true
  },
  "enableMultiStackedCertificates" : false,
  "autoRenewalStartTime" : null,
  "pendingChanges" : [ 
	"/cps-api/enrollments/1/changes/2"
   ],
  "maxAllowedSanNames" : 100,
  "maxAllowedWildcardSanNames" : 100
}, {
  "location" : "/cps-api/enrollments/2",
  "ra" : "lets-encrypt",
  "validationType" : "dv",
  "certificateType" : "san",
  "certificateChainType" : "default",
  "networkConfiguration" : {
    "geography" : "core",
    "secureNetwork" : "enhanced-tls",
    "mustHaveCiphers" : "ak-akamai-default-2017q3",
    "preferredCiphers" : "ak-akamai-default-2017q3",
    "disallowedTlsVersions" : [ "TLSv1", "TLSv1_1" ],
    "sniOnly" : true,
    "quicEnabled" : false,
    "dnsNameSettings" : {
      "cloneDnsNames" : true,
      "dnsNames" : [ "jmm.20210504-dsa12.faden.me" ]
    },
    "ocspStapling" : "on",
    "clientMutualAuthentication" : null
  },
  "signatureAlgorithm" : "SHA-256",
  "changeManagement" : false,
  "csr" : {
    "cn" : "jmm.20210504-dsa12.faden.me",
    "c" : "US",
    "st" : "MA",
    "l" : "Cambridge",
    "o" : "Akamai Technologies, Inc.",
    "ou" : null,
    "sans" : [ "jmm.20210504-dsa12.faden.me" ],
    "preferredTrustChain" : null
  },
  "org" : {
    "name" : "Akamai Technologies, Inc.",
    "addressLineOne" : "150 Broadway",
    "addressLineTwo" : null,
    "city" : "Cambridge",
    "region" : "MA",
    "postalCode" : "02142",
    "country" : "US",
    "phone" : "617-444-3000"
  },
  "adminContact" : {
    "firstName" : "R3",
    "lastName" : "D3",
    "phone" : "8577068086",
    "email" : "rd3@nomail-akamai.com",
    "addressLineOne" : null,
    "addressLineTwo" : null,
    "city" : null,
    "country" : null,
    "organizationName" : null,
    "postalCode" : null,
    "region" : null,
    "title" : null
  },
  "techContact" : {
    "firstName" : "R4",
    "lastName" : "D4",
    "phone" : "617-444-3000",
    "email" : "rd4@akamai.com",
    "addressLineOne" : null,
    "addressLineTwo" : null,
    "city" : null,
    "country" : null,
    "organizationName" : null,
    "postalCode" : null,
    "region" : null,
    "title" : null
  },
  "thirdParty" : null,
  "enableMultiStackedCertificates" : false,
  "autoRenewalStartTime" : null,
  "pendingChanges" : [ 
     "/cps-api/enrollments/2/changes/2"
   ],
  "maxAllowedSanNames" : 100,
  "maxAllowedWildcardSanNames" : 25
},
{
  "location" : "/cps-api/enrollments/3",
  "ra" : "third-party",
  "validationType" : "third-party",
  "certificateType" : "third-party",
  "certificateChainType" : "default",
  "networkConfiguration" : {
    "geography" : "core",
    "secureNetwork" : "enhanced-tls",
    "mustHaveCiphers" : "ak-akamai-2020q1",
    "preferredCiphers" : "ak-akamai-2020q1",
    "disallowedTlsVersions" : [ "TLSv1", "TLSv1_1" ],
    "sniOnly" : true,
    "quicEnabled" : false,
    "dnsNameSettings" : {
      "cloneDnsNames" : true,
      "dnsNames" : [ "san1-submishr-ghj1-mediatest.com", "san2-submishr-ghj1-mediatest.com", "submishr-ghj1-mediatest.com" ]
    },
    "ocspStapling" : "on",
    "clientMutualAuthentication" : null
  },
  "signatureAlgorithm" : null,
  "changeManagement" : false,
  "csr" : {
    "cn" : "submishr-ghj1-mediatest.com",
    "c" : "IN",
    "st" : "karnataka",
    "l" : "Bangalore",
    "o" : "Akamai",
    "ou" : "",
    "sans" : [ "san1-submishr-ghj1-mediatest.com", "san2-submishr-ghj1-mediatest.com", "submishr-ghj1-mediatest.com" ]
  },
  "org" : {
    "name" : "Akamai",
    "addressLineOne" : "EGL",
    "addressLineTwo" : "Bangalore",
    "city" : "Bangalore",
    "region" : "karnataka",
    "postalCode" : "560071",
    "country" : "IN",
    "phone" : "34234353453"
  },
  "adminContact" : {
    "firstName" : "DevQA",
    "lastName" : "Tester",
    "phone" : "6173000033",
    "email" : "devqa@tester.com",
    "addressLineOne" : null,
    "addressLineTwo" : null,
    "city" : null,
    "country" : null,
    "organizationName" : null,
    "postalCode" : null,
    "region" : null,
    "title" : null
  },
  "techContact" : {
    "firstName" : "John",
    "lastName" : "Doe",
    "phone" : "111000111",
    "email" : "john@example.com",
    "addressLineOne" : null,
    "addressLineTwo" : null,
    "city" : null,
    "country" : null,
    "organizationName" : null,
    "postalCode" : null,
    "region" : null,
    "title" : null
  },
  "thirdParty" : {
    "excludeSans" : false
  },
  "enableMultiStackedCertificates" : true,
  "autoRenewalStartTime" : null,
  "pendingChanges" : [ "/cps-api/enrollments/3/changes/30" ],
  "maxAllowedSanNames" : 100,
  "maxAllowedWildcardSanNames" : 100
}
]}`,
			expectedPath: "/cps/v2/enrollments?contractId=Contract-123",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.enrollments.v9+json",
			},
			expectedResponse: &ListEnrollmentsResponse{Enrollments: []Enrollment{
				{
					AdminContact: &Contact{
						AddressLineOne:   "EGL",
						City:             "BLR",
						Country:          "IN",
						Email:            "rd1@akamai.com",
						FirstName:        "R1",
						LastName:         "D1",
						OrganizationName: "Akamai",
						Phone:            "4356",
						PostalCode:       "71",
						Region:           "KA",
					},
					CertificateChainType: "default",
					CertificateType:      "third-party",
					ChangeManagement:     false,
					CSR: &CSR{
						C:  "IN",
						CN: "res-sqa2-3pss-4111-10-1-3CV382-ui.com",
						L:  "BLR",
						O:  "Akamai",
						OU: "ETG",
						SANS: []string{"san1.res-sqa2-3pss-4111-10-1-3CV382-ui.com",
							"res-sqa2-3pss-4111-10-1-3CV382-ui.com"},
						ST: "KA",
					},
					EnableMultiStackedCertificates: false,
					Location:                       "/cps-api/enrollments/1",
					MaxAllowedSanNames:             100,
					MaxAllowedWildcardSanNames:     100,
					NetworkConfiguration: &NetworkConfiguration{
						DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
						DNSNameSettings: &DNSNameSettings{
							CloneDNSNames: true,
							DNSNames: []string{"res-sqa2-3pss-4111-10-1-3CV382-ui.com",
								"san1.res-sqa2-3pss-4111-10-1-3CV382-ui.com"},
						},
						Geography:        "core",
						MustHaveCiphers:  "ak-akamai-2020q1",
						OCSPStapling:     "on",
						PreferredCiphers: "ak-akamai-2020q1",
						QuicEnabled:      false,
						SecureNetwork:    "standard-tls",
						SNIOnly:          true,
					},
					Org: &Org{
						AddressLineOne: "EGL",
						City:           "BLR",
						Country:        "IN",
						Name:           "Akamai",
						Phone:          "12",
						PostalCode:     "71",
						Region:         "KA",
					},
					PendingChanges: []string{"/cps-api/enrollments/1/changes/2"},
					RA:             "third-party",
					TechContact: &Contact{
						AddressLineOne:   "150 Broadway",
						City:             "Cambridge",
						Country:          "US",
						Email:            "rd2@akamai.com",
						FirstName:        "R2",
						LastName:         "D2",
						OrganizationName: "Akamai Technologies",
						Phone:            "6456",
						PostalCode:       "02142",
						Region:           "Massachusetts",
					},
					ThirdParty:     &ThirdParty{ExcludeSANS: true},
					ValidationType: "third-party",
				},
				{
					AdminContact: &Contact{
						Email:     "rd3@nomail-akamai.com",
						FirstName: "R3",
						LastName:  "D3",
						Phone:     "8577068086",
					},
					CertificateChainType: "default",
					CertificateType:      "san",
					ChangeManagement:     false,
					CSR: &CSR{
						C:    "US",
						CN:   "jmm.20210504-dsa12.faden.me",
						L:    "Cambridge",
						O:    "Akamai Technologies, Inc.",
						SANS: []string{"jmm.20210504-dsa12.faden.me"},
						ST:   "MA",
					},
					EnableMultiStackedCertificates: false,
					Location:                       "/cps-api/enrollments/2",
					MaxAllowedSanNames:             100,
					MaxAllowedWildcardSanNames:     25,
					NetworkConfiguration: &NetworkConfiguration{
						DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
						DNSNameSettings: &DNSNameSettings{
							CloneDNSNames: true,
							DNSNames:      []string{"jmm.20210504-dsa12.faden.me"},
						},
						Geography:        "core",
						MustHaveCiphers:  "ak-akamai-default-2017q3",
						OCSPStapling:     "on",
						PreferredCiphers: "ak-akamai-default-2017q3",
						QuicEnabled:      false,
						SecureNetwork:    "enhanced-tls",
						SNIOnly:          true,
					},
					Org: &Org{
						AddressLineOne: "150 Broadway",
						City:           "Cambridge",
						Country:        "US",
						Name:           "Akamai Technologies, Inc.",
						Phone:          "617-444-3000",
						PostalCode:     "02142",
						Region:         "MA",
					},
					PendingChanges: []string{"/cps-api/enrollments/2/changes/2"},
					RA:             "lets-encrypt",
					TechContact: &Contact{
						Email:     "rd4@akamai.com",
						FirstName: "R4",
						LastName:  "D4",
						Phone:     "617-444-3000",
					},
					ValidationType:     "dv",
					SignatureAlgorithm: "SHA-256",
				},
				{
					AdminContact: &Contact{
						Email:     "devqa@tester.com",
						FirstName: "DevQA",
						LastName:  "Tester",
						Phone:     "6173000033",
					},
					CertificateChainType: "default",
					CertificateType:      "third-party",
					ChangeManagement:     false,
					CSR: &CSR{
						C:  "IN",
						CN: "submishr-ghj1-mediatest.com",
						L:  "Bangalore",
						O:  "Akamai",
						SANS: []string{
							"san1-submishr-ghj1-mediatest.com",
							"san2-submishr-ghj1-mediatest.com",
							"submishr-ghj1-mediatest.com"},
						ST: "karnataka",
					},
					EnableMultiStackedCertificates: true,
					Location:                       "/cps-api/enrollments/3",
					MaxAllowedSanNames:             100,
					MaxAllowedWildcardSanNames:     100,
					NetworkConfiguration: &NetworkConfiguration{
						DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
						DNSNameSettings: &DNSNameSettings{
							CloneDNSNames: true,
							DNSNames: []string{
								"san1-submishr-ghj1-mediatest.com",
								"san2-submishr-ghj1-mediatest.com",
								"submishr-ghj1-mediatest.com"},
						},
						Geography:        "core",
						MustHaveCiphers:  "ak-akamai-2020q1",
						OCSPStapling:     "on",
						PreferredCiphers: "ak-akamai-2020q1",
						QuicEnabled:      false,
						SecureNetwork:    "enhanced-tls",
						SNIOnly:          true,
					},
					Org: &Org{
						AddressLineOne: "EGL",
						AddressLineTwo: "Bangalore",
						City:           "Bangalore",
						Country:        "IN",
						Name:           "Akamai",
						Phone:          "34234353453",
						PostalCode:     "560071",
						Region:         "karnataka",
					},
					PendingChanges: []string{"/cps-api/enrollments/3/changes/30"},
					RA:             "third-party",
					TechContact: &Contact{
						Email:     "john@example.com",
						FirstName: "John",
						LastName:  "Doe",
						Phone:     "111000111",
					},
					ThirdParty:     &ThirdParty{ExcludeSANS: false},
					ValidationType: "third-party",
				},
			}},
		},
		"500 internal server error": {
			params:         ListEnrollmentsRequest{ContractID: "1"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments?contractId=1",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.enrollments.v9+json",
			},
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListEnrollments(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetEnrollment(t *testing.T) {
	tests := map[string]struct {
		params           GetEnrollmentRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedHeaders  map[string]string
		expectedResponse *Enrollment
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         GetEnrollmentRequest{EnrollmentID: 1},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "location": "/cps-api/enrollments/1",
    "ra": "third-party",
    "validationType": "third-party",
    "certificateType": "third-party",
    "certificateChainType": "default",
    "networkConfiguration": {
        "geography": "core",
        "secureNetwork": "enhanced-tls",
        "mustHaveCiphers": "ak-akamai-default",
        "preferredCiphers": "ak-akamai-default-interim",
        "disallowedTlsVersions": [
            "TLSv1"
        ],
        "sniOnly": true,
        "quicEnabled": false,
        "dnsNameSettings": {
            "cloneDnsNames": false,
            "dnsNames": [
                "san1.example.com"
            ]
        },
        "ocspStapling": "on",
        "clientMutualAuthentication": {
            "setId": "Custom_CPS-6134b_B-3-1AHBENT.xml",
            "authenticationOptions": {
                "sendCaListToClient": false,
                "ocsp": {
                    "enabled": false
                }
            }
        }
    },
    "signatureAlgorithm": null,
    "changeManagement": true,
    "csr": {
        "cn": "www.example.com",
        "c": "US",
        "st": "MA",
        "l": "Cambridge",
        "o": "Akamai",
        "ou": "WebEx",
        "sans": [
            "www.example.com"
        ]
    },
    "org": {
        "name": "Akamai Technologies",
        "addressLineOne": "150 Broadway",
        "addressLineTwo": null,
        "city": "Cambridge",
        "region": "MA",
        "postalCode": "02142",
        "country": "US",
        "phone": "617-555-0111"
    },
    "adminContact": {
        "firstName": "R1",
        "lastName": "D1",
        "phone": "617-555-0111",
        "email": "r1d1@akamai.com",
        "addressLineOne": "150 Broadway",
        "addressLineTwo": null,
        "city": "Cambridge",
        "country": "US",
        "organizationName": "Akamai",
        "postalCode": "02142",
        "region": "MA",
        "title": "Administrator"
    },
    "techContact": {
        "firstName": "R2",
        "lastName": "D2",
        "phone": "617-555-0111",
        "email": "r2d2@akamai.com",
        "addressLineOne": "150 Broadway",
        "addressLineTwo": null,
        "city": "Cambridge",
        "country": "US",
        "organizationName": "Akamai",
        "postalCode": "02142",
        "region": "MA",
        "title": "Technical Engineer"
    },
    "thirdParty": {
        "excludeSans": false
    },
    "enableMultiStackedCertificates": false,
    "autoRenewalStartTime": null,
    "pendingChanges": [
        "/cps-api/enrollments/1/changes/2"
    ],
    "maxAllowedSanNames": 100,
    "maxAllowedWildcardSanNames": 100
}`,
			expectedPath: "/cps/v2/enrollments/1",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.enrollment.v9+json",
			},
			expectedResponse: &Enrollment{
				AdminContact: &Contact{
					AddressLineOne:   "150 Broadway",
					City:             "Cambridge",
					Country:          "US",
					Email:            "r1d1@akamai.com",
					FirstName:        "R1",
					LastName:         "D1",
					OrganizationName: "Akamai",
					Phone:            "617-555-0111",
					PostalCode:       "02142",
					Region:           "MA",
					Title:            "Administrator",
				},
				CertificateChainType: "default",
				CertificateType:      "third-party",
				ChangeManagement:     true,
				CSR: &CSR{
					C:    "US",
					CN:   "www.example.com",
					L:    "Cambridge",
					O:    "Akamai",
					OU:   "WebEx",
					SANS: []string{"www.example.com"},
					ST:   "MA",
				},
				EnableMultiStackedCertificates: false,
				Location:                       "/cps-api/enrollments/1",
				MaxAllowedSanNames:             100,
				MaxAllowedWildcardSanNames:     100,
				NetworkConfiguration: &NetworkConfiguration{
					ClientMutualAuthentication: &ClientMutualAuthentication{
						AuthenticationOptions: &AuthenticationOptions{
							OCSP:               &OCSP{BoolPtr(false)},
							SendCAListToClient: BoolPtr(false),
						},
						SetID: "Custom_CPS-6134b_B-3-1AHBENT.xml",
					},
					DisallowedTLSVersions: []string{"TLSv1"},
					DNSNameSettings: &DNSNameSettings{
						CloneDNSNames: false,
						DNSNames:      []string{"san1.example.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-default",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-default-interim",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				},
				Org: &Org{
					AddressLineOne: "150 Broadway",
					City:           "Cambridge",
					Country:        "US",
					Name:           "Akamai Technologies",
					Phone:          "617-555-0111",
					PostalCode:     "02142",
					Region:         "MA",
				},
				PendingChanges: []string{"/cps-api/enrollments/1/changes/2"},
				RA:             "third-party",
				TechContact: &Contact{
					AddressLineOne:   "150 Broadway",
					City:             "Cambridge",
					Country:          "US",
					Email:            "r2d2@akamai.com",
					FirstName:        "R2",
					LastName:         "D2",
					OrganizationName: "Akamai",
					Phone:            "617-555-0111",
					PostalCode:       "02142",
					Region:           "MA",
					Title:            "Technical Engineer",
				},
				ThirdParty:     &ThirdParty{ExcludeSANS: false},
				ValidationType: "third-party",
			},
		},
		"500 internal server error": {
			params:         GetEnrollmentRequest{EnrollmentID: 1},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.enrollment.v9+json",
			},
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetEnrollment(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateEnrollment(t *testing.T) {
	tests := map[string]struct {
		request          CreateEnrollmentRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateEnrollmentResponse
		withError        error
	}{
		"202 accepted": {
			request: CreateEnrollmentRequest{
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				ContractID:      "ctr-1",
				DeployNotAfter:  "12-12-2021",
				DeployNotBefore: "12-07-2020",
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
	"enrollment": "/cps-api/enrollments/1",
	"changes": ["/cps-api/enrollments/1/changes/10002"]
}`,
			expectedPath: "/cps/v2/enrollments?contractId=ctr-1&deploy-not-after=12-12-2021&deploy-not-before=12-07-2020",
			expectedResponse: &CreateEnrollmentResponse{
				Enrollment: "/cps-api/enrollments/1",
				Changes:    []string{"/cps-api/enrollments/1/changes/10002"},
				ID:         1,
			},
		},
		"202 accepted allow duplicate cn": {
			request: CreateEnrollmentRequest{
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				ContractID:       "ctr-1",
				DeployNotAfter:   "12-12-2021",
				DeployNotBefore:  "12-07-2020",
				AllowDuplicateCN: true,
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
	"enrollment": "/cps-api/enrollments/1",
	"changes": ["/cps-api/enrollments/1/changes/10002"]
}`,
			expectedPath: "/cps/v2/enrollments?allow-duplicate-cn=true&contractId=ctr-1&deploy-not-after=12-12-2021&deploy-not-before=12-07-2020",
			expectedResponse: &CreateEnrollmentResponse{
				Enrollment: "/cps-api/enrollments/1",
				Changes:    []string{"/cps-api/enrollments/1/changes/10002"},
				ID:         1,
			},
		},
		"500 internal server error": {
			request: CreateEnrollmentRequest{
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				ContractID:      "ctr-1",
				DeployNotAfter:  "12-12-2021",
				DeployNotBefore: "12-07-2020",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments?contractId=ctr-1&deploy-not-after=12-12-2021&deploy-not-before=12-07-2020",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request:   CreateEnrollmentRequest{},
			withError: ErrStructValidation,
		},
		"invalid location": {
			request: CreateEnrollmentRequest{
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				ContractID:      "ctr-1",
				DeployNotAfter:  "12-12-2021",
				DeployNotBefore: "12-07-2020",
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
	"enrollment": "abc",
	"changes": ["/cps-api/enrollments/1/changes/10002"]
}`,
			expectedPath: "/cps/v2/enrollments?contractId=ctr-1&deploy-not-after=12-12-2021&deploy-not-before=12-07-2020",
			withError:    ErrInvalidLocation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.enrollment-status.v1+json", r.Header.Get("Accept"))
				assert.Equal(t, "application/vnd.akamai.cps.enrollment.v11+json; charset=utf-8", r.Header.Get("Content-Type"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateEnrollment(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateEnrollment(t *testing.T) {
	tests := map[string]struct {
		request          UpdateEnrollmentRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateEnrollmentResponse
		withError        error
	}{
		"202 accepted": {
			request: UpdateEnrollmentRequest{
				EnrollmentID: 1,
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				DeployNotAfter:            "12-12-2021",
				DeployNotBefore:           "12-07-2020",
				RenewalDateCheckOverride:  BoolPtr(true),
				AllowCancelPendingChanges: BoolPtr(true),
				AllowStagingBypass:        BoolPtr(true),
				ForceRenewal:              BoolPtr(true),
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
	"enrollment": "/cps-api/enrollments/1",
	"changes": ["/cps-api/enrollments/1/changes/10002"]
}`,
			expectedPath: "/cps/v2/enrollments/1?allow-cancel-pending-changes=true&allow-staging-bypass=true&deploy-not-after=12-12-2021&deploy-not-before=12-07-2020&force-renewal=true&renewal-date-check-override=true",
			expectedResponse: &UpdateEnrollmentResponse{
				Enrollment: "/cps-api/enrollments/1",
				Changes:    []string{"/cps-api/enrollments/1/changes/10002"},
				ID:         1,
			},
		},
		"500 internal server error": {
			request: UpdateEnrollmentRequest{
				EnrollmentID: 1,
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				DeployNotAfter:  "12-12-2021",
				DeployNotBefore: "12-07-2020",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error updating enrollment",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1?deploy-not-after=12-12-2021&deploy-not-before=12-07-2020",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request:   UpdateEnrollmentRequest{},
			withError: ErrStructValidation,
		},
		"invalid location URL": {
			request: UpdateEnrollmentRequest{
				EnrollmentID: 1,
				Enrollment: Enrollment{
					AdminContact: &Contact{
						Email: "r1d1@akamai.com",
					},
					CertificateType: "third-party",
					CSR: &CSR{
						CN: "www.example.com",
					},
					NetworkConfiguration: &NetworkConfiguration{},
					Org:                  &Org{Name: "Akamai"},
					RA:                   "third-party",
					TechContact: &Contact{
						Email: "r2d2@akamai.com",
					},
					ValidationType: "third-party",
				},
				DeployNotAfter:            "12-12-2021",
				DeployNotBefore:           "12-07-2020",
				RenewalDateCheckOverride:  BoolPtr(true),
				AllowCancelPendingChanges: BoolPtr(true),
				AllowStagingBypass:        BoolPtr(true),
				ForceRenewal:              BoolPtr(true),
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
	"enrollment": "abc",
	"changes": ["/cps-api/enrollments/1/changes/10002"]
}`,
			expectedPath: "/cps/v2/enrollments/1?allow-cancel-pending-changes=true&allow-staging-bypass=true&deploy-not-after=12-12-2021&deploy-not-before=12-07-2020&force-renewal=true&renewal-date-check-override=true",
			withError:    ErrInvalidLocation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.enrollment-status.v1+json", r.Header.Get("Accept"))
				assert.Equal(t, "application/vnd.akamai.cps.enrollment.v11+json; charset=utf-8", r.Header.Get("Content-Type"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateEnrollment(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestRemoveEnrollment(t *testing.T) {
	tests := map[string]struct {
		request          RemoveEnrollmentRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveEnrollmentResponse
		withError        error
	}{
		"200 OK": {
			request: RemoveEnrollmentRequest{
				EnrollmentID:              1,
				AllowCancelPendingChanges: BoolPtr(true),
				DeployNotAfter:            "12-12-2021",
				DeployNotBefore:           "12-07-2021",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"enrollment": "/cps-api/enrollments/1",
	"changes": ["/cps-api/enrollments/1/changes/10002"]
}`,
			expectedPath: "/cps/v2/enrollments/1?allow-cancel-pending-changes=true&deploy-not-after=12-12-2021&deploy-not-before=12-07-2021",
			expectedResponse: &RemoveEnrollmentResponse{
				Enrollment: "/cps-api/enrollments/1",
				Changes:    []string{"/cps-api/enrollments/1/changes/10002"},
			},
		},
		"500 internal server error": {
			request: RemoveEnrollmentRequest{
				EnrollmentID:              1,
				AllowCancelPendingChanges: BoolPtr(true),
				DeployNotAfter:            "12-12-2021",
				DeployNotBefore:           "12-07-2021",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error removing enrollment",
    "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1?allow-cancel-pending-changes=true&deploy-not-after=12-12-2021&deploy-not-before=12-07-2021",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error removing enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request:   RemoveEnrollmentRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.enrollment-status.v1+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveEnrollment(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func BoolPtr(b bool) *bool {
	return &b
}
