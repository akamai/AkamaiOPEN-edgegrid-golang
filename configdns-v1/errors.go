package dns

import (
	"fmt"
)

type Error interface {
	error
	IsZoneNotFound() bool
	IsFailedToSave() bool
}

type ZoneNotFoundError struct {
	zoneName string
	err      error
}

type FailedToSaveError struct {
	err error
}

func IsZoneNotFound(err error) bool {
	_, ok := err.(*ZoneNotFoundError)
	return ok
}

func IsFailedToSave(err error) bool {
	_, ok := err.(*FailedToSaveError)
	return ok
}

func (e *ZoneNotFoundError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Zone \"%s\" not found, creating new zone.", e.zoneName)
}

func (e *FailedToSaveError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Unable to save record (%s)", e.err.Error())
}
