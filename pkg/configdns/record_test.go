package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDns_CreateRecord(t *testing.T) {
	tests := map[string]struct {
		body           RecordBody
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK": {
			responseStatus: http.StatusCreated,
			body: RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
			},
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
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
			body: RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
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
			err := client.CreateRecord(context.Background(), &test.body, "example.com")
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestDns_UpdateRecord(t *testing.T) {
	tests := map[string]struct {
		body           RecordBody
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			responseStatus: http.StatusOK,
			body: RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
			},
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
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
			body: RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
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
			err := client.UpdateRecord(context.Background(), &test.body, "example.com")
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestDns_DeleteRecord(t *testing.T) {
	tests := map[string]struct {
		body           RecordBody
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			responseStatus: http.StatusNoContent,
			body: RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
			},
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			responseBody: ``,
		},
		"500 internal server error": {
			body: RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
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
			err := client.DeleteRecord(context.Background(), &test.body, "example.com")
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestDns_RecordToMap(t *testing.T) {
	client := Client(session.Must(session.New()))

	data := client.RecordToMap(context.Background(), &RecordBody{
		Name:       "www.example.com",
		RecordType: "A",
		TTL:        300,
		Target:     []string{"10.0.0.2", "10.0.0.3"},
	})

	assert.Equal(t, map[string]interface{}{
		"active":     false,
		"name":       "www.example.com",
		"recordtype": "A",
		"target":     []string{"10.0.0.2", "10.0.0.3"},
		"ttl":        300,
	}, data)

	data = client.RecordToMap(context.Background(), &RecordBody{
		RecordType: "A",
		TTL:        300,
		Target:     []string{"10.0.0.2", "10.0.0.3"},
	})

	assert.Nil(t, data)
}

func TestDns_NewRecordBody(t *testing.T) {
	client := Client(session.Must(session.New()))

	toCopy := RecordBody{
		Name: "www.example.com",
	}

	newbody := client.NewRecordBody(context.Background(), toCopy)

	assert.Equal(t, toCopy, *newbody)
}
