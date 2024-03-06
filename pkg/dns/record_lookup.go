package dns

import (
	"context"
	"fmt"
	"net/http"

	"encoding/hex"
	"net"
	"strconv"
	"strings"
)

func (d *dns) FullIPv6(ctx context.Context, ip net.IP) string {
	logger := d.Log(ctx)
	logger.Debug("FullIPv6")

	dst := make([]byte, hex.EncodedLen(len(ip)))
	_ = hex.Encode(dst, ip)
	return string(dst[0:4]) + ":" +
		string(dst[4:8]) + ":" +
		string(dst[8:12]) + ":" +
		string(dst[12:16]) + ":" +
		string(dst[16:20]) + ":" +
		string(dst[20:24]) + ":" +
		string(dst[24:28]) + ":" +
		string(dst[28:])
}

func padValue(str string) string {
	newStr := strings.Replace(str, "m", "", -1)
	float, err := strconv.ParseFloat(newStr, 32)
	if err != nil {
		return "FAIL"
	}

	return fmt.Sprintf("%.2f", float)
}

func (d *dns) PadCoordinates(ctx context.Context, str string) string {
	logger := d.Log(ctx)
	logger.Debug("PadCoordinates")

	s := strings.Split(str, " ")
	if len(s) < 12 {
		return ""
	}

	latd, latm, lats, latDir, longd, longm, longs, longDir, altitude, size, horizPrecision, vertPrecision := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11]

	return latd + " " + latm + " " + lats + " " + latDir + " " + longd + " " + longm + " " + longs + " " + longDir + " " + padValue(altitude) + "m " + padValue(size) + "m " + padValue(horizPrecision) + "m " + padValue(vertPrecision) + "m"
}

func (d *dns) GetRecord(ctx context.Context, zone, name, recordType string) (*RecordBody, error) {
	logger := d.Log(ctx)
	logger.Debug("GetRecord")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, name, recordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecord request: %w", err)
	}

	var result RecordBody
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetRecordList(ctx context.Context, zone, _, recordType string) (*RecordSetResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetRecordList")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets?types=%s&showAll=true", zone, recordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecordList request: %w", err)
	}

	var result RecordSetResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRecordList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetRdata(ctx context.Context, zone, name, recordType string) ([]string, error) {
	logger := d.Log(ctx)
	logger.Debug("GetrData")

	records, err := d.GetRecordList(ctx, zone, name, recordType)
	if err != nil {
		return nil, err
	}

	var rData []string
	for _, r := range records.RecordSets {
		if r.Name == name {
			for _, i := range r.Rdata {
				str := i

				if recordType == "AAAA" {
					addr := net.ParseIP(str)
					result := d.FullIPv6(ctx, addr)
					str = result
				} else if recordType == "LOC" {
					str = d.PadCoordinates(ctx, str)
				}
				rData = append(rData, str)
			}
		}
	}
	return rData, nil
}

func (d *dns) ProcessRdata(ctx context.Context, rData []string, rType string) []string {
	logger := d.Log(ctx)
	logger.Debug("ProcessrData")

	var newRData []string
	for _, i := range rData {
		str := i
		if rType == "AAAA" {
			addr := net.ParseIP(str)
			result := d.FullIPv6(ctx, addr)
			str = result
		} else if rType == "LOC" {
			str = d.PadCoordinates(ctx, str)
		}
		newRData = append(newRData, str)
	}

	return newRData
}

func (d *dns) ParseRData(ctx context.Context, rType string, rData []string) map[string]interface{} {
	logger := d.Log(ctx)
	logger.Debug("ParserData")

	fieldMap := make(map[string]interface{}, 0)
	if len(rData) == 0 {
		return fieldMap
	}
	newRData := make([]string, 0, len(rData))
	fieldMap["target"] = newRData

	switch rType {
	case "AFSDB":
		resolveAFSDBType(rData, newRData, fieldMap)

	case "DNSKEY":
		resolveDNSKEYType(rData, fieldMap)

	case "DS":
		resolveDSType(rData, fieldMap)

	case "HINFO":
		resolveHINFOType(rData, fieldMap)
	/*
		// too many variations to calculate pri and increment
		case "MX":
			sort.Strings(rData)
			parts := strings.Split(rData[0], " ")
			fieldMap["priority"], _ = strconv.Atoi(parts[0])
			if len(rData) > 1 {
				parts = strings.Split(rData[1], " ")
				tpri, _ := strconv.Atoi(parts[0])
				fieldMap["priority_increment"] = tpri - fieldMap["priority"].(int)
			}
			for _, rContent := range rData {
				parts := strings.Split(rContent, " ")
				newrData = append(newrData, parts[1])
			}
			fieldMap["target"] = newrData
	*/

	case "NAPTR":
		resolveNAPTRType(rData, fieldMap)

	case "NSEC3":
		resolveNSEC3Type(rData, fieldMap)

	case "NSEC3PARAM":
		resolveNSEC3PARAMType(rData, fieldMap)

	case "RP":
		resolveRPType(rData, fieldMap)

	case "RRSIG":
		resolveRRSIGType(rData, fieldMap)

	case "SRV":
		resolveSRVType(rData, newRData, fieldMap)

	case "SSHFP":
		resolveSSHFPType(rData, fieldMap)

	case "SOA":
		resolveSOAType(rData, fieldMap)

	case "AKAMAITLC":
		resolveAKAMAITLCType(rData, fieldMap)

	case "SPF":
		resolveSPFType(rData, newRData, fieldMap)

	case "TXT":
		resolveTXTType(rData, newRData, fieldMap)

	case "AAAA":
		resolveAAAAType(ctx, d, rData, newRData, fieldMap)

	case "LOC":
		resolveLOCType(ctx, d, rData, newRData, fieldMap)

	case "CERT":
		resolveCERTType(rData, fieldMap)

	case "TLSA":
		resolveTLSAType(rData, fieldMap)

	case "SVCB":
		resolveSVCBType(rData, fieldMap)

	case "HTTPS":
		resolveHTTPSType(rData, fieldMap)

	default:
		for _, rContent := range rData {
			newRData = append(newRData, rContent)
		}
		fieldMap["target"] = newRData
	}

	return fieldMap
}

func resolveAFSDBType(rData, newRData []string, fieldMap map[string]interface{}) {
	parts := strings.Split(rData[0], " ")
	fieldMap["subtype"], _ = strconv.Atoi(parts[0])
	for _, rContent := range rData {
		parts = strings.Split(rContent, " ")
		newRData = append(newRData, parts[1])
	}
	fieldMap["target"] = newRData
}

func resolveDNSKEYType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["flags"], _ = strconv.Atoi(parts[0])
		fieldMap["protocol"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[2])
		key := parts[3]
		// key can have whitespace
		if len(parts) > 4 {
			i := 4
			for i < len(parts) {
				key += " " + parts[i]
			}
		}
		fieldMap["key"] = key
		break
	}
}

func resolveSVCBType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.SplitN(rContent, " ", 3)
		// has to be at least two fields.
		if len(parts) < 2 {
			break
		}
		fieldMap["svc_priority"], _ = strconv.Atoi(parts[0])
		fieldMap["target_name"] = parts[1]
		if len(parts) > 2 {
			fieldMap["svc_params"] = parts[2]
		}
		break
	}
}

func resolveDSType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["keytag"], _ = strconv.Atoi(parts[0])
		fieldMap["digest_type"], _ = strconv.Atoi(parts[2])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
		dig := parts[3]
		// digest can have whitespace
		if len(parts) > 4 {
			i := 4
			for i < len(parts) {
				dig += " " + parts[i]
			}
		}
		fieldMap["digest"] = dig
		break
	}
}

func resolveHINFOType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["hardware"] = parts[0]
		fieldMap["software"] = parts[1]
		break
	}
}

func resolveNAPTRType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["order"], _ = strconv.Atoi(parts[0])
		fieldMap["preference"], _ = strconv.Atoi(parts[1])
		fieldMap["flagsnaptr"] = parts[2]
		fieldMap["service"] = parts[3]
		fieldMap["regexp"] = parts[4]
		fieldMap["replacement"] = parts[5]
		break
	}
}

func resolveNSEC3Type(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["flags"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
		fieldMap["iterations"], _ = strconv.Atoi(parts[2])
		fieldMap["salt"] = parts[3]
		fieldMap["next_hashed_owner_name"] = parts[4]
		fieldMap["type_bitmaps"] = parts[5]
		break
	}
}

func resolveNSEC3PARAMType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["flags"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
		fieldMap["iterations"], _ = strconv.Atoi(parts[2])
		fieldMap["salt"] = parts[3]
		break
	}
}

func resolveRPType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["mailbox"] = parts[0]
		fieldMap["txt"] = parts[1]
		break
	}
}

func resolveRRSIGType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["type_covered"] = parts[0]
		fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
		fieldMap["labels"], _ = strconv.Atoi(parts[2])
		fieldMap["original_ttl"], _ = strconv.Atoi(parts[3])
		fieldMap["expiration"] = parts[4]
		fieldMap["inception"] = parts[5]
		fieldMap["signer"] = parts[7]
		fieldMap["keytag"], _ = strconv.Atoi(parts[6])
		sig := parts[8]
		// sig can have whitespace
		if len(parts) > 9 {
			i := 9
			for i < len(parts) {
				sig += " " + parts[i]
			}
		}
		fieldMap["signature"] = sig
		break
	}
}

func resolveSRVType(rData, newRData []string, fieldMap map[string]interface{}) {
	// pull out some fields
	parts := strings.Split(rData[0], " ")
	fieldMap["priority"], _ = strconv.Atoi(parts[0])
	fieldMap["weight"], _ = strconv.Atoi(parts[1])
	fieldMap["port"], _ = strconv.Atoi(parts[2])
	// populate target
	for _, rContent := range rData {
		parts = strings.Split(rContent, " ")
		newRData = append(newRData, parts[3])
	}
	fieldMap["target"] = newRData
}

func resolveSSHFPType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
		fieldMap["fingerprint_type"], _ = strconv.Atoi(parts[1])
		fieldMap["fingerprint"] = parts[2]
		break
	}
}

func resolveSOAType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["name_server"] = parts[0]
		fieldMap["email_address"] = parts[1]
		fieldMap["serial"], _ = strconv.Atoi(parts[2])
		fieldMap["refresh"], _ = strconv.Atoi(parts[3])
		fieldMap["retry"], _ = strconv.Atoi(parts[4])
		fieldMap["expiry"], _ = strconv.Atoi(parts[5])
		fieldMap["nxdomain_ttl"], _ = strconv.Atoi(parts[6])
		break
	}
}

func resolveAKAMAITLCType(rData []string, fieldMap map[string]interface{}) {
	parts := strings.Split(rData[0], " ")
	fieldMap["answer_type"] = parts[0]
	fieldMap["dns_name"] = parts[1]
}

func resolveSPFType(rData, newRData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		newRData = append(newRData, rContent)
	}
	fieldMap["target"] = newRData
}

func resolveTXTType(rData, newRData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		newRData = append(newRData, rContent)
	}
	fieldMap["target"] = newRData
}

func resolveAAAAType(ctx context.Context, d *dns, rData, newRData []string, fieldMap map[string]interface{}) {
	for _, i := range rData {
		str := i
		addr := net.ParseIP(str)
		result := d.FullIPv6(ctx, addr)
		str = result
		newRData = append(newRData, str)
	}
	fieldMap["target"] = newRData
}

func resolveLOCType(ctx context.Context, d *dns, rData, newRData []string, fieldMap map[string]interface{}) {
	for _, i := range rData {
		str := i
		str = d.PadCoordinates(ctx, str)
		newRData = append(newRData, str)
	}
	fieldMap["target"] = newRData
}

func resolveCERTType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		val, err := strconv.Atoi(parts[0])
		if err == nil {
			fieldMap["type_value"] = val
		} else {
			fieldMap["type_mnemonic"] = parts[0]
		}
		fieldMap["keytag"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[2])
		fieldMap["certificate"] = parts[3]
		break
	}
}

func resolveTLSAType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.Split(rContent, " ")
		fieldMap["usage"], _ = strconv.Atoi(parts[0])
		fieldMap["selector"], _ = strconv.Atoi(parts[1])
		fieldMap["match_type"], _ = strconv.Atoi(parts[2])
		fieldMap["certificate"] = parts[3]
		break
	}
}

func resolveHTTPSType(rData []string, fieldMap map[string]interface{}) {
	for _, rContent := range rData {
		parts := strings.SplitN(rContent, " ", 3)
		// has to be at least two fields.
		if len(parts) < 2 {
			break
		}
		fieldMap["svc_priority"], _ = strconv.Atoi(parts[0])
		fieldMap["target_name"] = parts[1]
		if len(parts) > 2 {
			fieldMap["svc_params"] = parts[2]
		}
		break
	}
}
