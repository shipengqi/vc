package vc

import (
	"fmt"
	"regexp"
	"strings"
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
	return &Constraints{constraints: gcs}, nil
}

func (c *Constraints) Check(ver Comparable) bool {
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

// func (c *Constraints) CheckString(ver string) bool {
// 	for _, v := range c.constraints {
// 		for _, v2 := range v {
// 			operatorsMap[v2.operator]()
// 		}
// 	}
//
// 	return false
// }

func parseConstraintGroup(group string, fn New, result *[]*constraint) error {
	if strings.Contains(group, " - ") {
		gs := strings.Split(group, " - ")
		if len(gs) > 1 {
			err := parseConstraintGroup(">="+gs[0], fn, result)
			if err != nil {
				return err
			}
			err = parseConstraintGroup("<="+gs[1], fn, result)
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
		group = strings.ReplaceAll(group, "x", "0")
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
	allowed := `\w\.\+-`
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
	valid := findConstraintRegex.MatchString(c)
	if !valid {
		return nil, ErrInvalidConstraint
	}
	var result []*constraint
	var ver, op string
	var err error
	found := findConstraintRegex.FindStringSubmatch(c)
	if len(found) == 0 {
		op = "="
		ver = c
	} else if len(found) == 2 {
		op = "="
		ver = found[1]
	} else if len(found) > 2 {
		op = found[1]
		if op == "" {
			op = "="
		}
		ver = found[2]
	} else {
		return nil, ErrInvalidConstraint
	}

	if op == "^" {
		result, err = parseCaretConstraint(ver, fn)
		if err != nil {
			return nil, err
		}
	} else if op == "~" {
		result, err = parseTildeConstraint(ver, fn)
		if err != nil {
			return nil, err
		}
	} else {
		if ver == "*" {
			ver = "0.0.0"
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
func parseCaretConstraint(ver string, fn New) ([]*constraint, error) {
	var result []*constraint
	ori, err := fn(ver)
	if err != nil {
		return nil, err
	}
	var max Comparable
	if ori.Major() > 0 {
		max = ori.IncMajor()
	} else if ori.Minor() > 0 {
		max = ori.IncMinor()
	} else {
		max = ori.IncPatch()
	}
	result = append(result,
		&constraint{version: ori.Version(), operator: ">=", com: ori},
		&constraint{version: max.Version(), operator: "<", com: max},
	)
	return result, nil
}

// ~2, ~2.x, ~2.x.x  -->  >=2.0.0, <3.0.0
// ~1.0.0            -->  >=1.0.0, <2.0.0
// ~1.2, ~1.2.x      -->  >=1.2.0, <1.3.0
// ~1.2.3,           -->  >=1.2.3, <1.3.0
// ~1.2.0            -->  >=1.2.0, <1.3.0
func parseTildeConstraint(ver string, fn New) ([]*constraint, error) {
	var result []*constraint
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
		&constraint{version: ori.Version(), operator: ">=", com: ori},
		&constraint{version: max.Version(), operator: "<", com: max},
	)
	return result, nil
}
