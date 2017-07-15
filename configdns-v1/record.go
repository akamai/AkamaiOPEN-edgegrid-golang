package dns

import "strings"

// Record represents a single DNS Record and it's settings
type Record struct {
	RecordType          string `json:"-"`
	Active              bool   `json:"active,omitempty"`
	Algorithm           int    `json:"algorithm,omitempty"`
	Contact             string `json:"contact,omitempty"`
	Digest              string `json:"digest,omitempty"`
	DigestType          int    `json:"digest_type,omitempty"`
	Expiration          string `json:"expiration,omitempty"`
	Expire              int    `json:"expire,omitempty"`
	Fingerprint         string `json:"fingerprint,omitempty"`
	FingerprintType     int    `json:"fingerprint_type,omitempty"`
	Flags               int    `json:"flags,omitempty"`
	Hardware            string `json:"hardware,omitempty"`
	Inception           string `json:"inception,omitempty"`
	Iterations          int    `json:"iterations,omitempty"`
	Key                 string `json:"key,omitempty"`
	Keytag              int    `json:"keytag,omitempty"`
	Labels              int    `json:"labels,omitempty"`
	Mailbox             string `json:"mailbox,omitempty"`
	Minimum             int    `json:"minimum,omitempty"`
	Name                string `json:"name,omitempty"`
	NextHashedOwnerName string `json:"next_hashed_owner_name,omitempty"`
	Order               int    `json:"order,omitempty"`
	OriginalTTL         int    `json:"original_ttl,omitempty"`
	Originserver        string `json:"originserver,omitempty"`
	Port                int    `json:"port,omitempty"`
	Preference          int    `json:"preference,omitempty"`
	Priority            int    `json:"priority,omitempty"`
	Protocol            int    `json:"protocol,omitempty"`
	Refresh             int    `json:"refresh,omitempty"`
	Regexp              string `json:"regexp,omitempty"`
	Replacement         string `json:"replacement,omitempty"`
	Retry               int    `json:"retry,omitempty"`
	Salt                string `json:"salt,omitempty"`
	Serial              int    `json:"serial,omitempty"`
	Service             string `json:"service,omitempty"`
	Signature           string `json:"signature,omitempty"`
	Signer              string `json:"signer,omitempty"`
	Software            string `json:"software,omitempty"`
	Subtype             int    `json:"subtype,omitempty"`
	Target              string `json:"target,omitempty"`
	TTL                 int    `json:"ttl,omitempty"`
	Txt                 string `json:"txt,omitempty"`
	TypeBitmaps         string `json:"type_bitmaps,omitempty"`
	TypeCovered         string `json:"type_coverered,omitempty"`
	Weight              uint   `json:"weight,omitempty"`
}

// Allows will validates if a the current record type allows a given field
func (record *Record) Allows(field string) bool {
	field = strings.ToLower(field)

	fieldMap := map[string]map[string]struct{}{
		"active": {
			"A":          {},
			"AAAA":       {},
			"AFSDB":      {},
			"CNAME":      {},
			"DNSKEY":     {},
			"DS":         {},
			"HINFO":      {},
			"LOC":        {},
			"MX":         {},
			"NAPTR":      {},
			"NS":         {},
			"NSEC3":      {},
			"NSEC3PARAM": {},
			"PTR":        {},
			"RP":         {},
			"RRSIG":      {},
			"SPF":        {},
			"SRV":        {},
			"SSHFP":      {},
			"TXT":        {},
		},
		"algorithm": {
			"DNSKEY":     {},
			"DS":         {},
			"NSEC3":      {},
			"NSEC3PARAM": {},
			"RRSIG":      {},
			"SSHFP":      {},
		},
		"contact":         {"SOA": {}},
		"digest":          {"DS": {}},
		"digesttype":      {"DS": {}},
		"expiration":      {"RRSIG": {}},
		"expire":          {"SOA": {}},
		"fingerprint":     {"SSHFP": {}},
		"fingerprinttype": {"SSHFP": {}},
		"flags": {
			"DNSKEY":     {},
			"NAPTR":      {},
			"NSEC3":      {},
			"NSEC3PARAM": {},
		},
		"hardware":  {"HINFO": {}},
		"inception": {"RRSIG": {}},
		"iterations": {
			"NSEC3":       {},
			"NSEC3PARAMS": {},
		},
		"key": {
			"DNSKEY": {},
			"DS":     {},
		},
		"keytag":  {"RRSIG": {}},
		"labels":  {"RRSIG": {}},
		"mailbox": {"RP": {}},
		"minimum": {"SOA": {}},
		"name": {
			"A":          {},
			"AAAA":       {},
			"AFSDB":      {},
			"CNAME":      {},
			"DNSKEY":     {},
			"DS":         {},
			"HINFO":      {},
			"LOC":        {},
			"MX":         {},
			"NAPTR":      {},
			"NS":         {},
			"NSEC3":      {},
			"NSEC3PARAM": {},
			"PTR":        {},
			"RP":         {},
			"RRSIG":      {},
			"SPF":        {},
			"SRV":        {},
			"SSHFP":      {},
			"TXT":        {},
		},
		"nexthashedownername": {"NSEC3": {}},
		"order":               {"NAPTR": {}},
		"originalttl":         {"RRSIG": {}},
		"originserver":        {"SOA": {}},
		"port":                {"SRV": {}},
		"preference":          {"NAPTR": {}},
		"priority": {
			"SRV": {},
			"MX":  {},
		},
		"protocol":    {"DNSKEY": {}},
		"refresh":     {"SOA": {}},
		"regexp":      {"NAPTR": {}},
		"replacement": {"NAPTR": {}},
		"retry":       {"SOA": {}},
		"salt": {
			"NSEC3":      {},
			"NSEC3PARAM": {},
		},
		"serial":    {"SOA": {}},
		"service":   {"NAPTR": {}},
		"signature": {"RRSIG": {}},
		"signer":    {"RRSIG": {}},
		"software":  {"HINFO": {}},
		"subtype":   {"AFSDB": {}},
		"targets": {
			"A":          {},
			"AAAA":       {},
			"AFSDB":      {},
			"CNAME":      {},
			"DNSKEY":     {},
			"DS":         {},
			"HINFO":      {},
			"LOC":        {},
			"MX":         {},
			"NAPTR":      {},
			"NS":         {},
			"NSEC3":      {},
			"NSEC3PARAM": {},
			"PTR":        {},
			"RP":         {},
			"RRSIG":      {},
			"SOA":        {},
			"SPF":        {},
			"SRV":        {},
			"SSHFP":      {},
			"TXT":        {},
		},
		"ttl": {
			"A":     {},
			"AAAA":  {},
			"AFSDB": {},
			"CNAME": {},
			"LOC":   {},
			"MX":    {},
			"NS":    {},
			"PTR":   {},
			"SPF":   {},
			"SRV":   {},
			"TXT":   {},
		},
		"txt":         {"RP": {}},
		"typebitmaps": {"NSEC3": {}},
		"typecovered": {"RRSIG": {}},
		"weight":      {"SRV": {}},
	}

	_, ok := fieldMap[field][strings.ToUpper(record.RecordType)]

	return ok
}
