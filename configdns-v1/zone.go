package dns

import (
	"fmt"
	"log"
	"time"

	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"strings"
)

var (
	recordTypes = []string{
		"A",
		"AAAA",
		"AFSDB",
		"CNAME",
		"DNSKEY",
		"DS",
		"HINFO",
		"LOC",
		"MX",
		"NAPTR",
		"NS",
		"NSEC3",
		"NSEC3PARAM",
		"PTR",
		"RP",
		"RRSIG",
		"SOA",
		"SPF",
		"SRV",
		"SSHFP",
		"TXT",
	}
)

// Zone represents a DNS zone
type Zone struct {
	Token string `json:"token"`
	Zone  struct {
		Name       string              `json:"name,omitempty"`
		A          RecordSet           `json:"a,omitempty"`
		AAAA       RecordSet           `json:"aaaa,omitempty"`
		Afsdb      RecordSet           `json:"afsdb,omitempty"`
		Cname      RecordSet           `json:"cname,omitempty"`
		Dnskey     RecordSet           `json:"dnskey,omitempty"`
		Ds         RecordSet           `json:"ds,omitempty"`
		Hinfo      RecordSet           `json:"hinfo,omitempty"`
		Loc        RecordSet           `json:"loc,omitempty"`
		Mx         RecordSet           `json:"mx,omitempty"`
		Naptr      NaptrRecordSet      `json:"naptr,omitempty"`
		Ns         NsRecordSet         `json:"ns,omitempty"`
		Nsec3      Nsec3RecordSet      `json:"nsec3,omitempty"`
		Nsec3param Nsec3paramRecordSet `json:"nsec3param,omitempty"`
		Ptr        RecordSet           `json:"ptr,omitempty"`
		Rp         RecordSet           `json:"rp,omitempty"`
		Rrsig      RecordSet           `json:"rrsig,omitempty"`
		Soa        *Record             `json:"soa,omitempty"`
		Spf        RecordSet           `json:"spf,omitempty"`
		Srv        SrvRecordSet        `json:"srv,omitempty"`
		Sshfp      RecordSet           `json:"sshfp,omitempty"`
		Txt        RecordSet           `json:"txt,omitempty"`
	} `json:"zone"`
}

// NewZone creates a new Zone
func NewZone(hostname string) *Zone {
	zone := &Zone{Token: "new"}
	zone.Zone.Name = hostname
	return zone
}

// Save updates the Zone
func (zone *Zone) Save() error {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/config-dns/v1/zones/"+zone.Zone.Name,
		zone,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		err := client.NewAPIError(res)
		return fmt.Errorf("Unable to save record (%s)", err.Error())
	}

	for {
		updatedZone, err := GetZone(zone.Zone.Name)
		if err != nil {
			return err
		}

		if updatedZone.Token != zone.Token {
			log.Printf("[TRACE] Token updated: old: %s, new: %s", zone.Token, updatedZone.Token)
			*zone = *updatedZone
			break
		}
		log.Println("[DEBUG] Token not updated, retrying...")
		time.Sleep(time.Second)
	}

	if err != nil {
		return fmt.Errorf(errorMap[ErrFailedToSave], err.Error())
	}

	log.Printf("[INFO] Zone Saved")

	return nil
}

func (zone *Zone) GetRecordType(name string) interface{} {
	name = strings.ToUpper(name)
	switch name {
	case "A":
		return zone.Zone.A
	case "AAAA":
		return zone.Zone.AAAA
	case "AFSDB":
		return zone.Zone.Afsdb
	case "CNAME":
		return zone.Zone.Cname
	case "DNSKEY":
		return zone.Zone.Dnskey
	case "DS":
		return zone.Zone.Ds
	case "HINFO":
		return zone.Zone.Hinfo
	case "LOC":
		return zone.Zone.Loc
	case "MX":
		return zone.Zone.Mx
	case "NAPTR":
		return zone.Zone.Naptr
	case "NS":
		return zone.Zone.Ns
	case "NSEC3":
		return zone.Zone.Nsec3
	case "NSEC3PARAM":
		return zone.Zone.Nsec3param
	case "PTR":
		return zone.Zone.Ptr
	case "RP":
		return zone.Zone.Rp
	case "RRSIG":
		return zone.Zone.Rrsig
	case "SPF":
		return zone.Zone.Spf
	case "SRV":
		return zone.Zone.Srv
	case "SSHFP":
		return zone.Zone.Sshfp
	case "TXT":
		return zone.Zone.Txt
	}

	return nil
}

func (zone *Zone) SetRecord(recordPtr interface{}) error {
	// CNAME "name" must be unique in the zone, if it's not we remove clashing records
	if recordPtr.(*Record).RecordType == "CNAME" {
		records, err := zone.RemoveRecordsByName(recordPtr.(*Record).Name, recordTypes)
		if err != nil {
			return err
		}

		if records > 0 {
			log.Printf(
				"[WARN] %d Record conflicts with CNAME \"%s\", record(s) removed.",
				records,
				recordPtr.(*Record).Name,
			)
		}
	} else if recordPtr.(*Record).Name != "" {
		records, err := zone.RemoveRecordsByName(recordPtr.(*Record).Name, []string{"CNAME"})
		if err != nil {
			return err
		}

		if records > 0 {
			log.Printf(
				"[WARN] %s Record \"%s\" conflicts with CNAME, CNAME Record removed.",
				recordPtr.(*Record).RecordType,
				recordPtr.(*Record).Name,
			)
		}
	}

	switch recordPtr.(*Record).RecordType {
	case /*recordPtr.(*Record).RecordType == */ "A":
		return zone.addARecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "AAAA":
		return zone.addAaaaRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "AFSDB":
		return zone.addAfsdbRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "CNAME":
		return zone.addCnameRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "DNSKEY":
		return zone.addDnskeyRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "DS":
		return zone.addDsRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "HINFO":
		return zone.addHinfoRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "LOC":
		return zone.addLocRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "MX":
		return zone.addMxRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*NaptrRecord).RecordType == */ "NAPTR":
		return zone.addNaptrRecord(recordPtr.(*NaptrRecord), true)
	case /*recordPtr.(*NsRecord).RecordType == */ "NS":
		return zone.addNsRecord(recordPtr.(*NsRecord), true)
	case /*recordPtr.(*Nsec3Record).RecordType == */ "NSEC3":
		return zone.addNsec3Record(recordPtr.(*Nsec3Record), true)
	case /*recordPtr.(*Nsec3paramRecord).RecordType == */ "NSEC3PARAM":
		return zone.addNsec3paramRecord(recordPtr.(*Nsec3paramRecord), true)
	case /*recordPtr.(*Record).RecordType == */ "PTR":
		return zone.addPtrRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "RP":
		return zone.addRpRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "RRSIG":
		return zone.addRrsigRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "SPF":
		return zone.addSpfRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*SrvRecord).RecordType == */ "SRV":
		return zone.addSrvRecord(recordPtr.(*SrvRecord), true)
	case /*recordPtr.(*Record).RecordType == */ "SSHFP":
		return zone.addSshfpRecord(recordPtr.(*Record), true)
	case /*recordPtr.(*Record).RecordType == */ "TXT":
		return zone.addTxtRecord(recordPtr.(*Record), true)
	}

	return nil
}

func (zone *Zone) AddRecord(recordPtr interface{}) error {
	// CNAME "name" must be unique in the zone, if it's not an error is returned
	if recordPtr.(*Record).RecordType == "CNAME" {
		records := zone.FindRecordsByName(recordPtr.(*Record).Name, recordTypes)
		if len(records) > 0 {
			return fmt.Errorf(
				"Existing Record(s) conflicts with CNAME \"%s\"",
				recordPtr.(*Record).Name,
			)
		}
	} else if recordPtr.(*Record).Name != "" {
		records := zone.FindRecordsByName(recordPtr.(*Record).Name, []string{"CNAME"})
		if len(records) > 0 {
			return fmt.Errorf(
				"%s Record \"%s\" conflicts with existing CNAME \"%s\"",
				recordPtr.(*Record).RecordType,
				recordPtr.(*Record).Name,
				records[1].(*Record).Name,
			)
		}
	}

	switch recordPtr.(*Record).RecordType {
	case /*recordPtr.(*Record).RecordType == */ "A":
		return zone.addARecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "AAAA":
		return zone.addAaaaRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "AFSDB":
		return zone.addAfsdbRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "CNAME":
		return zone.addCnameRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "DNSKEY":
		return zone.addDnskeyRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "DS":
		return zone.addDsRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "HINFO":
		return zone.addHinfoRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "LOC":
		return zone.addLocRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "MX":
		return zone.addMxRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*NaptrRecord).RecordType == */ "NAPTR":
		return zone.addNaptrRecord(recordPtr.(*NaptrRecord), false)
	case /*recordPtr.(*NsRecord).RecordType == */ "NS":
		return zone.addNsRecord(recordPtr.(*NsRecord), false)
	case /*recordPtr.(*Nsec3Record).RecordType == */ "NSEC3":
		return zone.addNsec3Record(recordPtr.(*Nsec3Record), false)
	case /*recordPtr.(*Nsec3paramRecord).RecordType == */ "NSEC3PARAM":
		return zone.addNsec3paramRecord(recordPtr.(*Nsec3paramRecord), false)
	case /*recordPtr.(*Record).RecordType == */ "PTR":
		return zone.addPtrRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "RP":
		return zone.addRpRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "RRSIG":
		return zone.addRrsigRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "SPF":
		return zone.addSpfRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*SrvRecord).RecordType == */ "SRV":
		return zone.addSrvRecord(recordPtr.(*SrvRecord), false)
	case /*recordPtr.(*Record).RecordType == */ "SSHFP":
		return zone.addSshfpRecord(recordPtr.(*Record), false)
	case /*recordPtr.(*Record).RecordType == */ "TXT":
		return zone.addTxtRecord(recordPtr.(*Record), false)
	}

	return nil
}

func (zone *Zone) RemoveRecord(recordPtr interface{}) error {
	switch recordPtr.(*Record).RecordType {
	case /*recordPtr.(*Record).RecordType == */ "A":
		return zone.removeARecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "AAAA":
		return zone.removeAaaaRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "AFSDB":
		return zone.removeAfsdbRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "CNAME":
		return zone.removeCnameRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "DNSKEY":
		return zone.removeDnskeyRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "DS":
		return zone.removeDsRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "HINFO":
		return zone.removeHinfoRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "LOC":
		return zone.removeLocRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "MX":
		return zone.removeMxRecord(recordPtr.(*Record))
	case /*recordPtr.(*NaptrRecord).RecordType == */ "NAPTR":
		return zone.removeNaptrRecord(recordPtr.(*NaptrRecord))
	case /*recordPtr.(*NsRecord).RecordType == */ "NS":
		return zone.removeNsRecord(recordPtr.(*NsRecord))
	case /*recordPtr.(*Nsec3Record).RecordType == */ "NSEC3":
		return zone.removeNsec3Record(recordPtr.(*Nsec3Record))
	case /*recordPtr.(*Nsec3paramRecord).RecordType == */ "NSEC3PARAM":
		return zone.removeNsec3paramRecord(recordPtr.(*Nsec3paramRecord))
	case /*recordPtr.(*Record).RecordType == */ "PTR":
		return zone.removePtrRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "RP":
		return zone.removeRpRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "RRSIG":
		return zone.removeRrsigRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "SPF":
		return zone.removeSpfRecord(recordPtr.(*Record))
	case /*recordPtr.(*SrvRecord).RecordType == */ "SRV":
		return zone.removeSrvRecord(recordPtr.(*SrvRecord))
	case /*recordPtr.(*Record).RecordType == */ "SSHFP":
		return zone.removeSshfpRecord(recordPtr.(*Record))
	case /*recordPtr.(*Record).RecordType == */ "TXT":
		return zone.removeTxtRecord(recordPtr.(*Record))
	}

	return nil
}

func (zone *Zone) RemoveRecordsByName(name string, filterRecordTypes []string) (count int, err error) {
	records := zone.FindRecordsByName(name, filterRecordTypes)
	for _, record := range records {
		err := zone.RemoveRecord(record)
		if err != nil {
			return count, err
		}
	}

	return len(records), nil
}

func (zone *Zone) FindRecordsByName(name string, filterRecordTypes []string) []interface{} {
	var records []interface{}

	name = strings.ToLower(name)
	if len(filterRecordTypes) == 0 {
		filterRecordTypes = recordTypes
	}

	for _, recordType := range filterRecordTypes {
		for _, record := range zone.GetRecordType(recordType).(RecordSet) {
			if strings.ToLower(record.Name) == name {
				records = append(records, record)
			}
		}
	}

	return records
}

func (zone *Zone) addARecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.A {
			if r.Name == record.Name {
				zone.Zone.A[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.A = append(zone.Zone.A, record)
	}

	return nil
}

func (zone *Zone) addAaaaRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.AAAA {
			if r.Name == record.Name {
				zone.Zone.AAAA[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.AAAA = append(zone.Zone.AAAA, record)
	}

	return nil
}

func (zone *Zone) addAfsdbRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Afsdb {
			if r.Name == record.Name {
				zone.Zone.Afsdb[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Afsdb = append(zone.Zone.Afsdb, record)
	}

	return nil
}

func (zone *Zone) addCnameRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Cname {
			if r.Name == record.Name {
				zone.Zone.Cname[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Cname = append(zone.Zone.Cname, record)
	}

	return nil
}

func (zone *Zone) addDnskeyRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Dnskey {
			if r.Name == record.Name {
				zone.Zone.Dnskey[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Dnskey = append(zone.Zone.Dnskey, record)
	}

	return nil
}

func (zone *Zone) addDsRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Ds {
			if r.Name == record.Name {
				zone.Zone.Ds[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Ds = append(zone.Zone.Ds, record)
	}

	return nil
}

func (zone *Zone) addHinfoRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Hinfo {
			if r.Name == record.Name {
				zone.Zone.Hinfo[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Hinfo = append(zone.Zone.Hinfo, record)
	}

	return nil
}

func (zone *Zone) addLocRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Loc {
			if r.Name == record.Name {
				zone.Zone.Loc[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Loc = append(zone.Zone.Loc, record)
	}

	return nil
}

func (zone *Zone) addMxRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Mx {
			if r.Name == record.Name {
				zone.Zone.Mx[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Mx = append(zone.Zone.Mx, record)
	}

	return nil
}

func (zone *Zone) addNaptrRecord(record *NaptrRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Naptr {
			if r.Name == record.Name {
				zone.Zone.Naptr[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Naptr = append(zone.Zone.Naptr, record)
	}

	return nil
}

func (zone *Zone) addNsRecord(record *NsRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Ns {
			if r.Name == record.Name {
				zone.Zone.Ns[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Ns = append(zone.Zone.Ns, record)
	}

	return nil
}

func (zone *Zone) addNsec3Record(record *Nsec3Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Nsec3 {
			if r.Name == record.Name {
				zone.Zone.Nsec3[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Nsec3 = append(zone.Zone.Nsec3, record)
	}

	return nil
}

func (zone *Zone) addNsec3paramRecord(record *Nsec3paramRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Nsec3param {
			if r.Name == record.Name {
				zone.Zone.Nsec3param[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Nsec3param = append(zone.Zone.Nsec3param, record)
	}

	return nil
}

func (zone *Zone) addPtrRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Ptr {
			if r.Name == record.Name {
				zone.Zone.Ptr[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Ptr = append(zone.Zone.Ptr, record)
	}

	return nil
}

func (zone *Zone) addRpRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Rp {
			if r.Name == record.Name {
				zone.Zone.Rp[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Rp = append(zone.Zone.Rp, record)
	}

	return nil
}

func (zone *Zone) addRrsigRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Rrsig {
			if r.Name == record.Name {
				zone.Zone.Rrsig[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Rrsig = append(zone.Zone.Rrsig, record)
	}

	return nil
}

func (zone *Zone) addSoaRecord(record *Record, replace bool) error {
	zone.Zone.Soa = record
	return nil
}

func (zone *Zone) addSpfRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Spf {
			if r.Name == record.Name {
				zone.Zone.Spf[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Spf = append(zone.Zone.Spf, record)
	}

	return nil
}

func (zone *Zone) addSrvRecord(record *SrvRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Srv {
			if r.Name == record.Name {
				zone.Zone.Srv[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Srv = append(zone.Zone.Srv, record)
	}

	return nil
}

func (zone *Zone) addSshfpRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Sshfp {
			if r.Name == record.Name {
				zone.Zone.Sshfp[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Sshfp = append(zone.Zone.Sshfp, record)
	}

	return nil
}

func (zone *Zone) addTxtRecord(record *Record, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Txt {
			if r.Name == record.Name {
				zone.Zone.Txt[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Txt = append(zone.Zone.Txt, record)
	}

	return nil
}

func (zone *Zone) removeARecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.A {
		if r == record {
			records := zone.Zone.A[:key]
			zone.Zone.A = append(records, zone.Zone.A[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("A Record not found")
	}

	return nil
}

func (zone *Zone) removeAaaaRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.AAAA {
		if r == record {
			records := zone.Zone.AAAA[:key]
			zone.Zone.AAAA = append(records, zone.Zone.AAAA[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("AAAA Record not found")
	}

	return nil
}

func (zone *Zone) removeAfsdbRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Afsdb {
		if r == record {
			records := zone.Zone.Afsdb[:key]
			zone.Zone.Afsdb = append(records, zone.Zone.Afsdb[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Afsdb Record not found")
	}

	return nil
}

func (zone *Zone) removeCnameRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Cname {
		if r == record {
			records := zone.Zone.Cname[:key]
			zone.Zone.Cname = append(records, zone.Zone.Cname[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Cname Record not found")
	}

	return nil
}

func (zone *Zone) removeDnskeyRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Dnskey {
		if r == record {
			records := zone.Zone.Dnskey[:key]
			zone.Zone.Dnskey = append(records, zone.Zone.Dnskey[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Dnskey Record not found")
	}

	return nil
}

func (zone *Zone) removeDsRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Ds {
		if r == record {
			records := zone.Zone.Ds[:key]
			zone.Zone.Ds = append(records, zone.Zone.Ds[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Ds Record not found")
	}

	return nil
}

func (zone *Zone) removeHinfoRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Hinfo {
		if r == record {
			records := zone.Zone.Hinfo[:key]
			zone.Zone.Hinfo = append(records, zone.Zone.Hinfo[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Hinfo Record not found")
	}

	return nil
}

func (zone *Zone) removeLocRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Loc {
		if r == record {
			records := zone.Zone.Loc[:key]
			zone.Zone.Loc = append(records, zone.Zone.Loc[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Loc Record not found")
	}

	return nil
}

func (zone *Zone) removeMxRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Mx {
		if r == record {
			records := zone.Zone.Mx[:key]
			zone.Zone.Mx = append(records, zone.Zone.Mx[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Mx Record not found")
	}

	return nil
}

func (zone *Zone) removeNaptrRecord(record *NaptrRecord) error {
	var found bool
	for key, r := range zone.Zone.Naptr {
		if r == record {
			records := zone.Zone.Naptr[:key]
			zone.Zone.Naptr = append(records, zone.Zone.Naptr[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Naptr Record not found")
	}

	return nil
}

func (zone *Zone) removeNsRecord(record *NsRecord) error {
	var found bool
	for key, r := range zone.Zone.Ns {
		if r == record {
			records := zone.Zone.Ns[:key]
			zone.Zone.Ns = append(records, zone.Zone.Ns[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Ns Record not found")
	}

	return nil
}

func (zone *Zone) removeNsec3Record(record *Nsec3Record) error {
	var found bool
	for key, r := range zone.Zone.Nsec3 {
		if r == record {
			records := zone.Zone.Nsec3[:key]
			zone.Zone.Nsec3 = append(records, zone.Zone.Nsec3[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Nsec3 Record not found")
	}

	return nil
}

func (zone *Zone) removeNsec3paramRecord(record *Nsec3paramRecord) error {
	var found bool
	for key, r := range zone.Zone.Nsec3param {
		if r == record {
			records := zone.Zone.Nsec3param[:key]
			zone.Zone.Nsec3param = append(records, zone.Zone.Nsec3param[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Nsec3param Record not found")
	}

	return nil
}

func (zone *Zone) removePtrRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Ptr {
		if r == record {
			records := zone.Zone.Ptr[:key]
			zone.Zone.Ptr = append(records, zone.Zone.Ptr[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Ptr Record not found")
	}

	return nil
}

func (zone *Zone) removeRpRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Rp {
		if r == record {
			records := zone.Zone.Rp[:key]
			zone.Zone.Rp = append(records, zone.Zone.Rp[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Rp Record not found")
	}

	return nil
}

func (zone *Zone) removeRrsigRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Rrsig {
		if r == record {
			records := zone.Zone.Rrsig[:key]
			zone.Zone.Rrsig = append(records, zone.Zone.Rrsig[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Rrsig Record not found")
	}

	return nil
}

func (zone *Zone) removeSoaRecord(record *Record) error {
	zone.Zone.Soa = record
	return nil
}

func (zone *Zone) removeSpfRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Spf {
		if r == record {
			records := zone.Zone.Spf[:key]
			zone.Zone.Spf = append(records, zone.Zone.Spf[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Spf Record not found")
	}

	return nil
}

func (zone *Zone) removeSrvRecord(record *SrvRecord) error {
	var found bool
	for key, r := range zone.Zone.Srv {
		if r == record {
			records := zone.Zone.Srv[:key]
			zone.Zone.Srv = append(records, zone.Zone.Srv[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Srv Record not found")
	}

	return nil
}

func (zone *Zone) removeSshfpRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Sshfp {
		if r == record {
			records := zone.Zone.Sshfp[:key]
			zone.Zone.Sshfp = append(records, zone.Zone.Sshfp[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Sshfp Record not found")
	}

	return nil
}

func (zone *Zone) removeTxtRecord(record *Record) error {
	var found bool
	for key, r := range zone.Zone.Txt {
		if r == record {
			records := zone.Zone.Txt[:key]
			zone.Zone.Txt = append(records, zone.Zone.Txt[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("Txt Record not found")
	}

	return nil
}

func (zone *Zone) PreMarshalJSON() error {
	zone.Zone.Soa.Serial = int(time.Now().Unix())
	return nil
}
