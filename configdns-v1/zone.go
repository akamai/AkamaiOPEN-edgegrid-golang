package dns

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Zone represents a DNS zone
type Zone struct {
	Token string `json:"token"`
	Zone  struct {
		Name       string               `json:"name,omitempty"`
		A          RecordSet            `json:"a,omitempty"`
		AAAA       RecordSet            `json:"aaaa,omitempty"`
		Afsdb      RecordSet            `json:"afsdb,omitempty"`
		Cname      RecordSet            `json:"cname,omitempty"`
		Dnskey     RecordSet            `json:"dnskey,omitempty"`
		Ds         RecordSet            `json:"ds,omitempty"`
		Hinfo      RecordSet            `json:"hinfo,omitempty"`
		Loc        RecordSet            `json:"loc,omitempty"`
		Mx         RecordSet            `json:"mx,omitempty"`
		Naptr      RecordSet            `json:"naptr,omitempty"`
		Ns         RecordSet            `json:"ns,omitempty"`
		Nsec3      RecordSet            `json:"nsec3,omitempty"`
		Nsec3param RecordSet            `json:"nsec3param,omitempty"`
		Ptr        RecordSet            `json:"ptr,omitempty"`
		Rp         RecordSet            `json:"rp,omitempty"`
		Rrsig      RecordSet            `json:"rrsig,omitempty"`
		Soa        *Record              `json:"soa,omitempty"`
		Spf        RecordSet            `json:"spf,omitempty"`
		Srv        RecordSet            `json:"srv,omitempty"`
		Sshfp      RecordSet            `json:"sshfp,omitempty"`
		Txt        RecordSet            `json:"txt,omitempty"`
		Records    map[string]RecordSet `json:"-"`
	} `json:"zone"`
}

// NewZone creates a new Zone
func NewZone(hostname string) Zone {
	zone := Zone{Token: "new"}
	zone.Zone.Name = hostname
	return zone
}

// Save updates the Zone
func (zone *Zone) Save() error {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/edgegrid-dns/v1/zones/"+zone.Zone.Name,
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

func (zone *Zone) PreMarshalJSON() error {
	zone.Zone.Records = make(map[string]RecordSet)
	zone.Zone.Records["A"] = zone.Zone.A
	zone.Zone.Records["AAAA"] = zone.Zone.AAAA
	zone.Zone.Records["AFSDB"] = zone.Zone.Afsdb
	zone.Zone.Records["CNAME"] = zone.Zone.Cname
	zone.Zone.Records["DNSKEY"] = zone.Zone.Dnskey
	zone.Zone.Records["DS"] = zone.Zone.Ds
	zone.Zone.Records["HINFO"] = zone.Zone.Hinfo
	zone.Zone.Records["LOC"] = zone.Zone.Loc
	zone.Zone.Records["MX"] = zone.Zone.Mx
	zone.Zone.Records["NAPTR"] = zone.Zone.Naptr
	zone.Zone.Records["NS"] = zone.Zone.Ns
	zone.Zone.Records["NSEC3"] = zone.Zone.Nsec3
	zone.Zone.Records["NSEC3PARAM"] = zone.Zone.Nsec3param
	zone.Zone.Records["PTR"] = zone.Zone.Ptr
	zone.Zone.Records["RP"] = zone.Zone.Rp
	zone.Zone.Records["RRSIG"] = zone.Zone.Rrsig
	zone.Zone.Records["SOA"] = []*Record{zone.Zone.Soa}
	zone.Zone.Records["SPF"] = zone.Zone.Spf
	zone.Zone.Records["SRV"] = zone.Zone.Srv
	zone.Zone.Records["SSHFP"] = zone.Zone.Sshfp
	zone.Zone.Records["TXT"] = zone.Zone.Txt

	zone.Zone.Soa.Serial = int(time.Now().Unix())

	return nil
}

func (zone *Zone) PostUnmarshalJSON() error {
	zone.Zone.A = zone.Zone.Records["A"]
	zone.Zone.AAAA = zone.Zone.Records["AAAA"]
	zone.Zone.Afsdb = zone.Zone.Records["AFSDB"]
	zone.Zone.Cname = zone.Zone.Records["CNAME"]
	zone.Zone.Dnskey = zone.Zone.Records["DNSKEY"]
	zone.Zone.Ds = zone.Zone.Records["DS"]
	zone.Zone.Hinfo = zone.Zone.Records["HINFO"]
	zone.Zone.Loc = zone.Zone.Records["LOC"]
	zone.Zone.Mx = zone.Zone.Records["MX"]
	zone.Zone.Naptr = zone.Zone.Records["NAPTR"]
	zone.Zone.Ns = zone.Zone.Records["NS"]
	zone.Zone.Nsec3 = zone.Zone.Records["NSEC3"]
	zone.Zone.Nsec3param = zone.Zone.Records["NSEC3PARAM"]
	zone.Zone.Ptr = zone.Zone.Records["PTR"]
	zone.Zone.Rp = zone.Zone.Records["RP"]
	zone.Zone.Rrsig = zone.Zone.Records["RRSIG"]
	zone.Zone.Soa = zone.Zone.Records["SOA"][0]
	zone.Zone.Spf = zone.Zone.Records["SPF"]
	zone.Zone.Srv = zone.Zone.Records["SRV"]
	zone.Zone.Sshfp = zone.Zone.Records["SSHFP"]
	zone.Zone.Txt = zone.Zone.Records["TXT"]

	return nil
}

func (zone *Zone) fixupCnames(record *Record) {
	if record.RecordType == "CNAME" {
		names := make(map[string]string, len(zone.Zone.Records["CNAME"]))
		for _, record := range zone.Zone.Records["CNAME"] {
			names[strings.ToUpper(record.Name)] = record.Name
		}

		for recordType, records := range zone.Zone.Records {
			if recordType == "CNAME" {
				continue
			}

			newRecords := RecordSet{}
			for _, record := range records {
				if _, ok := names[record.Name]; !ok {
					newRecords = append(newRecords, record)
				} else {
					log.Printf(
						"[WARN] %s Record conflicts with CNAME \"%s\", %[1]s Record ignored.",
						recordType,
						names[strings.ToUpper(record.Name)],
					)
				}
			}
			zone.Zone.Records[recordType] = newRecords
		}
	} else if record.Name != "" {
		name := strings.ToLower(record.Name)

		newRecords := RecordSet{}
		for _, cname := range zone.Zone.Records["CNAME"] {
			if strings.ToLower(cname.Name) != name {
				newRecords = append(newRecords, cname)
			} else {
				log.Printf(
					"[WARN] %s Record \"%s\" conflicts with existing CNAME \"%s\", removing CNAME",
					record.RecordType,
					record.Name,
					cname.Name,
				)
			}
		}

		zone.Zone.Records["CNAME"] = newRecords
	}
}
