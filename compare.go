package vc

import (
	"strconv"
	"strings"
)

// Comparable An implementation of Comparable interface can be compared with constraints.
type Comparable interface {
	// Version converts major,minor and patch to a string.
	Version() string
	// Major returns the major version.
	Major() uint64
	// Minor returns the minor version.
	Minor() uint64
	// Patch returns the patch version.
	Patch() uint64
	// Prerelease returns the prerelease version.
	Prerelease() string
	// IncMajor produces the next major version.
	IncMajor() Comparable
	// IncMinor produces the next minor version.
	IncMinor() Comparable
	// IncPatch produces the next patch version.
	IncPatch() Comparable
}

// Lt tests if one version is less than another one.
func Lt(v1, v2 Comparable) bool {
	return Compare(v1, v2) < 0
}

// Gt tests if one version is greater than another one.
func Gt(v1, v2 Comparable) bool {
	return Compare(v1, v2) > 0
}

// Eq tests if two versions are equal to each other.
// Note, versions can be equal with different metadata since metadata
// is not considered part of the comparable version.
func Eq(v1, v2 Comparable) bool {
	return Compare(v1, v2) == 0
}

// Compare compares a Comparable to another one. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version.
//
// Prerelease is lower than the version without a prerelease.
// Compare always takes into account prerelease.
// If you want to work with ranges using typical range syntax that
// skip prerelease if the range is not looking for them use constraints.
func Compare(v1, v2 Comparable) int {
	// Compare the major, minor, and patch version for differences. If a
	// difference is found return the comparison.
	if d := compareSegment(v1.Major(), v2.Major()); d != 0 {
		return d
	}
	if d := compareSegment(v1.Minor(), v2.Minor()); d != 0 {
		return d
	}
	if d := compareSegment(v1.Patch(), v2.Patch()); d != 0 {
		return d
	}

	// At this point the major, minor, and patch versions are the same.
	pre1 := v1.Prerelease()
	pre2 := v2.Prerelease()

	if pre1 == "" && pre2 == "" {
		return 0
	}
	if pre1 == "" {
		return 1
	}
	if pre2 == "" {
		return -1
	}

	return comparePrerelease(pre1, pre2)
}

func compareSegment(v, o uint64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}
	return 0
}

func comparePrerelease(v, o string) int {
	// split the prerelease versions by their part. The separator, per the spec,
	// is a .
	sparts := strings.Split(v, ".")
	oparts := strings.Split(o, ".")

	// Find the longer length of the parts to know how many loop iterations to
	// go through.
	slen := len(sparts)
	olen := len(oparts)

	l := slen
	if olen > slen {
		l = olen
	}

	// Iterate over each part of the prerelease to compare the differences.
	for i := 0; i < l; i++ {
		// Since the lentgh of the parts can be different we need to create
		// a placeholder. This is to avoid out of bounds issues.
		stemp := ""
		if i < slen {
			stemp = sparts[i]
		}

		otemp := ""
		if i < olen {
			otemp = oparts[i]
		}

		d := comparePrePart(stemp, otemp)
		if d != 0 {
			return d
		}
	}

	// Reaching here means two versions are of equal value but have different
	// metadata (the part following a +). They are not identical in string form
	// but the version comparison finds them to be equal.
	return 0
}

func comparePrePart(s, o string) int {
	// Fastpath if they are equal
	if s == o {
		return 0
	}

	// When s or o are empty we can use the other in an attempt to determine
	// the response.
	if s == "" {
		if o != "" {
			return -1
		}
		return 1
	}

	if o == "" {
		if s != "" {
			return 1
		}
		return -1
	}

	// When comparing strings "99" is greater than "103". To handle
	// cases like this we need to detect numbers and compare them. According
	// to the semver spec, numbers are always positive. If there is a - at the
	// start like -99 this is to be evaluated as an alphanum. numbers always
	// have precedence over alphanum. Parsing as Uints because negative numbers
	// are ignored.

	oi, n1 := strconv.ParseUint(o, 10, 64)
	si, n2 := strconv.ParseUint(s, 10, 64)

	// The case where both are strings compare the strings
	if n1 != nil && n2 != nil {
		if s > o {
			return 1
		}
		return -1
	} else if n1 != nil {
		// o is a string and s is a number
		return -1
	} else if n2 != nil {
		// s is a string and o is a number
		return 1
	}
	// Both are numbers
	if si > oi {
		return 1
	}
	return -1
}
