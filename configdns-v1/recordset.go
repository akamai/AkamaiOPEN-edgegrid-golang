package dns

// RecordSet represents a collection of records
type RecordSet []*Record

type NaptrRecordSet []*NaptrRecord
type NsRecordSet []*NsRecord
type Nsec3RecordSet []*Nsec3Record
type Nsec3paramRecordSet []*Nsec3paramRecord
type SrvRecordSet []*SrvRecord
