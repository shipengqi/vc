package vc

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// semVerRegex is the regular expression used to parse a semantic version.
	semverReg string = `^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
		`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
		`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`
	allowedNum   = "0123456789"
	allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-" + allowedNum
)

var _ Comparable = &Semver{}

// The compiled version of the regex created at init() is cached here, so it
// only needs to be created once.
var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile(semverReg)
}

type Semver struct {
	major, minor, patch uint64
	pre                 string
	metadata            string
	original            string
}

// NewSemverStr parses a given version and returns an instance of Semver or
// an error if unable to parse the version. If the version is SemVer-ish it
// attempts to convert it to Semver.
func NewSemverStr(ver string) (*Semver, error) {
	m := versionRegex.FindStringSubmatch(ver)
	if m == nil {
		return nil, ErrInvalidSemVer
	}

	sv := &Semver{
		metadata: m[8],
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

	if sv.metadata != "" {
		if err = validateMetadata(sv.metadata); err != nil {
			return nil, err
		}
	}

	return sv, nil
}

// NewSemver creates a new instance of Semver with each of the parts passed in as
// arguments instead of parsing a version string.
func NewSemver(major, minor, patch uint64, pre, metadata string) *Semver {
	v := Semver{
		major:    major,
		minor:    minor,
		patch:    patch,
		pre:      pre,
		metadata: metadata,
		original: "",
	}

	v.original = v.String()

	return &v
}

// String converts a Semver object to a string.
// Note, if the original version contained a leading v this version will not.
// See the Original() method to retrieve the original value. Semantic Versions
// don't contain a leading v per the spec. Instead, it's optional on
// implementation.
func (v *Semver) String() string {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, "%d.%d.%d", v.major, v.minor, v.patch)
	if v.pre != "" {
		_, _ = fmt.Fprintf(&buf, "-%s", v.pre)
	}
	if v.metadata != "" {
		_, _ = fmt.Fprintf(&buf, "+%s", v.metadata)
	}

	return buf.String()
}

// Version converts major,minor and patch to a string.
func (v *Semver) Version() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "%d.%d.%d", v.major, v.minor, v.patch)
	return buf.String()
}

// Original returns the original value passed in to be parsed.
func (v *Semver) Original() string {
	return v.original
}

// Major returns the major version.
func (v *Semver) Major() uint64 {
	return v.major
}

// Minor returns the minor version.
func (v *Semver) Minor() uint64 {
	return v.minor
}

// Patch returns the patch version.
func (v *Semver) Patch() uint64 {
	return v.patch
}

// Prerelease returns the prerelease version.
func (v *Semver) Prerelease() string {
	return v.pre
}

// Metadata returns the metadata on the version.
func (v *Semver) Metadata() string {
	return v.metadata
}

// IncMajor produces the next major version.
// Sets patch to 0.
// Sets minor to 0.
// Increments major number.
// Unsets metadata.
// Unsets prerelease status.
func (v *Semver) IncMajor() Comparable {
	vNext := &Semver{}
	vNext.metadata = ""
	vNext.pre = ""
	vNext.patch = 0
	vNext.minor = 0
	vNext.major = v.major + 1
	vNext.original = v.originalVPrefix() + "" + vNext.String()
	return vNext
}

// IncMinor produces the next minor version.
// Sets patch to 0.
// Increments minor number.
// Unsets metadata.
// Unsets prerelease status.
func (v *Semver) IncMinor() Comparable {
	vNext := &Semver{}
	vNext.metadata = ""
	vNext.pre = ""
	vNext.patch = 0
	vNext.minor = v.minor + 1
	vNext.major = v.major
	vNext.original = v.originalVPrefix() + "" + vNext.String()
	return vNext
}

// IncPatch produces the next patch version.
// If the current version does not have prerelease/metadata information,
// it unsets metadata and prerelease values, increments patch number.
// If the current version has any of prerelease or metadata information,
// it unsets both values and keeps current patch value
func (v *Semver) IncPatch() Comparable {
	vNext := &Semver{}
	// according to http://semver.org/#spec-item-9
	// Pre-release versions have a lower precedence than the associated normal version.
	// according to http://semver.org/#spec-item-10
	// Build metadata SHOULD be ignored when determining version precedence.
	if v.pre != "" {
		vNext.metadata = ""
		vNext.pre = ""
		vNext.patch = v.patch
	} else {
		vNext.metadata = ""
		vNext.pre = ""
		vNext.patch = v.patch + 1
	}
	vNext.minor = v.minor
	vNext.major = v.major
	vNext.original = v.originalVPrefix() + "" + vNext.String()
	return vNext
}

// Lt tests if one version is less than another one.
func (v *Semver) Lt(o *Semver) bool {
	return v.Compare(o) < 0
}

// Gt tests if one version is greater than another one.
func (v *Semver) Gt(o *Semver) bool {
	return v.Compare(o) > 0
}

// Eq tests if two versions are equal to each other.
// Note, versions can be equal with different metadata since metadata
// is not considered part of the comparable version.
func (v *Semver) Eq(o *Semver) bool {
	return v.Compare(o) == 0
}

// Compare compares this version to another Semver. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version.
//
// Versions are compared by X.Y.Z. Build metadata is ignored. Prerelease is
// lower than the version without a prerelease. Compare always takes into account
// prereleases. If you want to work with ranges using typical range syntaxes that
// skip prereleases if the range is not looking for them use constraints.
func (v *Semver) Compare(o *Semver) int {
	return Compare(v, o)
}

// originalVPrefix returns the original 'v' prefix if any.
func (v *Semver) originalVPrefix() string {
	// Note, only lowercase v is supported as a prefix by the parser.
	if v.original != "" && v.original[:1] == "v" {
		return "v"
	}
	return ""
}
