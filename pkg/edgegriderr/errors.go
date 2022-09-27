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
		keys := make([]string, 0, len(validationErrors))
		for k := range validationErrors {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var s strings.Builder
		for _, key := range keys {
			if validationErrors[key] == nil {
				continue
			}
			errs, ok := validationErrors[key].(validation.Errors)

			if !ok {
				fmt.Fprintf(&s, "%s%s: %s\n", strings.Repeat("\t", indentSize), key, validationErrors[key].Error())
			}
			if _, err := strconv.Atoi(key); err != nil {
				fmt.Fprintf(&s, "%s", parser(errs, key, indentSize))
				continue
			}
			fmt.Fprintf(&s, "%s%s[%s]: {\n%s%s}\n", strings.Repeat("\t", indentSize), indexedFieldName, key, parser(errs, "", indentSize+1), strings.Repeat("\t", indentSize))
		}

		return s.String()
	}

	return parser
}
