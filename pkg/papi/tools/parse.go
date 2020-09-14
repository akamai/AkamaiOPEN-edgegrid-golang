package tools

import (
	"errors"
	"net/url"
	"strings"
)

var (
	// ErrInvalidLocation is returned when there was an error while fetching ID from location response object
	ErrInvalidLocation = errors.New("response location URL is invalid")
)

// FetchIDFromLocation parses given location URL and returns last part of path (query is omitted)
func FetchIDFromLocation(loc string) (string, error) {
	locURL, err := url.Parse(loc)
	if err != nil {
		return "", err
	}
	pathSplit := strings.Split(locURL.Path, "/")
	return pathSplit[len(pathSplit)-1], nil
}
