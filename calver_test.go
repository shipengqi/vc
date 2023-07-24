package vc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCalVer(t *testing.T) {
	v := NewCalVer(2023, 7, 23, "")
	assert.Equal(t, "2023.07.23", v.String())
	v = NewCalVer(2023, 07, 03, "alpha.1")
	assert.Equal(t, "2023.07.03-alpha.1", v.String())
	v = NewCalVer(2023, 7, 3, "alpha.1")
	assert.Equal(t, "2023.07.03-alpha.1", v.String())
	v = NewCalVer(2023, 11, 3, "alpha.1")
	assert.Equal(t, "2023.11.03-alpha.1", v.String())
}

func TestNewCalVerStr(t *testing.T) {
	tests := []struct {
		version string
		err     bool
	}{
		{"2023.7.23", false},
		{"2023.07.3", false},
		{"2023.07.03", false},
		{"2023.11.13", false},
		{"2023.117.23", true},
		{"202.7.23", true},
		{"2023.7.202", true},
	}

	for _, tc := range tests {
		_, err := NewCalVerStr(tc.version)
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
