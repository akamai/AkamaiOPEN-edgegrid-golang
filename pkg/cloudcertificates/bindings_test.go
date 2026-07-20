package cloudcertificates

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListLineageBindings(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListLineageBindingsRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListLineageBindingsResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - list bindings": {
			params:         ListLineageBindingsRequest{LineageID: 12345},
			expectedPath:   "/ccm/v2/lineages/12345/bindings",
			responseStatus: http.StatusOK,
			responseBody: `{
				"bindings": [
					{ "active": true,  "hostname": "www.example.com",     "network": ["PRODUCTION"] },
					{ "active": false, "hostname": "api.example.com",     "network": ["PRODUCTION"] },
					{ "active": true,  "hostname": "staging.example.com", "network": ["STAGING", "PRODUCTION"] }
				],
				"nextCursor": "12350",
				"totalCount": 87
			}`,
			expectedResponse: &ListLineageBindingsResponse{
				Bindings: []LineageBinding{
					{Active: true, Hostname: "www.example.com", Networks: []string{"PRODUCTION"}},
					{Active: false, Hostname: "api.example.com", Networks: []string{"PRODUCTION"}},
					{Active: true, Hostname: "staging.example.com", Networks: []string{"STAGING", "PRODUCTION"}},
				},
				NextCursor: ptr.To("12350"),
				TotalCount: 87,
			},
		},
		"200 OK - list bindings, empty": {
			params:         ListLineageBindingsRequest{LineageID: 500033},
			expectedPath:   "/ccm/v2/lineages/500033/bindings",
			responseStatus: http.StatusOK,
			responseBody: `{
				"bindings": [],
				"totalCount": 0
			}`,
			expectedResponse: &ListLineageBindingsResponse{
				Bindings:   []LineageBinding{},
				TotalCount: 0,
			},
		},
		"200 OK - filtered by network, pageSize, after and sort": {
			params: ListLineageBindingsRequest{
				LineageID: 12345,
				Network:   TargetNetworkStaging,
				PageSize:  2,
				After:     "12300",
				Sort:      SortOrderDescending,
			},
			expectedPath:   "/ccm/v2/lineages/12345/bindings?after=12300&network=STAGING&pageSize=2&sort=DESC",
			responseStatus: http.StatusOK,
			responseBody: `{
				"bindings": [
					{ "active": true, "hostname": "staging.example.com", "network": ["STAGING"] }
				],
				"totalCount": 1
			}`,
			expectedResponse: &ListLineageBindingsResponse{
				Bindings: []LineageBinding{
					{Active: true, Hostname: "staging.example.com", Networks: []string{"STAGING"}},
				},
				TotalCount: 1,
			},
		},
		"200 OK - maximum page size": {
			params:         ListLineageBindingsRequest{LineageID: 12345, PageSize: MaxListLineageBindingsPageSize},
			expectedPath:   fmt.Sprintf("/ccm/v2/lineages/12345/bindings?pageSize=%d", MaxListLineageBindingsPageSize),
			responseStatus: http.StatusOK,
			responseBody:   `{"bindings": [], "totalCount": 0}`,
			expectedResponse: &ListLineageBindingsResponse{
				Bindings:   []LineageBinding{},
				TotalCount: 0,
			},
		},
		"404 lineage not found": {
			params:         ListLineageBindingsRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999/bindings",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891080",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageBindings, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891080",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrListLineageBindings)
			},
		},
		"403 forbidden - missing ACMI_CCM_READ_ONLY": {
			params:         ListLineageBindingsRequest{LineageID: 12345},
			expectedPath:   "/ccm/v2/lineages/12345/bindings",
			responseStatus: http.StatusForbidden,
			responseBody: `{
				"type": "/error-types/forbidden",
				"title": "Forbidden.",
				"status": 403,
				"detail": "User lacks ACMI_CCM_READ_ONLY on this lineage.",
				"instance": "/error-types/forbidden?traceId=1234567891081"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageBindings, &Error{
					Type:     "/error-types/forbidden",
					Title:    "Forbidden.",
					Status:   http.StatusForbidden,
					Detail:   "User lacks ACMI_CCM_READ_ONLY on this lineage.",
					Instance: "/error-types/forbidden?traceId=1234567891081",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrForbidden)
				assert.ErrorIs(t, err, ErrListLineageBindings)
			},
		},
		"500 internal server error": {
			params:         ListLineageBindingsRequest{LineageID: 12345},
			expectedPath:   "/ccm/v2/lineages/12345/bindings",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891085"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageBindings, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891085",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrListLineageBindings)
			},
		},
		"validation error - invalid network value": {
			params: ListLineageBindingsRequest{LineageID: 12345, Network: "not-a-real-network"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineage bindings: struct validation: Network: value "+
					"'not-a-real-network' is invalid. Must be either 'STAGING' or 'PRODUCTION'")
				assert.ErrorIs(t, err, ErrListLineageBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid sort value": {
			params: ListLineageBindingsRequest{LineageID: 12345, Sort: "not-a-real-sort"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineage bindings: struct validation: Sort: value "+
					"'not-a-real-sort' is invalid. Must be either 'ASC' or 'DESC'")
				assert.ErrorIs(t, err, ErrListLineageBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"400 invalid or expired after cursor": {
			params:         ListLineageBindingsRequest{LineageID: 12345, After: "invalid-or-expired-cursor"},
			expectedPath:   "/ccm/v2/lineages/12345/bindings?after=invalid-or-expired-cursor",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/invalid-cursor",
				"title": "Invalid or expired pagination cursor.",
				"status": 400,
				"detail": "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
				"instance": "/error-types/invalid-cursor?traceId=1234567891084"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageBindings, &Error{
					Type:     "/error-types/invalid-cursor",
					Title:    "Invalid or expired pagination cursor.",
					Status:   http.StatusBadRequest,
					Detail:   "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
					Instance: "/error-types/invalid-cursor?traceId=1234567891084",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInvalidCursor)
				assert.ErrorIs(t, err, ErrListLineageBindings)
			},
		},
		"validation error - missing LineageID": {
			params: ListLineageBindingsRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineage bindings: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrListLineageBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - PageSize exceeds maximum": {
			params: ListLineageBindingsRequest{LineageID: 12345, PageSize: 101},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, fmt.Sprintf("listing lineage bindings: struct validation: PageSize: must be no greater than %d", MaxListLineageBindingsPageSize))
				assert.ErrorIs(t, err, ErrListLineageBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - PageSize is negative": {
			params: ListLineageBindingsRequest{LineageID: 12345, PageSize: -1},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineage bindings: struct validation: PageSize: must be no less than 0")
				assert.ErrorIs(t, err, ErrListLineageBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListLineageBindings(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
