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
                "addressLineTwo": "Death Star",
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
                "addressLineTwo": "Death Star",
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

  dns_names1 := []string{
    "san2.example.com",
    "san1.example.com",
  }

  dns_names2 := []string{
    "san1.example.com",
    "san2.example.com",
  }

  Config.NewConfig(config)
  e, err := GetEnrollments()

  assert.NoError(t, err)
  assert.IsType(t, &Enrollments{}, e)

  // First Enrollment
  e1 := e.Enrollments[0]
  assert.IsType(t, Enrollment{}, e1)
  assert.Equal(t, "/cps/v2/enrollments/10002", e1.Location)
  assert.Equal(t, "third-party", e1.RootAuth)
  assert.Equal(t, "third-party", e1.ValidationType)
  assert.Equal(t, "third-party", e1.CertificateType)
  assert.Equal(t, "", e1.SignatureAlgorithm)
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

  // Org
  org1 := e1.Org
  assert.IsType(t, Org{}, org1)
  assert.Equal(t, "Akamai Technologies", org1.Name)
  assert.Equal(t, "150 Broadway", org1.AddressLineOne)
  assert.Equal(t, "", org1.AddressLineTwo)
  assert.Equal(t, "Cambridge", org1.City)
  assert.Equal(t, "MA", org1.Region)
  assert.Equal(t, "02142", org1.PostalCode)
  assert.Equal(t, "US", org1.Country)
  assert.Equal(t, "617-555-0111", org1.Phone)

  // TechContact
  tc1 := e1.TechContact
  assert.IsType(t, Contact{}, tc1)
  assert.Equal(t, "R2", tc1.FirstName)
  assert.Equal(t, "D2", tc1.LastName)
  assert.Equal(t, "617-555-0111", tc1.Phone)
  assert.Equal(t, "r2d2@akamai.com", tc1.Email)
  assert.Equal(t, "150 Broadway", tc1.AddressLineOne)
  assert.Equal(t, "", tc1.AddressLineTwo)
  assert.Equal(t, "Cambridge", tc1.City)
  assert.Equal(t, "US", tc1.Country)
  assert.Equal(t, "Akamai", tc1.OrganizationName)
  assert.Equal(t, "02142", tc1.PostalCode)
  assert.Equal(t, "MA", tc1.Region)
  assert.Equal(t, "Astromech Droid", tc1.Title)

  // AdminContact
  ac1 := e1.AdminContact
  assert.IsType(t, Contact{}, ac1)
  assert.Equal(t, "Darth", ac1.FirstName)
  assert.Equal(t, "Vader", ac1.LastName)
  assert.Equal(t, "617-555-0123", ac1.Phone)
  assert.Equal(t, "vader@example.com", ac1.Email)
  assert.Equal(t, "666 Evil Way", ac1.AddressLineOne)
  assert.Equal(t, "Death Star", ac1.AddressLineTwo)
  assert.Equal(t, "Cambridge", ac1.City)
  assert.Equal(t, "US", ac1.Country)
  assert.Equal(t, "Dark Side", ac1.OrganizationName)
  assert.Equal(t, "02142", ac1.PostalCode)
  assert.Equal(t, "MA", ac1.Region)
  assert.Equal(t, "Lord", ac1.Title)

  // NetworkConfiguration
  nc1 := e1.NetworkConfiguration
  assert.IsType(t, NetworkConfiguration{}, nc1)
  assert.Equal(t, "core", nc1.Geography)
  assert.Equal(t, "enhanced-tls", nc1.SecureNetwork)
  assert.Equal(t, "ak-akamai-default2016q3", nc1.MustHaveCiphers)
  assert.Equal(t, "ak-akamai-default", nc1.PreferredCiphers)
  assert.Equal(t, []string{}, nc1.DisallowedTlsVersions)

  // SNI
  nc_sni1 := nc1.SNI
  assert.IsType(t, SNI{}, nc_sni1)
  assert.Equal(t, false, nc_sni1.CloneDnsNames)
  assert.Equal(t, dns_names1, nc_sni1.DnsNames)

  // Second Enrollment
  e2 := e.Enrollments[1]
  assert.IsType(t, Enrollment{}, e2)
  assert.Equal(t, "/cps/v2/enrollments/10003", e2.Location)
  assert.Equal(t, "third-party", e2.RootAuth)
  assert.Equal(t, "third-party", e2.ValidationType)
  assert.Equal(t, "third-party", e2.CertificateType)
  assert.Equal(t, "", e2.SignatureAlgorithm)
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

  // Org
  org2 := e2.Org
  assert.IsType(t, Org{}, org2)
  assert.Equal(t, "Akamai Technologies", org2.Name)
  assert.Equal(t, "150 Broadway", org2.AddressLineOne)
  assert.Equal(t, "", org2.AddressLineTwo)
  assert.Equal(t, "Cambridge", org2.City)
  assert.Equal(t, "MA", org2.Region)
  assert.Equal(t, "02142", org2.PostalCode)
  assert.Equal(t, "US", org2.Country)
  assert.Equal(t, "617-555-0111", org2.Phone)

  // TechContact
  tc2 := e2.TechContact
  assert.IsType(t, Contact{}, tc2)
  assert.Equal(t, "R2", tc2.FirstName)
  assert.Equal(t, "D2", tc2.LastName)
  assert.Equal(t, "617-555-0111", tc2.Phone)
  assert.Equal(t, "r2d2@akamai.com", tc2.Email)
  assert.Equal(t, "150 Broadway", tc2.AddressLineOne)
  assert.Equal(t, "", tc2.AddressLineTwo)
  assert.Equal(t, "Cambridge", tc2.City)
  assert.Equal(t, "US", tc2.Country)
  assert.Equal(t, "Akamai", tc2.OrganizationName)
  assert.Equal(t, "02142", tc2.PostalCode)
  assert.Equal(t, "MA", tc2.Region)
  assert.Equal(t, "Astromech Droid", tc2.Title)

  // AdminContact
  ac2 := e2.AdminContact
  assert.IsType(t, Contact{}, ac2)
  assert.Equal(t, "Darth", ac2.FirstName)
  assert.Equal(t, "Vader", ac2.LastName)
  assert.Equal(t, "617-555-0123", ac2.Phone)
  assert.Equal(t, "vader@example.com", ac2.Email)
  assert.Equal(t, "666 Evil Way", ac2.AddressLineOne)
  assert.Equal(t, "Death Star", ac2.AddressLineTwo)
  assert.Equal(t, "Cambridge", ac2.City)
  assert.Equal(t, "US", ac2.Country)
  assert.Equal(t, "Dark Side", ac2.OrganizationName)
  assert.Equal(t, "02142", ac2.PostalCode)
  assert.Equal(t, "MA", ac2.Region)
  assert.Equal(t, "Lord", ac2.Title)

  // NetworkConfiguration
  nc2 := e2.NetworkConfiguration
  assert.IsType(t, NetworkConfiguration{}, nc2)
  assert.Equal(t, "core", nc2.Geography)
  assert.Equal(t, "enhanced-tls", nc2.SecureNetwork)
  assert.Equal(t, "ak-akamai-default2016q3", nc2.MustHaveCiphers)
  assert.Equal(t, "ak-akamai-default", nc2.PreferredCiphers)
  assert.Equal(t, []string{}, nc2.DisallowedTlsVersions)

  // SNI
  nc_sni2 := nc2.SNI
  assert.IsType(t, SNI{}, nc_sni2)
  assert.Equal(t, false, nc_sni2.CloneDnsNames)
  assert.Equal(t, dns_names2, nc_sni2.DnsNames)
}

func TestCPS_GetEnrollmentsError(t *testing.T) {
  defer gock.Off()
  mock := gock.New("https://test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/cps/v2/enrollments")
	mock.
		Get("/cps/v2/enrollments").
		HeaderPresent("Authorization").
		Reply(401).
		SetHeader("Content-Type", "application/json").
		BodyString(`{}`)

  Config.NewConfig(config)
  _, err := GetEnrollments()
  assert.Error(t, err)
}

func TestCPS_GetEnrollmentsParseError(t *testing.T) {
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
          "location": "/cps/v2/enrollments/10002",
        }
      ]
    }`)

  Config.NewConfig(config)
  _, err := GetEnrollments()
  assert.Error(t, err)
}

func TestCPS_NewEnrollment(t *testing.T) {
  e := NewEnrollment()
  assert.IsType(t, &Enrollment{}, e)
}

func TestCPS_GetEnrollment(t *testing.T) {
  defer gock.Off()
  mock := gock.New("https://test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/cps/v2/enrollments/10002")
	mock.
		Get("/cps/v2/enrollments/10002").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
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
        "location": "/cps/v2/enrollments/10002",
        "ra": "third-party",
        "adminContact": {
            "city": "Cambridge",
            "organizationName": "Dark Side",
            "firstName": "Darth",
            "addressLineTwo": "Death Star",
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

  dns_names := []string{
    "san2.example.com",
    "san1.example.com",
  }

  Config.NewConfig(config)
  _, err := GetEnrollment("")
  assert.Error(t, err)

  e1, err1 := GetEnrollment("10002")
  assert.NoError(t, err1)
  assert.IsType(t, &Enrollment{}, e1)
  assert.Equal(t, "/cps/v2/enrollments/10002", e1.Location)
  assert.Equal(t, "third-party", e1.RootAuth)
  assert.Equal(t, "third-party", e1.ValidationType)
  assert.Equal(t, "third-party", e1.CertificateType)
  assert.Equal(t, "", e1.SignatureAlgorithm)
  assert.Equal(t, true, e1.ChangeManagement)
  assert.Equal(t, false, e1.EnableMultiStackedCertificates)
  assert.Equal(t, pending, e1.PendingChanges)

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

  // Org
  org1 := e1.Org
  assert.IsType(t, Org{}, org1)
  assert.Equal(t, "Akamai Technologies", org1.Name)
  assert.Equal(t, "150 Broadway", org1.AddressLineOne)
  assert.Equal(t, "", org1.AddressLineTwo)
  assert.Equal(t, "Cambridge", org1.City)
  assert.Equal(t, "MA", org1.Region)
  assert.Equal(t, "02142", org1.PostalCode)
  assert.Equal(t, "US", org1.Country)
  assert.Equal(t, "617-555-0111", org1.Phone)

  // TechContact
  tc1 := e1.TechContact
  assert.IsType(t, Contact{}, tc1)
  assert.Equal(t, "R2", tc1.FirstName)
  assert.Equal(t, "D2", tc1.LastName)
  assert.Equal(t, "617-555-0111", tc1.Phone)
  assert.Equal(t, "r2d2@akamai.com", tc1.Email)
  assert.Equal(t, "150 Broadway", tc1.AddressLineOne)
  assert.Equal(t, "", tc1.AddressLineTwo)
  assert.Equal(t, "Cambridge", tc1.City)
  assert.Equal(t, "US", tc1.Country)
  assert.Equal(t, "Akamai", tc1.OrganizationName)
  assert.Equal(t, "02142", tc1.PostalCode)
  assert.Equal(t, "MA", tc1.Region)
  assert.Equal(t, "Astromech Droid", tc1.Title)

  // AdminContact
  ac1 := e1.AdminContact
  assert.IsType(t, Contact{}, ac1)
  assert.Equal(t, "Darth", ac1.FirstName)
  assert.Equal(t, "Vader", ac1.LastName)
  assert.Equal(t, "617-555-0123", ac1.Phone)
  assert.Equal(t, "vader@example.com", ac1.Email)
  assert.Equal(t, "666 Evil Way", ac1.AddressLineOne)
  assert.Equal(t, "Death Star", ac1.AddressLineTwo)
  assert.Equal(t, "Cambridge", ac1.City)
  assert.Equal(t, "US", ac1.Country)
  assert.Equal(t, "Dark Side", ac1.OrganizationName)
  assert.Equal(t, "02142", ac1.PostalCode)
  assert.Equal(t, "MA", ac1.Region)
  assert.Equal(t, "Lord", ac1.Title)

  // NetworkConfiguration
  nc1 := e1.NetworkConfiguration
  assert.IsType(t, NetworkConfiguration{}, nc1)
  assert.Equal(t, "core", nc1.Geography)
  assert.Equal(t, "enhanced-tls", nc1.SecureNetwork)
  assert.Equal(t, "ak-akamai-default2016q3", nc1.MustHaveCiphers)
  assert.Equal(t, "ak-akamai-default", nc1.PreferredCiphers)
  assert.Equal(t, []string{}, nc1.DisallowedTlsVersions)

  // SNI
  nc_sni1 := nc1.SNI
  assert.IsType(t, SNI{}, nc_sni1)
  assert.Equal(t, false, nc_sni1.CloneDnsNames)
  assert.Equal(t, dns_names, nc_sni1.DnsNames)
}

func TestCPS_GetEnrollmentError(t *testing.T) {
  defer gock.Off()
  mock := gock.New("https://test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/cps/v2/enrollments/10002")
	mock.
		Get("/cps/v2/enrollments/10002").
		HeaderPresent("Authorization").
		Reply(401).
		SetHeader("Content-Type", "application/json").
		BodyString(`{}`)

  Config.NewConfig(config)
  _, err := GetEnrollment("10002")
  assert.Error(t, err)
}

func TestCPS_GetEnrollmentParseError(t *testing.T) {
  defer gock.Off()
  mock := gock.New("https://test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/cps/v2/enrollments/10002")
	mock.
		Get("/cps/v2/enrollments/10002").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
      "location": "/cps/v2/enrollments/10002",
    }`)

  Config.NewConfig(config)
  _, err := GetEnrollment("10002")
  assert.Error(t, err)
}
