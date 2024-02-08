// Package errs provides utilities for working with errors during JSON data unmarshalling.
// It includes functions for unescaping HTML content and checking if a string contains HTML or XML data.
package errs

import (
	"strings"

	"golang.org/x/net/html"
)

// UnescapeContent unescapes HTML content.
func UnescapeContent(content string) string {
	//check if the content is HTML or XML
	if isHTML(content) {
		return html.UnescapeString(content)
	}
	return content
}

func isHTML(data string) bool {
	_, err := html.Parse(strings.NewReader(data))
	return err == nil
}
