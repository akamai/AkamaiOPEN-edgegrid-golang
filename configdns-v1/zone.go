package dns

import (
	"fmt"
	"log"
	"time"

	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
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
		Name       string    `json:"name,omitempty"`
		A          RecordSet `json:"a,omitempty"`
		Aaaa       RecordSet `json:"aaaa,omitempty"`
		Afsdb      RecordSet `json:"afsdb,omitempty"`
		Cname      RecordSet `json:"cname,omitempty"`
		Dnskey     RecordSet `json:"dnskey,omitempty"`
		Ds         RecordSet `json:"ds,omitempty"`
		Hinfo      RecordSet `json:"hinfo,omitempty"`
		Loc        RecordSet `json:"loc,omitempty"`
		Mx         RecordSet `json:"mx,omitempty"`
		Naptr      RecordSet `json:"naptr,omitempty"`
		Ns         RecordSet `json:"ns,omitempty"`
		Nsec3      RecordSet `json:"nsec3,omitempty"`
		Nsec3param RecordSet `json:"nsec3param,omitempty"`
		Ptr        RecordSet `json:"ptr,omitempty"`
		Rp         RecordSet `json:"rp,omitempty"`
		Rrsig      RecordSet `json:"rrsig,omitempty"`
		Soa        SoaRecord `json:"soa,omitempty"`
		Spf        RecordSet `json:"spf,omitempty"`
		Srv        RecordSet `json:"srv,omitempty"`
		Sshfp      RecordSet `json:"sshfp,omitempty"`
		Txt        RecordSet `json:"txt,omitempty"`
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
		return &FailedToSaveError{err: err}
	}

	log.Printf("[INFO] Zone Saved")

	return nil
}

func (zone *Zone) AddRecord(recordPtr interface{}) error {
	switch recordPtr.(type) {
	case ARecord:
		return zone.addARecord(recordPtr.(*ARecord), true)
	case AaaaRecord:
		return zone.addAaaaRecord(recordPtr.(*AaaaRecord), true)
	case AfsdbRecord:
		return zone.addAfsdbRecord(recordPtr.(*AfsdbRecord), true)
	case CnameRecord:
		return zone.addCnameRecord(recordPtr.(*CnameRecord), true)
	case DnskeyRecord:
		return zone.addDnskeyRecord(recordPtr.(*DnskeyRecord), true)
	case DsRecord:
		return zone.addDsRecord(recordPtr.(*DsRecord), true)
	case HinfoRecord:
		return zone.addHinfoRecord(recordPtr.(*HinfoRecord), true)
	case LocRecord:
		return zone.addLocRecord(recordPtr.(*LocRecord), true)
	case MxRecord:
		return zone.addMxRecord(recordPtr.(*MxRecord), true)
	case NaptrRecord:
		return zone.addNaptrRecord(recordPtr.(*NaptrRecord), true)
	case NsRecord:
		return zone.addNsRecord(recordPtr.(*NsRecord), true)
	case Nsec3Record:
		return zone.addNsec3Record(recordPtr.(*Nsec3Record), true)
	case Nsec3paramRecord:
		return zone.addNsec3paramRecord(recordPtr.(*Nsec3paramRecord), true)
	case PtrRecord:
		return zone.addPtrRecord(recordPtr.(*PtrRecord), true)
	case RpRecord:
		return zone.addRpRecord(recordPtr.(*RpRecord), true)
	case RrsigRecord:
		return zone.addRrsigRecord(recordPtr.(*RrsigRecord), true)
	case SpfRecord:
		return zone.addSpfRecord(recordPtr.(*SpfRecord), true)
	case SrvRecord:
		return zone.addSrvRecord(recordPtr.(*SrvRecord), true)
	case SshfpRecord:
		return zone.addSshfpRecord(recordPtr.(*SshfpRecord), true)
	case TxtRecord:
		return zone.addTxtRecord(recordPtr.(*TxtRecord), true)
	}

	return nil
}

func (zone *Zone) RemoveRecord(recordPtr interface{}) error {
	switch recordPtr.(type) {
	case ARecord:
		return zone.removeARecord(recordPtr.(*ARecord))
	case AaaaRecord:
		return zone.removeAaaaRecord(recordPtr.(*AaaaRecord))
	case AfsdbRecord:
		return zone.removeAfsdbRecord(recordPtr.(*AfsdbRecord))
	case CnameRecord:
		return zone.removeCnameRecord(recordPtr.(*CnameRecord))
	case DnskeyRecord:
		return zone.removeDnskeyRecord(recordPtr.(*DnskeyRecord))
	case DsRecord:
		return zone.removeDsRecord(recordPtr.(*DsRecord))
	case HinfoRecord:
		return zone.removeHinfoRecord(recordPtr.(*HinfoRecord))
	case LocRecord:
		return zone.removeLocRecord(recordPtr.(*LocRecord))
	case MxRecord:
		return zone.removeMxRecord(recordPtr.(*MxRecord))
	case NaptrRecord:
		return zone.removeNaptrRecord(recordPtr.(*NaptrRecord))
	case NsRecord:
		return zone.removeNsRecord(recordPtr.(*NsRecord))
	case Nsec3Record:
		return zone.removeNsec3Record(recordPtr.(*Nsec3Record))
	case Nsec3paramRecord:
		return zone.removeNsec3paramRecord(recordPtr.(*Nsec3paramRecord))
	case PtrRecord:
		return zone.removePtrRecord(recordPtr.(*PtrRecord))
	case RpRecord:
		return zone.removeRpRecord(recordPtr.(*RpRecord))
	case RrsigRecord:
		return zone.removeRrsigRecord(recordPtr.(*RrsigRecord))
	case SpfRecord:
		return zone.removeSpfRecord(recordPtr.(*SpfRecord))
	case SrvRecord:
		return zone.removeSrvRecord(recordPtr.(*SrvRecord))
	case SshfpRecord:
		return zone.removeSshfpRecord(recordPtr.(*SshfpRecord))
	case TxtRecord:
		return zone.removeTxtRecord(recordPtr.(*TxtRecord))
	}

	return nil
}

func (zone *Zone) addARecord(record *ARecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.A {
			if r.(ARecord).Name == record.Name {
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

func (zone *Zone) addAaaaRecord(record *AaaaRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Aaaa {
			if r.(AaaaRecord).Name == record.Name {
				zone.Zone.Aaaa[key] = record
				found = true
			}
		}
	}

	if !found {
		zone.Zone.Aaaa = append(zone.Zone.Aaaa, record)
	}

	return nil
}

func (zone *Zone) addAfsdbRecord(record *AfsdbRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Afsdb {
			if r.(AfsdbRecord).Name == record.Name {
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

func (zone *Zone) addCnameRecord(record *CnameRecord, replace bool) error {
	// CNAME records are unique, so always replace the record
	zone.Zone.Cname = RecordSet{record}
	return nil
}

func (zone *Zone) addDnskeyRecord(record *DnskeyRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Dnskey {
			if r.(DnskeyRecord).Name == record.Name {
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

func (zone *Zone) addDsRecord(record *DsRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Ds {
			if r.(DsRecord).Name == record.Name {
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

func (zone *Zone) addHinfoRecord(record *HinfoRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Hinfo {
			if r.(HinfoRecord).Name == record.Name {
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

func (zone *Zone) addLocRecord(record *LocRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Loc {
			if r.(LocRecord).Name == record.Name {
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

func (zone *Zone) addMxRecord(record *MxRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Mx {
			if r.(MxRecord).Name == record.Name {
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
			if r.(NaptrRecord).Name == record.Name {
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
			if r.(NsRecord).Name == record.Name {
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
			if r.(Nsec3Record).Name == record.Name {
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
			if r.(Nsec3paramRecord).Name == record.Name {
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

func (zone *Zone) addPtrRecord(record *PtrRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Ptr {
			if r.(PtrRecord).Name == record.Name {
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

func (zone *Zone) addRpRecord(record *RpRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Rp {
			if r.(RpRecord).Name == record.Name {
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

func (zone *Zone) addRrsigRecord(record *RrsigRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Rrsig {
			if r.(RrsigRecord).Name == record.Name {
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

func (zone *Zone) addSoaRecord(record *SoaRecord, replace bool) error {
	zone.Zone.Soa = *record
	return nil
}

func (zone *Zone) addSpfRecord(record *SpfRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Spf {
			if r.(SpfRecord).Name == record.Name {
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
			if r.(SrvRecord).Name == record.Name {
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

func (zone *Zone) addSshfpRecord(record *SshfpRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Sshfp {
			if r.(SshfpRecord).Name == record.Name {
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

func (zone *Zone) addTxtRecord(record *TxtRecord, replace bool) error {
	var found bool
	if replace == true {
		for key, r := range zone.Zone.Txt {
			if r.(TxtRecord).Name == record.Name {
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

func (zone *Zone) removeARecord(record *ARecord) error {
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

func (zone *Zone) removeAaaaRecord(record *AaaaRecord) error {
	var found bool
	for key, r := range zone.Zone.Aaaa {
		if r == record {
			records := zone.Zone.Aaaa[:key]
			zone.Zone.Aaaa = append(records, zone.Zone.Aaaa[key+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("AAAA Record not found")
	}

	return nil
}

func (zone *Zone) removeAfsdbRecord(record *AfsdbRecord) error {
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

func (zone *Zone) removeCnameRecord(record *CnameRecord) error {
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

func (zone *Zone) removeDnskeyRecord(record *DnskeyRecord) error {
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

func (zone *Zone) removeDsRecord(record *DsRecord) error {
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

func (zone *Zone) removeHinfoRecord(record *HinfoRecord) error {
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

func (zone *Zone) removeLocRecord(record *LocRecord) error {
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

func (zone *Zone) removeMxRecord(record *MxRecord) error {
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

func (zone *Zone) removePtrRecord(record *PtrRecord) error {
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

func (zone *Zone) removeRpRecord(record *RpRecord) error {
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

func (zone *Zone) removeRrsigRecord(record *RrsigRecord) error {
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

func (zone *Zone) removeSoaRecord(record *SoaRecord) error {
	zone.Zone.Soa = *record
	return nil
}

func (zone *Zone) removeSpfRecord(record *SpfRecord) error {
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

func (zone *Zone) removeSshfpRecord(record *SshfpRecord) error {
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

func (zone *Zone) removeTxtRecord(record *TxtRecord) error {
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
