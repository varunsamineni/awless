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
		if err != nil {
			return fmt.Errorf("all of %v", n.rules)
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
	return fmt.Sprintf("all of %v", n.rules)
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
		return fmt.Errorf("only one of %v", n.rules)
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
	return fmt.Sprintf("only one of %v", n.rules)
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
		if err := r.Validate(input); err == nil {
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
	return fmt.Sprintf("at least one of %v", n.rules)
}

type opt struct {
	optionals []string
}

func Opt(s ...string) Rule {
	o := opt{}
	o.optionals = append(o.optionals, s...)
	return o
}

func (n opt) Validate(input []string) error {
	return nil
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
	return fmt.Sprintf("optionals: %s", strings.Join(n.optionals, ","))
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

func build(rules []Rule) (d defaultRule) {
	for _, n := range rules {
		d.rules = append(d.rules, n)
	}
	return
}

type defaultRule struct {
	rules []Rule
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
