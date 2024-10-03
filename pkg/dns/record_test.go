package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNS_CreateRecord(t *testing.T) {
	tests := map[string]struct {
		params         CreateRecordRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK": {
			params: CreateRecordRequest{
				Record: &RecordBody{
					Name:       "www.example.com",
					RecordType: "A",
					TTL:        300,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				},
				Zone: "example.com",
			},
			responseStatus: http.StatusCreated,
			expectedPath:   "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			responseBody: `
			{
				"name": "www.example.com",
				"type": "A",
				"ttl": 300,
				"rdata": [
					"10.0.0.2",
					"10.0.0.3"
				]
			}`,
		},
		"500 internal server error": {
			params: CreateRecordRequest{
				Record: &RecordBody{
					Name:       "www.example.com",
					RecordType: "A",
					TTL:        300,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				},
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.CreateRecord(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestDNS_UpdateRecord(t *testing.T) {
	tests := map[string]struct {
		params         UpdateRecordRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			params: UpdateRecordRequest{
				Record: &RecordBody{
					Name:       "www.example.com",
					RecordType: "A",
					TTL:        300,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				},
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			responseBody: `
			{
				"name": "www.example.com",
				"type": "A",
				"ttl": 300,
				"rdata": [
					"10.0.0.2",
					"10.0.0.3"
				]
			}`,
		},
		"500 internal server error": {
			params: UpdateRecordRequest{
				Record: &RecordBody{
					Name:       "www.example.com",
					RecordType: "A",
					TTL:        300,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				},
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.UpdateRecord(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestDNS_DeleteRecord(t *testing.T) {
	tests := map[string]struct {
		params         DeleteRecordRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			params: DeleteRecordRequest{
				Name:       "www.example.com",
				RecordType: "A",
				Zone:       "example.com",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			responseBody:   ``,
		},
		"500 internal server error": {
			params: DeleteRecordRequest{
				Name:       "www.example.com",
				RecordType: "A",
				Zone:       "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteRecord(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
