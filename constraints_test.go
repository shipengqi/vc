package vc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConstraint(t *testing.T) {
	tests := []struct {
		in  string
		err bool
	}{
		{">=1.2.3", false},
		{"1.0", false},
		{"foo", true},
		{"<=1.2", false},
		{">=1.2", false},
		{"v1.2", false},
		{"=1.5", false},
		{">1.3", false},
		{"<1.4.1", false},
		{"<40.50.10", false},
		{"<40.50.10 || =1.5", false},
		{"<40.50.10 || =1.5 <1.9", false},
		{"<40.50.10 || =1.5 <1.9 || =2.0", false},
	}

	for _, tc := range tests {
		c, _ := NewConstraint(tc.in, func(ver string) (Comparable, error) {
			return NewSemver(ver)
		})
		for k, cs := range c.constraints {
			for _, csv := range cs {
				t.Logf("==========> %d, %v", k, csv)
			}
		}
	}
}

func TestParseConstraint(t *testing.T) {
	tests := []struct {
		in       string
		expected *constraint
		err      bool
	}{
		{">=", nil, true},
		{"=", nil, true},
		{">= 1.0.0", nil, true},
		{"1.0.0 >=", nil, true},
		{">= 1.0.0 >=", nil, true},
		{"=<1.2", nil, true},
		{"=>1.2", nil, true},
		{"~>1.5", nil, true},
		{"1.0", &constraint{version: "1.0", operator: "="}, false},
		{"v1.2", &constraint{version: "v1.2", operator: "="}, false},
		{"foo", &constraint{version: "foo", operator: "="}, false},
		{"<1.4.1", &constraint{version: "1.4.1", operator: "<"}, false},
		{"<=1.2", &constraint{version: "1.2", operator: "<="}, false},
		{">1.3", &constraint{version: "1.3", operator: ">"}, false},
		{">=1.2", &constraint{version: "1.2", operator: ">="}, false},
		{"=1.5", &constraint{version: "1.5", operator: "="}, false},
		{"!1.5", &constraint{version: "1.5", operator: "!"}, false},
		{"!=1.5", &constraint{version: "1.5", operator: "!="}, false},
		{"~1.5", &constraint{version: "1.5", operator: "~"}, false},
		{"^1.5", &constraint{version: "1.5", operator: "^"}, false},
	}
	for _, tc := range tests {
		c, err := parseConstraint(tc.in, func(ver string) (Comparable, error) {
			return NewSemver(ver)
		})
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, c)
		}
	}
}

func TestParseConstraintGroup(t *testing.T) {
	tests := []struct {
		in       string
		expected []string
	}{
		{"=1.5", []string{"=1.5"}},
		{"=1.5.0", []string{"=1.5"}},
		{"<40.50.10 =1.5.0", []string{"<40.50.10", "=1.5"}},
		{"<40.50.10 =1.5.0 <1.9.0", []string{"<40.50.10", "=1.5", "<1.9"}},
		{"<40.50.10 >=1.9.0", []string{"<40.50.10", ">=1.9"}},
		{"1.0.0 - 1.5.0 <40.50.10 >=1.9.0", []string{"1.0", "1.5", "<40.50.10", ">=1.9"}},
		{"1.0.0 - 1.5.0 1.1.0 - 1.2.0 <40.50.10 >=1.9.0", []string{"1.0", "1.5", "1.1", "1.2", "<40.50.10", ">=1.9"}},
		{"~1.5", []string{"=1.5"}},
		{"~1.5.1", []string{"=1.5"}},
		{"~1.5.0", []string{"=1.5"}},
		{"~1.0", []string{"=1.5"}},
		{"~0.1", []string{"=1.5"}},
		{"~0.1.2", []string{"=1.5"}},
		{"~0.0.2", []string{"=1.5"}},
		{"^1.5.1", []string{"=1.5"}},
		{"^1.5", []string{"=1.5"}},
		{"^1.5.0", []string{"=1.5"}},
		{"^1.0", []string{"=1.5"}},
		{"^0.1", []string{"=1.5"}},
		{"^0.1.2", []string{"=1.5"}},
		{"^0.0.2", []string{"=1.5"}},
		{"^0.1.x", []string{"=1.5"}},
		{"^1.x.x", []string{"=1.5"}},
	}
	for _, tc := range tests {
		var res []*constraint
		_ = parseConstraintGroup(tc.in, func(ver string) (Comparable, error) {
			return NewSemver(ver)
		}, &res)

		t.Log(res[0], len(res))
		// assert.Equal(t, tc.expected, res)
	}
}

func TestConstraintsCheck(t *testing.T) {
	tests := []struct {
		ver   string
		con   string
		valid bool
	}{
		{"1.2.3", ">=1.2.3", true},
		{"1.2.4", ">=1.2.3", true},
		{"1.3.0", ">=1.2.3", true},
		{"2.0.0", ">=1.2.3", true},
		{"1.2.1", ">=1.2.3", false},
		{"1.1.0", ">=1.2.3", false},
		{"0.1.0", ">=1.2.3", false},
		{"1.2.4", ">1.2.3", true},
		{"1.2.3", ">1.2.3", false},
		{"1.3.0", "<=1.2.3", false},
		{"2.0.0", "<=1.2.3", false},
		{"1.2.1", "<=1.2.3", true},
		{"1.1.0", "<=1.2.3", true},
		{"2.0.0", "<=1.2.3", false},
		{"2.0.0", "<1.2.3", false},
	}

	for _, tc := range tests {
		c, _ := NewConstraint(tc.con, func(ver string) (Comparable, error) {
			return NewSemver(ver)
		})
		ver, _ := NewSemver(tc.ver)
		got := c.Check(ver)
		t.Log(tc.ver, tc.con)
		assert.Equal(t, tc.valid, got)
	}
}
