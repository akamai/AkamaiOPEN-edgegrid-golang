package dns

const (
	ErrZoneNotFound = iota
	ErrFailedToSave
)

var (
	errorMap = map[int]string{
		ErrZoneNotFound: "Zone \"%s\" not found, creating new zone.",
		ErrFailedToSave: "Unable to save record (%s)",
	}
)
