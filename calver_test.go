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
	v = NewCalVer(23, 7, 0, "")
	assert.Equal(t, "23.07.00", v.String())
	v = NewCalVer(23, 7, 0, "alpha.1")
	assert.Equal(t, "23.07.00-alpha.1", v.String())
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
		{"2023.11.13-alpha", false},
		{"23.7.23", false},
		{"23.07.3", false},
		{"23.11.13-alpha", false},
		{"2023.11.13-dev+1", true},
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

func TestCalVerParts(t *testing.T) {
	v, err := NewSemverStr("2023.2.3-beta.1")
	assert.NoError(t, err)

	assert.Equal(t, uint64(2023), v.Major())
	assert.Equal(t, uint64(2), v.Minor())
	assert.Equal(t, uint64(3), v.Patch())
	assert.Equal(t, "beta.1", v.Prerelease())
	assert.Equal(t, "", v.Metadata())
}

func TestCalVerGt(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"2022.2.3", "2023.2.3", false},
		{"2023.1.3", "2023.2.3", false},
		{"2023.2.1", "2023.2.3", false},
		{"2024.2.3", "2023.2.3", true},
		{"2023.5.1", "2023.2.3", true},
		{"2023.2.1", "2023.2.01-beta", true},
		{"2023.2.2", "2023.02.01-beta", true},
	}

	for _, tc := range tests {
		v1, err := NewCalVerStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewCalVerStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Gt(v2)
		assert.Equal(t, a, tc.expected)
	}
}

func TestCalVerLt(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"2023.2.3", "2022.2.3", false},
		{"2023.2.3", "2023.1.3", false},
		{"2023.2.3", "2023.2.1", false},
		{"2023.2.3", "2024.2.3", true},
		{"2023.2.3", "2023.5.1", true},
		{"2023.2.01-beta", "2023.2.1", true},
		{"2023.02.01-beta", "2023.2.2", true},
	}

	for _, tc := range tests {
		v1, err := NewCalVerStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewCalVerStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Lt(v2)
		assert.Equal(t, a, tc.expected)
	}
}

func TestCalVerEq(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"2023.2.3", "2023.5.1", false},
		{"2023.2.01-beta", "2023.2.1-beta", true},
		{"2023.02.01-beta", "2023.2.1-beta", true},
	}

	for _, tc := range tests {
		v1, err := NewCalVerStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewCalVerStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Eq(v2)
		assert.Equal(t, a, tc.expected)
	}
}

func TestCalVerInc(t *testing.T) {
	tests := []struct {
		v1               string
		expected         string
		how              string
		expectedOriginal string
	}{
		{"2023.2.3", "2023.2.4", "patch", "2023.02.04"},
		{"2023.2.4", "2023.2.5", "patch", "2023.02.05"},
		{"2023.2.3", "2023.3.0", "minor", "2023.03.00"},
		{"2023.2.4", "2023.3.0", "minor", "2023.03.00"},
		{"2023.2.3", "2024.0.0", "major", "2024.00.00"},
		{"2023.2.4", "2024.0.0", "major", "2024.00.00"},
		{"2023.2.3-beta", "2023.2.3", "patch", "2023.02.03"},
		{"2023.2.4-beta", "2023.2.4", "patch", "2023.02.04"},
		{"2023.2.3-beta", "2023.3.0", "minor", "2023.03.00"},
		{"2023.2.4-beta", "2023.3.0", "minor", "2023.03.00"},
		{"2023.2.3-beta", "2023.0.0", "major", "2024.00.00"},
		{"2023.2.4-beta", "2023.0.0", "major", "2024.00.00"},
	}

	for _, tc := range tests {
		v1, err := NewCalVerStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}
		var v2 Comparable
		switch tc.how {
		case "patch":
			v2 = v1.IncPatch()
		case "minor":
			v2 = v1.IncMinor()
		case "major":
			v2 = v1.IncMajor()
		}

		sem, ok := v2.(*CalVer)
		assert.True(t, ok)
		a := sem.Original()
		assert.Equal(t, tc.expectedOriginal, a)
	}
}
