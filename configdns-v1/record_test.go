package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecord_ContainsHelper(t *testing.T) {
	tm1 := []string{
		"test1",
		"test2",
		"test3",
	}

	assert.Equal(t, contains(tm1, "test1"), true)
	assert.Equal(t, contains(tm1, "test2"), true)
	assert.Equal(t, contains(tm1, "test3"), true)
	assert.Equal(t, contains(tm1, "test4"), false)
}

func TestRecord_ARecord(t *testing.T) {
	a := NewARecord()
	f := []string{
		"name",
		"ttl",
		"active",
		"target",
	}
	assert.Equal(t, a.fieldMap, f)
	assert.Equal(t, a.fieldMap, a.GetAllowedFields())
	assert.Equal(t, a.SetField("name", "test1"), nil)
	assert.Equal(t, a.SetField("doesntExist", "test1"), &RecordError{fieldName: "doesntExist"})
	a.SetField("ttl", 900)
	a.SetField("active", true)
	a.SetField("target", "test2")
	assert.Equal(t, a.ToMap(), map[string]interface{}{
		"name":   "test1",
		"ttl":    900,
		"active": true,
		"target": "test2",
	})
}
