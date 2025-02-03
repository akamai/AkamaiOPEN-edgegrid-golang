// Package edgegriderr is used for parsing validation errors to make them more readable.
// It formats error(s) in a way, that the return value is one formatted error type, consisting of all the errors that occurred
// in human-readable form. It is important to provide all the validation errors to the function.
// Usage example:
//
//	error := edgegriderr.ParseValidationErrors(validation.Errors{
//		"Validation1": validation.Validate(...),
//		"Validation2": validation.Validate(...),
//	})
package edgegriderr

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ParseValidationErrors parses validation errors into easily readable form
// The output error is formatted with indentations and struct field indexing for collections
func ParseValidationErrors(e validation.Errors) error {
	if e.Filter() == nil {
		return nil
	}

	parser := validationErrorsParser()
	return fmt.Errorf("%s", strings.TrimSuffix(parser(e, "", 0), "\n"))
}

// validationErrorsParser returns a function that parses validation errors
// Returned function takes validation.Errors, field to be indexed (empty at the beginning) and index size at start as parameters
func validationErrorsParser() func(validation.Errors, string, int) string {
	var parser func(validation.Errors, string, int) string

	parser = func(validationErrors validation.Errors, indexedFieldName string, indentSize int) string {
		keys := getSortedKeys(validationErrors)

		var s strings.Builder
		for _, key := range keys {
			if validationErrors[key] == nil {
				continue
			}
			errs, ok := validationErrors[key].(validation.Errors)
			if !ok {
				fmt.Fprintf(&s, "%s%s: %s\n", indent(indentSize), key, validationErrors[key].Error())
				continue
			}

			if _, err := strconv.Atoi(key); err != nil {
				if hasNumericKeys(errs) {
					fmt.Fprintf(&s, "%s", parser(errs, key, indentSize))
				} else {
					fmt.Fprintf(&s, "%s%s: {\n%s%s}\n", indent(indentSize), key, parser(errs, key, indentSize+1), indent(indentSize))
				}
				continue
			}

			fmt.Fprintf(&s, "%s%s[%s]: {\n%s%s}\n", indent(indentSize), indexedFieldName, key, parser(errs, "", indentSize+1), indent(indentSize))
		}

		return s.String()
	}

	return parser
}

func getSortedKeys(errs validation.Errors) []string {
	keys := make([]string, 0, len(errs))
	for k := range errs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func indent(size int) string {
	return strings.Repeat("\t", size)
}

func hasNumericKeys(errs validation.Errors) bool {
	for key := range errs {
		if _, err := strconv.Atoi(key); err != nil {
			return false
		}
	}
	return true
}
