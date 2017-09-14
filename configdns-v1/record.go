package dns

import (
	// "fmt"
	"strings"
)

// type ValidationFailedError struct {
// 	fieldName string
// 	err       error
// }
//
// type Error interface {
// 	error
// 	IsZoneNotFound() bool
// 	IsFailedToSave() bool
// 	IsValidationFailed() bool
// }
//
// func (e ValidationFailedError) Error() string {
// 	if e == nil {
// 		return "<nil>"
// 	}
// 	return fmt.Sprintf("Validation Failed - Field is not allowed for this type: %s", e.fieldName)
// }

// All record types (below) must implement the DNSRecord interface
// This allows the record to be used dynamically in slices - see the Zone struct definition in zone.go
type DNSRecord interface {
	// Allows will validates if a the current record type allows a given field
	Allows(field string) bool
	// Get the list of allowed fields for this type
	GetAllowedFields() []string
	// Set a field on the struct, which check for valid fields
	SetField(name string, value interface{}) error
}

// RecordSet represents a collection of Records
type RecordSet []DNSRecord

// baseRecord represents the basic settings required for all types of DNS Records
type baseRecord struct {
	fieldMap            []string `json:"-"`
	RecordType          string   `json:"-"`
	Active              bool     `json:"active,omitempty"`
	Algorithm           int      `json:"algorithm,omitempty"`
	Contact             string   `json:"contact,omitempty"`
	Digest              string   `json:"digest,omitempty"`
	DigestType          int      `json:"digest_type,omitempty"`
	Expiration          string   `json:"expiration,omitempty"`
	Expire              int      `json:"expire,omitempty"`
	Fingerprint         string   `json:"fingerprint,omitempty"`
	FingerprintType     int      `json:"fingerprint_type,omitempty"`
	Flags               int      `json:"flags,omitempty"`
	Hardware            string   `json:"hardware,omitempty"`
	Inception           string   `json:"inception,omitempty"`
	Iterations          int      `json:"iterations,omitempty"`
	Key                 string   `json:"key,omitempty"`
	Keytag              int      `json:"keytag,omitempty"`
	Labels              int      `json:"labels,omitempty"`
	Mailbox             string   `json:"mailbox,omitempty"`
	Minimum             int      `json:"minimum,omitempty"`
	Name                string   `json:"name,omitempty"`
	NextHashedOwnerName string   `json:"next_hashed_owner_name,omitempty"`
	Order               int      `json:"order,omitempty"`
	OriginalTTL         int      `json:"original_ttl,omitempty"`
	Originserver        string   `json:"originserver,omitempty"`
	Port                int      `json:"port,omitempty"`
	Preference          int      `json:"preference,omitempty"`
	Priority            int      `json:"priority,omitempty"`
	Protocol            int      `json:"protocol,omitempty"`
	Refresh             int      `json:"refresh,omitempty"`
	Regexp              string   `json:"regexp,omitempty"`
	Replacement         string   `json:"replacement,omitempty"`
	Retry               int      `json:"retry,omitempty"`
	Salt                string   `json:"salt,omitempty"`
	Serial              int      `json:"serial,omitempty"`
	Service             string   `json:"service,omitempty"`
	Signature           string   `json:"signature,omitempty"`
	Signer              string   `json:"signer,omitempty"`
	Software            string   `json:"software,omitempty"`
	Subtype             int      `json:"subtype,omitempty"`
	Target              string   `json:"target,omitempty"`
	TTL                 int      `json:"ttl,omitempty"`
	Txt                 string   `json:"txt,omitempty"`
	TypeBitmaps         string   `json:"type_bitmaps,omitempty"`
	TypeCovered         string   `json:"type_covered,omitempty"`
	Weight              uint     `json:"weight,omitempty"`
}

type ARecord struct {
	baseRecord
}

func NewARecord() ARecord {
	return ARecord{
		baseRecord{
			RecordType: "A",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record ARecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record ARecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record ARecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type AaaaRecord struct {
	baseRecord
}

func NewAaaaRecord() AaaaRecord {
	return AaaaRecord{
		baseRecord{
			RecordType: "AAAA",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record AaaaRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record AaaaRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record AaaaRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type AfsdbRecord struct {
	baseRecord
}

func NewAfsdbRecord() AfsdbRecord {
	return AfsdbRecord{
		baseRecord{
			RecordType: "AFSDB",
			fieldMap: []string{
				"active",
				"name",
				"subtype",
				"targets",
				"ttl",
			},
		},
	}
}

func (record AfsdbRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record AfsdbRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record AfsdbRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "subtype":
			record.Subtype = value.(int)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type CnameRecord struct {
	baseRecord
}

func NewCnameRecord() CnameRecord {
	return CnameRecord{
		baseRecord{
			RecordType: "CNAME",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record CnameRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record CnameRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record CnameRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type DnskeyRecord struct {
	baseRecord
}

func NewDnskeyRecord() DnskeyRecord {
	return DnskeyRecord{
		baseRecord{
			RecordType: "DNSKEY",
			fieldMap: []string{
				"active",
				"algorithm",
				"flags",
				"key",
				"name",
				"protocol",
				"targets",
			},
		},
	}
}

func (record DnskeyRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record DnskeyRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record DnskeyRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "flags":
			record.Flags = value.(int)
		case "key":
			record.Key = value.(string)
		case "name":
			record.Name = value.(string)
		case "protocol":
			record.Protocol = value.(int)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type DsRecord struct {
	baseRecord
}

func NewDsRecord() DsRecord {
	return DsRecord{
		baseRecord{
			RecordType: "DS",
			fieldMap: []string{
				"active",
				"algorithm",
				"digest",
				"digesttype",
				"key",
				"name",
				"targets",
			},
		},
	}
}

func (record DsRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record DsRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record DsRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "digest":
			record.Digest = value.(string)
		case "digesttype":
			record.DigestType = value.(int)
		case "key":
			record.Key = value.(string)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type HinfoRecord struct {
	baseRecord
}

func NewHinfoRecord() HinfoRecord {
	return HinfoRecord{
		baseRecord{
			RecordType: "HINFO",
			fieldMap: []string{
				"active",
				"hardware",
				"name",
				"software",
				"targets",
			},
		},
	}
}

func (record HinfoRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record HinfoRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record HinfoRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "hardware":
			record.Hardware = value.(string)
		case "name":
			record.Name = value.(string)
		case "software":
			record.Software = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type LocRecord struct {
	baseRecord
}

func NewLocRecord() LocRecord {
	return LocRecord{
		baseRecord{
			RecordType: "LOC",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record LocRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record LocRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record LocRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type MxRecord struct {
	baseRecord
}

func NewMxRecord() MxRecord {
	return MxRecord{
		baseRecord{
			RecordType: "MX",
			fieldMap: []string{
				"active",
				"name",
				"priority",
				"targets",
				"ttl",
			},
		},
	}
}

func (record MxRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record MxRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record MxRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "priority":
			record.Priority = value.(int)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type NaptrRecord struct {
	baseRecord
}

func NewNaptrRecord() NaptrRecord {
	return NaptrRecord{
		baseRecord{
			RecordType: "NAPTR",
			fieldMap: []string{
				"active",
				"flags",
				"name",
				"order",
				"preference",
				"regexp",
				"replacement",
				"service",
				"targets",
			},
		},
	}
}

func (record NaptrRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record NaptrRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record NaptrRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "flags":
			record.Flags = value.(int)
		case "name":
			record.Name = value.(string)
		case "order":
			record.Order = value.(int)
		case "preference":
			record.Preference = value.(int)
		case "regexp":
			record.Regexp = value.(string)
		case "replacement":
			record.Replacement = value.(string)
		case "service":
			record.Service = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type NsRecord struct {
	baseRecord
}

func NewNsRecord() NsRecord {
	return NsRecord{
		baseRecord{
			RecordType: "NS",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record NsRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record NsRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record NsRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type Nsec3Record struct {
	baseRecord
}

func NewNsec3Record() Nsec3Record {
	return Nsec3Record{
		baseRecord{
			RecordType: "NSEC3",
			fieldMap: []string{
				"active",
				"algorithm",
				"flags",
				"iterations",
				"name",
				"nexthashedownername",
				"salt",
				"targets",
				"typebitmaps",
			},
		},
	}
}

func (record Nsec3Record) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record Nsec3Record) GetAllowedFields() []string {
	return record.fieldMap
}

func (record Nsec3Record) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "flags":
			record.Flags = value.(int)
		case "iterations":
			record.Iterations = value.(int)
		case "name":
			record.Name = value.(string)
		case "nexthashedownername":
			record.NextHashedOwnerName = value.(string)
		case "salt":
			record.Salt = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "typebitmaps":
			record.TypeBitmaps = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type Nsec3paramRecord struct {
	baseRecord
}

func NewNsec3paramRecord() Nsec3paramRecord {
	return Nsec3paramRecord{
		baseRecord{
			RecordType: "NSEC3PARAM",
			fieldMap: []string{
				"active",
				"algorithm",
				"flags",
				"iterations",
				"name",
				"salt",
				"targets",
			},
		},
	}
}

func (record Nsec3paramRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record Nsec3paramRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record Nsec3paramRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "flags":
			record.Flags = value.(int)
		case "iterations":
			record.Iterations = value.(int)
		case "name":
			record.Name = value.(string)
		case "salt":
			record.Salt = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type PtrRecord struct {
	baseRecord
}

func NewPtrRecord() PtrRecord {
	return PtrRecord{
		baseRecord{
			RecordType: "PTR",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record PtrRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record PtrRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record PtrRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type RpRecord struct {
	baseRecord
}

func NewRpRecord() RpRecord {
	return RpRecord{
		baseRecord{
			RecordType: "RP",
			fieldMap: []string{
				"active",
				"mailbox",
				"name",
				"targets",
				"txt",
			},
		},
	}
}

func (record RpRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record RpRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record RpRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "mailbox":
			record.Mailbox = value.(string)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "txt":
			record.Txt = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type RrsigRecord struct {
	baseRecord
}

func NewRrsigRecord() RrsigRecord {
	return RrsigRecord{
		baseRecord{
			RecordType: "RRSIG",
			fieldMap: []string{
				"active",
				"algorithm",
				"expiration",
				"inception",
				"keytag",
				"labels",
				"name",
				"originalttl",
				"service",
				"signature",
				"signer",
				"targets",
				"typecovered",
			},
		},
	}
}

func (record RrsigRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record RrsigRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record RrsigRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "expiration":
			record.Expiration = value.(string)
		case "inception":
			record.Inception = value.(string)
		case "keytag":
			record.Keytag = value.(int)
		case "labels":
			record.Labels = value.(int)
		case "name":
			record.Name = value.(string)
		case "originalttl":
			record.OriginalTTL = value.(int)
		case "service":
			record.Service = value.(string)
		case "signature":
			record.Signature = value.(string)
		case "signer":
			record.Signer = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "typecovered":
			record.TypeCovered = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type SoaRecord struct {
	baseRecord
}

func NewSoaRecord() SoaRecord {
	return SoaRecord{
		baseRecord{
			RecordType: "SOA",
			fieldMap: []string{
				"contact",
				"expire",
				"minimum",
				"originserver",
				"refresh",
				"retry",
				"serial",
				"targets",
			},
		},
	}
}

func (record SoaRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record SoaRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record SoaRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "contact":
			record.Contact = value.(string)
		case "expire":
			record.Expire = value.(int)
		case "minimum":
			record.Minimum = value.(int)
		case "originserver":
			record.Originserver = value.(string)
		case "refresh":
			record.Refresh = value.(int)
		case "retry":
			record.Retry = value.(int)
		case "serial":
			record.Serial = value.(int)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type SpfRecord struct {
	baseRecord
}

func NewSpfRecord() SpfRecord {
	return SpfRecord{
		baseRecord{
			RecordType: "SPF",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record SpfRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record SpfRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record SpfRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type SrvRecord struct {
	baseRecord
}

func NewSrvRecord() SrvRecord {
	return SrvRecord{
		baseRecord{
			RecordType: "SRV",
			fieldMap: []string{
				"active",
				"name",
				"port",
				"priority",
				"targets",
				"ttl",
				"weight",
			},
		},
	}
}

func (record SrvRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record SrvRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record SrvRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "port":
			record.Port = value.(int)
		case "priority":
			record.Priority = value.(int)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "weight":
			record.Weight = value.(uint)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type SshfpRecord struct {
	baseRecord
}

func NewSshfpRecord() SshfpRecord {
	return SshfpRecord{
		baseRecord{
			RecordType: "SSHFP",
			fieldMap: []string{
				"active",
				"algorithm",
				"fingerprint",
				"fingerprinttype",
				"name",
				"targets",
			},
		},
	}
}

func (record SshfpRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record SshfpRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record SshfpRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "fingerprint":
			record.Fingerprint = value.(string)
		case "fingerprinttype":
			record.FingerprintType = value.(int)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

type TxtRecord struct {
	baseRecord
}

func NewTxtRecord() TxtRecord {
	return TxtRecord{
		baseRecord{
			RecordType: "TXT",
			fieldMap: []string{
				"active",
				"name",
				"targets",
				"ttl",
			},
		},
	}
}

func (record TxtRecord) Allows(field string) bool {
	return contains(record.fieldMap, field)
}

func (record TxtRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record TxtRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "active":
			record.Active = value.(bool)
		case "name":
			record.Name = value.(string)
		case "targets":
			// targets should come through as single value
			record.Target = value.(string)
		case "ttl":
			record.TTL = value.(int)
		}
		return nil
	}
	return &ValidationFailedError{fieldName: name}
}

func contains(fieldMap []string, field string) bool {
	field = strings.ToLower(field)

	for _, r := range fieldMap {
		if r == field {
			return true
		}
	}

	return false
}
