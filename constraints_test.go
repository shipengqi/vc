package vc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConstraint(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		err      bool
	}{
		{"foo", "", true},
		{">=1.2.3", "1.2.3", false},
		{"1.0", "1.2.3", false},
		{"<=1.2", "1.2.3", false},
		{">=1.2", "1.2.3", false},
		{"v1.2", "1.2.3", false},
		{"=1.5", "1.2.3", false},
		{">1.3", "1.2.3", false},
		{"<1.4.1", "1.2.3", false},
		{"<40.50.10", "1.2.3", false},
		{"<40.50.10 || =1.5", "1.2.3", false},
		{"<40.50.10 || =1.5 <1.9", "1.2.3", false},
		{"<40.50.10 || =1.5 <1.9 || =2.0", "1.2.3", false},
		{"<40.50.10 || =1.5 <1.9 || 2.0 - 3.0", "1.2.3", false},
	}

	for _, tc := range tests {
		c, _ := NewConstraint(tc.in, func(ver string) (Comparable, error) {
			return NewSemverStr(ver)
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
			return NewSemverStr(ver)
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
			return NewSemverStr(ver)
		}, &res)

		t.Log(res[0], len(res))
		// assert.Equal(t, tc.expected, res)
	}
}

func TestConstraintCheck(t *testing.T) {
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
		{"4.1.0", "!=4.1.0", false},
		{"4.1.1", "!=4.1.0", true},
		{"4.1.0", "!=4.1", false},
		{"4.1.1", "!=4.1", true},
		{"5.1.0-alpha.1", "!=4.1", true},
		{"4.1.0", "!=4.1-alpha", true},
		{"5.1.0", "!=4.1", true},
		{"11.1.0", ">11", true},
		{"11.1.0", ">11.1", false},
		{"11.1.1", ">11.1", true},
		{"11.2.1", ">11.1", true},
		{"11.1.2", ">=11", true},
		{"11.1.2", ">=11.1", true},
		{"11.0.2", ">=11.1", false},
		{"4.1.0", ">=1.1", true},
		{"1.1.0", ">=1.1", true},
		{"0.0.9", ">=1.1", false},
		{"0.0.1-alpha", ">=0", true},
		{"0.0.1-alpha", ">=0.0", true},
		{"0.0.1-alpha", ">=0-0", true},
		{"0.0.1-alpha", ">=0.0-0", true},
		{"0.0.0-alpha", ">=0", false},
		{"0.0.0-alpha", ">=0-0", true},
		{"0.0.0-alpha", ">=0.0.0-0", true},
		{"1.2.3", ">=0.0.0-0", true},
		{"3.4.5-beta.1", ">=0.0.0-0", true},
		{"0.0.0-alpha", "<0", true},
		{"0.0.0-alpha", "<0-z", true},
		{"0", ">=0", true},
		{"1", "=0", false},
		{"1", "*", true},
		{"4.5.6", "*", true},
		{"1.2.3-alpha.1", "*", true},
		{"1.2.3-alpha.1", "*-0", true},
		{"1", "2.*", false},
		{"3.4.5", "2.*", false},
		{"2.1.1", "2.*", true},
		{"2.1.1", "2.1.*", true},
		{"2.2.1", "2.1.*", false},
		{"1", "2", false},
		{"3.4.5", "2", false},
		{"2.1.1", "2", false},
		{"2.1.1", "2.1", false},
		{"2.2.1", "2.1", false},
		{"1.2.4", "~1.2.3", true},
		{"1.3.4", "~1.2.3", false},
		{"1.2.4", "~1.2", true},
		{"1.3.4", "~1.2", false},
		{"1.2.4", "~1", true},
		{"2.3.4", "~1", false},
		{"0.2.5", "~0.2.3", true},
		{"0.3.5", "~0.2.3", false},
		{"1.2.3-beta.4", "~1.2.3-beta.2", true},
		{"1.2.4-beta.2", "~1.2.3-beta.2", true},
		{"1.3.4-beta.2", "~1.2.3-beta.2", false},
		{"1.8.9", "^1.2.3", true},
		{"2.8.9", "^1.2.3", false},
		{"1.2.1", "^1.2.3", false},
		{"2.1.0", "^1.1.0", false},
		{"2.2.1", "^1.2.0", false},
		{"1.2.1-alpha.1", "^1.2.0", true},
		{"1.2.1-alpha.1", "^1.2.0-alpha.0", true},
		{"1.2.1-alpha.0", "^1.2.0-alpha.0", true},
		{"1.2.0-alpha.1", "^1.2.0-alpha.2", false},
		{"1.8.9", "^1.2", true},
		{"2.8.9", "^1.2", false},
		{"1.8.9", "^1", true},
		{"2.8.9", "^1", false},
		{"0.2.5", "^0.2.3", true},
		{"0.5.6", "^0.2.3", false},
		{"0.2.5", "^0.2", true},
		{"0.5.6", "^0.2", false},
		{"0.0.3", "^0.0.3", true},
		{"0.0.4", "^0.0.3", false},
		{"0.0.3", "^0.0", true},
		{"0.1.4", "^0.0", false},
		{"1.0.4", "^0.0", false},
		{"0.2.3", "^0", true},
		{"1.1.4", "^0", false},
		{"0.2.3-beta.4", "^0.2.3-beta.2", true},
		{"0.2.4-beta.2", "^0.2.3-beta.2", true},
		{"0.3.4-beta.2", "^0.2.3-beta.2", false},
		{"0.2.3-beta.2", "^0.2.3-beta.2", true},
	}

	for _, tc := range tests {
		c, err := NewConstraint(tc.con, func(ver string) (Comparable, error) {
			return NewSemverStr(ver)
		})
		assert.NoError(t, err)
		ver, err := NewSemverStr(tc.ver)
		assert.NoError(t, err)
		got := c.Check(ver)
		t.Log(tc.ver, tc.con)
		assert.Equal(t, tc.valid, got)
	}
}

func TestConstraintsCheck(t *testing.T) {
	tests := []struct {
		con   string
		ver   string
		valid bool
	}{
		{">1.1 <2", "1.1.1", true},
		{">1.1 <2", "1.2.1", true},
		{">1.1 <3", "4.3.2", false},
		{">=1.1 <2 !=1.2.3", "1.2.3", false},
		{">1.1 <2", "1.2.1", true},
		{">=1.1    <2    !=1.2.3", "1.2.3", false},
		{">1.1 <3", "4.3.2", false},
		{">1.1    <3", "4.3.2", false},
		{">=1.1 <2 !=1.2.3 || >3", "4.1.2", true},
		{">=1.1 <2 !=1.2.3 || >3", "3.1.2", true},
		{">=1.1 <2 !=1.2.3 || >=3", "3.0.0", true},
		{">=1.1 <2 !=1.2.3 || >3", "3.0.0", false},
		{">=1.1 <2 !=1.2.3 || >3", "1.2.3", false},
		{">=1.1 <2 !=1.2.3", "1.2.3", false},
		{"1.1 - 2", "1.1.1", true},
		{"1.5.0 - 4.5", "3.7.0", true},
		{"1.0.0 - 2.0.0 <=2.0.0", "1.5.0", true},
	}

	for _, tc := range tests {
		c, err := NewConstraint(tc.con, func(s string) (Comparable, error) {
			return NewSemverStr(s)
		})
		assert.NoError(t, err)

		v, err := NewSemverStr(tc.ver)
		assert.NoError(t, err)

		a := c.Check(v)
		if a != tc.valid {
			t.Errorf("Constraint '%s' failing with '%s'", tc.con, tc.ver)
		}
	}
}
