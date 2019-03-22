package dnsv2

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type RecordSetResponse struct {
	Metadata struct {
		ShowAll       bool `json:"showAll"`
		TotalElements int  `json:"totalElements"`
	} `json:"metadata"`
	Recordsets []struct {
		Name  string   `json:"name"`
		Type  string   `json:"type"`
		TTL   int      `json:"ttl"`
		Rdata []string `json:"rdata"`
	} `json:"recordsets"`
}

/*
{
  "metadata": {
    "zone": "example.com",
    "types": [
      "A"
    ],
    "page": 1,
    "pageSize": 25,
    "totalElements": 2
  },
  "recordsets": [
    {
      "name": "www.example.com",
      "type": "A",
      "ttl": 300,
      "rdata": [
        "10.0.0.2",
        "10.0.0.3"
      ]
    },
    {
      "name": "mail.example.com",
      "type": "A",
      "ttl": 300,
      "rdata": [
        "192.168.0.1",
        "192.168.0.2"
      ]
    }
  ]
}

*/

func NewRecordSetResponse(name string) *RecordSetResponse {
	recordset := &RecordSetResponse{}
	return recordset
}

func GetRecordList(zone string, name string, record_type string) (*RecordSetResponse, error) {
	records := NewRecordSetResponse(name)

	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-dns/v2/zones/"+zone+"/recordsets?types="+record_type+"&showAll=true",
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}
	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, &ZoneError{zoneName: name}
	} else {
		err = client.BodyJSON(res, records)
		if err != nil {
			return nil, err
		}
		return records, nil
	}
}

func GetRdata(zone string, name string, record_type string) ([]string, error) {
	records, err := GetRecordList(zone, name, record_type)
	if err != nil {
		return nil, err
	}

	var arrLength int
	for _, c := range records.Recordsets {
		if c.Name == name {
			arrLength = len(c.Rdata)
		}
	}

	rdata := make([]string, 0, arrLength)

	for _, r := range records.Recordsets {
		if r.Name == name {
			for _, i := range r.Rdata {
				rdata = append(rdata, i)
			}
		}
	}
	return rdata, nil
}
