package botman

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

// Test Get RecategorizedAkamaiDefinedBot List
func TestBotman_GetRecategorizedAkamaiDefinedBotList(t *testing.T) {

	tests := map[string]struct {
		params           GetRecategorizedAkamaiDefinedBotListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRecategorizedAkamaiDefinedBotListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetRecategorizedAkamaiDefinedBotListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"recategorizedBots": [
		{"botId":"b85e3eaa-d334-466d-857e-33308ce416be", "customBotCategoryId":"39cbadc6-c5ef-42d1-9290-7895f24316ad"},
		{"botId":"69acad64-7459-4c1d-9bad-672600150127", "customBotCategoryId":"5eb700c8-275d-4866-a271-b6fa25e1fdc5"},
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
		{"botId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "customBotCategoryId":"87a152a9-8af0-4c4f-9c37-a895fe7ca6b4"},
		{"botId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "customBotCategoryId":"b61a3017-bff4-41b0-9396-be378d4f07c1"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots",
			expectedResponse: &GetRecategorizedAkamaiDefinedBotListResponse{
				Bots: []RecategorizedAkamaiDefinedBotResponse{
					{BotID: "b85e3eaa-d334-466d-857e-33308ce416be", CategoryID: "39cbadc6-c5ef-42d1-9290-7895f24316ad"},
					{BotID: "69acad64-7459-4c1d-9bad-672600150127", CategoryID: "5eb700c8-275d-4866-a271-b6fa25e1fdc5"},
					{BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
					{BotID: "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", CategoryID: "87a152a9-8af0-4c4f-9c37-a895fe7ca6b4"},
					{BotID: "4d64d85a-a07f-485a-bbac-24c60658a1b8", CategoryID: "b61a3017-bff4-41b0-9396-be378d4f07c1"},
				},
			},
		},
		"200 OK One Record": {
			params: GetRecategorizedAkamaiDefinedBotListRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"recategorizedBots":[
		{"botId":"b85e3eaa-d334-466d-857e-33308ce416be", "customBotCategoryId":"39cbadc6-c5ef-42d1-9290-7895f24316ad"},
		{"botId":"69acad64-7459-4c1d-9bad-672600150127", "customBotCategoryId":"5eb700c8-275d-4866-a271-b6fa25e1fdc5"},
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
		{"botId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "customBotCategoryId":"87a152a9-8af0-4c4f-9c37-a895fe7ca6b4"},
		{"botId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "customBotCategoryId":"b61a3017-bff4-41b0-9396-be378d4f07c1"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots",
			expectedResponse: &GetRecategorizedAkamaiDefinedBotListResponse{
				Bots: []RecategorizedAkamaiDefinedBotResponse{
					{BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
				},
			},
		},
		"500 internal server error": {
			params: GetRecategorizedAkamaiDefinedBotListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
   "type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error fetching data",
   "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots",
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
		"Missing ConfigID": {
			params: GetRecategorizedAkamaiDefinedBotListRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetRecategorizedAkamaiDefinedBotListRequest{
				ConfigID: 43253,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
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
			result, err := client.GetRecategorizedAkamaiDefinedBotList(
				session.ContextWithOptions(
					context.Background(),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Get RecategorizedAkamaiDefinedBot
func TestBotman_GetRecategorizedAkamaiDefinedBot(t *testing.T) {
	tests := map[string]struct {
		params           GetRecategorizedAkamaiDefinedBotRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecategorizedAkamaiDefinedBotResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"}`,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			expectedResponse: &RecategorizedAkamaiDefinedBotResponse{
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
		},
		"500 internal server error": {
			params: GetRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
		"Missing ConfigID": {
			params: GetRecategorizedAkamaiDefinedBotRequest{
				Version: 15,
				BotID:   "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing BotID": {
			params: GetRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "BotID")
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
			result, err := client.GetRecategorizedAkamaiDefinedBot(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create RecategorizedAkamaiDefinedBot
func TestBotman_CreateRecategorizedAkamaiDefinedBot(t *testing.T) {

	tests := map[string]struct {
		params           CreateRecategorizedAkamaiDefinedBotRequest
		prop             *CreateRecategorizedAkamaiDefinedBotRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecategorizedAkamaiDefinedBotResponse
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    15,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			responseStatus: http.StatusCreated,
			responseBody:   `{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"}`,
			expectedResponse: &RecategorizedAkamaiDefinedBotResponse{
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			expectedPath: "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots",
		},
		"500 internal server error": {
			params: CreateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    15,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/recategorized-akamai-defined-bots",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: CreateRecategorizedAkamaiDefinedBotRequest{
				Version:    15,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: CreateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing BotID": {
			params: CreateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "BotID")
			},
		},
		"Missing CategoryID": {
			params: CreateRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CategoryID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateRecategorizedAkamaiDefinedBot(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update RecategorizedAkamaiDefinedBot
func TestBotman_UpdateRecategorizedAkamaiDefinedBot(t *testing.T) {
	tests := map[string]struct {
		params           UpdateRecategorizedAkamaiDefinedBotRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecategorizedAkamaiDefinedBotResponse
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    10,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"}`,
			expectedResponse: &RecategorizedAkamaiDefinedBotResponse{
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			expectedPath: "/appsec/v1/configs/43253/versions/10/recategorized-akamai-defined-bots/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: UpdateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    10,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/recategorized-akamai-defined-bots/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating zone",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: UpdateRecategorizedAkamaiDefinedBotRequest{
				Version:    15,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing CategoryID": {
			params: UpdateRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CategoryID")
			},
		},
		"Missing BotID": {
			params: UpdateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "BotID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRecategorizedAkamaiDefinedBot(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Remove RecategorizedAkamaiDefinedBot
func TestBotman_RemoveRecategorizedAkamaiDefinedBot(t *testing.T) {
	tests := map[string]struct {
		params           RemoveRecategorizedAkamaiDefinedBotRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: RemoveRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  10,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/10/recategorized-akamai-defined-bots/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: RemoveRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  10,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/recategorized-akamai-defined-bots/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error deleting match target",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: RemoveRecategorizedAkamaiDefinedBotRequest{
				Version: 15,
				BotID:   "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: RemoveRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing BotID": {
			params: RemoveRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "BotID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.RemoveRecategorizedAkamaiDefinedBot(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
