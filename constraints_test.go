package vc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConstraint(t *testing.T) {
	tests := []struct {
		in       string
		expected []*constraint
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
		{"foo", nil, true},
		{"1.0", []*constraint{
			{version: "1.0", operator: "=", original: "1.0"},
		}, false},
		{"v1.2", []*constraint{{
			version: "v1.2", operator: "=", original: "v1.2",
		}}, false},
		{"<1.4.1", []*constraint{
			{version: "1.4.1", operator: "<", original: "<1.4.1"},
		}, false},
		{"<=1.2", []*constraint{
			{version: "1.2", operator: "<=", original: "<=1.2"},
		}, false},
		{">1.3", []*constraint{
			{version: "1.3", operator: ">", original: ">1.3"},
		}, false},
		{">=1.2", []*constraint{
			{version: "1.2", operator: ">=", original: ">=1.2"},
		}, false},
		{"=1.5", []*constraint{
			{version: "1.5", operator: "=", original: "=1.5"},
		}, false},
		{"!1.5", []*constraint{
			{version: "1.5", operator: "!", original: "!1.5"},
		}, false},
		{"!=1.5", []*constraint{
			{version: "1.5", operator: "!=", original: "!=1.5"},
		}, false},
		{"~1.5", []*constraint{
			{version: "1.5.0", operator: ">=", original: "~1.5"},
			{version: "1.6.0", operator: "<", original: "~1.5"},
		}, false},
		{"^1.5", []*constraint{
			{version: "1.5.0", operator: ">=", original: "^1.5"},
			{version: "2.0.0", operator: "<", original: "^1.5"},
		}, false},
		{"^0", []*constraint{
			{version: "0.0.0", operator: ">=", original: "^0"},
			{version: "1.0.0", operator: "<", original: "^0"},
		}, false},
		{"^0.0", []*constraint{
			{version: "0.0.0", operator: ">=", original: "^0.0"},
			{version: "0.1.0", operator: "<", original: "^0.0"},
		}, false},
	}
	for _, tc := range tests {
		c, err := parseConstraint(tc.in, func(ver string) (Comparable, error) {
			return NewSemverStr(ver)
		})
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expected), len(c))
			for k, com := range c {
				assert.Equal(t, tc.expected[k].version, com.version)
				assert.Equal(t, tc.expected[k].operator, com.operator)
				assert.Equal(t, tc.expected[k].original, com.original)
			}

		}
	}
}

func TestParseConstraintGroup(t *testing.T) {
	tests := []struct {
		in       string
		expected []*constraint
		err      bool
	}{
		{"<40.50.10 =1.5.0", []*constraint{
			{version: "40.50.10", operator: "<", original: "<40.50.10"},
			{version: "1.5.0", operator: "=", original: "=1.5.0"},
		}, false},
		{"<40.50.10 =1.5.0 <1.9.0", []*constraint{
			{version: "40.50.10", operator: "<", original: "<40.50.10"},
			{version: "1.5.0", operator: "=", original: "=1.5.0"},
			{version: "1.9.0", operator: "<", original: "<1.9.0"},
		}, false},
		{"1.0.0 - 1.5.0 1.1.0 - 1.2.0 <40.50.10 >=1.9.0", nil, true},
		{"1.0.0 - 1.5.0 <40.50.10 >=1.9.0", []*constraint{
			{version: "1.0.0", operator: ">=", original: ">=1.0.0"},
			{version: "1.5.0", operator: "<=", original: "<=1.5.0"},
			{version: "40.50.10", operator: "<", original: "<40.50.10"},
			{version: "1.9.0", operator: ">=", original: ">=1.9.0"},
		}, false},
		{"~1.5", []*constraint{
			{version: "1.5.0", operator: ">=", original: "~1.5"},
			{version: "1.6.0", operator: "<", original: "~1.5"},
		}, false},
		{"~1.5.1", []*constraint{
			{version: "1.5.1", operator: ">=", original: "~1.5.1"},
			{version: "1.6.0", operator: "<", original: "~1.5.1"},
		}, false},
		{"~1.5.0", []*constraint{
			{version: "1.5.0", operator: ">=", original: "~1.5.0"},
			{version: "1.6.0", operator: "<", original: "~1.5.0"},
		}, false},
		{"~1.0", []*constraint{
			{version: "1.0.0", operator: ">=", original: "~1.0"},
			{version: "2.0.0", operator: "<", original: "~1.0"},
		}, false},
		{"~0.1", []*constraint{
			{version: "0.1.0", operator: ">=", original: "~0.1"},
			{version: "0.2.0", operator: "<", original: "~0.1"},
		}, false},
		{"~0.1.2", []*constraint{
			{version: "0.1.2", operator: ">=", original: "~0.1.2"},
			{version: "0.2.0", operator: "<", original: "~0.1.2"},
		}, false},
		{"~0.0.2", []*constraint{
			{version: "0.0.2", operator: ">=", original: "~0.0.2"},
			{version: "0.1.0", operator: "<", original: "~0.0.2"},
		}, false},
		{"~1.x", []*constraint{
			{version: "1.0.0", operator: ">=", original: "~1.x"},
			{version: "2.0.0", operator: "<", original: "~1.x"},
		}, false},
		{"~1.2.x", []*constraint{
			{version: "1.2.0", operator: ">=", original: "~1.2.x"},
			{version: "1.3.0", operator: "<", original: "~1.2.x"},
		}, false},
		{"^1.5", []*constraint{
			{version: "1.5.0", operator: ">=", original: "^1.5"},
			{version: "2.0.0", operator: "<", original: "^1.5"},
		}, false},
		{"^1.5.1", []*constraint{
			{version: "1.5.1", operator: ">=", original: "^1.5.1"},
			{version: "2.0.0", operator: "<", original: "^1.5.1"},
		}, false},
		{"^1.5.0", []*constraint{
			{version: "1.5.0", operator: ">=", original: "^1.5.0"},
			{version: "2.0.0", operator: "<", original: "^1.5.0"},
		}, false},
		{"^1.0", []*constraint{
			{version: "1.0.0", operator: ">=", original: "^1.0"},
			{version: "2.0.0", operator: "<", original: "^1.0"},
		}, false},
		{"^0.1", []*constraint{
			{version: "0.1.0", operator: ">=", original: "^0.1"},
			{version: "0.2.0", operator: "<", original: "^0.1"},
		}, false},
		{"^0.1.2", []*constraint{
			{version: "0.1.2", operator: ">=", original: "^0.1.2"},
			{version: "0.2.0", operator: "<", original: "^0.1.2"},
		}, false},
		{"^0.0.2", []*constraint{
			{version: "0.0.2", operator: ">=", original: "^0.0.2"},
			{version: "0.0.3", operator: "<", original: "^0.0.2"},
		}, false},
		{"^0.1.x", []*constraint{
			{version: "0.1.0", operator: ">=", original: "^0.1.x"},
			{version: "0.2.0", operator: "<", original: "^0.1.x"},
		}, false},
		{"^1.x.x", []*constraint{
			{version: "1.0.0", operator: ">=", original: "^1.x.x"},
			{version: "2.0.0", operator: "<", original: "^1.x.x"},
		}, false},
		{">=1.2.x", []*constraint{
			{version: "1.2.0", operator: ">=", original: ">=1.2.x"},
			{version: "1.3.0", operator: "<", original: ">=1.2.x"},
		}, false},
		{"<=2.x", []*constraint{
			{version: "2.0.0", operator: ">=", original: "<=2.x"},
			{version: "3.0.0", operator: "<", original: "<=2.x"},
		}, false},
		{"*", []*constraint{
			{version: "0.0.0", operator: ">=", original: ">=0.0.0"},
		}, false},
	}
	for _, tc := range tests {
		var c []*constraint
		err := parseConstraintGroup(tc.in, func(ver string) (Comparable, error) {
			return NewSemverStr(ver)
		}, &c)
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expected), len(c))
			for k, com := range c {
				assert.Equal(t, tc.expected[k].version, com.version)
				assert.Equal(t, tc.expected[k].operator, com.operator)
				assert.Equal(t, tc.expected[k].original, com.original)
			}

		}
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

func TestConstraintsCheckString(t *testing.T) {
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
		a, err := c.CheckString(tc.ver)
		assert.NoError(t, err)
		if a != tc.valid {
			t.Errorf("Constraint '%s' failing with '%s'", tc.con, tc.ver)
		}
	}
}
