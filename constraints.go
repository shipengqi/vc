package vc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shipengqi/golib/strutil"
)

const (
	OperatorGte   = ">="
	OperatorGt    = ">"
	OperatorLte   = "<="
	OperatorLt    = "<"
	OperatorEq    = "="
	OperatorRange = " - "
	OperatorCaret = "^"
	OperatorTilde = "~"
)

const (
	VersionAll      = "*"
	VersionX        = "x"
	VersionMinimum  = "0.0.0"
	VersionAllAlias = OperatorGte + VersionMinimum
)

var operatorsMap map[string]operation

type operation func(v Comparable, c *constraint) bool

// Constraints is one or more constraint that a version can be
// checked against.
type Constraints struct {
	constraints [][]*constraint
	newfn       New
}

// New a function to generate a Comparable instance.
type New func(string) (Comparable, error)

// NewConstraint returns a Constraints instance that a Comparable instance can
// be checked against. If there is a parse error it will be returned.
func NewConstraint(c string, fn New) (*Constraints, error) {
	if strutil.IsEmpty(c) {
		return nil, ErrInvalidConstraint
	}
	groups := strings.Split(c, "||")
	gcs := make([][]*constraint, len(groups))

	for k, v := range groups {
		var result []*constraint
		err := parseConstraintGroup(v, fn, &result)
		if err != nil {
			return nil, err
		}
		gcs[k] = result
	}
	return &Constraints{constraints: gcs, newfn: fn}, nil
}

func (c *Constraints) Check(ver Comparable) bool {
	return c.check(ver)
}

func (c *Constraints) CheckString(ver string) (bool, error) {
	com, err := c.newfn(ver)
	if err != nil {
		return false, err
	}
	return c.check(com), nil
}

func (c *Constraints) check(ver Comparable) bool {
	for _, v := range c.constraints {
		joy := true
		for _, v2 := range v {
			if !operatorsMap[v2.operator](ver, v2) {
				joy = false
				break
			}
		}
		if joy {
			return true
		}
	}

	return false
}

func parseConstraintGroup(group string, fn New, result *[]*constraint) error {
	group = strings.TrimSpace(group)
	if strings.Contains(group, OperatorRange) {
		gs := strings.Split(group, OperatorRange)
		// Contains more than one " - "
		if len(gs) > 2 {
			return ErrInvalidConstraint
		}
		if len(gs) > 1 {
			err := parseConstraintGroup(OperatorGte+gs[0], fn, result)
			if err != nil {
				return err
			}
			err = parseConstraintGroup(OperatorLte+gs[1], fn, result)
			if err != nil {
				return err
			}
		}
	} else if strings.Contains(group, " ") {
		gs := strings.Split(group, " ")
		for _, gv := range gs {
			if gv != "" {
				err := parseConstraintGroup(gv, fn, result)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if group == VersionAll || strings.HasPrefix(group, VersionAll) {
			group = VersionAllAlias
		}
		cons, err := parseConstraint(group, fn)
		if err != nil {
			return err
		}
		*result = append(*result, cons...)
	}
	return nil
}

// Used to find individual constraints within a multi-constraint string
var findConstraintRegex *regexp.Regexp

func init() {
	operatorsMap = map[string]operation{
		"=":  constraintEqual,
		"!":  constraintNotEqual,
		"!=": constraintNotEqual,
		">":  constraintGreaterThan,
		"<":  constraintLessThan,
		">=": constraintGreaterThanEqual,
		"<=": constraintLessThanEqual,
	}

	ops := `\^|>=|<=|!=|!|>|<|~|=`
	allowed := `\w\.\+\*-`
	findConstraintRegex = regexp.MustCompile(fmt.Sprintf(
		`^(%s)?([%s]+)$`, ops, allowed))
}

// An individual constraint
type constraint struct {
	// The original value
	original string
	// The original version (e.g., 4.x from != 4.x)
	version string
	// The original operator for the constraint
	operator string
	com      Comparable
}

func parseConstraint(c string, fn New) ([]*constraint, error) {
	// replace x to 0
	// c = strings.ReplaceAll(c, "x", "0")

	valid := findConstraintRegex.MatchString(c)
	if !valid {
		return nil, ErrInvalidConstraint
	}
	var result []*constraint
	var ver, op string
	var err error
	found := findConstraintRegex.FindStringSubmatch(c)
	if len(found) == 0 {
		op = OperatorEq
		ver = c
	} else if len(found) == 2 {
		op = OperatorEq
		ver = found[1]
	} else if len(found) > 2 {
		op = found[1]
		if op == "" {
			op = OperatorEq
		}
		ver = found[2]
	} else {
		return nil, ErrInvalidConstraint
	}

	if op == OperatorCaret {
		result, err = parseCaretConstraint(c, ver, fn)
		if err != nil {
			return nil, err
		}
	} else if op == OperatorTilde {
		result, err = parseTildeConstraint(c, ver, fn)
		if err != nil {
			return nil, err
		}
	} else if strings.Contains(ver, VersionAll) || strings.Contains(ver, "x") {
		result, err = parseStarConstraint(c, ver, fn)
		if err != nil {
			return nil, err
		}
	} else {
		if ver == VersionAll {
			op = OperatorGte
			ver = VersionMinimum
		}
		com, err := fn(ver)
		if err != nil {
			return nil, err
		}
		result = append(result, &constraint{original: c, version: ver, operator: op, com: com})
	}
	return result, nil
}

func constraintEqual(ver Comparable, c *constraint) bool {
	return Compare(ver, c.com) == 0
}

func constraintNotEqual(ver Comparable, c *constraint) bool {
	return Compare(ver, c.com) != 0
}

func constraintGreaterThan(ver Comparable, c *constraint) bool {
	return Compare(ver, c.com) > 0
}

func constraintLessThan(ver Comparable, c *constraint) bool {
	return Compare(ver, c.com) < 0
}

func constraintGreaterThanEqual(ver Comparable, c *constraint) bool {
	return Compare(ver, c.com) >= 0
}

func constraintLessThanEqual(ver Comparable, c *constraint) bool {
	return Compare(ver, c.com) <= 0
}

// ^1.2.3  -->  >=1.2.3 <2.0.0
// ^1.2    -->  >=1.2.0 <2.0.0
// ^1      -->  >=1.0.0 <2.0.0
// ^0.2.3  -->  >=0.2.3 <0.3.0
// ^0.2    -->  >=0.2.0 <0.3.0
// ^0.0.3  -->  >=0.0.3 <0.0.4
// ^0.0    -->  >=0.0.0 <0.1.0
// ^0      -->  >=0.0.0 <1.0.0
// ^0.x    -->  >=0.0.0 <1.0.0
// ^1.x    -->  >=1.0.0 <2.0.0
func parseCaretConstraint(original, ver string, fn New) ([]*constraint, error) {
	var result []*constraint
	ver = strings.ReplaceAll(ver, VersionX, "0")

	ori, err := fn(ver)
	if err != nil {
		return nil, err
	}
	var max Comparable
	if ori.Major() > 0 {
		max = ori.IncMajor()
	} else if ori.Minor() > 0 {
		max = ori.IncMinor()
	} else if ori.Patch() > 0 {
		max = ori.IncPatch()
	} else {
		// version is ^0.0.0
		vs := strings.Split(ver, ".")
		if len(vs) == 1 {
			max = ori.IncMajor()
		} else if len(vs) == 2 {
			max = ori.IncMinor()
		} else {
			max = ori.IncPatch()
		}
	}
	result = append(result,
		&constraint{version: ori.Version(), operator: OperatorGte, com: ori, original: original},
		&constraint{version: max.Version(), operator: OperatorLt, com: max, original: original},
	)
	return result, nil
}

// ~2, ~2.x, ~2.x.x  -->  >=2.0.0, <3.0.0
// ~1.0.0            -->  >=1.0.0, <2.0.0
// ~1.2, ~1.2.x      -->  >=1.2.0, <1.3.0
// ~1.2.3,           -->  >=1.2.3, <1.3.0
// ~1.2.0            -->  >=1.2.0, <1.3.0
func parseTildeConstraint(original, ver string, fn New) ([]*constraint, error) {
	var result []*constraint
	ver = strings.ReplaceAll(ver, VersionX, "0")
	ori, err := fn(ver)
	if err != nil {
		return nil, err
	}
	var max Comparable
	if ori.Minor() == 0 && ori.Patch() == 0 {
		max = ori.IncMajor()
	} else {
		max = ori.IncMinor()
	}
	result = append(result,
		&constraint{version: ori.Version(), operator: OperatorGte, com: ori, original: original},
		&constraint{version: max.Version(), operator: OperatorLt, com: max, original: original},
	)
	return result, nil
}

// 2.*    -->  >=2.0.0, <3.0.0
// 2.1.*  -->  >=2.1.0, <2.2.0
func parseStarConstraint(original, ver string, fn New) ([]*constraint, error) {
	var result []*constraint
	var minorall, patchall bool
	vs := strings.Split(ver, ".")
	if len(vs) == 1 {
		vs = append(vs, "0", "0")
	} else if len(vs) == 2 {
		if vs[1] == VersionAll || vs[1] == VersionX {
			vs[1] = "0"
			minorall = true
		}
	} else if len(vs) > 2 {
		if vs[1] == VersionAll || vs[1] == VersionX {
			minorall = true
			vs[1] = "0"
		} else if vs[2] == VersionAll || vs[2] == VersionX {
			vs[2] = "0"
			patchall = true
		}
	}
	ver = strings.Join(vs, ".")
	ori, err := fn(ver)
	if err != nil {
		return nil, err
	}
	var max Comparable
	if minorall {
		max = ori.IncMajor()
	} else if patchall {
		max = ori.IncMinor()
	}
	result = append(result,
		&constraint{version: ori.Version(), operator: OperatorGte, com: ori, original: original},
		&constraint{version: max.Version(), operator: OperatorLt, com: max, original: original},
	)
	return result, nil
}
