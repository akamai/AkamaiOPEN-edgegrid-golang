package dnsv2

import (
	"strings"
	"time"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// All record types (below) must implement the DNSRecord interface
// This allows the record to be used dynamically in slices - see the Zone struct definition in zone.go
//
// The record types implemented and their fields are as defined here
// https://developer.akamai.com/api/luna/config-dns/data.html

type RecordBody struct {
	Name     string   `json:"name,omitempty"`
	RecordType     string   `json:"type,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Target   []string   `json:"rdata,omitempty"`
	Subtype  int      `json:"subtype,omitempty"` //AfsdbRecord
	Flags     int      `json:"flags,omitempty"` //DnskeyRecord Nsec3paramRecord
	Protocol  int      `json:"protocol,omitempty"` //DnskeyRecord
	Algorithm int      `json:"algorithm,omitempty"` //DnskeyRecord DsRecord Nsec3paramRecord RrsigRecord SshfpRecord
	Key       string   `json:"key,omitempty"` //DnskeyRecord
  Keytag     int      `json:"keytag,omitempty"` //DsRecord RrsigRecord
  DigestType int      `json:"digest_type,omitempty"` //DsRecord
  Digest     string   `json:"digest,omitempty"` //DsRecord
	Hardware string   `json:"hardware,omitempty"` //HinfoRecord
	Software string   `json:"software,omitempty"` //HinfoRecord
	Priority int      `json:"priority,omitempty"` //MxRecord SrvRecord
	Order       uint16   `json:"order,omitempty"` //NaptrRecord
	Preference  uint16   `json:"preference,omitempty"` //NaptrRecord
	FlagsNaptr  string   `json:"flags,omitempty"` //NaptrRecord
	Service     string   `json:"service,omitempty"` //NaptrRecord
	Regexp      string   `json:"regexp,omitempty"` //NaptrRecord
	Replacement string   `json:"replacement,omitempty"` //NaptrRecord
	Iterations          int      `json:"iterations,omitempty"` //Nsec3Record Nsec3paramRecord
	Salt                string   `json:"salt,omitempty"`  //Nsec3Record Nsec3paramRecord
	NextHashedOwnerName string   `json:"next_hashed_owner_name,omitempty"` //Nsec3Record
	TypeBitmaps         string   `json:"type_bitmaps,omitempty"` //Nsec3Record
	Mailbox  string   `json:"mailbox,omitempty"` //RpRecord
	Txt      string   `json:"txt,omitempty"` //RpRecord
	TypeCovered string   `json:"type_covered,omitempty"` //RrsigRecord
	OriginalTTL int      `json:"original_ttl,omitempty"` //RrsigRecord
	Expiration  string   `json:"expiration,omitempty"` //RrsigRecord
	Inception   string   `json:"inception,omitempty"` //RrsigRecord
	Signer      string   `json:"signer,omitempty"` //RrsigRecord
	Signature   string   `json:"signature,omitempty"` //RrsigRecord
	Labels      int      `json:"labels,omitempty"` //RrsigRecord
	Weight   uint16   `json:"weight,omitempty"` //SrvRecord
	Port     uint16   `json:"port,omitempty"` //SrvRecord
	FingerprintType int      `json:"fingerprint_type,omitempty"` //SshfpRecord
	Fingerprint     string   `json:"fingerprint,omitempty"` //SshfpRecord
}

func NewRecordBody(params RecordBody) *RecordBody {
	recordbody := &RecordBody{Name: params.Name }
	return recordbody
}

func (record *RecordBody ) Save(zone string) error {

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/config-dns/v2/zones/"+zone+"/names/"+record.Name+"/types/"+record.RecordType,
		record,
	)
	if err != nil {
		return err
	}
	res, err := client.Do(Config, req)

	// Network error
	if err != nil {
		return &ZoneError{
			zoneName:         zone,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return &ZoneError{zoneName: zone, apiErrorMessage: err.Detail, err: err}
	}

	return nil
}


func (record *RecordBody ) Delete(zone string) error {


	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		"/config-dns/v2/zones/"+zone+"/names/"+record.Name+"/types/"+record.RecordType,
		record,
	)
	if err != nil {
		return err
	}
	res, err := client.Do(Config, req)

	// Network error
	if err != nil {
		return &ZoneError{
			zoneName:         zone,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return &ZoneError{zoneName: zone, apiErrorMessage: err.Detail, err: err}
	}

	return nil
}


type SoaRecord struct {
	fieldMap       []string `json:"-"`
	originalSerial uint     `json:"-"`
	TTL            int      `json:"ttl,omitempty"`
	Originserver   string   `json:"originserver,omitempty"`
	Contact        string   `json:"contact,omitempty"`
	Serial         uint     `json:"serial,omitempty"`
	Refresh        int      `json:"refresh,omitempty"`
	Retry          int      `json:"retry,omitempty"`
	Expire         int      `json:"expire,omitempty"`
	Minimum        uint     `json:"minimum,omitempty"`
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
			v, ok := value.(int)
			if ok {
				record.TTL = v
				return nil
			}
		case "originserver":
			v, ok := value.(string)
			if ok {
				record.Originserver = v
				return nil
			}
		case "contact":
			v, ok := value.(string)
			if ok {
				record.Contact = v
				return nil
			}
		case "serial":
			v, ok := value.(uint)
			if ok {
				record.Serial = v
				return nil
			}
		case "refresh":
			v, ok := value.(int)
			if ok {
				record.Refresh = v
				return nil
			}
		case "retry":
			v, ok := value.(int)
			if ok {
				record.Retry = v
				return nil
			}
		case "expire":
			v, ok := value.(int)
			if ok {
				record.Expire = v
				return nil
			}
		case "minimum":
			v, ok := value.(uint)
			if ok {
				record.Minimum = v
				return nil
			}
		}
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


func contains(fieldMap []string, field string) bool {
	field = strings.ToLower(field)

	for _, r := range fieldMap {
		if r == field {
			return true
		}
	}

	return false
}
