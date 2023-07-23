package vc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSemverStr(t *testing.T) {
	tests := []struct {
		version string
		err     bool
	}{
		{"1.2.3", false},
		{"1.2.3-alpha.01", true},
		{"1.2.3+test.01", false},
		{"1.2.3-alpha.-1", false},
		{"v1.2.3", false},
		{"1.0", false},
		{"v1.0", false},
		{"1", false},
		{"v1", false},
		{"1.2.beta", true},
		{"v1.2.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"v1.2-5", false},
		{"1.2-beta.5", false},
		{"v1.2-beta.5", false},
		{"\n1.2", true},
		{"\nv1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"v1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hypen", false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", false},
		{"1.2.3-rc1-with-hypen", false},
		{"v1.2.3-rc1-with-hypen", false},
		{"1.2.3.4", true},
		{"v1.2.3.4", true},
		{"1.2.2147483648", false},
		{"1.2147483648.3", false},
		{"2147483648.3.0", false},

		// Due to having 4 parts these should produce an error. See
		// https://github.com/Masterminds/semver/issues/185 for the reason for
		// these tests.
		{"12.3.4.1234", true},
		{"12.23.4.1234", true},
		{"12.3.34.1234", true},

		// The SemVer spec in a pre-release expects to allow [0-9A-Za-z-].
		{"20221209-update-renovatejson-v4", false},
	}

	for _, tc := range tests {
		_, err := NewSemverStr(tc.version)
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestNewSemver(t *testing.T) {
	// v0.1.2
	v := NewSemver(0, 1, 2, "", "")
	assert.Equal(t, "0.1.2", v.String())

	// v1.2.3-alpha.1+foo.bar
	v = NewSemver(1, 2, 3, "alpha.1", "foo.bar")
	assert.Equal(t, "1.2.3-alpha.1+foo.bar", v.String())
}

func TestOriginal(t *testing.T) {
	tests := []string{
		"1.2.3",
		"v1.2.3",
		"1.0",
		"v1.0",
		"1",
		"v1",
		"1.2-5",
		"v1.2-5",
		"1.2-beta.5",
		"v1.2-beta.5",
		"1.2.0-x.Y.0+metadata",
		"v1.2.0-x.Y.0+metadata",
		"1.2.0-x.Y.0+metadata-width-hypen",
		"v1.2.0-x.Y.0+metadata-width-hypen",
		"1.2.3-rc1-with-hypen",
		"v1.2.3-rc1-with-hypen",
	}

	for _, tc := range tests {
		v, err := NewSemverStr(tc)
		assert.NoError(t, err)

		o := v.Original()
		assert.Equal(t, tc, o)
	}
}

func TestParts(t *testing.T) {
	v, err := NewSemverStr("1.2.3-beta.1+build.123")
	assert.NoError(t, err)

	assert.Equal(t, uint64(1), v.Major())
	assert.Equal(t, uint64(2), v.Minor())
	assert.Equal(t, uint64(3), v.Patch())
	assert.Equal(t, "beta.1", v.Prerelease())
	assert.Equal(t, "build.123", v.Metadata())
}

func TestCoerceString(t *testing.T) {
	tests := []struct {
		version  string
		expected string
	}{
		{"1.2.3", "1.2.3"},
		{"v1.2.3", "1.2.3"},
		{"1.0", "1.0.0"},
		{"v1.0", "1.0.0"},
		{"1", "1.0.0"},
		{"v1", "1.0.0"},
		{"1.2-5", "1.2.0-5"},
		{"v1.2-5", "1.2.0-5"},
		{"1.2-beta.5", "1.2.0-beta.5"},
		{"v1.2-beta.5", "1.2.0-beta.5"},
		{"1.2.0-x.Y.0+metadata", "1.2.0-x.Y.0+metadata"},
		{"v1.2.0-x.Y.0+metadata", "1.2.0-x.Y.0+metadata"},
		{"1.2.0-x.Y.0+metadata-width-hypen", "1.2.0-x.Y.0+metadata-width-hypen"},
		{"v1.2.0-x.Y.0+metadata-width-hypen", "1.2.0-x.Y.0+metadata-width-hypen"},
		{"1.2.3-rc1-with-hypen", "1.2.3-rc1-with-hypen"},
		{"v1.2.3-rc1-with-hypen", "1.2.3-rc1-with-hypen"},
	}

	for _, tc := range tests {
		v, err := NewSemverStr(tc.version)
		assert.NoError(t, err)

		s := v.String()
		assert.Equal(t, tc.expected, s)
	}
}

func TestSemverGt(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.5.1", false},
		{"2.2.3", "1.5.1", true},
		{"3.2.0-beta", "3.2.0-beta", false},
		{"3.2.0-beta.1", "3.2.0-beta.5", false},
		{"3.2.0-beta.4", "3.2.0-beta.2", true},
		{"7.43.0-SNAPSHOT.99", "7.43.0-SNAPSHOT.103", false},
		{"7.43.0-SNAPSHOT.FOO", "7.43.0-SNAPSHOT.103", true},
		{"7.43.0-SNAPSHOT.99", "7.43.0-SNAPSHOT.BAR", false},
	}

	for _, tc := range tests {
		v1, err := NewSemverStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemverStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Gt(v2)
		assert.Equal(t, a, tc.expected)
	}
}

func TestSemverLt(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.5.1", true},
		{"2.2.3", "1.5.1", false},
		{"3.2.0-beta", "3.2.0-beta", false},
	}

	for _, tc := range tests {
		v1, err := NewSemverStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemverStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Lt(v2)
		assert.Equal(t, a, tc.expected)
	}
}

func TestSemverEq(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.5.1", false},
		{"2.2.3", "1.5.1", false},
		{"3.2.0-beta", "3.2.0-beta", true},
		{"3.2.0-beta+foo", "3.2.0-beta+bar", true},
	}

	for _, tc := range tests {
		v1, err := NewSemverStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemverStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Eq(v2)
		assert.Equal(t, a, tc.expected)
	}
}
func TestCompare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.5.1", -1},
		{"2.2.3", "1.5.1", 1},
		{"2.2.3", "2.2.2", 1},
		{"3.2.0-beta", "3.2.0-beta", 0},
		{"1.3.0", "1.1.4", 1},
		{"4.2.0", "4.2.0-beta", 1},
		{"4.2.0-beta", "4.2.0", -1},
		{"4.2.0-alpha", "4.2.0-beta", -1},
		{"4.2.0-alpha", "4.2.0-alpha", 0},
		{"4.2.0-beta.2", "4.2.0-beta.1", 1},
		{"4.2.0-beta2", "4.2.0-beta1", 1},
		{"4.2.0-beta", "4.2.0-beta.2", -1},
		{"4.2.0-beta", "4.2.0-beta.foo", -1},
		{"4.2.0-beta.2", "4.2.0-beta", 1},
		{"4.2.0-beta.foo", "4.2.0-beta", 1},
		{"1.2.0+bar", "1.2.0+baz", 0},
		{"1.0.0-beta.4", "1.0.0-beta.-2", -1},
		{"1.0.0-beta.-2", "1.0.0-beta.-3", -1},
		{"1.0.0-beta.-3", "1.0.0-beta.5", 1},
	}

	for _, tc := range tests {
		v1, err := NewSemverStr(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemverStr(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Compare(v2)
		assert.Equal(t, a, tc.expected)
	}
}

func TestInc(t *testing.T) {
	tests := []struct {
		v1               string
		expected         string
		how              string
		expectedOriginal string
	}{
		{"1.2.3", "1.2.4", "patch", "1.2.4"},
		{"v1.2.4", "1.2.5", "patch", "v1.2.5"},
		{"1.2.3", "1.3.0", "minor", "1.3.0"},
		{"v1.2.4", "1.3.0", "minor", "v1.3.0"},
		{"1.2.3", "2.0.0", "major", "2.0.0"},
		{"v1.2.4", "2.0.0", "major", "v2.0.0"},
		{"1.2.3+meta", "1.2.4", "patch", "1.2.4"},
		{"1.2.3-beta+meta", "1.2.3", "patch", "1.2.3"},
		{"v1.2.4-beta+meta", "1.2.4", "patch", "v1.2.4"},
		{"1.2.3-beta+meta", "1.3.0", "minor", "1.3.0"},
		{"v1.2.4-beta+meta", "1.3.0", "minor", "v1.3.0"},
		{"1.2.3-beta+meta", "2.0.0", "major", "2.0.0"},
		{"v1.2.4-beta+meta", "2.0.0", "major", "v2.0.0"},
	}

	for _, tc := range tests {
		v1, err := NewSemverStr(tc.v1)
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

		// a := v2.String()
		// assert.Equal(t, tc.expected, a)

		sem, ok := v2.(*Semver)
		assert.True(t, ok)
		a := sem.Original()
		assert.Equal(t, tc.expectedOriginal, a)
	}
}

func TestOriginalVPrefix(t *testing.T) {
	tests := []struct {
		version string
		vprefix string
	}{
		{"1.2.3", ""},
		{"v1.2.4", "v"},
	}

	for _, tc := range tests {
		v1, _ := NewSemverStr(tc.version)
		a := v1.originalVPrefix()
		assert.Equal(t, tc.vprefix, a)
	}
}
