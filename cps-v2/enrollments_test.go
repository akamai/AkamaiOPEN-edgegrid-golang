package cps

import (
  "testing"
  "github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestCPS_NewEnrollments(t *testing.T) {
  e := NewEnrollments()
  assert.IsType(t, &Enrollments{}, e)
}

func TestCPS_GetEnrollments(t *testing.T) {
  defer gock.Off()
  mock := gock.New("https://test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/cps/v2/enrollments")
	mock.
		Get("/cps/v2/enrollments").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
    "enrollments": [
        {
            "networkConfiguration": {
                "preferredCiphers": "ak-akamai-default",
                "sni": {
                    "cloneDnsNames": false,
                    "dnsNames": [
                        "san2.example.com",
                        "san1.example.com"
                    ]
                },
                "mustHaveCiphers": "ak-akamai-default2016q3",
                "disallowedTlsVersions": [],
                "secureNetwork": "enhanced-tls",
                "geography": "core"
            },
            "enableMultiStackedCertificates": false,
            "pendingChanges": [],
            "thirdParty": {
                "excludeSans": false
            },
            "validationType": "third-party",
            "techContact": {
                "city": "Cambridge",
                "organizationName": "Akamai",
                "firstName": "R2",
                "addressLineTwo": null,
                "title": "Astromech Droid",
                "addressLineOne": "150 Broadway",
                "lastName": "D2",
                "region": "MA",
                "phone": "617-555-0111",
                "country": "US",
                "postalCode": "02142",
                "email": "r2d2@akamai.com"
            },
            "changeManagement": true,
            "location": "/cps/v2/enrollments/10002",
            "ra": "third-party",
            "adminContact": {
                "city": "Cambridge",
                "organizationName": "Dark Side",
                "firstName": "Darth",
                "addressLineTwo": null,
                "title": "Lord",
                "addressLineOne": "666 Evil Way",
                "lastName": "Vader",
                "region": "MA",
                "phone": "617-555-0123",
                "country": "US",
                "postalCode": "02142",
                "email": "vader@example.com"
            },
            "certificateChainType": "default",
            "org": {
                "city": "Cambridge",
                "name": "Akamai Technologies",
                "addressLineTwo": null,
                "addressLineOne": "150 Broadway",
                "country": "US",
                "region": "MA",
                "phone": "617-555-0111",
                "postalCode": "02142"
            },
            "certificateType": "third-party",
            "csr": {
                "c": "US",
                "cn": "www.example.com",
                "l": "Cambridge",
                "o": "Akamai",
                "st": "MA",
                "ou": "WebEx",
                "sans": [
                    "san1.example.com",
                    "san2.example.com",
                    "san3.example.com",
                    "san4.example.com",
                    "www.example.com"
                ]
            },
            "signatureAlgorithm": null
        },
        {
            "networkConfiguration": {
                "preferredCiphers": "ak-akamai-default",
                "sni": {
                    "cloneDnsNames": false,
                    "dnsNames": [
                        "san1.example.com",
                        "san2.example.com"
                    ]
                },
                "mustHaveCiphers": "ak-akamai-default2016q3",
                "disallowedTlsVersions": [],
                "secureNetwork": "enhanced-tls",
                "geography": "core"
            },
            "enableMultiStackedCertificates": false,
            "pendingChanges": [
                "/cps/v2/enrollments/10003/changes/10004"
            ],
            "thirdParty": {
                "excludeSans": false
            },
            "validationType": "third-party",
            "techContact": {
                "city": "Cambridge",
                "organizationName": "Akamai",
                "firstName": "R2",
                "addressLineTwo": null,
                "title": "Astromech Droid",
                "addressLineOne": "150 Broadway",
                "lastName": "D2",
                "region": "MA",
                "phone": "617-555-0111",
                "country": "US",
                "postalCode": "02142",
                "email": "r2d2@akamai.com"
            },
            "changeManagement": true,
            "location": "/cps/v2/enrollments/10003",
            "ra": "third-party",
            "adminContact": {
                "city": "Cambridge",
                "organizationName": "Dark Side",
                "firstName": "Darth",
                "addressLineTwo": null,
                "title": "Lord",
                "addressLineOne": "666 Evil Way",
                "lastName": "Vader",
                "region": "MA",
                "phone": "617-555-0123",
                "country": "US",
                "postalCode": "02142",
                "email": "vader@example.com"
            },
            "certificateChainType": "default",
            "org": {
                "city": "Cambridge",
                "name": "Akamai Technologies",
                "addressLineTwo": null,
                "addressLineOne": "150 Broadway",
                "country": "US",
                "region": "MA",
                "phone": "617-555-0111",
                "postalCode": "02142"
            },
            "certificateType": "third-party",
            "csr": {
                "c": "US",
                "cn": "www.example2.com",
                "l": "Cambridge",
                "o": "Akamai",
                "st": "MA",
                "ou": "WebEx",
                "sans": [
                    "san1.example.com",
                    "san2.example.com",
                    "san3.example.com",
                    "san4.example.com",
                    "www.example.com"
                ]
            },
            "signatureAlgorithm": null
        }
    ]
}`)

  sans := []string{
    "san1.example.com",
    "san2.example.com",
    "san3.example.com",
    "san4.example.com",
    "www.example.com",
  }

  pending := []string{
    "/cps/v2/enrollments/10003/changes/10004",
  }

  Config.NewConfig(config)
  e := NewEnrollments()
  err := e.GetEnrollments()

  assert.NoError(t, err)
  assert.IsType(t, &Enrollments{}, e)

  // First Enrollment
  e1 := e.Enrollments[0]
  assert.IsType(t, Enrollment{}, e1)
  assert.Equal(t, "/cps/v2/enrollments/10002", e1.Location)
  assert.Equal(t, "third-party", e1.RootAuth)
  assert.Equal(t, "third-party", e1.ValidationType)
  assert.Equal(t, "third-party", e1.CertificateType)
  assert.Nil(t, e1.SignatureAlgorithm)
  assert.Equal(t, true, e1.ChangeManagement)
  assert.Equal(t, false, e1.EnableMultiStackedCertificates)
  assert.Equal(t, []string{}, e1.PendingChanges)

  // CSR
  csr1 := e1.CSR
  assert.IsType(t, CSR{}, csr1)
  assert.Equal(t, "US", csr1.Country)
  assert.Equal(t, "www.example.com", csr1.CommonName)
  assert.Equal(t, "Cambridge", csr1.Location)
  assert.Equal(t, "Akamai", csr1.Org)
  assert.Equal(t, "MA", csr1.State)
  assert.Equal(t, "WebEx", csr1.OrgUnit)
  assert.Equal(t, sans, csr1.SANS)

  // Second Enrollment
  e2 := e.Enrollments[1]
  assert.IsType(t, Enrollment{}, e2)
  assert.Equal(t, "/cps/v2/enrollments/10003", e2.Location)
  assert.Equal(t, "third-party", e2.RootAuth)
  assert.Equal(t, "third-party", e2.ValidationType)
  assert.Equal(t, "third-party", e2.CertificateType)
  assert.Nil(t, e2.SignatureAlgorithm)
  assert.Equal(t, true, e2.ChangeManagement)
  assert.Equal(t, false, e2.EnableMultiStackedCertificates)
  assert.Equal(t, pending, e2.PendingChanges)

  // CSR
  csr2 := e2.CSR
  assert.IsType(t, CSR{}, csr2)
  assert.Equal(t, "US", csr2.Country)
  assert.Equal(t, "www.example2.com", csr2.CommonName)
  assert.Equal(t, "Cambridge", csr2.Location)
  assert.Equal(t, "Akamai", csr2.Org)
  assert.Equal(t, "MA", csr2.State)
  assert.Equal(t, "WebEx", csr2.OrgUnit)
  assert.Equal(t, sans, csr2.SANS)
}
