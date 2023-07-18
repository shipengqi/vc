package vc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const semverReg = "^v?(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"

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

func NewSemver(ver string) (*Semver, error) {
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

// From the spec, "Identifiers MUST comprise only
// ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty.
// Numeric identifiers MUST NOT include leading zeroes.". These segments can
// be dot separated.
func validatePrerelease(p string) error {
	eparts := strings.Split(p, ".")
	for _, p := range eparts {
		if containsOnly(p, num) {
			if len(p) > 1 && p[0] == '0' {
				return ErrSegmentStartsZero
			}
		} else if !containsOnly(p, allowed) {
			return ErrInvalidPrerelease
		}
	}

	return nil
}

// From the spec, "Build metadata MAY be denoted by
// appending a plus sign and a series of dot separated identifiers immediately
// following the patch or pre-release version. Identifiers MUST comprise only
// ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty."
func validateMetadata(m string) error {
	eparts := strings.Split(m, ".")
	for _, p := range eparts {
		if !containsOnly(p, allowed) {
			return ErrInvalidMetadata
		}
	}
	return nil
}
