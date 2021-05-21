package cps

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

var (
	// ErrInvalidLocation is returned when there was an error while fetching ID from location response object
	ErrInvalidLocation = errors.New("location URL is invalid")
)

// GetIDFromLocation parse the link and returns the id
func GetIDFromLocation(location string) (int, error) {
	locURL, err := url.Parse(location)
	if err != nil {
		return 0, err
	}
	pathSplit := strings.Split(locURL.Path, "/")
	idStr := pathSplit[len(pathSplit)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
