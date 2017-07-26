package dns

// RecordSet represents a collection of Records
type RecordSet []*Record

// NaptrRecordSet represents a collection of NaptrRecords
type NaptrRecordSet []*NaptrRecord

// NsRecordSet represents a collection of NsRecords
type NsRecordSet []*NsRecord

// Nsec3RecordSet represents a collection of Nsec3Records
type Nsec3RecordSet []*Nsec3Record

// Nsec3paramRecordSet represents a collection of Nsec3paramRecords
type Nsec3paramRecordSet []*Nsec3paramRecord

// SrvRecordSet represents a collection of SrvRecords
type SrvRecordSet []*SrvRecord
