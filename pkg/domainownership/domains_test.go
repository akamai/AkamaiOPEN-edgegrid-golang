package domainownership

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDomains(t *testing.T) {
	tests := map[string]struct {
		params           ListDomainsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListDomainsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - no arguments, multiple various results": {
			params:         ListDomainsRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "metadata": {
        "hasPrevious": false,
        "hasNext": true,
        "page": 1,
        "pageSize": 10,
        "totalItems": 11
    },
    "domains": [
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom1.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:27:19Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom2.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:58Z", 
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-2",
                    "target": "cname-target-2"
                },
                "txtRecord": {
                    "name": "txt-name-2",
                    "value": "txt-value-2"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom3.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:50Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-3",
                    "target": "cname-target-3"
                },
                "txtRecord": {
                    "name": "txt-name-3",
                    "value": "txt-value-3"
                },
                "httpFile": {
                    "path": "http-file-path-3",
                    "content": "http-file-content-3",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-3",
                    "to": "http-redirect-to-3"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom4.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:50Z", 
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-4",
                    "target": "cname-target-4"
                },
                "txtRecord": {
                    "name": "txt-name-4",
                    "value": "txt-value-4"
                },
                "httpFile": {
                    "path": "http-file-path-4",
                    "content": "http-file-content-4",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-4",
                    "to": "http-redirect-to-4"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom5.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:50Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-5",
                    "target": "cname-target-5"
                },
                "txtRecord": {
                    "name": "txt-name-5",
                    "value": "txt-value-5"
                },
                "httpFile": {
                    "path": "http-file-path-5",
                    "content": "http-file-content-5",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-5",
                    "to": "http-redirect-to-5"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom6.test",
            "validationScope": "WILDCARD",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:50Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-6",
                    "target": "cname-target-6"
                },
                "txtRecord": {
                    "name": "txt-name-6",
                    "value": "txt-value-6"
                },
                "httpFile": {
                    "path": "http-file-path-6",
                    "content": "http-file-content-6",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-6",
                    "to": "http-redirect-to-6"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom7.test",
            "validationScope": "DOMAIN",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:50Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-7",
                    "target": "cname-target-7"
                },
                "txtRecord": {
                    "name": "txt-name-7",
                    "value": "txt-value-7"
                },
                "httpFile": {
                    "path": "http-file-path-7",
                    "content": "http-file-content-7",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-7",
                    "to": "http-redirect-to-7"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom8.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:29Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-8",
                    "target": "cname-target-8"
                },
                "txtRecord": {
                    "name": "txt-name-8",
                    "value": "txt-value-8"
                },
                "httpFile": {
                    "path": "http-file-path-8",
                    "content": "http-file-content-8",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-8",
                    "to": "http-redirect-to-8"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom9.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:29Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-9",
                    "target": "cname-target-9"
                },
                "txtRecord": {
                    "name": "txt-name-9",
                    "value": "txt-value-9"
                },
                "httpFile": {
                    "path": "http-file-path-9",
                    "content": "http-file-content-9",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-9",
                    "to": "http-redirect-to-9"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom10.test",
            "validationScope": "DOMAIN",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:29Z",
           "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-10",
                    "target": "cname-target-10"
                },
                "txtRecord": {
                    "name": "txt-name-10",
                    "value": "txt-value-10"
                },
                "httpFile": {
                    "path": "http-file-path-10",
                    "content": "http-file-content-10",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-10",
                    "to": "http-redirect-to-10"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ],
    "links": [
        {
            "rel": "self",
            "href": "/domain-validation-service/api/v1/domains?page=1&pageSize=10"
        },
        {
            "rel": "next",
            "href": "/domain-validation-service/api/v1/domains?page=2&pageSize=10"
        }
    ]
}`,
			expectedPath: "/domain-validation/v1/domains",
			expectedResponse: &ListDomainsResponse{
				Metadata: Metadata{
					HasPrevious: false,
					HasNext:     true,
					Page:        1,
					PageSize:    10,
					TotalItems:  11,
				},
				Domains: []DomainItem{
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom1.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom2.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:58Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-2",
								Target: "cname-target-2",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-2",
								Value: "txt-value-2",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom3.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:50Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-3",
								Target: "cname-target-3",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-3",
								Value: "txt-value-3",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-3",
								Content:     "http-file-content-3",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-3",
								To:   "http-redirect-to-3",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom4.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:50Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-4",
								Target: "cname-target-4",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-4",
								Value: "txt-value-4",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-4",
								Content:     "http-file-content-4",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-4",
								To:   "http-redirect-to-4",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom5.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:50Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-5",
								Target: "cname-target-5",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-5",
								Value: "txt-value-5",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-5",
								Content:     "http-file-content-5",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-5",
								To:   "http-redirect-to-5",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom6.test",
						ValidationScope:         "WILDCARD",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:50Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-6",
								Target: "cname-target-6",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-6",
								Value: "txt-value-6",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-6",
								Content:     "http-file-content-6",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-6",
								To:   "http-redirect-to-6",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom7.test",
						ValidationScope:         "DOMAIN",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:50Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-7",
								Target: "cname-target-7",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-7",
								Value: "txt-value-7",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-7",
								Content:     "http-file-content-7",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-7",
								To:   "http-redirect-to-7",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom8.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:29Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-8",
								Target: "cname-target-8",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-8",
								Value: "txt-value-8",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-8",
								Content:     "http-file-content-8",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-8",
								To:   "http-redirect-to-8",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom9.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:29Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-9",
								Target: "cname-target-9",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-9",
								Value: "txt-value-9",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-9",
								Content:     "http-file-content-9",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-9",
								To:   "http-redirect-to-9",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom10.test",
						ValidationScope:         "DOMAIN",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:29Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-10",
								Target: "cname-target-10",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-10",
								Value: "txt-value-10",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-10",
								Content:     "http-file-content-10",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-10",
								To:   "http-redirect-to-10",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "/domain-validation-service/api/v1/domains?page=1&pageSize=10",
					},
					{
						Rel:  "next",
						Href: "/domain-validation-service/api/v1/domains?page=2&pageSize=10",
					},
				},
			},
		},
		"200 OK - explicit page and pageSize": {
			params:         ListDomainsRequest{Page: 1, PageSize: 10},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "metadata": {
        "hasPrevious": false,
        "hasNext": false,
        "page": 1,
        "pageSize": 10,
        "totalItems": 1
    },
    "domains": [
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom1.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:27:19Z",
             "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ],
    "links": [
        {
            "rel": "self",
            "href": "/domain-validation-service/api/v1/domains?page=1&pageSize=10"
        }
    ]
}`,
			expectedPath: "/domain-validation/v1/domains?page=1&pageSize=10",
			expectedResponse: &ListDomainsResponse{
				Metadata: Metadata{
					HasPrevious: false,
					HasNext:     false,
					Page:        1,
					PageSize:    10,
					TotalItems:  1,
				},
				Domains: []DomainItem{
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom1.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "/domain-validation-service/api/v1/domains?page=1&pageSize=10",
					},
				},
			},
		},
		"200 OK - explicit paginate, page and pageSize": {
			params: ListDomainsRequest{
				Paginate: ptr.To(true),
				Page:     1,
				PageSize: 10,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "metadata": {
        "hasPrevious": false,
        "hasNext": false,
        "page": 1,
        "pageSize": 10,
        "totalItems": 1
    },
    "domains": [
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom9.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:26:29Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-9",
                    "target": "cname-target-9"
                },
                "txtRecord": {
                    "name": "txt-name-9",
                    "value": "txt-value-9"
                },
                "httpFile": {
                    "path": "http-file-path-9",
                    "content": "http-file-content-9",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-9",
                    "to": "http-redirect-to-9"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ],
    "links": [
        {
            "rel": "self",
            "href": "/domain-validation-service/api/v1/domains?page=1&pageSize=10"
        }
    ]
}`,
			expectedPath: "/domain-validation/v1/domains?page=1&pageSize=10&paginate=true",
			expectedResponse: &ListDomainsResponse{
				Metadata: Metadata{
					HasPrevious: false,
					HasNext:     false,
					Page:        1,
					PageSize:    10,
					TotalItems:  1,
				},
				Domains: []DomainItem{
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom9.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:26:29Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-9",
								Target: "cname-target-9",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-9",
								Value: "txt-value-9",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-9",
								Content:     "http-file-content-9",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-9",
								To:   "http-redirect-to-9",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "/domain-validation-service/api/v1/domains?page=1&pageSize=10",
					},
				},
			},
		},
		"200 OK - only paginate": {
			params: ListDomainsRequest{
				Paginate: ptr.To(true),
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "metadata": {
        "hasPrevious": false,
        "hasNext": false,
        "page": 1,
        "pageSize": 10,
        "totalItems": 1
    },
    "domains": [
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom1.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:27:19Z",
             "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ],
    "links": [
        {
            "rel": "self",
            "href": "/domain-validation-service/api/v1/domains?page=1&pageSize=10"
        }
    ]
}`,
			expectedPath: "/domain-validation/v1/domains?paginate=true",
			expectedResponse: &ListDomainsResponse{
				Metadata: Metadata{
					HasPrevious: false,
					HasNext:     false,
					Page:        1,
					PageSize:    10,
					TotalItems:  1,
				},
				Domains: []DomainItem{
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom1.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "/domain-validation-service/api/v1/domains?page=1&pageSize=10",
					},
				},
			},
		},
		"200 OK - only pageSize": {
			params: ListDomainsRequest{
				PageSize: 10,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "metadata": {
        "hasPrevious": false,
        "hasNext": false,
        "page": 1,
        "pageSize": 10,
        "totalItems": 1
    },
    "domains": [
        {
            "accountId": "1-ACCOUN",
            "domainName": "dom1.test",
            "validationScope": "HOST",
            "domainStatus": "REQUEST_ACCEPTED",
            "validationRequestedBy": "someuser",
            "validationRequestedDate": "2025-08-04T13:27:19Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ],
    "links": [
        {
            "rel": "self",
            "href": "/domain-validation-service/api/v1/domains?page=1&pageSize=10"
        }
    ]
}`,
			expectedPath: "/domain-validation/v1/domains?pageSize=10",
			expectedResponse: &ListDomainsResponse{
				Metadata: Metadata{
					HasPrevious: false,
					HasNext:     false,
					Page:        1,
					PageSize:    10,
					TotalItems:  1,
				},
				Domains: []DomainItem{
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom1.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "/domain-validation-service/api/v1/domains?page=1&pageSize=10",
					},
				},
			},
		},
		"200 OK - only page": {
			params: ListDomainsRequest{
				Page: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
		{
		    "metadata": {
		        "hasPrevious": false,
		        "hasNext": false,
		        "page": 1,
		        "pageSize": 10,
		        "totalItems": 1
		    },
		    "domains": [
		        {
		            "accountId": "1-ACCOUN",
		            "domainName": "dom1.test",
		            "validationScope": "HOST",
		            "domainStatus": "REQUEST_ACCEPTED",
		            "validationRequestedBy": "someuser",
		            "validationRequestedDate": "2025-08-04T13:27:19Z",
		            "validationChallenge": {
                		"cnameRecord": {
                    		"name": "cname-name-1",
                    		"target": "cname-target-1"
               		 	},
                		"txtRecord": {
                    		"name": "txt-name-1",
                    		"value": "txt-value-1"
                		},
                		"httpFile": {
                    		"path": "http-file-path-1",
                    		"content": "http-file-content-1",
                    		"contentType": "text/plain"
                		},
                		"httpRedirect": {
                    		"from": "http-redirect-from-1",
                    		"to": "http-redirect-to-1"
                		},
                		"expirationDate": "2025-08-05T13:27:19Z"
            		}
		        }
		    ],
		    "links": [
		        {
		            "rel": "self",
		            "href": "/domain-validation-service/api/v1/domains?page=1&pageSize=10"
		        }
		    ]
		}`,
			expectedPath: "/domain-validation/v1/domains?page=1",
			expectedResponse: &ListDomainsResponse{
				Metadata: Metadata{
					HasPrevious: false,
					HasNext:     false,
					Page:        1,
					PageSize:    10,
					TotalItems:  1,
				},
				Domains: []DomainItem{
					{
						AccountID:               "1-ACCOUN",
						DomainName:              "dom1.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationRequestedBy:   "someuser",
						ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "/domain-validation-service/api/v1/domains?page=1&pageSize=10",
					},
				},
			},
		},
		"validation - page or pageSize without paging": {
			params: ListDomainsRequest{
				Paginate: ptr.To(false),
				Page:     1,
				PageSize: 10,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list domains: struct validation:\nPage: must be 0 when Paginate is false\nPageSize: must be 0 when Paginate is false", err.Error())
			},
		},
		"validation - pageSize too small": {
			params: ListDomainsRequest{
				PageSize: 1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list domains: struct validation:\nPageSize: must be no less than 10", err.Error())
			},
		},
		"validation - pageSize too big": {
			params: ListDomainsRequest{
				PageSize: 1001,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list domains: struct validation:\nPageSize: must be no greater than 1000", err.Error())
			},
		},
		"500 internal server error": {
			params:         ListDomainsRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListDomains(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetDomain(t *testing.T) {
	tests := map[string]struct {
		params           GetDomainRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetDomainResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - only required arguments - not validated": {
			params: GetDomainRequest{
				DomainName:      "dom1.test",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "1-ACCOUN",
    "domainName": "dom1.test",
    "validationScope": "HOST",
    "domainStatus": "REQUEST_ACCEPTED",
    "validationRequestedBy": "someuser",
    "validationRequestedDate": "2025-08-04T13:27:19Z",
    "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
}`,
			expectedPath: "/domain-validation/v1/domains/dom1.test?validationScope=HOST",
			expectedResponse: &GetDomainResponse{
				AccountID:               "1-ACCOUN",
				DomainName:              "dom1.test",
				ValidationScope:         "HOST",
				DomainStatus:            "REQUEST_ACCEPTED",
				ValidationRequestedBy:   "someuser",
				ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
				ValidationChallenge: &ValidationChallenge{
					CnameRecord: CnameRecord{
						Name:   "cname-name-1",
						Target: "cname-target-1",
					},
					TXTRecord: TXTRecord{
						Name:  "txt-name-1",
						Value: "txt-value-1",
					},
					HTTPFile: &HTTPFile{
						Path:        "http-file-path-1",
						Content:     "http-file-content-1",
						ContentType: "text/plain",
					},
					HTTPRedirect: &HTTPRedirect{
						From: "http-redirect-from-1",
						To:   "http-redirect-to-1",
					},
					ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
				},
			},
		},
		"200 OK - only required arguments - not validated and minimal challenges": {
			params: GetDomainRequest{
				DomainName:      "dom1.test",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "1-ACCOUN",
    "domainName": "dom1.test",
    "validationScope": "HOST",
    "domainStatus": "REQUEST_ACCEPTED",
    "validationRequestedBy": "someuser",
    "validationRequestedDate": "2025-08-04T13:27:19Z",
    "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
}`,
			expectedPath: "/domain-validation/v1/domains/dom1.test?validationScope=HOST",
			expectedResponse: &GetDomainResponse{
				AccountID:               "1-ACCOUN",
				DomainName:              "dom1.test",
				ValidationScope:         "HOST",
				DomainStatus:            "REQUEST_ACCEPTED",
				ValidationRequestedBy:   "someuser",
				ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
				ValidationChallenge: &ValidationChallenge{
					CnameRecord: CnameRecord{
						Name:   "cname-name-1",
						Target: "cname-target-1",
					},
					TXTRecord: TXTRecord{
						Name:  "txt-name-1",
						Value: "txt-value-1",
					},
					HTTPFile: &HTTPFile{
						Path:        "http-file-path-1",
						Content:     "http-file-content-1",
						ContentType: "text/plain",
					},
					HTTPRedirect: &HTTPRedirect{
						From: "http-redirect-from-1",
						To:   "http-redirect-to-1",
					},
					ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
				},
			},
		},
		"200 OK - with status history - not validated": {
			params: GetDomainRequest{
				DomainName:                 "dom1.test",
				ValidationScope:            ValidationScopeHost,
				IncludeDomainStatusHistory: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "1-ACCOUN",
    "domainName": "dom1.test",
    "validationScope": "HOST",
    "domainStatus": "VALIDATION_IN_PROGRESS",
    "validationRequestedBy": "someuser",
    "validationRequestedDate": "2025-08-04T13:27:19Z",
    "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            },
    "domainStatusHistory": [
        {
            "domainStatus": "REQUEST_ACCEPTED",
            "modifiedUser": "someuser",
            "modifiedDate": "2025-08-04T11:49:53Z"
        },
        {
            "domainStatus": "VALIDATION_IN_PROGRESS",
            "message": "DNS verification failed.",
            "modifiedUser": "someuser",
            "modifiedDate": "2025-08-04T11:50:53Z"
        }
    ]
}`,
			expectedPath: "/domain-validation/v1/domains/dom1.test?includeDomainStatusHistory=true&validationScope=HOST",
			expectedResponse: &GetDomainResponse{
				AccountID:               "1-ACCOUN",
				DomainName:              "dom1.test",
				ValidationScope:         "HOST",
				DomainStatus:            "VALIDATION_IN_PROGRESS",
				ValidationRequestedBy:   "someuser",
				ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
				ValidationChallenge: &ValidationChallenge{
					CnameRecord: CnameRecord{
						Name:   "cname-name-1",
						Target: "cname-target-1",
					},
					TXTRecord: TXTRecord{
						Name:  "txt-name-1",
						Value: "txt-value-1",
					},
					HTTPFile: &HTTPFile{
						Path:        "http-file-path-1",
						Content:     "http-file-content-1",
						ContentType: "text/plain",
					},
					HTTPRedirect: &HTTPRedirect{
						From: "http-redirect-from-1",
						To:   "http-redirect-to-1",
					},
					ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
				},
				DomainStatusHistory: []DomainStatusHistory{
					{
						DomainStatus: "REQUEST_ACCEPTED",
						ModifiedUser: "someuser",
						ModifiedDate: test.NewTimeFromString(t, "2025-08-04T11:49:53Z"),
					},
					{
						DomainStatus: "VALIDATION_IN_PROGRESS",
						Message:      ptr.To("DNS verification failed."),
						ModifiedUser: "someuser",
						ModifiedDate: test.NewTimeFromString(t, "2025-08-04T11:50:53Z"),
					},
				},
			},
		},
		"200 OK - validated": {
			params: GetDomainRequest{
				DomainName:      "dom1.test",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "1-ACCOUN",
	"domainName": "dom1.test",
	"validationScope": "HOST",
	"domainStatus": "VALIDATED",
	"validationMethod": "SYSTEM",
	"validationRequestedBy": "someuser",
	"validationRequestedDate": "2025-08-04T13:27:19Z",
	"validationCompletedDate": "2025-08-05T11:56:07Z"
}`,
			expectedPath: "/domain-validation/v1/domains/dom1.test?validationScope=HOST",
			expectedResponse: &GetDomainResponse{
				AccountID:               "1-ACCOUN",
				DomainName:              "dom1.test",
				ValidationScope:         "HOST",
				DomainStatus:            "VALIDATED",
				ValidationMethod:        ptr.To("SYSTEM"),
				ValidationRequestedBy:   "someuser",
				ValidationRequestedDate: test.NewTimeFromString(t, "2025-08-04T13:27:19Z"),
				ValidationCompletedDate: ptr.To(test.NewTimeFromString(t, "2025-08-05T11:56:07Z")),
			},
		},
		"validation - no arguments": {
			params: GetDomainRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get domain: struct validation:\nDomainName: cannot be blank\nValidationScope: cannot be blank", err.Error())
			},
		},
		"validation - incorrect ValidationScope ": {
			params: GetDomainRequest{
				DomainName:      "dom1.test",
				ValidationScope: ValidationScope("incorrect"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get domain: struct validation:\nValidationScope: value 'incorrect' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'", err.Error())
			},
		},
		"500 internal server error": {
			params: GetDomainRequest{
				DomainName:      "dom1.test",
				ValidationScope: ValidationScopeDomain,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains/dom1.test?validationScope=DOMAIN",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetDomain(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestSearchDomains(t *testing.T) {
	tests := map[string]struct {
		params              SearchDomainsRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *SearchDomainsResponse
		withError           func(*testing.T, error)
	}{
		"200 OK - several different elements in search without details": {
			params: SearchDomainsRequest{
				Body: SearchDomainsBody{Domains: []Domain{
					{
						DomainName:      "dom1.test",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "dom2.test",
						ValidationScope: ValidationScopeDomain,
					},
				}},
			},
			expectedRequestBody: `{"domains":[{"domainName":"dom1.test","validationScope":"HOST"},{"domainName":"dom2.test","validationScope":"DOMAIN"}]}`,
			responseStatus:      http.StatusOK,
			responseBody: `
{
	"domains": [
		{
			"domainName": "dom1.test",
			"validationScope": "HOST",
			"domainStatus": "REQUEST_ACCEPTED",
			"validationLevel": "FQDN"
		},
		{
			"domainName": "dom2.test",
			"validationScope": "HOST",
			"domainStatus": "VALIDATED",
			"validationLevel": "FQDN"
		}
	]
}`,
			expectedPath: "/domain-validation/v1/domains/search",
			expectedResponse: &SearchDomainsResponse{
				Domains: []SearchDomainItem{
					{
						DomainName:      "dom1.test",
						ValidationScope: "HOST",
						DomainStatus:    "REQUEST_ACCEPTED",
						ValidationLevel: "FQDN",
					},
					{
						DomainName:      "dom2.test",
						ValidationScope: "HOST",
						DomainStatus:    "VALIDATED",
						ValidationLevel: "FQDN",
					},
				},
			},
		},
		"200 OK - several different elements in search with details": {
			params: SearchDomainsRequest{
				IncludeAll: true,
				Body: SearchDomainsBody{Domains: []Domain{
					{
						DomainName:      "dom1.test",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "dom2.test",
						ValidationScope: ValidationScopeDomain,
					},
				}},
			},
			expectedRequestBody: `{"domains":[{"domainName":"dom1.test","validationScope":"HOST"},{"domainName":"dom2.test","validationScope":"DOMAIN"}]}`,
			responseStatus:      http.StatusOK,
			responseBody: `
{
	"domains": [
		{
			"accountId": "1-ACCOUN",
			"domainName": "dom1.test",
			"validationScope": "HOST",
			"domainStatus": "REQUEST_ACCEPTED",
			"validationLevel": "FQDN",
			"validationMethod": "DNS_TXT",
			"validationRequestedBy": "someuser",
			"validationRequestedDate": "2025-08-04T13:27:19Z",
			"validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
		},
		{
			"accountId": "1-ACCOUN",
			"domainName": "dom2.test",
			"validationScope": "HOST",
			"domainStatus": "VALIDATED",
			"validationLevel": "FQDN",
			"validationMethod": "SYSTEM",
			"validationRequestedBy": "someuser",
			"validationRequestedDate": "2025-08-04T13:27:19Z",
			"validationCompletedDate": "2025-08-05T11:56:07Z"
		}
	]
}`,
			expectedPath: "/domain-validation/v1/domains/search?includeAll=true",
			expectedResponse: &SearchDomainsResponse{
				Domains: []SearchDomainItem{
					{
						AccountID:               ptr.To("1-ACCOUN"),
						DomainName:              "dom1.test",
						ValidationScope:         "HOST",
						DomainStatus:            "REQUEST_ACCEPTED",
						ValidationLevel:         "FQDN",
						ValidationMethod:        ptr.To("DNS_TXT"),
						ValidationRequestedBy:   ptr.To("someuser"),
						ValidationRequestedDate: ptr.To(test.NewTimeFromString(t, "2025-08-04T13:27:19Z")),
						ValidationChallenge: &ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{
						AccountID:               ptr.To("1-ACCOUN"),
						DomainName:              "dom2.test",
						ValidationScope:         "HOST",
						DomainStatus:            "VALIDATED",
						ValidationLevel:         "FQDN",
						ValidationMethod:        ptr.To("SYSTEM"),
						ValidationRequestedBy:   ptr.To("someuser"),
						ValidationRequestedDate: ptr.To(test.NewTimeFromString(t, "2025-08-04T13:27:19Z")),
						ValidationCompletedDate: ptr.To(test.NewTimeFromString(t, "2025-08-05T11:56:07Z")),
					},
				},
			},
		},
		"validation - no arguments": {
			params: SearchDomainsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "search domains: struct validation:\nBody: {\n\tDomains: cannot be blank\n}", err.Error())
			},
		},
		"validation - empty domain": {
			params: SearchDomainsRequest{
				Body: SearchDomainsBody{Domains: []Domain{
					{},
				}},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "search domains: struct validation:\nBody: {\n\tDomains[0]: {\n\t\tDomainName: cannot be blank\n\t\tValidationScope: cannot be blank\n\t}\n}", err.Error())
			},
		},
		"validation - incorrect ValidationScope": {
			params: SearchDomainsRequest{
				Body: SearchDomainsBody{Domains: []Domain{
					{
						DomainName:      "dom1.test",
						ValidationScope: ValidationScope("incorrect"),
					},
				}},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "search domains: struct validation:\nBody: {\n\tDomains[0]: {\n\t\tValidationScope: value 'incorrect' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'\n\t}\n}", err.Error())
			},
		},
		"500 internal server error": {
			params: SearchDomainsRequest{
				Body: SearchDomainsBody{
					Domains: []Domain{
						{
							DomainName:      "dom1.test",
							ValidationScope: ValidationScopeDomain,
						},
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains/search",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.SearchDomains(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestAddDomains(t *testing.T) {
	tests := map[string]struct {
		request             AddDomainsRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *AddDomainsResponse
		withError           func(*testing.T, error)
	}{
		"207 All Success": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "sample2.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			responseStatus: http.StatusMultiStatus,
			responseBody: `{
    "errors": [],
    "successes": [
        {
            "domainName": "sample1.com",
            "domainStatus": "REQUEST_ACCEPTED",
            "accountId": "A-CCT5678",
            "validationScope": "HOST",
            "validationMethod": "DNS_TXT",
            "validationRequestedBy": "someone",
            "validationRequestedDate": "2024-02-06T06:01:45Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "domainName": "sample2.com",
            "domainStatus": "REQUEST_ACCEPTED",
            "accountId": "A-CCT7890",
            "validationScope": "HOST",
            "validationRequestedBy": "someone",
            "validationRequestedDate": "2024-02-06T06:01:45Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-2",
                    "target": "cname-target-2"
                },
                "txtRecord": {
                    "name": "txt-name-2",
                    "value": "txt-value-2"
                },
                "httpFile": {
                    "path": "http-file-path-2",
                    "content": "http-file-content-2",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-2",
                    "to": "http-redirect-to-2"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ]
}`,
			expectedPath:        "/domain-validation/v1/domains",
			expectedRequestBody: `{"domains":[{"domainName":"sample1.com","validationScope":"HOST"},{"domainName":"sample2.com","validationScope":"HOST"}]}`,
			expectedResponse: &AddDomainsResponse{
				Errors: []AddDomainError{},
				Successes: []AddDomainSuccess{
					{
						DomainName:              "sample1.com",
						DomainStatus:            "REQUEST_ACCEPTED",
						AccountID:               "A-CCT5678",
						ValidationScope:         "HOST",
						ValidationMethod:        ptr.To("DNS_TXT"),
						ValidationRequestedBy:   "someone",
						ValidationRequestedDate: test.NewTimeFromString(t, "2024-02-06T06:01:45Z"),
						ValidationChallenge: ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{

						DomainName:              "sample2.com",
						DomainStatus:            "REQUEST_ACCEPTED",
						AccountID:               "A-CCT7890",
						ValidationScope:         "HOST",
						ValidationRequestedBy:   "someone",
						ValidationRequestedDate: test.NewTimeFromString(t, "2024-02-06T06:01:45Z"),
						ValidationChallenge: ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-2",
								Target: "cname-target-2",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-2",
								Value: "txt-value-2",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-2",
								Content:     "http-file-content-2",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-2",
								To:   "http-redirect-to-2",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
			},
		},
		"207 Partial Success": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "sample2.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "sample3.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "sample4.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "sample5.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			responseStatus: http.StatusMultiStatus,
			responseBody: `{
    "errors": [
        {
            "domainName": "sample3.com",
            "detail": "Domain already exists.",
            "title": "Internal Server Error",
            "type": "internal-server-error",
            "validationScope": "HOST"
        },
        {
            "domainName": "sample4.com",
            "detail": "Supernet domain has been validated and is ready for use.",
            "title": "Internal Server Error",
            "type": "internal-server-error",
            "validationScope": "HOST"
        },
        {
            "domainName": "sample5.com",
            "detail": "Domain is already in use within the system. You cannot use this domain.",
            "title": "Internal Server Error",
            "type": "internal-server-error",
            "validationScope": "HOST"
        }
    ],
    "successes": [
        {
            "domainName": "sample1.com",
            "domainStatus": "REQUEST_ACCEPTED",
            "accountId": "A-CCT5678",
            "validationScope": "HOST",
            "validationMethod": "DNS_TXT",
            "validationRequestedBy": "someone",
            "validationRequestedDate": "2024-02-06T06:01:45Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-1",
                    "target": "cname-target-1"
                },
                "txtRecord": {
                    "name": "txt-name-1",
                    "value": "txt-value-1"
                },
                "httpFile": {
                    "path": "http-file-path-1",
                    "content": "http-file-content-1",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-1",
                    "to": "http-redirect-to-1"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        },
        {
            "domainName": "sample2.com",
            "domainStatus": "REQUEST_ACCEPTED",
            "accountId": "A-CCT7890",
            "validationScope": "HOST",
            "validationRequestedBy": "someone",
            "validationRequestedDate": "2024-02-06T06:01:45Z",
            "validationChallenge": {
                "cnameRecord": {
                    "name": "cname-name-2",
                    "target": "cname-target-2"
                },
                "txtRecord": {
                    "name": "txt-name-2",
                    "value": "txt-value-2"
                },
                "httpFile": {
                    "path": "http-file-path-2",
                    "content": "http-file-content-2",
                    "contentType": "text/plain"
                },
                "httpRedirect": {
                    "from": "http-redirect-from-2",
                    "to": "http-redirect-to-2"
                },
                "expirationDate": "2025-08-05T13:27:19Z"
            }
        }
    ]
}`,
			expectedPath:        "/domain-validation/v1/domains",
			expectedRequestBody: `{"domains":[{"domainName":"sample1.com","validationScope":"HOST"},{"domainName":"sample2.com","validationScope":"HOST"},{"domainName":"sample3.com","validationScope":"HOST"},{"domainName":"sample4.com","validationScope":"HOST"},{"domainName":"sample5.com","validationScope":"HOST"}]}`,
			expectedResponse: &AddDomainsResponse{
				Errors: []AddDomainError{
					{
						DomainName:      "sample3.com",
						Detail:          "Domain already exists.",
						Title:           "Internal Server Error",
						Type:            "internal-server-error",
						ValidationScope: "HOST",
					},
					{
						DomainName:      "sample4.com",
						Detail:          "Supernet domain has been validated and is ready for use.",
						Title:           "Internal Server Error",
						Type:            "internal-server-error",
						ValidationScope: "HOST",
					},
					{
						DomainName:      "sample5.com",
						Detail:          "Domain is already in use within the system. You cannot use this domain.",
						Title:           "Internal Server Error",
						Type:            "internal-server-error",
						ValidationScope: "HOST",
					},
				},
				Successes: []AddDomainSuccess{
					{
						DomainName:              "sample1.com",
						DomainStatus:            "REQUEST_ACCEPTED",
						AccountID:               "A-CCT5678",
						ValidationScope:         "HOST",
						ValidationMethod:        ptr.To("DNS_TXT"),
						ValidationRequestedBy:   "someone",
						ValidationRequestedDate: test.NewTimeFromString(t, "2024-02-06T06:01:45Z"),
						ValidationChallenge: ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-1",
								Target: "cname-target-1",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-1",
								Value: "txt-value-1",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-1",
								Content:     "http-file-content-1",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-1",
								To:   "http-redirect-to-1",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
					{

						DomainName:              "sample2.com",
						DomainStatus:            "REQUEST_ACCEPTED",
						AccountID:               "A-CCT7890",
						ValidationScope:         "HOST",
						ValidationRequestedBy:   "someone",
						ValidationRequestedDate: test.NewTimeFromString(t, "2024-02-06T06:01:45Z"),
						ValidationChallenge: ValidationChallenge{
							CnameRecord: CnameRecord{
								Name:   "cname-name-2",
								Target: "cname-target-2",
							},
							TXTRecord: TXTRecord{
								Name:  "txt-name-2",
								Value: "txt-value-2",
							},
							HTTPFile: &HTTPFile{
								Path:        "http-file-path-2",
								Content:     "http-file-content-2",
								ContentType: "text/plain",
							},
							HTTPRedirect: &HTTPRedirect{
								From: "http-redirect-from-2",
								To:   "http-redirect-to-2",
							},
							ExpirationDate: test.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
						},
					},
				},
			},
		},
		"validation - empty domain": {
			request: AddDomainsRequest{
				Domains: []Domain{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains: cannot be blank", err.Error())
			},
		},
		"validation - domain Name not supplied": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains[1]: {\n\tDomainName: cannot be blank\n}\nHint: Domain must: not be empty, not begin with '*', not begin or end with whitespace, and not exceed 200 characters", err.Error())
			},
		},
		"validation - validation scope not supplied": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName: "sample1.com",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains[0]: {\n\tValidationScope: cannot be blank\n}", err.Error())
			},
		},
		"validation - domain Name cannot start with `*`": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "*sample1.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains[0]: {\n\tDomainName: domain '*sample1.com': invalid name format\n}\nHint: Domain must: not be empty, not begin with '*', not begin or end with whitespace, and not exceed 200 characters", err.Error())
			},
		},
		"validation - domain Name does not match the allowed format": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "*example.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains[0]: {\n\tDomainName: domain '*example.com': invalid name format\n}\nHint: Domain must: not be empty, not begin with '*', not begin or end with whitespace, and not exceed 200 characters", err.Error())
			},
		},
		"validation - domain Name greater than 200 characters": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com" + strings.Repeat("a", 190),
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains[0]: {\n\tDomainName: domain 'sample1.comaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa': cannot exceed 200 characters\n}\nHint: Domain must: not be empty, not begin with '*', not begin or end with whitespace, and not exceed 200 characters", err.Error())
			},
		},
		"validation - incorrect ValidationScope": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScope("incorrect"),
					},
					{
						DomainName:      "sample2.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "add domains: struct validation:\nDomains[0]: {\n\tValidationScope: value 'incorrect' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'\n}", err.Error())
			},
		},
		"500 internal server error": {
			request: AddDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeDomain,
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.AddDomains(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeleteDomain(t *testing.T) {
	tests := map[string]struct {
		params              DeleteDomainRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params: DeleteDomainRequest{
				DomainName:      "sample1.com",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/domain-validation/v1/domains/sample1.com?validationScope=HOST",
		},
		"validation errors": {
			params: DeleteDomainRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domain: struct validation:\nDomainName: cannot be blank\nValidationScope: cannot be blank", err.Error())
			},
		},
		"validation errors - DomainName missing ": {
			params: DeleteDomainRequest{
				ValidationScope: ValidationScopeHost,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domain: struct validation:\nDomainName: cannot be blank", err.Error())
			},
		},
		"validation errors - ValidationScope missing ": {
			params: DeleteDomainRequest{
				DomainName: "sample1.com",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domain: struct validation:\nValidationScope: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: DeleteDomainRequest{
				DomainName:      "sample1.com",
				ValidationScope: ValidationScopeHost,
			},
			expectedPath:   "/domain-validation/v1/domains/sample1.com?validationScope=HOST",
			responseStatus: http.StatusNotFound,
			responseBody: `
						{
						    "detail": "Domain is not found.",
							"instance": "55f55b02-bfac-4654-91f6-f72626839bb3",
				  			"status": 404,
				  			"title": "Not Found",
				  			"type": "not-found"
						}
						`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:    "Not Found",
					Type:     "not-found",
					Status:   http.StatusNotFound,
					Instance: "55f55b02-bfac-4654-91f6-f72626839bb3",
					Detail:   "Domain is not found.",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: DeleteDomainRequest{
				DomainName:      "sample1.com",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
						{
						   "type": "internal_error",
						   "title": "Internal Server Error",
						   "detail": "Error making request",
						   "status": 500
						}
						`,
			expectedPath: "/domain-validation/v1/domains/sample1.com?validationScope=HOST",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.DeleteDomain(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDeleteDomains(t *testing.T) {
	tests := map[string]struct {
		params              DeleteDomainsRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			responseStatus:      http.StatusNoContent,
			expectedPath:        "/domain-validation/v1/domains",
			expectedRequestBody: `{"domains":[{"domainName":"sample1.com","validationScope":"HOST"}]}`,
		},
		"validation errors - empty params": {
			params: DeleteDomainsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: Domains: cannot be blank", err.Error())
			},
		},
		"validation errors - empty domains": {
			params: DeleteDomainsRequest{
				Domains: []Domain{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: Domains: cannot be blank", err.Error())
			},
		},
		"validation errors - empty domain": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: Domains[0]: {\n\tDomainName: cannot be blank\n\tValidationScope: cannot be blank\n}", err.Error())
			},
		},
		"validation errors - DomainName missing ": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: Domains[0]: {\n\tDomainName: cannot be blank\n}", err.Error())
			},
		},
		"validation errors - ValidationScope missing": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{
						DomainName: "sample1.com",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: Domains[0]: {\n\tValidationScope: cannot be blank\n}", err.Error())
			},
		},
		"validation errors - invalid ValidationScope": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: "foo",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: Domains[0]: {\n\tValidationScope: value 'foo' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'\n}", err.Error())
			},
		},
		"400 Domain Not Found": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			expectedPath:   "/domain-validation/v1/domains",
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
					"type": "bad-request",
					"title": "Bad Request",
					"instance": "12345-c988-463e-9382-6c15e80868c0",
					"status": 400,
					"detail": "Oops, something wasn't right. Please correct the errors.",
					"errors": [
						{
							"type": "error-types/invalid",
							"title": "Invalid Check",
							"detail": "Domain is not found.",
							"problemId": "0030b473-0bd9-40eb-8afa-d6a08c9be687",
							"field": "domains[0].domainName"
						}
					],
					"problemId": "85588d59-c988-463e-9382-6c15e80868c0"
				}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:     "Bad Request",
					Type:      "bad-request",
					Status:    http.StatusBadRequest,
					Instance:  "12345-c988-463e-9382-6c15e80868c0",
					Detail:    "Oops, something wasn't right. Please correct the errors.",
					ProblemID: "85588d59-c988-463e-9382-6c15e80868c0",
					Errors: []ErrorDetail{
						{
							Type:      "error-types/invalid",
							Title:     "Invalid Check",
							Detail:    "Domain is not found.",
							ProblemID: "0030b473-0bd9-40eb-8afa-d6a08c9be687",
							Field:     "domains[0].domainName",
						},
					},
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: DeleteDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
						{
						   "type": "internal_error",
						   "title": "Internal Server Error",
						   "detail": "Error making request",
						   "status": 500
						}
						`,
			expectedPath: "/domain-validation/v1/domains",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.DeleteDomains(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
