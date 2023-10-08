# vc
Compare version strings to find greater, equal or lesser. 
Provides the ability to work with Semantic Versions (http://semver.org) in Go. 

Inspired by [Masterminds/semver](https://github.com/Masterminds/semver).

[![test](https://github.com/shipengqi/vc/actions/workflows/test.yaml/badge.svg)](https://github.com/shipengqi/vc/actions/workflows/test.yaml)
[![Codecov](https://codecov.io/gh/shipengqi/vc/branch/main/graph/badge.svg?token=SMU4SI304O)](https://codecov.io/gh/shipengqi/vc)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipengqi/vc)](https://goreportcard.com/report/github.com/shipengqi/vc)
[![Release](https://img.shields.io/github/release/shipengqi/vc.svg)](https://github.com/shipengqi/vc/releases)
[![License](https://img.shields.io/github/license/shipengqi/vc)](https://github.com/shipengqi/vc/blob/main/LICENSE)

## Getting Started

## Semver Versions

```go
v, err := vc.NewSemverStr("0.1.2")

v := NewSemver(0, 1, 2, "", "")
```

## Calendar Versions

```go
v, err := vc.NewCalVerStr("2023.07.05")

v := NewCalVer(2023, 7, 5, "")
```

## Constraints

```go
con, _ := NewConstraint(">=1.1 <2", func(s string) (Comparable, error) {
  return NewSemverStr(s)
})
v, _ := NewSemverStr("1.2.3")

a := con.Check(v)
if a != tc.valid {
    fmt.Printf("Constraint '%s' failing with '%s'", tc.con, tc.ver)
}

// Or
a, _ = con.CheckString("1.2.3")
if a != tc.valid {
    fmt.Printf("Constraint '%s' failing with '%s'", tc.con, tc.ver)
}
```

### Basic Comparisons

There are two elements to the comparisons. First, a comparison string is a list
of space separated **AND** comparisons. These are then separated by `||` (OR)
comparisons. For example, `">=1.2 <3.0.0 || >=4.2.3"` is looking for a
comparison that's greater than or equal to `1.2` and less than `3.0.0` or is
greater than or equal to `4.2.3`.

The basic comparisons are:

* `=`: equal (aliased to no operator)
* `!=`: not equal
* `>`: greater than
* `<`: less than
* `>=`: greater than or equal to
* `<=`: less than or equal to

### Hyphen Range Comparisons

There are multiple methods to handle ranges and the first is hyphens ranges.
These look like:

* `1.2 - 1.4.5` which is equivalent to `>=1.2.0 <=1.4.5`
* `2.3.4 - 4.5` which is equivalent to `>=2.3.4 <=4.5.0`

### Wildcards In Comparisons

The `x`, `X`, and `*` characters can be used as a wildcard character. This works
for all comparison operators. When used on the `=` operator it falls
back to the patch level comparison (see tilde below). For example,

* `1.2.x` is equivalent to `>=1.2.0, <1.3.0`
* `>=1.2.x` is equivalent to `>=1.2.0`
* `<=2.x` is equivalent to `<=3.0.0`
* `*` is equivalent to `>= 0.0.0`

### Tilde Range Comparisons (Patch)

The tilde (`~`) comparison operator is for patch level ranges when a minor
version is specified and major level changes when the minor number is missing.
For example,

* `~1.2.3` is equivalent to `>=1.2.3, <1.3.0`
* `~1` is equivalent to `>=1, <2`
* `~2.3` is equivalent to `>=2.3, <2.4`
* `~1.2.x` is equivalent to `>=1.2.0, <1.3.0`
* `~1.x` is equivalent to `>=1, <2`

### Caret Range Comparisons (Major)

The caret (`^`) comparison operator is for major level changes once a stable
(1.0.0) release has occurred. Prior to a 1.0.0 release the minor versions acts
as the API stability level. This is useful when comparisons of API versions as a
major change is API breaking. For example,

* `^1.2.3` is equivalent to `>= 1.2.3, < 2.0.0`
* `^1.2.x` is equivalent to `>= 1.2.0, < 2.0.0`
* `^2.3` is equivalent to `>= 2.3, < 3`
* `^2.x` is equivalent to `>= 2.0.0, < 3`
* `^0.2.3` is equivalent to `>=0.2.3 <0.3.0`
* `^0.2` is equivalent to `>=0.2.0 <0.3.0`
* `^0.0.3` is equivalent to `>=0.0.3 <0.0.4`
* `^0.0` is equivalent to `>=0.0.0 <0.1.0`
* `^0` is equivalent to `>=0.0.0 <1.0.0`

### Comparable Interface

An implementation of `Comparable` interface can be compared with constraints.

```go
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
```

## 🔋 JetBrains OS licenses

`vc` had been being developed with **IntelliJ IDEA** under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=vc" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo." width="250" align="middle"></a>
