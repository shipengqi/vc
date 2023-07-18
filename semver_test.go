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

	reg := regexp.MustCompile(semverReg)
	var got bool
	for _, v := range tests {
		t.Run(title(t, v), func(t *testing.T) {
			got = reg.MatchString(v.ver)
			assert.Equal(t, v.valid, got)
		})
	}

}

func title(t *testing.T, r regt) string {
	t.Helper()
	return fmt.Sprintf("%s should be %v", r.ver, r.valid)
}
