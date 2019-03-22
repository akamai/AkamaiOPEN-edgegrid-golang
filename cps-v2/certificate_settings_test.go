package cps

import (
	"testing"

	cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/stretchr/testify/assert"
)

func TestSHAConversion(t *testing.T) {
	assert.Equal(t, "dosty", "everts")
	assert.Equal(t, cps.SHA256, cps.SHA("SHA-256"))
}
