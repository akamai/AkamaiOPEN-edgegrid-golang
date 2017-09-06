package dns

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var (
	config = edgegrid.Config{
		Host:         "akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
	}
)

func TestGetZoneSimple(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
	mock.
		Get("/config-dns/v1/zones/example.com").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
			"token": "a184671d5307a388180fbf7f11dbdf46",
			"zone": {
				"name": "example.com",
				"soa": {
					"contact": "hostmaster.akamai.com.",
					"expire": 604800,
					"minimum": 180,
					"originserver": "use4.akamai.com.",
					"refresh": 900,
					"retry": 300,
					"serial": 1271354824,
					"ttl": 900
				},
				"ns": [
					{
						"active": true,
						"name": "",
						"target": "use4.akam.net.",
						"ttl": 3600
					},
					{
						"active": true,
						"name": "",
						"target": "use3.akam.net.",
						"ttl": 3600
					}
				],
				"a": [
					{
						"active": true,
						"name": "www",
						"target": "1.2.3.4",
						"ttl": 30
					}
				]
			}
		}`)

	Init(config)
	zone, err := GetZone("example.com")

	assert.NoError(t, err)

	assert.IsType(t, &Zone{}, zone)
	assert.Equal(t, "a184671d5307a388180fbf7f11dbdf46", zone.Token)
	assert.Equal(t, "example.com", zone.Zone.Name)

	assert.IsType(t, &Record{}, zone.Zone.Soa)
	assert.Equal(t, "hostmaster.akamai.com.", zone.Zone.Soa.Contact)
	assert.Equal(t, 604800, zone.Zone.Soa.Expire)
	assert.Equal(t, 180, zone.Zone.Soa.Minimum)
	assert.Equal(t, "use4.akamai.com.", zone.Zone.Soa.Originserver)
	assert.Equal(t, 900, zone.Zone.Soa.Refresh)
	assert.Equal(t, 300, zone.Zone.Soa.Retry)
	assert.Equal(t, 1271354824, zone.Zone.Soa.Serial)
	assert.Equal(t, 900, zone.Zone.Soa.TTL)

	assert.IsType(t, NsRecordSet{}, zone.Zone.Ns)
	assert.Len(t, zone.Zone.Ns, 2)

	assert.Equal(t, true, zone.Zone.Ns[0].Active)
	assert.Equal(t, stringPointer(""), zone.Zone.Ns[0].Name)
	assert.Equal(t, "use4.akam.net.", zone.Zone.Ns[0].Target)
	assert.Equal(t, 3600, zone.Zone.Ns[0].TTL)

	assert.Equal(t, true, zone.Zone.Ns[1].Active)
	assert.Equal(t, stringPointer(""), zone.Zone.Ns[1].Name)
	assert.Equal(t, "use3.akam.net.", zone.Zone.Ns[1].Target)
	assert.Equal(t, 3600, zone.Zone.Ns[1].TTL)

	assert.IsType(t, RecordSet{}, zone.Zone.A)
	assert.Len(t, zone.Zone.A, 1)

	assert.Equal(t, true, zone.Zone.A[0].Active)
	assert.Equal(t, "www", zone.Zone.A[0].Name)
	assert.Equal(t, "1.2.3.4", zone.Zone.A[0].Target)
	assert.Equal(t, 30, zone.Zone.A[0].TTL)
}

func TestGetZone(t *testing.T) {
	defer gock.Off()

	tests := testGetZoneCompleteProvider()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				records := zone.GetRecordType(test.recordType)
				assert.IsType(t, test.expectedType, records)
				assert.Equal(
					t,
					len(test.expectedRecords.([]Record)),
					len(records.(RecordSet)),
				)

				for key, record := range test.expectedRecords.([]Record) {
					assert.ObjectsAreEqual(record, records.(RecordSet)[key])
				}
			}
		})
	}
}

func TestGetZoneNaptrRecords(t *testing.T) {
	defer gock.Off()

	tests := recordTests{
		{
			name:       "NAPTR Records",
			recordType: "NAPTR",
			responseBody: `{
				"zone": {
					"naptr": [
						{
							"active": true,
							"flags": "S",
							"name": "naptrrecord",
							"order": 0,
							"preference": 10,
							"regexp": "!^.*$!sip:customer-service@example.com!",
							"replacement": ".",
							"service": "SIP+D2U",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: NaptrRecordSet{},
			expectedRecords: []NaptrRecord{
				NaptrRecord{
					Active:      true,
					Flags:       "S",
					Name:        "naptrrecord",
					Order:       0,
					Preference:  10,
					Regexp:      "!^.*$!sip:customer-service@example.com!",
					Replacement: ".",
					Service:     "SIP+D2U",
					TTL:         3600,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				records := zone.GetRecordType(test.recordType)
				assert.IsType(t, test.expectedType, records)
				assert.Equal(
					t,
					len(test.expectedRecords.([]NaptrRecord)),
					len(records.(NaptrRecordSet)),
				)

				for key, record := range test.expectedRecords.([]NaptrRecord) {
					assert.ObjectsAreEqual(record, records.(NaptrRecordSet)[key])
				}
			}
		})
	}
}

func TestGetZoneNsRecords(t *testing.T) {
	defer gock.Off()

	tests := recordTests{
		{
			name:       "NS Records",
			recordType: "NS",
			responseBody: `{
				"zone": {
					"ns": [
						{
							"active": true,
							"name": null,
							"target": "use4.akam.net.",
							"ttl": 3600
						},
						{
							"active": true,
							"name": null,
							"target": "use3.akam.net.",
							"ttl": 3600
						},
						{
							"active": true,
							"name": "five",
							"target": "use4.akam.net.",
							"ttl": 172800
						}
					]
				}
			}`,
			expectedType: NsRecordSet{},
			expectedRecords: []NsRecord{
				NsRecord{
					Active: true,
					Target: "use4.akam.net.",
					TTL:    3600,
				},
				NsRecord{
					Active: true,
					Target: "us34.akam.net.",
					TTL:    3600,
				},
				NsRecord{
					Active: true,
					Name:   stringPointer("five"),
					Target: "use4.akam.net.",
					TTL:    172800,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				records := zone.GetRecordType(test.recordType)
				assert.IsType(t, test.expectedType, records)
				assert.Equal(
					t,
					len(test.expectedRecords.([]NsRecord)),
					len(records.(NsRecordSet)),
				)

				for key, record := range test.expectedRecords.([]NsRecord) {
					assert.ObjectsAreEqual(record, records.(NsRecordSet)[key])
				}
			}
		})
	}
}

func TestGetZoneNsec3Records(t *testing.T) {
	defer gock.Off()

	tests := recordTests{
		{
			name:       "NSEC3 Records",
			recordType: "NSEC3",
			responseBody: `{
				"zone": {
					"nsec3": [
						{
							"active": true,
							"algorithm": 1,
							"flags": 0,
							"iterations": 1,
							"name": "qdeo8lqu4l81uo67oolpo9h0nv9l13dh",
							"next_hashed_owner_name": "R2NUSMGFSEUHT195P59KOU2AI30JR96P",
							"salt": "EBD1E0942543A01B",
							"ttl": 7200,
							"type_bitmaps": "CNAME RRSIG"
						}
					]
				}
			}`,
			expectedType: Nsec3RecordSet{},
			expectedRecords: []Nsec3Record{
				Nsec3Record{
					Active:              true,
					Algorithm:           1,
					Flags:               0,
					Iterations:          1,
					Name:                "qdeo8lqu4l81uo67oolpo9h0nv9l13dh",
					NextHashedOwnerName: "R2NUSMGFSEUHT195P59KOU2AI30JR96P",
					Salt:                "EBD1E0942543A01B",
					TTL:                 7200,
					TypeBitmaps:         "CNAME RRSIG",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				records := zone.GetRecordType(test.recordType)
				assert.IsType(t, test.expectedType, records)
				assert.Equal(
					t,
					len(test.expectedRecords.([]Nsec3Record)),
					len(records.(Nsec3RecordSet)),
				)

				for key, record := range test.expectedRecords.([]Nsec3Record) {
					assert.ObjectsAreEqual(record, records.(Nsec3RecordSet)[key])
				}
			}
		})
	}
}

func TestGetZoneNsec3paramRecords(t *testing.T) {
	defer gock.Off()

	tests := recordTests{
		{
			name:       "NSEC3PARAM Records",
			recordType: "NSEC3PARAM",
			responseBody: `{
				"zone": {
					"nsec3param": [
						{
							"active": true,
							"algorithm": 1,
							"flags": 0,
							"iterations": 1,
							"name": "qnsec3param",
							"salt": "EBD1E0942543A01B",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: Nsec3paramRecordSet{},
			expectedRecords: []Nsec3paramRecord{
				Nsec3paramRecord{
					Active:     true,
					Algorithm:  1,
					Flags:      0,
					Iterations: 1,
					Name:       "qnsec3param",
					Salt:       "EBD1E0942543A01B",
					TTL:        7200,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				records := zone.GetRecordType(test.recordType)
				assert.IsType(t, test.expectedType, records)
				assert.Equal(
					t,
					len(test.expectedRecords.([]Nsec3paramRecord)),
					len(records.(Nsec3paramRecordSet)),
				)

				for key, record := range test.expectedRecords.([]Nsec3paramRecord) {
					assert.ObjectsAreEqual(record, records.(Nsec3paramRecordSet)[key])
				}
			}
		})
	}
}

func TestGetZoneSrvRecords(t *testing.T) {
	defer gock.Off()

	tests := recordTests{
		{
			name:       "SRV Records",
			recordType: "SRV",
			responseBody: `{
				"zone": {
					"srv": [
						{
							"active": true,
							"name": "srv",
							"port": 522,
							"priority": 10,
							"target": "target.akamai.com.",
							"ttl": 7200,
							"weight": 0
						}
					]
				}
			}`,
			expectedType: SrvRecordSet{},
			expectedRecords: []SrvRecord{
				SrvRecord{
					Active:   true,
					Name:     "srv",
					Port:     522,
					Priority: 10,
					Target:   "target.akamai.com.",
					TTL:      7200,
					Weight:   0,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				records := zone.GetRecordType(test.recordType)
				assert.IsType(t, test.expectedType, records)
				assert.Equal(
					t,
					len(test.expectedRecords.([]SrvRecord)),
					len(records.(SrvRecordSet)),
				)

				for key, record := range test.expectedRecords.([]SrvRecord) {
					assert.ObjectsAreEqual(record, records.(SrvRecordSet)[key])
				}
			}
		})
	}
}

type recordTests []struct {
	name            string
	recordType      string
	responseBody    string
	expectedType    interface{}
	expectedRecord  Record
	expectedRecords interface{}
}

func testGetZoneCompleteProvider() recordTests {
	return recordTests{
		{
			name:       "A Records",
			recordType: "A",
			responseBody: `{
				"zone": {
					"a": [
						{
							"active": true,
							"name": "arecord",
							"target": "1.2.3.5",
							"ttl": 3600
						},
						{
							"active": true,
							"name": "origin",
							"target": "1.2.3.9",
							"ttl": 3600
						},
						{
							"active": true,
							"name": "arecord",
							"target": "1.2.3.4",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "arecord",
					Target: "1.2.3.5",
					TTL:    3600,
				},
				Record{
					Active: true,
					Name:   "origin",
					Target: "1.2.3.9",
					TTL:    3600,
				},
				Record{
					Active: true,
					Name:   "arecord",
					Target: "1.2.3.4",
					TTL:    3600,
				},
			},
		},
		{
			name:       "AAAA Records",
			recordType: "AAAA",
			responseBody: `{
				"zone": {
					"aaaa": [
						{
							"active": true,
							"name": "ipv6record",
							"target": "2001:0db8::ff00:0042:8329",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "ipv6record",
					Target: "2001:0db8::ff00:0042:8329",
					TTL:    3600,
				},
			},
		},
		{
			name:       "AFSDB Records",
			recordType: "AFSDB",
			responseBody: `{
				"zone": {
					"afsdb": [
						{
							"active": true,
							"name": "afsdb",
							"subtype": 1,
							"target": "example.com.",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:  true,
					Name:    "afsdb",
					Subtype: 1,
					Target:  "example.com.",
					TTL:     7200,
				},
			},
		},
		{
			name:       "CNAME Records",
			recordType: "CNAME",
			responseBody: `{
				"zone": {
					"cname": [
						{
							"active": true,
							"name": "redirect",
							"target": "arecord.example.com.",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "redirect",
					Target: "arecord.example.com.",
					TTL:    3600,
				},
			},
		},
		{
			name:       "DNSKEY Records",
			recordType: "DNSKEY",
			responseBody: `{
				"zone": {
					"dnskey": [
						{
							"active": true,
							"algorithm": 3,
							"flags": 257,
							"key": "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ==",
							"name": "dnskey",
							"protocol": 7,
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:    true,
					Algorithm: 3,
					Flags:     257,
					Key:       "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ==",
					Name:      "dnskey",
					Protocol:  7,
					TTL:       7200,
				},
			},
		},
		{
			name:       "DS Records",
			recordType: "DS",
			responseBody: `{
				"zone": {
					"ds": [
						{
							"active": true,
							"algorithm": 7,
							"digest": "909FF0B4DD66F91F56524C4F968D13083BE42380",
							"digest_type": 1,
							"keytag": 30336,
							"name": "ds",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:     true,
					Algorithm:  7,
					Digest:     "909FF0B4DD66F91F56524C4F968D13083BE42380",
					DigestType: 1,
					Keytag:     30336,
					Name:       "ds",
					TTL:        7200,
				},
			},
		},
		{
			name:       "HINFO Records",
			recordType: "HINFO",
			responseBody: `{
				"zone": {
					"hinfo": [
						{
							"active": true,
							"hardware": "INTEL-386",
							"name": "hinfo",
							"software": "UNIX",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:   true,
					Hardware: "INTEL-386",
					Name:     "hinfo",
					Software: "UNIX",
					TTL:      7200,
				},
			},
		},
		{
			name:       "LOC Records",
			recordType: "LOC",
			responseBody: `{
				"zone": {
					"loc": [
						{
							"active": true,
							"name": "location",
							"target": "51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "location",
					Target: "51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m",
					TTL:    7200,
				},
			},
		},
		{
			name:       "MX Records",
			recordType: "MX",
			responseBody: `{
				"zone": {
					"mx": [
						{
							"active": true,
							"name": "four",
							"priority": 10,
							"target": "mx1.akamai.com.",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:   true,
					Name:     "four",
					Priority: 10,
					Target:   "mx1.akamai.com.",
					TTL:      7200,
				},
			},
		},
		{
			name:       "PTR Records",
			recordType: "PTR",
			responseBody: `{
				"zone": {
					"ptr": [
						{
							"active": true,
							"name": "ptr",
							"target": "ptr.example.com.",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "ptr",
					Target: "ptr.example.com.",
					TTL:    7200,
				},
			},
		},
		{
			name:       "RP Records",
			recordType: "RP",
			responseBody: `{
				"zone": {
					"rp": [
						{
							"active": true,
							"mailbox": "admin.example.com.",
							"name": "rp",
							"ttl": 7200,
							"txt": "txt.example.com."
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:  true,
					Mailbox: "admin.example.com.",
					Name:    "rp",
					TTL:     7200,
					Txt:     "txt.example.com.",
				},
			},
		},
		{
			name:       "RRSIG Records",
			recordType: "RRSIG",
			responseBody: `{
				"zone": {
					"rrsig": [
						{
							"active": true,
							"algorithm": 7,
							"expiration": "20120318104101",
							"inception": "20120315094101",
							"keytag": 63761,
							"labels": 3,
							"name": "arecord",
							"original_ttl": 3600,
							"signature": "toCy19QnAb86vRlQjf5 ... z1doJdHEr8PiI+Is9Eafxh+4Idcw8Ysv",
							"signer": "example.com.",
							"ttl": 7200,
							"type_covered": "A"
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:      true,
					Algorithm:   7,
					Expiration:  "20120318104101",
					Inception:   "20120315094101",
					Keytag:      63761,
					Labels:      3,
					Name:        "arecord",
					OriginalTTL: 3600,
					Signature:   "toCy19QnAb86vRlQjf5 ... z1doJdHEr8PiI+Is9Eafxh+4Idcw8Ysv",
					Signer:      "example.com.",
					TTL:         7200,
					TypeCovered: "A",
				},
			},
		},
		{
			name:       "SPF Records",
			recordType: "SPF",
			responseBody: `{
				"zone": {
					"spf": [
						{
							"active": true,
							"name": "spf",
							"target": "v=spf a -all",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "spf",
					Target: "v=spf a -all",
					TTL:    7200,
				},
			},
		},
		{
			name:       "SSHFP Records",
			recordType: "SSHFP",
			responseBody: `{
				"zone": {
					"sshfp": [
						{
							"active": true,
							"algorithm": 2,
							"fingerprint": "123456789ABCDEF67890123456789ABCDEF67890",
							"fingerprint_type": 1,
							"name": "host",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active:          true,
					Algorithm:       2,
					Fingerprint:     "123456789ABCDEF67890123456789ABCDEF67890",
					FingerprintType: 1,
					Name:            "host",
					TTL:             3600,
				},
			},
		},
		{
			name:       "TXT Records",
			recordType: "TXT",
			responseBody: `{
				"zone": {
					"txt": [
						{
							"active": true,
							"name": "text",
							"target": "Hello world!",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: RecordSet{},
			expectedRecords: []Record{
				Record{
					Active: true,
					Name:   "text",
					Target: "Hello world!",
					TTL:    7200,
				},
			},
		},
	}
}

func stringPointer(s string) *string {
	return &s
}
