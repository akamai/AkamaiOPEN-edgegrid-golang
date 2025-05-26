package botman

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get CustomBotCategoryItemSequence
func TestBotman_GetCustomBotCategoryItemSequence(t *testing.T) {
	tests := map[string]struct {
		params           GetCustomBotCategoryItemSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomBotCategoryItemSequenceResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "f4f0cb20-eddb-4421-93d9-90954e509d5f",
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"sequence":["fake3f89-e179-4892-89cf-d5e623ba9dc7","fake85df-e399-43e8-bb0f-c0d980a88e4f","fake9b8-4fd5-430e-a061-1c61df1d2ac2"]}`,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/custom-bot-categories/f4f0cb20-eddb-4421-93d9-90954e509d5f/custom-bot-category-item-sequence",
			expectedResponse: &GetCustomBotCategoryItemSequenceResponse{
				Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"},
			},
		},
		"500 internal server error": {
			params: GetCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "f4f0cb20-eddb-4421-93d9-90954e509d5f",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-bot-categories/f4f0cb20-eddb-4421-93d9-90954e509d5f/custom-bot-category-item-sequence",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"missing required params - validation error": {
			params: GetCustomBotCategoryItemSequenceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: CategoryID: cannot be blank\nConfigID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCustomBotCategoryItemSequence(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update CustomBotCategoryItemSequence.
func TestBotman_UpdateCustomBotCategoryItemSequence(t *testing.T) {
	tests := map[string]struct {
		params           UpdateCustomBotCategoryItemSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateCustomBotCategoryItemSequenceResponse
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "f4f0cb20-eddb-4421-93d9-90954e509d5f",
				Sequence:   UUIDSequence{Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"sequence":["fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"]}`,
			expectedResponse: &UpdateCustomBotCategoryItemSequenceResponse{
				Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"},
			},
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-bot-categories/f4f0cb20-eddb-4421-93d9-90954e509d5f/custom-bot-category-item-sequence",
		},
		"500 internal server error": {
			params: UpdateCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "f4f0cb20-eddb-4421-93d9-90954e509d5f",
				Sequence:   UUIDSequence{Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error updating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-bot-categories/f4f0cb20-eddb-4421-93d9-90954e509d5f/custom-bot-category-item-sequence",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"missing required params - validation error": {
			params: UpdateCustomBotCategoryItemSequenceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: CategoryID: cannot be blank\nConfigID: cannot be blank\nSequence: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateCustomBotCategoryItemSequence(
				session.ContextWithOptions(
					context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
