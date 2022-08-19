package edgegriderr

import (
	"fmt"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidationErrorsParsing(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		assert.Equal(t, nil, ParseValidationErrors(nil))
	})

	tests := map[string]struct {
		input    validation.Errors
		expected string
	}{
		"single error": {
			input: validation.Errors{
				"Error": fmt.Errorf("oops"),
			},
			expected: `
Error: oops`,
		},
		"multiple errors": {
			input: validation.Errors{
				"Error1": fmt.Errorf("oops"),
				"Error2": fmt.Errorf("oops"),
			},
			expected: `
Error1: oops
Error2: oops`,
		},
		"nested error": {
			input: validation.Errors{
				"Error": validation.Errors{
					"0": validation.Errors{
						"NestedError1": fmt.Errorf("oops"),
					},
				},
			},
			expected: `
Error[0]: {
	NestedError1: oops
}`,
		},
		"multiple nested errors": {
			input: validation.Errors{
				"Error": validation.Errors{
					"0": validation.Errors{
						"NestedError1": fmt.Errorf("oops"),
					},
					"1": validation.Errors{
						"Error1": validation.Errors{
							"Error1-1": validation.Errors{
								"0": validation.Errors{
									"NestedError1": fmt.Errorf("oops"),
								},
								"1": validation.Errors{
									"NestedError1": fmt.Errorf("oops"),
								},
							},
							"Error1-2": validation.Errors{
								"NestedError1": fmt.Errorf("oops"),
							},
						},
					},
				},
			},
			expected: `
Error[0]: {
	NestedError1: oops
}
Error[1]: {
	Error1: {
		Error1-1[0]: {
			NestedError1: oops
		}
		Error1-1[1]: {
			NestedError1: oops
		}
		Error1-2: {
			NestedError1: oops
		}
	}
}`,
		},
		"errors on multiple levels": {
			input: validation.Errors{
				"Error1": validation.Errors{
					"0": validation.Errors{
						"NestedError": validation.Errors{
							"0": validation.Errors{
								"DoubleNestedError1": fmt.Errorf("oops"),
								"DoubleNestedError2": fmt.Errorf("oops"),
							},
							"1": validation.Errors{
								"DoubleNestedError1": fmt.Errorf("oops"),
							},
						},
					},
					"1": validation.Errors{
						"NestedError1": fmt.Errorf("oops"),
					},
				},
				"Error2": fmt.Errorf("oops"),
			},
			expected: `
Error1[0]: {
	NestedError[0]: {
		DoubleNestedError1: oops
		DoubleNestedError2: oops
	}
	NestedError[1]: {
		DoubleNestedError1: oops
	}
}
Error1[1]: {
	NestedError1: oops
}
Error2: oops`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, strings.TrimPrefix(test.expected, "\n"), ParseValidationErrors(test.input).Error())
		})
	}
}
