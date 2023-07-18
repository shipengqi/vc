package vc

import "errors"

var (
	// ErrInvalidSemVer is returned a version is found to be invalid when
	// being parsed.
	ErrInvalidSemVer = errors.New("invalid semantic version")

	// ErrParseSemVer is returned a version is found to be invalid when
	// being parsed.
	ErrParseSemVer = errors.New("invalid semantic version")
)
