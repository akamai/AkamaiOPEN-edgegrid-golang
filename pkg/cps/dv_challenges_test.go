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

func TestGetChangeLetsEncryptChallenges(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DVArray
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "dv": [
        {
            "status": "Awaiting user",
            "error": "The domain is not ready for validation.",
            "validationStatus": "RESPONSE_ERROR",
            "requestTimestamp": "2018-09-05T15:55:49Z",
            "validatedTimestamp": "2018-09-05T17:53:22Z",
            "expires": "2018-09-06T17:55:17Z",
            "challenges": [
                {
                    "type": "dns-01",
                    "status": "pending",
                    "error": null,
                    "token": "cGBnw-3YO7rUhq61EuuHqcGrYkaQWALAgi8szTqRoHA",
                    "responseBody": "0yVISDJjpXR7BXzR5QgfA51tt-I6aKremGnPwK_lvH4",
                    "fullPath": "_acme-challenge.www.cps-example-dv.com.",
                    "redirectFullPath": "",
                    "validationRecords": []
                }
            ],
            "domain": "www.cps-example-dv.com"
        }
    ]
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/lets-encrypt-challenges",
			expectedResponse: &DVArray{DV: []DV{
				{
					Challenges: []Challenge{
						{
							FullPath:          "_acme-challenge.www.cps-example-dv.com.",
							ResponseBody:      "0yVISDJjpXR7BXzR5QgfA51tt-I6aKremGnPwK_lvH4",
							Status:            "pending",
							Token:             "cGBnw-3YO7rUhq61EuuHqcGrYkaQWALAgi8szTqRoHA",
							Type:              "dns-01",
							ValidationRecords: []ValidationRecord{},
						},
					},
					Domain:             "www.cps-example-dv.com",
					Error:              "The domain is not ready for validation.",
					Expires:            "2018-09-06T17:55:17Z",
					RequestTimestamp:   "2018-09-05T15:55:49Z",
					Status:             "Awaiting user",
					ValidatedTimestamp: "2018-09-05T17:53:22Z",
					ValidationStatus:   "RESPONSE_ERROR",
				},
			}},
		},
		"500 internal server error": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/lets-encrypt-challenges",
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
		"validation error": {
			params: GetChangeRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.dv-challenges.v2+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeLetsEncryptChallenges(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAcknowledgeDVChallenges(t *testing.T) {
	tests := map[string]struct {
		params           AcknowledgementRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DVArray
		withError        func(*testing.T, error)
	}{
		"204 no content": {
			params: AcknowledgementRequest{
				EnrollmentID:    1,
				ChangeID:        2,
				Acknowledgement: Acknowledgement{"acknowledge"},
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/cps/v2/enrollments/1/changes/2/input/update/lets-encrypt-challenges-completed",
		},
		"500 internal server error": {
			params: AcknowledgementRequest{
				EnrollmentID:    1,
				ChangeID:        2,
				Acknowledgement: Acknowledgement{"acknowledge"},
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/cps/v2/enrollments/1/changes/2/input/update/lets-encrypt-challenges-completed",
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
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
		"validation error": {
			params: AcknowledgementRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.change-id.v1+json", r.Header.Get("Accept"))
				assert.Equal(t, "application/vnd.akamai.cps.acknowledgement.v1+json", r.Header.Get("Content-Type"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				require.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.AcknowledgeDVChallenges(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
