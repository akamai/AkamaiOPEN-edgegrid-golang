package dns

import (
	"strings"
	"time"
)

// All record types (below) must implement the DNSRecord interface
// This allows the record to be used dynamically in slices - see the Zone struct definition in zone.go
//
// The record types implemented and their fields are as defined here
// https://developer.akamai.com/api/luna/config-dns/data.html
type DNSRecord interface {
	// Get the list of allowed fields for this type
	GetAllowedFields() []string
	// Set a field on the struct, which check for valid fields
	SetField(name string, value interface{}) error
	// Translate struct properties to a map
	ToMap() map[string]interface{}
}

type ARecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewARecord() *ARecord {
	return &ARecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *ARecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *ARecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
	}
	return &RecordError{fieldName: name}
}

func (record *ARecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type AaaaRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewAaaaRecord() *AaaaRecord {
	return &AaaaRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *AaaaRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *AaaaRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *AaaaRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type AfsdbRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
	Subtype  int      `json:"subtype,omitempty"`
}

func NewAfsdbRecord() *AfsdbRecord {
	return &AfsdbRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
			"subtype",
		},
	}
}

func (record *AfsdbRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *AfsdbRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		case "subtype":
			record.Subtype = value.(int)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *AfsdbRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    record.Name,
		"ttl":     record.TTL,
		"active":  record.Active,
		"target":  record.Target,
		"subtype": record.Subtype,
	}
}

type CnameRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewCnameRecord() *CnameRecord {
	return &CnameRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *CnameRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *CnameRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *CnameRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type DnskeyRecord struct {
	fieldMap  []string `json:"-"`
	Name      string   `json:"name,omitempty"`
	TTL       int      `json:"ttl,omitempty"`
	Active    bool     `json:"active,omitempty"`
	Flags     int      `json:"flags,omitempty"`
	Protocol  int      `json:"protocol,omitempty"`
	Algorithm int      `json:"algorithm,omitempty"`
	Key       string   `json:"key,omitempty"`
}

func NewDnskeyRecord() *DnskeyRecord {
	return &DnskeyRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"flags",
			"protocol",
			"algorithm",
			"key",
		},
	}
}

func (record *DnskeyRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *DnskeyRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "flags":
			record.Flags = value.(int)
		case "protocol":
			record.Protocol = value.(int)
		case "algorithm":
			record.Algorithm = value.(int)
		case "key":
			record.Key = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *DnskeyRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":      record.Name,
		"ttl":       record.TTL,
		"active":    record.Active,
		"flags":     record.Flags,
		"protocol":  record.Protocol,
		"algorithm": record.Algorithm,
		"key":       record.Key,
	}
}

type DsRecord struct {
	fieldMap   []string `json:"-"`
	Name       string   `json:"name,omitempty"`
	TTL        int      `json:"ttl,omitempty"`
	Active     bool     `json:"active,omitempty"`
	Keytag     int      `json:"keytag,omitempty"`
	Algorithm  int      `json:"algorithm,omitempty"`
	DigestType int      `json:"digest_type,omitempty"`
	Digest     string   `json:"digest,omitempty"`
}

func NewDsRecord() *DsRecord {
	return &DsRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"keytag",
			"algorithm",
			"digesttype",
			"digest",
		},
	}
}

func (record *DsRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *DsRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "keytag":
			record.Keytag = value.(int)
		case "algorithm":
			record.Algorithm = value.(int)
		case "digesttype":
			record.DigestType = value.(int)
		case "digest":
			record.Digest = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *DsRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       record.Name,
		"ttl":        record.TTL,
		"active":     record.Active,
		"keytag":     record.Keytag,
		"algorithm":  record.Algorithm,
		"digesttype": record.DigestType,
		"digest":     record.DigestType,
	}
}

type HinfoRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Hardware string   `json:"hardware,omitempty"`
	Software string   `json:"software,omitempty"`
}

func NewHinfoRecord() *HinfoRecord {
	return &HinfoRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"hardware",
			"software",
		},
	}
}

func (record *HinfoRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *HinfoRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "hardware":
			record.Hardware = value.(string)
		case "software":
			record.Software = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *HinfoRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":     record.Name,
		"ttl":      record.TTL,
		"active":   record.Active,
		"hardware": record.Hardware,
		"software": record.Software,
	}
}

type LocRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewLocRecord() *LocRecord {
	return &LocRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *LocRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *LocRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *LocRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type MxRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
	Priority int      `json:"priority,omitempty"`
}

func NewMxRecord() *MxRecord {
	return &MxRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
			"priority",
		},
	}
}

func (record *MxRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *MxRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		case "priority":
			record.Priority = value.(int)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *MxRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":     record.Name,
		"ttl":      record.TTL,
		"active":   record.Active,
		"target":   record.Target,
		"priority": record.Priority,
	}
}

type NaptrRecord struct {
	fieldMap    []string `json:"-"`
	Name        string   `json:"name,omitempty"`
	TTL         int      `json:"ttl,omitempty"`
	Active      bool     `json:"active,omitempty"`
	Order       int      `json:"order,omitempty"`
	Preference  int      `json:"preference,omitempty"`
	Flags       int      `json:"flags,omitempty"`
	Service     string   `json:"service,omitempty"`
	Regexp      string   `json:"regexp,omitempty"`
	Replacement string   `json:"replacement,omitempty"`
}

func NewNaptrRecord() *NaptrRecord {
	return &NaptrRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"order",
			"preference",
			"flags",
			"service",
			"regexp",
			"replacement",
		},
	}
}

func (record *NaptrRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *NaptrRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "order":
			record.Order = value.(int)
		case "preference":
			record.Preference = value.(int)
		case "flags":
			record.Flags = value.(int)
		case "service":
			record.Service = value.(string)
		case "regexp":
			record.Regexp = value.(string)
		case "replacement":
			record.Replacement = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *NaptrRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        record.Name,
		"ttl":         record.TTL,
		"active":      record.Active,
		"order":       record.Order,
		"preference":  record.Preference,
		"flags":       record.Flags,
		"service":     record.Service,
		"regexp":      record.Regexp,
		"replacement": record.Replacement,
	}
}

type NsRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewNsRecord() *NsRecord {
	return &NsRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *NsRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *NsRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *NsRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type Nsec3Record struct {
	fieldMap            []string `json:"-"`
	Name                string   `json:"name,omitempty"`
	TTL                 int      `json:"ttl,omitempty"`
	Active              bool     `json:"active,omitempty"`
	Algorithm           int      `json:"algorithm,omitempty"`
	Flags               int      `json:"flags,omitempty"`
	Iterations          int      `json:"iterations,omitempty"`
	Salt                string   `json:"salt,omitempty"`
	NextHashedOwnerName string   `json:"next_hashed_owner_name,omitempty"`
	TypeBitmaps         string   `json:"type_bitmaps,omitempty"`
}

func NewNsec3Record() *Nsec3Record {
	return &Nsec3Record{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"algorithm",
			"flags",
			"iterations",
			"salt",
			"nexthashedownername",
			"typebitmaps",
		},
	}
}

func (record *Nsec3Record) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *Nsec3Record) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "flags":
			record.Flags = value.(int)
		case "iterations":
			record.Iterations = value.(int)
		case "salt":
			record.Salt = value.(string)
		case "nexthashedownername":
			record.NextHashedOwnerName = value.(string)
		case "typebitmaps":
			record.TypeBitmaps = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *Nsec3Record) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":                record.Name,
		"ttl":                 record.TTL,
		"active":              record.Active,
		"algorithm":           record.Algorithm,
		"flags":               record.Flags,
		"iterations":          record.Iterations,
		"salt":                record.Salt,
		"nexthashedownername": record.NextHashedOwnerName,
		"typebitmaps":         record.TypeBitmaps,
	}
}

type Nsec3paramRecord struct {
	fieldMap   []string `json:"-"`
	Name       string   `json:"name,omitempty"`
	TTL        int      `json:"ttl,omitempty"`
	Active     bool     `json:"active,omitempty"`
	Algorithm  int      `json:"algorithm,omitempty"`
	Flags      int      `json:"flags,omitempty"`
	Iterations int      `json:"iterations,omitempty"`
	Salt       string   `json:"salt,omitempty"`
}

func NewNsec3paramRecord() *Nsec3paramRecord {
	return &Nsec3paramRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"algorithm",
			"flags",
			"iterations",
			"salt",
		},
	}
}

func (record *Nsec3paramRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *Nsec3paramRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "flags":
			record.Flags = value.(int)
		case "iterations":
			record.Iterations = value.(int)
		case "salt":
			record.Salt = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *Nsec3paramRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       record.Name,
		"ttl":        record.TTL,
		"active":     record.Active,
		"algorithm":  record.Algorithm,
		"flags":      record.Flags,
		"iterations": record.Iterations,
		"salt":       record.Salt,
	}
}

type PtrRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewPtrRecord() *PtrRecord {
	return &PtrRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *PtrRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *PtrRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *PtrRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type RpRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Mailbox  string   `json:"mailbox,omitempty"`
	Txt      string   `json:"txt,omitempty"`
}

func NewRpRecord() *RpRecord {
	return &RpRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"mailbox",
			"txt",
		},
	}
}

func (record *RpRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *RpRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "mailbox":
			record.Mailbox = value.(string)
		case "txt":
			record.Txt = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *RpRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    record.Name,
		"ttl":     record.TTL,
		"active":  record.Active,
		"mailbox": record.Mailbox,
		"txt":     record.Txt,
	}
}

type RrsigRecord struct {
	fieldMap    []string `json:"-"`
	Name        string   `json:"name,omitempty"`
	TTL         int      `json:"ttl,omitempty"`
	Active      bool     `json:"active,omitempty"`
	TypeCovered string   `json:"type_covered,omitempty"`
	Algorithm   int      `json:"algorithm,omitempty"`
	OriginalTTL int      `json:"original_ttl,omitempty"`
	Expiration  string   `json:"expiration,omitempty"`
	Inception   string   `json:"inception,omitempty"`
	Keytag      int      `json:"keytag,omitempty"`
	Signer      string   `json:"signer,omitempty"`
	Signature   string   `json:"signature,omitempty"`
	Labels      int      `json:"labels,omitempty"`
}

func NewRrsigRecord() *RrsigRecord {
	return &RrsigRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"typecovered",
			"algorithm",
			"originalttl",
			"expiration",
			"inception",
			"keytag",
			"signer",
			"signature",
			"labels",
		},
	}
}

func (record *RrsigRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *RrsigRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "typecovered":
			record.TypeCovered = value.(string)
		case "algorithm":
			record.Algorithm = value.(int)
		case "originalttl":
			record.OriginalTTL = value.(int)
		case "expiration":
			record.Expiration = value.(string)
		case "inception":
			record.Inception = value.(string)
		case "keytag":
			record.Keytag = value.(int)
		case "signer":
			record.Signer = value.(string)
		case "signature":
			record.Signature = value.(string)
		case "labels":
			record.Labels = value.(int)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *RrsigRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        record.Name,
		"ttl":         record.TTL,
		"active":      record.Active,
		"typecovered": record.TypeCovered,
		"algorithm":   record.Algorithm,
		"originalttl": record.OriginalTTL,
		"expiration":  record.Expiration,
		"inception":   record.Inception,
		"keytag":      record.Keytag,
		"signer":      record.Signer,
		"signature":   record.Signature,
		"labels":      record.Labels,
	}
}

type SoaRecord struct {
	fieldMap     []string `json:"-"`
	TTL          int      `json:"ttl,omitempty"`
	Originserver string   `json:"originserver,omitempty"`
	Contact      string   `json:"contact,omitempty"`
	Serial       int      `json:"serial,omitempty"`
	Refresh      int      `json:"refresh,omitempty"`
	Retry        int      `json:"retry,omitempty"`
	Expire       int      `json:"expire,omitempty"`
	Minimum      int      `json:"minimum,omitempty"`
}

func NewSoaRecord() *SoaRecord {
	r := &SoaRecord{
		fieldMap: []string{
			"ttl",
			"originserver",
			"contact",
			"serial",
			"refresh",
			"retry",
			"expire",
			"minimum",
		},
	}
	r.SetField("serial", int(time.Now().Unix()))
	return r
}

func (record *SoaRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *SoaRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "ttl":
			record.TTL = value.(int)
		case "originserver":
			record.Originserver = value.(string)
		case "contact":
			record.Contact = value.(string)
		case "serial":
			record.Serial = value.(int)
		case "refresh":
			record.Refresh = value.(int)
		case "retry":
			record.Retry = value.(int)
		case "expire":
			record.Expire = value.(int)
		case "minimum":
			record.Minimum = value.(int)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *SoaRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ttl":          record.TTL,
		"originserver": record.Originserver,
		"contact":      record.Contact,
		"serial":       record.Serial,
		"refresh":      record.Refresh,
		"retry":        record.Retry,
		"expire":       record.Expire,
		"minimum":      record.Minimum,
	}
}

type SpfRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewSpfRecord() *SpfRecord {
	return &SpfRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *SpfRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *SpfRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *SpfRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
}

type SrvRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
	Priority int      `json:"priority,omitempty"`
	Weight   uint     `json:"weight,omitempty"`
	Port     int      `json:"port,omitempty"`
}

func NewSrvRecord() *SrvRecord {
	return &SrvRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
			"priority",
			"weight",
			"port",
		},
	}
}

func (record *SrvRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *SrvRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		case "priority":
			record.Priority = value.(int)
		case "weight":
			record.Weight = value.(uint)
		case "port":
			record.Port = value.(int)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *SrvRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":     record.Name,
		"ttl":      record.TTL,
		"active":   record.Active,
		"target":   record.Target,
		"priority": record.Priority,
		"weight":   record.Weight,
		"port":     record.Port,
	}
}

type SshfpRecord struct {
	fieldMap        []string `json:"-"`
	Name            string   `json:"name,omitempty"`
	TTL             int      `json:"ttl,omitempty"`
	Active          bool     `json:"active,omitempty"`
	Algorithm       int      `json:"algorithm,omitempty"`
	FingerprintType int      `json:"fingerprint_type,omitempty"`
	Fingerprint     string   `json:"fingerprint,omitempty"`
}

func NewSshfpRecord() *SshfpRecord {
	return &SshfpRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"algorithm",
			"fingerprinttype",
			"fingerprint",
		},
	}
}

func (record *SshfpRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *SshfpRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "algorithm":
			record.Algorithm = value.(int)
		case "fingerprinttype":
			record.FingerprintType = value.(int)
		case "fingerprint":
			record.Fingerprint = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *SshfpRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":            record.Name,
		"ttl":             record.TTL,
		"active":          record.Active,
		"algorithm":       record.Algorithm,
		"fingerprinttype": record.FingerprintType,
		"fingerprint":     record.Fingerprint,
	}
}

type TxtRecord struct {
	fieldMap []string `json:"-"`
	Name     string   `json:"name,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   string   `json:"target,omitempty"`
}

func NewTxtRecord() *TxtRecord {
	return &TxtRecord{
		fieldMap: []string{
			"name",
			"ttl",
			"active",
			"target",
		},
	}
}

func (record *TxtRecord) GetAllowedFields() []string {
	return record.fieldMap
}

func (record *TxtRecord) SetField(name string, value interface{}) error {
	if contains(record.fieldMap, name) {
		switch name {
		case "name":
			record.Name = value.(string)
		case "ttl":
			record.TTL = value.(int)
		case "active":
			record.Active = value.(bool)
		case "target":
			record.Target = value.(string)
		}
		return nil
	}
	return &RecordError{fieldName: name}
}

func (record *TxtRecord) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"ttl":    record.TTL,
		"active": record.Active,
		"target": record.Target,
	}
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
