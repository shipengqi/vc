package vc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type regt struct {
	ver   string
	valid bool
}

type semvert struct {
	err   bool
	ver   string
	pre   string
	meta  string
	major uint64
	minor uint64
	patch uint64
}

func TestReg(t *testing.T) {
	tests := []regt{
		{"0.0.4", true},
		{"1.2.3", true},
		{"v0.0.4", true},
		{"v1.2.3", true},
		{"10.20.30", true},
		{"1.1.2-prerelease+meta", true},
		{"1.1.2+meta", true},
		{"1.1.2+meta-valid", true},
		{"v1.1.2-prerelease+meta", true},
		{"v1.1.2+meta", true},
		{"v1.1.2+meta-valid", true},
		{"1.0.0-alpha", true},
		{"1.0.0-beta", true},
		{"1.0.0-alpha.beta", true},
		{"1.0.0-alpha.beta.1", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-alpha0.valid", true},
		{"1.0.0-alpha.0valid", true},
		{"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", true},
		{"v1.0.0-alpha0.valid", true},
		{"v1.0.0-alpha.0valid", true},
		{"v1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", true},
		{"1.0.0-rc.1+build.1", true},
		{"2.0.0-rc.1+build.123", true},
		{"1.2.3-beta", true},
		{"10.2.3-DEV-SNAPSHOT", true},
		{"1.2.3-SNAPSHOT-123", true},
		{"v1.2.3-beta", true},
		{"v10.2.3-DEV-SNAPSHOT", true},
		{"v1.2.3-SNAPSHOT-123", true},
		{"1.0.0", true},
		{"2.0.0", true},
		{"1.1.7", true},
		{"2.0.0+build.1848", true},
		{"2.0.1-alpha.1227", true},
		{"1.0.0-alpha+beta", true},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12+788", true},
		{"1.2.3----R-S.12.9.1--.12+meta", true},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12", true},
		{"1.0.0+0.build.1-rc.10000aaa-kk-0.1", true},
		{"99999999999999999999999.999999999999999999.99999999999999999", true},
		{"1.0.0-0A.is.legal", true},
		{"1", false},
		{"1.2", false},
		{"1.2.3-0123", false},
		{"1.2.3-0123.0123", false},
		{"c1.2.3-0123", false},
		{"v1.2.3-0123.0123", false},
		{"1.1.2+.123", false},
		{"+invalid", false},
		{"-invalid", false},
		{"-invalid+invalid", false},
		{"-invalid.01", false},
		{"alpha", false},
		{"alpha.beta", false},
		{"alpha.beta.1", false},
		{"alpha.1", false},
		{"alpha+beta", false},
		{"alpha_beta", false},
		{"alpha.", false},
		{"alpha..", false},
		{"beta", false},
		{"1.0.0-alpha_beta", false},
		{"-alpha.", false},
		{"1.0.0-alpha..", false},
		{"1.0.0-alpha..1", false},
		{"1.0.0-alpha...1", false},
		{"1.0.0-alpha....1", false},
		{"1.0.0-alpha.....1", false},
		{"1.0.0-alpha......1", false},
		{"1.0.0-alpha.......1", false},
		{"01.1.1", false},
		{"1.01.1", false},
		{"1.1.01", false},
		{"1.2", false},
		{"1.2.3.DEV", false},
		{"1.2-SNAPSHOT", false},
		{"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788", false},
		{"1.2-RC-SNAPSHOT", false},
		{"-1.0.3-gamma+b7718", false},
		{"+justmeta", false},
		{"9.8.7+meta+meta", false},
		{"9.8.7-whatever+meta+meta", false},
		{"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12", false},
	}

	reg := regexp.MustCompile(strictSemverReg)
	var got bool
	for _, v := range tests {
		t.Run(title(t, v), func(t *testing.T) {
			got = reg.MatchString(v.ver)
			assert.Equal(t, v.valid, got)
		})
	}
}

func TestNewSemver(t *testing.T) {
	tests := []semvert{
		{true, "foo", "", "", 0, 0, 0},
		{false, "0.0.4", "", "", 0, 0, 4},
		{false, "1.2.3", "", "", 1, 2, 3},
		{false, "v0.0.4", "", "", 0, 0, 4},
		{false, "10.20.30", "", "", 10, 20, 30},
		{false, "1.1.2-prerelease+meta", "prerelease", "meta", 1, 1, 2},
		{false, "v1.1.2-prerelease+meta", "prerelease", "meta", 1, 1, 2},
		{false, "1.1.2+meta", "", "meta", 1, 1, 2},
		{false, "1.1.2+meta-valid", "", "meta-valid", 1, 1, 2},
		{false, "1.1.2-alpha", "alpha", "", 1, 1, 2},
		{false, "1.1.2-beta", "beta", "", 1, 1, 2},
		{false, "1.1.2-alpha.beta", "alpha.beta", "", 1, 1, 2},
		{false, "1.1.2-alpha.beta.1", "alpha.beta.1", "", 1, 1, 2},
		{false, "1.1.2-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", "alpha-a.b-c-somethinglong", "build.1-aef.1-its-okay", 1, 1, 2},
		{false, "1.1.2-rc.1+build.1", "rc.1", "build.1", 1, 1, 2},
		{false, "1.1.2-DEV-SNAPSHOT", "DEV-SNAPSHOT", "", 1, 1, 2},
		{false, "1.1.2-SNAPSHOT-123", "SNAPSHOT-123", "", 1, 1, 2},
		{false, "1.1.2+build.1848", "", "build.1848", 1, 1, 2},
		{false, "1.1.2----RC-SNAPSHOT.12.9.1--.12+788", "---RC-SNAPSHOT.12.9.1--.12", "788", 1, 1, 2},
		{false, "1.1.2----R-S.12.9.1--.12+meta", "---R-S.12.9.1--.12", "meta", 1, 1, 2},
		{false, "1.1.2----RC-SNAPSHOT.12.9.1--.12", "---RC-SNAPSHOT.12.9.1--.12", "", 1, 1, 2},
		{false, "1", "", "", 1, 0, 0},
		{false, "1.2", "", "", 1, 2, 0},
		{false, "1.1----R-S.12.9.1--.12+meta", "---R-S.12.9.1--.12", "meta", 1, 1, 0},
		{false, "1.1----RC-SNAPSHOT.12.9.1--.12", "---RC-SNAPSHOT.12.9.1--.12", "", 1, 1, 0},
		{false, "1----R-S.12.9.1--.12+meta", "---R-S.12.9.1--.12", "meta", 1, 0, 0},
		{false, "1----RC-SNAPSHOT.12.9.1--.12", "---RC-SNAPSHOT.12.9.1--.12", "", 1, 0, 0},
	}

	var got *Semver
	var err error
	for _, v := range tests {
		t.Run(fmt.Sprintf("%s should be %v", v.ver, v.err), func(t *testing.T) {
			got, err = NewSemver(v.ver)
			if v.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v.major, got.Major())
				assert.Equal(t, v.minor, got.Minor())
				assert.Equal(t, v.patch, got.Patch())
				assert.Equal(t, v.pre, got.Prerelease())
				assert.Equal(t, v.meta, got.Metadata())
			}
		})
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
		v1, err := NewSemver(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemver(tc.v2)
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
		v1, err := NewSemver(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemver(tc.v2)
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
		v1, err := NewSemver(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemver(tc.v2)
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
		v1, err := NewSemver(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewSemver(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.Compare(v2)
		assert.Equal(t, a, tc.expected)
	}
}
func title(t *testing.T, r regt) string {
	t.Helper()
	return fmt.Sprintf("%s should be %v", r.ver, r.valid)
}
