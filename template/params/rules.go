package params

import (
	"errors"
	"fmt"
	"strings"
)

type Rule interface {
	Validate(input []string) error
	Required() []string
	Optionals() []string
	Missing(input []string) []string
	String() string
}

type allOf struct {
	defaultRule
}

func AllOf(rules ...Rule) Rule {
	return allOf{build(rules)}
}

func (n allOf) Validate(input []string) (err error) {
	for _, r := range n.rules {
		err = r.Validate(input)
		if err != optErr && err != nil {
			return errors.New(n.String())
		}
	}
	return nil
}

func (n allOf) Missing(input []string) (miss []string) {
	for _, r := range n.rules {
		miss = append(miss, r.Missing(input)...)
	}
	return
}

func (n allOf) Required() (all []string) {
	for _, r := range n.rules {
		all = append(all, r.Required()...)
	}
	return
}

func (n allOf) String() string {
	return n.rules.join(" + ")
}

type onlyOneOf struct {
	defaultRule
}

func OnlyOneOf(rules ...Rule) Rule {
	return onlyOneOf{build(rules)}
}

func (n onlyOneOf) Validate(input []string) error {
	if len(n.rules) == 0 {
		return nil
	}
	var pass int
	for _, r := range n.rules {
		if err := r.Validate(input); err == nil {
			pass++
		}
	}
	if pass != 1 {
		return fmt.Errorf("only %s", n.rules)
	}
	return nil
}

func (n onlyOneOf) Missing(input []string) (miss []string) {
	if err := n.Validate(input); err != nil && len(n.rules) > 0 {
		miss = append(miss, n.rules[0].Missing(input)...)
	}
	return
}

func (n onlyOneOf) Required() (all []string) {
	if len(n.rules) > 0 {
		all = append(all, n.rules[0].Required()...)
	}
	return
}

func (n onlyOneOf) String() string {
	return "(" + n.rules.join(" | ") + ")"
}

type atLeastOneOf struct {
	defaultRule
}

func AtLeastOneOf(rules ...Rule) Rule {
	return atLeastOneOf{build(rules)}
}

func (n atLeastOneOf) Validate(input []string) error {
	if len(n.rules) == 0 {
		return nil
	}
	var pass int
	for _, r := range n.rules {
		if err := r.Validate(input); err == nil || err == optErr {
			pass++
		}
	}
	if pass < 1 {
		return fmt.Errorf("at least one of %v", n.rules)
	}
	return nil
}

func (n atLeastOneOf) Missing(input []string) (miss []string) {
	if err := n.Validate(input); err != nil && len(n.rules) > 0 {
		miss = append(miss, n.rules[0].Missing(input)...)
	}
	return
}

func (n atLeastOneOf) Required() (all []string) {
	if len(n.rules) > 0 {
		all = append(all, n.rules[0].Required()...)
	}
	return
}

func (n atLeastOneOf) String() string {
	return "(" + n.rules.join(" / ") + ")"
}

type opt struct {
	optionals []string
}

func Opt(s ...string) Rule {
	o := opt{}
	o.optionals = append(o.optionals, s...)
	return o
}

var optErr = errors.New("opt err")

func (n opt) Validate(input []string) error {
	return optErr
}

func (n opt) Missing(input []string) (miss []string) {
	return
}

func (n opt) Required() []string {
	return []string{}
}

func (n opt) Optionals() []string {
	return n.optionals
}

func (n opt) String() string {
	return "[" + strings.Join(n.optionals, " ") + "]"
}

type Key string

func (n Key) Validate(input []string) error {
	if s := string(n); !contains(input, s) {
		return errors.New(s)
	}
	return nil
}

func (n Key) Missing(input []string) (miss []string) {
	if s := string(n); !contains(input, s) {
		miss = append(miss, s)
	}
	return
}

func (n Key) Required() []string {
	return []string{string(n)}
}

func (n Key) Optionals() []string {
	return []string{}
}

func (n Key) String() string {
	return string(n)
}

type none struct{}

func None() Rule {
	return none{}
}

func (n none) Validate(input []string) error   { return nil }
func (n none) Required() []string              { return []string{} }
func (n none) Optionals() []string             { return []string{} }
func (n none) Missing(input []string) []string { return []string{} }
func (n none) String() string                  { return "none" }

func build(rules []Rule) (d defaultRule) {
	for _, n := range rules {
		d.rules = append(d.rules, n)
	}
	return
}

type defaultRule struct {
	rules rules
}

type rules []Rule

func (rs rules) join(sep string) string {
	var arr []string
	for _, r := range rs {
		arr = append(arr, fmt.Sprint(r))
	}
	return strings.Join(arr, sep)
}

func (r defaultRule) Optionals() (o []string) {
	for _, r := range r.rules {
		switch v := r.(type) {
		case opt:
			o = append(o, v.optionals...)
		}
	}
	return
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if s == a {
			return true
		}
	}
	return false
}
