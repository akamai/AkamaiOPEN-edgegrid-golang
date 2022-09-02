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

func TestGetDVHistory(t *testing.T) {
	tests := map[string]struct {
		request          GetDVHistoryRequest
		responseBody     string
		responseStatus   int
		withError        error
		expectedPath     string
		expectedResponse *GetDVHistoryResponse
		expectedHeaders  map[string]string
	}{
		"200 ok": {
			request: GetDVHistoryRequest{EnrollmentID: 28926},
			responseBody: `
{
   "results": [
       {
           "domain": "bartdtest.sqa-il.com",
           "domainHistory": [
               {
                   "domain": "bartdtest.sqa-il.com",
                   "responseBody": "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs.xI50JalgaT8I71x4FcarzkEx1UemXcsDqndvWaQqDgc",
                   "fullPath": "http://bartdtest.sqa-il.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                   "token": "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                   "status": "Error",
                   "error": "Expired authorization",
                   "validationStatus": "EXPIRED",
                   "requestTimestamp": "2022-05-18T20:12:05Z",
                   "validatedTimestamp": "2022-05-25T20:14:05Z",
                   "expires": "2022-05-25T20:12:05Z",
                   "redirectFullPath": "http://dcv.akamai.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                   "validationRecords": [],
                   "challenges": [
                       {
                           "type": "http-01",
                           "status": "pending",
                           "error": null,
                           "token": "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                           "responseBody": "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs.xI50JalgaT8I71x4FcarzkEx1UemXcsDqndvWaQqDgc",
                           "fullPath": "http://bartdtest.sqa-il.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                           "redirectFullPath": "http://dcv.akamai.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                           "validationRecords": []
                       },
                       {
                           "type": "dns-01",
                           "status": "pending",
                           "error": null,
                           "token": "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
                           "responseBody": "qyxqdksaqsLTbWIt0im_ob0wYRUVH_Nfe91rmTD3bn0",
                           "fullPath": "_acme-challenge.bartdtest.sqa-il.com.",
                           "redirectFullPath": "",
                           "validationRecords": []
                       }
                   ]
               },
               {
                   "domain": "bartdtest.sqa-il.com",
                   "responseBody": "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI.QzVNu9F4w2DPQOMaqyiwLtnih04pcfunDeZx-LK3h24",
                   "fullPath": "http://bartdtest.sqa-il.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                   "token": "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                   "status": "Error",
                   "error": "Expired authorization",
                   "validationStatus": "EXPIRED",
                   "requestTimestamp": "2022-05-25T20:14:35Z",
                   "validatedTimestamp": "2022-06-01T20:15:13Z",
                   "expires": "2022-06-01T20:14:35Z",
                   "redirectFullPath": "http://dcv.akamai.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                   "validationRecords": [],
                   "challenges": [
                       {
                           "type": "http-01",
                           "status": "pending",
                           "error": null,
                           "token": "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                           "responseBody": "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI.QzVNu9F4w2DPQOMaqyiwLtnih04pcfunDeZx-LK3h24",
                           "fullPath": "http://bartdtest.sqa-il.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                           "redirectFullPath": "http://dcv.akamai.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                           "validationRecords": []
                       },
                       {
                           "type": "dns-01",
                           "status": "pending",
                           "error": null,
                           "token": "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
                           "responseBody": "fNnU1Y9bCKG1jET8px-yE5cSd9HXMg-n6N1rCL0BqdE",
                           "fullPath": "_acme-challenge.bartdtest.sqa-il.com.",
                           "redirectFullPath": "",
                           "validationRecords": []
                       }
                   ]
               }
           ]
       }
   ]
}
`,
			responseStatus: http.StatusOK,
			expectedPath:   "/cps/v2/enrollments/28926/dv-history",
			expectedResponse: &GetDVHistoryResponse{Results: []HistoryResult{{
				Domain: "bartdtest.sqa-il.com",
				DomainHistory: []DomainHistory{
					{
						Domain:             "bartdtest.sqa-il.com",
						ResponseBody:       "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs.xI50JalgaT8I71x4FcarzkEx1UemXcsDqndvWaQqDgc",
						FullPath:           "http://bartdtest.sqa-il.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
						Token:              "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
						Status:             "Error",
						Error:              "Expired authorization",
						ValidationRecords:  []ValidationRecord{},
						ValidationStatus:   "EXPIRED",
						RequestTimestamp:   "2022-05-18T20:12:05Z",
						ValidatedTimestamp: "2022-05-25T20:14:05Z",
						Expires:            "2022-05-25T20:12:05Z",
						RedirectFullPath:   "http://dcv.akamai.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
						Challenges: []Challenge{
							{
								Type:              "http-01",
								Status:            "pending",
								Token:             "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
								ResponseBody:      "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs.xI50JalgaT8I71x4FcarzkEx1UemXcsDqndvWaQqDgc",
								FullPath:          "http://bartdtest.sqa-il.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
								RedirectFullPath:  "http://dcv.akamai.com/.well-known/acme-challenge/zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
								ValidationRecords: []ValidationRecord{},
							},
							{
								Type:              "dns-01",
								Status:            "pending",
								Token:             "zVylMi1AXeE6XbLprYQayt3iojtMSuhELHwhJ2t2pfs",
								ResponseBody:      "qyxqdksaqsLTbWIt0im_ob0wYRUVH_Nfe91rmTD3bn0",
								FullPath:          "_acme-challenge.bartdtest.sqa-il.com.",
								RedirectFullPath:  "",
								ValidationRecords: []ValidationRecord{},
							},
						},
					},
					{
						Domain:             "bartdtest.sqa-il.com",
						ResponseBody:       "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI.QzVNu9F4w2DPQOMaqyiwLtnih04pcfunDeZx-LK3h24",
						FullPath:           "http://bartdtest.sqa-il.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
						Token:              "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
						Status:             "Error",
						Error:              "Expired authorization",
						ValidationRecords:  []ValidationRecord{},
						ValidationStatus:   "EXPIRED",
						RequestTimestamp:   "2022-05-25T20:14:35Z",
						ValidatedTimestamp: "2022-06-01T20:15:13Z",
						Expires:            "2022-06-01T20:14:35Z",
						RedirectFullPath:   "http://dcv.akamai.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
						Challenges: []Challenge{
							{
								Type:              "http-01",
								Status:            "pending",
								Token:             "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
								ResponseBody:      "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI.QzVNu9F4w2DPQOMaqyiwLtnih04pcfunDeZx-LK3h24",
								FullPath:          "http://bartdtest.sqa-il.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
								RedirectFullPath:  "http://dcv.akamai.com/.well-known/acme-challenge/7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
								ValidationRecords: []ValidationRecord{},
							},
							{
								Type:              "dns-01",
								Status:            "pending",
								Token:             "7LwF89FciEFJZb_CUO0xogHQEh-r2iwN6R4BNHvLSoI",
								ResponseBody:      "fNnU1Y9bCKG1jET8px-yE5cSd9HXMg-n6N1rCL0BqdE",
								FullPath:          "_acme-challenge.bartdtest.sqa-il.com.",
								RedirectFullPath:  "",
								ValidationRecords: []ValidationRecord{},
							},
						},
					},
				},
			}}},
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.dv-history.v1+json",
			},
		},
		"404 not found": {
			request: GetDVHistoryRequest{EnrollmentID: 28926000},
			responseBody: `
{
    "type": "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found",
    "title": "Not Found",
    "instance": "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found?id=e925fbec2ef84012a4bee2782c5b0715"
}
`,
			responseStatus: http.StatusNotFound,
			withError: &Error{
				Type:       "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found",
				Title:      "Not Found",
				Instance:   "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found?id=e925fbec2ef84012a4bee2782c5b0715",
				StatusCode: 404,
			},
			expectedPath: "/cps/v2/enrollments/28926000/dv-history",
		},
		"missing enrollmentID": {
			request:   GetDVHistoryRequest{},
			withError: ErrStructValidation,
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
			result, err := client.GetDVHistory(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetCertificateHistory(t *testing.T) {
	tests := map[string]struct {
		request          GetCertificateHistoryRequest
		responseBody     string
		responseStatus   int
		withError        error
		expectedPath     string
		expectedResponse *GetCertificateHistoryResponse
		expectedHeaders  map[string]string
	}{
		"200 ok": {
			request: GetCertificateHistoryRequest{EnrollmentID: 28926},
			responseBody: `
{
    "certificates": [
        {
            "type": "third-party",
            "deploymentStatus": "active",
            "stagingStatus": "active",
            "slots": [
                757836
            ],
            "geography": "core",
            "ra": "third-party",
            "primaryCertificate": {
                "certificate": "-----BEGIN CERTIFICATE-----\nMIIDpDCCAgygAwIBAgIQG7UFcE+swJPyIGzHEPrWTzANBgkqhkiG9w0BAQsFADCB\njTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMTEwLwYDVQQLDCh3emFn\ncmFqY0BrcmstbXAwNmMgKFdvamNpZWNoIFphZ3JhamN6dWspMTgwNgYDVQQDDC9t\na2NlcnQgd3phZ3JhamNAa3JrLW1wMDZjIChXb2pjaWVjaCBaYWdyYWpjenVrKTAe\nFw0yMjA3MjAwNjE0MjlaFw0yNDEwMjAwNjE0MjlaMGsxCzAJBgNVBAYTAlBMMQ0w\nCwYDVQQIEwR0ZXN0MQ0wCwYDVQQHEwR0ZXN0MQ0wCwYDVQQKEwR0ZXN0MQ0wCwYD\nVQQLEwR0ZXN0MSAwHgYDVQQDExdjcHMud3phLXRlc3QwMDEuYWthdGVzdDBZMBMG\nByqGSM49AgEGCCqGSM49AwEHA0IABLAKxSxwaZgqpFVHhFi6feDQ6A8Q0S71p/mP\nuwkUD1zNvmyKPzDflAWDGTyocC8aCzGCHiFdt6CRhCy25RwDK2mjbDBqMA4GA1Ud\nDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAfBgNVHSMEGDAWgBTpiXDW\nLsgNHer9dSJuuzAt6LBpWzAiBgNVHREEGzAZghdjcHMud3phLXRlc3QwMDEuYWth\ndGVzdDANBgkqhkiG9w0BAQsFAAOCAYEAiXmq64svercHBIFbHfNPrRns3ccrOy0U\n+zHks6uSGEO2S/wnAhDtpK2D103SnddRK4WXDF3w0F/GCqrXGmWbKQvzukaUTBbl\neO4Wy36qwc/SypvjzZsPI2q2E6e7uONcB2cL6UA18aEO/w9X4jTfWS5zMV1zPO3N\naozEFRXrziFYGrPAJ9o2RldScHU9stl0TcIiuDWzUqPxvseGBuEMFNtVko590cJA\nl36muto6uC0XV8RtEaMAfbHe1yC64AJd+DzQP47ORQSR3L2+jA8oDQ9pJ+meiBA6\nAt6AJpI1MsekXleRHnZxacUMVYzdk4c472xst+0ueHyMdaIbWgii54csrDfy2vPw\nZwDNm437wJQvqg4RcUrOd5IoM37UCDfyisU9csY4yMXFxwwKFQtIk4Bn+lGbWdjC\nF+NSS+ujtHl0d5rg12QXegbWtFIol+E/ntxG4uS97dpldO2+cQCMM8RGsXA75teL\n4+HRwu0IEa7aaZZAVDZUA2U6wtAmhFM3\n-----END CERTIFICATE-----",
                "trustChain": "",
                "expiry": "2024-10-20T06:14:29Z",
                "keyAlgorithm": "ECDSA"
            },
            "multiStackedCertificates": []
        }
    ]
}
`,
			responseStatus: http.StatusOK,
			expectedPath:   "/cps/v2/enrollments/28926/history/certificates",
			expectedResponse: &GetCertificateHistoryResponse{
				Certificates: []HistoryCertificate{
					{
						Type:             "third-party",
						DeploymentStatus: "active",
						StagingStatus:    "active",
						Slots:            []int{757836},
						Geography:        "core",
						RA:               "third-party",
						PrimaryCertificate: CertificateObject{
							Certificate:  "-----BEGIN CERTIFICATE-----\nMIIDpDCCAgygAwIBAgIQG7UFcE+swJPyIGzHEPrWTzANBgkqhkiG9w0BAQsFADCB\njTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMTEwLwYDVQQLDCh3emFn\ncmFqY0BrcmstbXAwNmMgKFdvamNpZWNoIFphZ3JhamN6dWspMTgwNgYDVQQDDC9t\na2NlcnQgd3phZ3JhamNAa3JrLW1wMDZjIChXb2pjaWVjaCBaYWdyYWpjenVrKTAe\nFw0yMjA3MjAwNjE0MjlaFw0yNDEwMjAwNjE0MjlaMGsxCzAJBgNVBAYTAlBMMQ0w\nCwYDVQQIEwR0ZXN0MQ0wCwYDVQQHEwR0ZXN0MQ0wCwYDVQQKEwR0ZXN0MQ0wCwYD\nVQQLEwR0ZXN0MSAwHgYDVQQDExdjcHMud3phLXRlc3QwMDEuYWthdGVzdDBZMBMG\nByqGSM49AgEGCCqGSM49AwEHA0IABLAKxSxwaZgqpFVHhFi6feDQ6A8Q0S71p/mP\nuwkUD1zNvmyKPzDflAWDGTyocC8aCzGCHiFdt6CRhCy25RwDK2mjbDBqMA4GA1Ud\nDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAfBgNVHSMEGDAWgBTpiXDW\nLsgNHer9dSJuuzAt6LBpWzAiBgNVHREEGzAZghdjcHMud3phLXRlc3QwMDEuYWth\ndGVzdDANBgkqhkiG9w0BAQsFAAOCAYEAiXmq64svercHBIFbHfNPrRns3ccrOy0U\n+zHks6uSGEO2S/wnAhDtpK2D103SnddRK4WXDF3w0F/GCqrXGmWbKQvzukaUTBbl\neO4Wy36qwc/SypvjzZsPI2q2E6e7uONcB2cL6UA18aEO/w9X4jTfWS5zMV1zPO3N\naozEFRXrziFYGrPAJ9o2RldScHU9stl0TcIiuDWzUqPxvseGBuEMFNtVko590cJA\nl36muto6uC0XV8RtEaMAfbHe1yC64AJd+DzQP47ORQSR3L2+jA8oDQ9pJ+meiBA6\nAt6AJpI1MsekXleRHnZxacUMVYzdk4c472xst+0ueHyMdaIbWgii54csrDfy2vPw\nZwDNm437wJQvqg4RcUrOd5IoM37UCDfyisU9csY4yMXFxwwKFQtIk4Bn+lGbWdjC\nF+NSS+ujtHl0d5rg12QXegbWtFIol+E/ntxG4uS97dpldO2+cQCMM8RGsXA75teL\n4+HRwu0IEa7aaZZAVDZUA2U6wtAmhFM3\n-----END CERTIFICATE-----",
							Expiry:       "2024-10-20T06:14:29Z",
							KeyAlgorithm: "ECDSA"},
						MultiStackedCertificates: []CertificateObject{},
					},
				},
			},
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.certificate-history.v2+json",
			},
		},
		"404 not found": {
			request: GetCertificateHistoryRequest{EnrollmentID: 28926000},
			responseBody: `
{
    "type": "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found",
    "title": "Not Found",
    "instance": "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found?id=c3f6bd2c01bc4e138bb5490b0fdb6f5d"
}
`,
			responseStatus: http.StatusNotFound,
			withError: &Error{
				Type:       "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found",
				Title:      "Not Found",
				Instance:   "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found?id=c3f6bd2c01bc4e138bb5490b0fdb6f5d",
				StatusCode: 404,
			},
			expectedPath: "/cps/v2/enrollments/28926000/history/certificates",
		},
		"missing enrollmentID": {
			request:   GetCertificateHistoryRequest{},
			withError: ErrStructValidation,
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
			result, err := client.GetCertificateHistory(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetChangeHistory(t *testing.T) {
	tests := map[string]struct {
		request          GetChangeHistoryRequest
		responseBody     string
		responseStatus   int
		withError        error
		expectedPath     string
		expectedResponse *GetChangeHistoryResponse
		expectedHeaders  map[string]string
	}{
		"200 ok": {
			request: GetChangeHistoryRequest{EnrollmentID: 28926},
			responseBody: `
{
    "changes": [
        {
            "action": "new-certificate",
            "actionDescription": "Create New Certificate",
            "status": "completed",
            "lastUpdated": "2022-07-21T21:40:00Z",
            "createdBy": "wzagrajc",
            "createdOn": "2022-07-18T12:05:41Z",
            "ra": "third-party",
            "primaryCertificate": {
                "certificate": "-----BEGIN CERTIFICATE-----\nMIIDpDCCAgygAwIBAgIQG7UFcE+swJPyIGzHEPrWTzANBgkqhkiG9w0BAQsFADCB\njTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMTEwLwYDVQQLDCh3emFn\ncmFqY0BrcmstbXAwNmMgKFdvamNpZWNoIFphZ3JhamN6dWspMTgwNgYDVQQDDC9t\na2NlcnQgd3phZ3JhamNAa3JrLW1wMDZjIChXb2pjaWVjaCBaYWdyYWpjenVrKTAe\nFw0yMjA3MjAwNjE0MjlaFw0yNDEwMjAwNjE0MjlaMGsxCzAJBgNVBAYTAlBMMQ0w\nCwYDVQQIEwR0ZXN0MQ0wCwYDVQQHEwR0ZXN0MQ0wCwYDVQQKEwR0ZXN0MQ0wCwYD\nVQQLEwR0ZXN0MSAwHgYDVQQDExdjcHMud3phLXRlc3QwMDEuYWthdGVzdDBZMBMG\nByqGSM49AgEGCCqGSM49AwEHA0IABLAKxSxwaZgqpFVHhFi6feDQ6A8Q0S71p/mP\nuwkUD1zNvmyKPzDflAWDGTyocC8aCzGCHiFdt6CRhCy25RwDK2mjbDBqMA4GA1Ud\nDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAfBgNVHSMEGDAWgBTpiXDW\nLsgNHer9dSJuuzAt6LBpWzAiBgNVHREEGzAZghdjcHMud3phLXRlc3QwMDEuYWth\ndGVzdDANBgkqhkiG9w0BAQsFAAOCAYEAiXmq64svercHBIFbHfNPrRns3ccrOy0U\n+zHks6uSGEO2S/wnAhDtpK2D103SnddRK4WXDF3w0F/GCqrXGmWbKQvzukaUTBbl\neO4Wy36qwc/SypvjzZsPI2q2E6e7uONcB2cL6UA18aEO/w9X4jTfWS5zMV1zPO3N\naozEFRXrziFYGrPAJ9o2RldScHU9stl0TcIiuDWzUqPxvseGBuEMFNtVko590cJA\nl36muto6uC0XV8RtEaMAfbHe1yC64AJd+DzQP47ORQSR3L2+jA8oDQ9pJ+meiBA6\nAt6AJpI1MsekXleRHnZxacUMVYzdk4c472xst+0ueHyMdaIbWgii54csrDfy2vPw\nZwDNm437wJQvqg4RcUrOd5IoM37UCDfyisU9csY4yMXFxwwKFQtIk4Bn+lGbWdjC\nF+NSS+ujtHl0d5rg12QXegbWtFIol+E/ntxG4uS97dpldO2+cQCMM8RGsXA75teL\n4+HRwu0IEa7aaZZAVDZUA2U6wtAmhFM3\n-----END CERTIFICATE-----",
                "trustChain": null,
                "csr": "-----BEGIN CERTIFICATE REQUEST-----\nMIIBJTCBzQIBADBrMQswCQYDVQQGEwJQTDENMAsGA1UECAwEdGVzdDENMAsGA1UE\nBwwEdGVzdDENMAsGA1UECgwEdGVzdDENMAsGA1UECwwEdGVzdDEgMB4GA1UEAwwX\nY3BzLnd6YS10ZXN0MDAxLmFrYXRlc3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNC\nAASwCsUscGmYKqRVR4RYun3g0OgPENEu9af5j7sJFA9czb5sij8w35QFgxk8qHAv\nGgsxgh4hXbegkYQstuUcAytpoAAwCgYIKoZIzj0EAwIDRwAwRAIgZdTQC7pGlBQj\n5QHwPy/XFpXhgkssPkPGiyU+Ooauq6ACICCb/nP+DaEs193RFm6MBpScJAxf1F0r\ne9qtJn4af/6h\n-----END CERTIFICATE REQUEST-----\n",
                "keyAlgorithm": "ECDSA"
            },
            "multiStackedCertificates": [],
            "primaryCertificateOrderDetails": null,
            "businessCaseId": null
        }
    ]
}
`,
			responseStatus: http.StatusOK,
			expectedPath:   "/cps/v2/enrollments/28926/history/changes",
			expectedResponse: &GetChangeHistoryResponse{
				Changes: []ChangeHistory{{
					Action:            "new-certificate",
					ActionDescription: "Create New Certificate",
					Status:            "completed",
					LastUpdated:       "2022-07-21T21:40:00Z",
					CreatedBy:         "wzagrajc",
					CreatedOn:         "2022-07-18T12:05:41Z",
					RA:                "third-party",
					PrimaryCertificate: CertificateChangeHistory{
						Certificate:  "-----BEGIN CERTIFICATE-----\nMIIDpDCCAgygAwIBAgIQG7UFcE+swJPyIGzHEPrWTzANBgkqhkiG9w0BAQsFADCB\njTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMTEwLwYDVQQLDCh3emFn\ncmFqY0BrcmstbXAwNmMgKFdvamNpZWNoIFphZ3JhamN6dWspMTgwNgYDVQQDDC9t\na2NlcnQgd3phZ3JhamNAa3JrLW1wMDZjIChXb2pjaWVjaCBaYWdyYWpjenVrKTAe\nFw0yMjA3MjAwNjE0MjlaFw0yNDEwMjAwNjE0MjlaMGsxCzAJBgNVBAYTAlBMMQ0w\nCwYDVQQIEwR0ZXN0MQ0wCwYDVQQHEwR0ZXN0MQ0wCwYDVQQKEwR0ZXN0MQ0wCwYD\nVQQLEwR0ZXN0MSAwHgYDVQQDExdjcHMud3phLXRlc3QwMDEuYWthdGVzdDBZMBMG\nByqGSM49AgEGCCqGSM49AwEHA0IABLAKxSxwaZgqpFVHhFi6feDQ6A8Q0S71p/mP\nuwkUD1zNvmyKPzDflAWDGTyocC8aCzGCHiFdt6CRhCy25RwDK2mjbDBqMA4GA1Ud\nDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAfBgNVHSMEGDAWgBTpiXDW\nLsgNHer9dSJuuzAt6LBpWzAiBgNVHREEGzAZghdjcHMud3phLXRlc3QwMDEuYWth\ndGVzdDANBgkqhkiG9w0BAQsFAAOCAYEAiXmq64svercHBIFbHfNPrRns3ccrOy0U\n+zHks6uSGEO2S/wnAhDtpK2D103SnddRK4WXDF3w0F/GCqrXGmWbKQvzukaUTBbl\neO4Wy36qwc/SypvjzZsPI2q2E6e7uONcB2cL6UA18aEO/w9X4jTfWS5zMV1zPO3N\naozEFRXrziFYGrPAJ9o2RldScHU9stl0TcIiuDWzUqPxvseGBuEMFNtVko590cJA\nl36muto6uC0XV8RtEaMAfbHe1yC64AJd+DzQP47ORQSR3L2+jA8oDQ9pJ+meiBA6\nAt6AJpI1MsekXleRHnZxacUMVYzdk4c472xst+0ueHyMdaIbWgii54csrDfy2vPw\nZwDNm437wJQvqg4RcUrOd5IoM37UCDfyisU9csY4yMXFxwwKFQtIk4Bn+lGbWdjC\nF+NSS+ujtHl0d5rg12QXegbWtFIol+E/ntxG4uS97dpldO2+cQCMM8RGsXA75teL\n4+HRwu0IEa7aaZZAVDZUA2U6wtAmhFM3\n-----END CERTIFICATE-----",
						CSR:          "-----BEGIN CERTIFICATE REQUEST-----\nMIIBJTCBzQIBADBrMQswCQYDVQQGEwJQTDENMAsGA1UECAwEdGVzdDENMAsGA1UE\nBwwEdGVzdDENMAsGA1UECgwEdGVzdDENMAsGA1UECwwEdGVzdDEgMB4GA1UEAwwX\nY3BzLnd6YS10ZXN0MDAxLmFrYXRlc3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNC\nAASwCsUscGmYKqRVR4RYun3g0OgPENEu9af5j7sJFA9czb5sij8w35QFgxk8qHAv\nGgsxgh4hXbegkYQstuUcAytpoAAwCgYIKoZIzj0EAwIDRwAwRAIgZdTQC7pGlBQj\n5QHwPy/XFpXhgkssPkPGiyU+Ooauq6ACICCb/nP+DaEs193RFm6MBpScJAxf1F0r\ne9qtJn4af/6h\n-----END CERTIFICATE REQUEST-----\n",
						KeyAlgorithm: "ECDSA",
					},
					MultiStackedCertificates: []CertificateChangeHistory{},
				}},
			},
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.change-history.v5+json",
			},
		},
		"404 not found": {
			request: GetChangeHistoryRequest{EnrollmentID: 28926000},
			responseBody: `
{
    "type": "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found",
    "title": "Not Found",
    "instance": "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found?id=723cacad466e43688f4a8bd08639ba4e"
}
`,
			responseStatus: http.StatusNotFound,
			withError: &Error{
				Type:       "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found",
				Title:      "Not Found",
				Instance:   "https://akaa-wb66l66toq4ewuc4-haxhlepvmnlgidlc.luna-dev.akamaiapis.net/cps/v2/error-types/not-found?id=723cacad466e43688f4a8bd08639ba4e",
				StatusCode: 404,
			},
			expectedPath: "/cps/v2/enrollments/28926000/history/changes",
		},
		"missing enrollmentID": {
			request:   GetChangeHistoryRequest{},
			withError: ErrStructValidation,
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
			result, err := client.GetChangeHistory(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
