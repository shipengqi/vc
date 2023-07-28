package vc

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const calVerReg = `^([0-9]{4}|[0-9]{2})(\.[0-9]{1,2})?(\.[0-9]{1,2})?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`

// The compiled version of the regex created at init() is cached here, so it
// only needs to be created once.
var calVersionRegex *regexp.Regexp

func init() {
	calVersionRegex = regexp.MustCompile(calVerReg)
}

var _ Comparable = &CalVer{}

type CalVer struct {
	major, minor, patch uint64
	pre                 string
	original            string
}

// NewCalVerStr parses a given version and returns an instance of CalVer or
// an error if unable to parse the version. If the version is SemVer-ish it
// attempts to convert it to CalVer.
func NewCalVerStr(ver string) (*CalVer, error) {
	m := calVersionRegex.FindStringSubmatch(ver)
	if m == nil {
		return nil, ErrInvalidCalVer
	}

	sv := &CalVer{
		pre:      m[5],
		original: ver,
	}

	var err error
	sv.major, err = strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing version segment: %s", err)
	}
	if m[2] != "" {
		sv.minor, err = strconv.ParseUint(strings.TrimPrefix(m[2], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing version segment: %s", err)
		}
	} else {
		sv.minor = 0
	}

	if m[3] != "" {
		sv.patch, err = strconv.ParseUint(strings.TrimPrefix(m[3], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing version segment: %s", err)
		}
	} else {
		sv.patch = 0
	}

	// Perform some basic due diligence on the extra parts to ensure they are
	// valid.
	if sv.pre != "" {
		if err = validatePrerelease(sv.pre); err != nil {
			return nil, err
		}
	}
	return sv, nil
}

// NewCalVer creates a new instance of CalVer with each of the parts passed in as
// arguments instead of parsing a version string.
func NewCalVer(major, minor, patch uint64, pre string) *CalVer {
	v := CalVer{
		major:    major,
		minor:    minor,
		patch:    patch,
		pre:      pre,
		original: "",
	}
	v.original = v.String()

	return &v
}

// String converts a CalVer object to a string.
// Note, if the original version contained a leading v this version will not.
// See the Original() method to retrieve the original value. Semantic Versions
// don't contain a leading v per the spec. Instead, it's optional on
// implementation.
func (v *CalVer) String() string {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, "%d.%02d.%02d", v.major, v.minor, v.patch)
	if v.pre != "" {
		_, _ = fmt.Fprintf(&buf, "-%s", v.pre)
	}
	return buf.String()
}

// Version converts major,minor and patch to a string.
func (v *CalVer) Version() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "%d.%02d.%02d", v.major, v.minor, v.patch)
	return buf.String()
}

// Original returns the original value passed in to be parsed.
func (v *CalVer) Original() string {
	return v.original
}

// Major returns the major version.
func (v *CalVer) Major() uint64 {
	return v.major
}

// Minor returns the minor version.
func (v *CalVer) Minor() uint64 {
	return v.minor
}

// Patch returns the patch version.
func (v *CalVer) Patch() uint64 {
	return v.patch
}

// Prerelease returns the prerelease version.
func (v *CalVer) Prerelease() string {
	return v.pre
}

// Metadata returns empty string.
func (v *CalVer) Metadata() string {
	return ""
}

// IncMajor produces the next major version.
// Sets patch to 0.
// Sets minor to 0.
// Increments major number.
// Unsets prerelease status.
func (v *CalVer) IncMajor() Comparable {
	vNext := &CalVer{}
	vNext.pre = ""
	vNext.patch = 0
	vNext.minor = 0
	vNext.major = v.major + 1
	vNext.original = vNext.String()
	return vNext
}

// IncMinor produces the next minor version.
// Sets patch to 0.
// Increments minor number.
// Unsets prerelease status.
func (v *CalVer) IncMinor() Comparable {
	vNext := &CalVer{}
	vNext.pre = ""
	vNext.patch = 0
	vNext.minor = v.minor + 1
	vNext.major = v.major
	vNext.original = vNext.String()
	return vNext
}

// IncPatch produces the next patch version.
// If the current version does not have prerelease information,
// it unsets prerelease values, increments patch number.
// If the current version has any of prerelease information,
// it unsets both values and keeps current patch value
func (v *CalVer) IncPatch() Comparable {
	vNext := &CalVer{}
	if v.pre != "" {
		vNext.pre = ""
		vNext.patch = v.patch
	} else {
		vNext.pre = ""
		vNext.patch = v.patch + 1
	}
	vNext.minor = v.minor
	vNext.major = v.major
	vNext.original = vNext.String()
	return vNext
}

// Lt tests if one version is less than another one.
func (v *CalVer) Lt(o *CalVer) bool {
	return v.Compare(o) < 0
}

// Gt tests if one version is greater than another one.
func (v *CalVer) Gt(o *CalVer) bool {
	return v.Compare(o) > 0
}

// Eq tests if two versions are equal to each other.
func (v *CalVer) Eq(o *CalVer) bool {
	return v.Compare(o) == 0
}

// Compare compares this version to another CalVer. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version.
func (v *CalVer) Compare(o *CalVer) int {
	return Compare(v, o)
}
