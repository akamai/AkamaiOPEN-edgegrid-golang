package reportinggroups

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Parallel()
	sess, err := session.New()
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)

	tests := map[string]struct {
		response *http.Response
		expected *Error
	}{
		"Bad request 400": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`
				{
					"code": "bad.request",
					"title": "Bad Request",
					"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					"details": [
						{
						  "code": "invalid.data",
						  "message": "Group name is mandatory"
						}
					]
				}`)),
				Request: req,
			},
			expected: &Error{
				Code:       "bad.request",
				Title:      "Bad Request",
				IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				HTTPStatus: http.StatusBadRequest,
				Details: []SecondaryError{
					{
						Code:    "invalid.data",
						Message: "Group name is mandatory"},
				},
			},
		},
		"Unauthorized 401": {
			response: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body: io.NopCloser(strings.NewReader(`
				{
				  "code": "unauthorized",
				  "details": [
					{
					  "code": "invalid.data",
					  "message": "The request requires authentication"
					}
				  ],
				  "incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				  "title": "Unauthorized"
				}`)),
				Request: req,
			},
			expected: &Error{
				Code:       "unauthorized",
				Title:      "Unauthorized",
				IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				HTTPStatus: http.StatusUnauthorized,
				Details: []SecondaryError{
					{
						Code:    "invalid.data",
						Message: "The request requires authentication"},
				},
			},
		},
		"Forbidden 403": {
			response: &http.Response{
				StatusCode: http.StatusForbidden,
				Body: io.NopCloser(strings.NewReader(`
				{
				  "code": "forbidden",
				  "details": [
					{
					  "code": "invalid.data",
					  "message": "User is not authorized"
					}
				  ],
				  "incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				  "title": "Forbidden"
				}`)),
				Request: req,
			},
			expected: &Error{
				Code:       "forbidden",
				Title:      "Forbidden",
				IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				HTTPStatus: http.StatusForbidden,
				Details: []SecondaryError{
					{
						Code:    "invalid.data",
						Message: "User is not authorized"},
				},
			},
		},
		"Not found 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(`
				{
					"code": "not.found",
					"title": "Not Found",
					"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					"details": [
					  {
						"code": "invalid.data",
						"message": "The reporting group with id 000 is already deleted"
  					  }
					]
				}`)),
				Request: req,
			},
			expected: &Error{
				Code:       "not.found",
				Title:      "Not Found",
				IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				HTTPStatus: http.StatusNotFound,
				Details: []SecondaryError{
					{
						Code:    "invalid.data",
						Message: "The reporting group with id 000 is already deleted"},
				},
			},
		},
		"Method not allowed 405": {
			response: &http.Response{
				StatusCode: http.StatusMethodNotAllowed,
				Body: io.NopCloser(strings.NewReader(`
				{
					"code": "method.not.allowed",
					"title": "Method Not Allowed",
					"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					"details": [
					  {
						"code": "method.not.supported",
						"message": "Method not supported"
					  }
					]
				}`)),
				Request: req,
			},
			expected: &Error{
				Code:       "method.not.allowed",
				Title:      "Method Not Allowed",
				HTTPStatus: http.StatusMethodNotAllowed,
				IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				Details: []SecondaryError{
					{
						Code:    "method.not.supported",
						Message: "Method not supported"},
				},
			},
		},
		"Invalid response body": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Code:       "",
				Title:      "Failed to unmarshal error body. Reporting Groups API failed. Check details for more information.",
				IncidentID: "",
				HTTPStatus: http.StatusInternalServerError,
				Details: []SecondaryError{
					{
						Code:    "",
						Message: "test",
					},
				},
			},
		},
		"Empty response body": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			},
			expected: &Error{
				Code:       "",
				Title:      "Failed to unmarshal error body. Reporting Groups API failed. Check details for more information.",
				IncidentID: "",
				HTTPStatus: http.StatusInternalServerError,
				Details: []SecondaryError{
					{
						Code:    "",
						Message: "",
					},
				},
			},
		},
		"Internal server error 500": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(`
				{
					"code": "internal.server.error",
					"title": "Internal Server Error",
					"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
				}`)),
				Request: req,
			},
			expected: &Error{
				Code:       "internal.server.error",
				Title:      "Internal Server Error",
				IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
				HTTPStatus: http.StatusInternalServerError,
				Details:    nil,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*reportinggroups).Error(tc.response)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestIs(t *testing.T) {
	t.Parallel()
	someError := Error{
		Code:       "not.found",
		Title:      "Not Found",
		IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
		HTTPStatus: 404,
		Details: []SecondaryError{
			{
				Code:    "invalid.data",
				Message: "Group Name is mandatory",
			},
		}}

	tests := map[string]struct {
		target   Error
		source   *Error
		expected bool
	}{
		"different error code": {
			target:   Error{Code: "bad.request"},
			expected: false,
		},
		"different error title": {
			target:   Error{Title: "other error"},
			expected: false,
		},
		"same error title": {
			target:   Error{Title: "Not Found"},
			expected: true,
		},
		"same error code": {
			target:   Error{Code: "not.found"},
			expected: true,
		},
		"same error code and title": {
			target:   Error{Code: "not.found", Title: "Not Found"},
			expected: true,
		},
		"same error code but different title": {
			target:   Error{Code: "not.found", Title: "other error"},
			expected: false,
		},
		"same error title but different code": {
			target:   Error{Title: "Not Found", Code: "bad.request"},
			expected: false,
		},
		"both zero values match any": {
			target:   Error{},
			expected: true,
		},
		"same error code and title but different details": {
			target:   Error{Code: "not.found", Title: "Not Found", Details: []SecondaryError{{Code: "other.code", Message: "other message"}}},
			expected: false,
		},
		"same title but different status": {
			target:   Error{Title: "Not Found", HTTPStatus: 400},
			expected: false,
		},
		"same status": {
			target:   Error{HTTPStatus: 404},
			expected: true,
		},
		"different Details length": {
			target:   Error{Details: []SecondaryError{{Code: "invalid.data", Message: "Group Name is mandatory"}, {Code: "other.code", Message: "other message"}}},
			expected: false,
		},
		"same Details length but Code differs": {
			target:   Error{Details: []SecondaryError{{Code: "other.code", Message: "Group Name is mandatory"}}},
			expected: false,
		},
		"same Details length but message differs": {
			target:   Error{Details: []SecondaryError{{Code: "invalid.data", Message: "other message"}}},
			expected: false,
		},
		"equal Details structs": {
			target:   Error{Details: []SecondaryError{{Code: "invalid.data", Message: "Group Name is mandatory"}}},
			expected: true,
		},
		"zero value inside Details Code": {
			target:   Error{Details: []SecondaryError{{Message: "Group Name is mandatory"}}},
			expected: false,
		},
		"zero value inside Details Message": {
			target:   Error{Details: []SecondaryError{{Code: "invalid.data"}}},
			expected: false,
		},
		"zero values in Details": {
			target:   Error{Details: []SecondaryError{{}}},
			expected: false,
		},
		"nil details": {
			target:   Error{Details: nil},
			expected: true,
		},
		"empty Details slice": {
			target:   Error{Details: []SecondaryError{}},
			expected: true,
		},
		"with only HTTPStatus zero": {
			target:   Error{Code: "not.found", Title: "Not Found"},
			expected: true,
		},
		"all fields zero except HTTPStatus": {
			target:   Error{HTTPStatus: 404},
			expected: true,
		},
		"code set to empty string": {
			target:   Error{Code: ""},
			expected: true,
		},
		"empty source matches empty target": {
			source:   &Error{},
			target:   Error{},
			expected: true,
		},
		"empty source matches target with empty code": {
			source:   &Error{},
			target:   Error{Code: ""},
			expected: true,
		},
		"empty source matches target with empty details": {
			source:   &Error{},
			target:   Error{Details: []SecondaryError{}},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			source := test.source
			if source == nil {
				source = &someError
			}
			assert.Equal(t, test.expected, source.Is(&test.target))
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	e := &Error{
		Code:       "bad.request",
		Title:      "Bad Request",
		IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
		HTTPStatus: http.StatusBadRequest,
		Details: []SecondaryError{
			{
				Code:    "invalid.data",
				Message: "Group name is mandatory",
			},
		},
	}
	expected := `API error: 
{
	"code": "bad.request",
	"title": "Bad Request",
	"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566",
	"details": [
		{
			"code": "invalid.data",
			"message": "Group name is mandatory"
		}
	]
}`
	assert.EqualError(t, e, expected)
}
