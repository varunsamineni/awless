package params

import (
	"reflect"
	"strings"
	"testing"
)

func TestRuleOptionals(t *testing.T) {
	tcases := []struct {
		rules     Rule
		optionals []string
	}{
		{rules: Opt("1", "2", "3"), optionals: []string{"1", "2", "3"}},

		{rules: AllOf(Key("1"), Opt("2", "3")), optionals: []string{"2", "3"}},
		{rules: OnlyOneOf(Key("1"), Opt("2", "3")), optionals: []string{"2", "3"}},
		{rules: AtLeastOneOf(Key("1"), Opt("2", "3")), optionals: []string{"2", "3"}},
	}

	for _, tcase := range tcases {
		if got, want := tcase.rules.Optionals(), tcase.optionals; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestRuleMissing(t *testing.T) {
	tcases := []struct {
		rules   Rule
		in      []string
		missing []string
	}{
		{rules: AllOf()},
		{rules: OnlyOneOf()},
		{rules: AtLeastOneOf()},

		{rules: AllOf(Key("1")), missing: []string{"1"}},
		{rules: OnlyOneOf(Key("1")), missing: []string{"1"}},
		{rules: AtLeastOneOf(Key("1")), missing: []string{"1"}},

		{rules: AllOf(Key("2"), Key("1")), missing: []string{"2", "1"}},
		{rules: OnlyOneOf(Key("2"), Key("1")), missing: []string{"2"}},
		{rules: AtLeastOneOf(Key("2"), Key("1")), missing: []string{"2"}},

		{rules: AllOf(Key("1"), Key("2")), in: []string{"1", "2"}},
		{rules: AllOf(Key("1"), Key("2")), in: []string{"1"}, missing: []string{"2"}},
		{rules: OnlyOneOf(Key("1"), Key("2")), in: []string{"1"}},
		{rules: AtLeastOneOf(Key("2"), Key("1")), in: []string{"2"}},

		{rules: OnlyOneOf(
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), missing: []string{"5"}},
		{rules: OnlyOneOf(
			OnlyOneOf(Key("1"), Key("2")),
			Key("5"),
			AtLeastOneOf(Key("3"), Key("4"))), missing: []string{"1"}},
		{rules: OnlyOneOf(
			OnlyOneOf(Key("1"), Key("2")),
			Key("5")), missing: []string{"1"}},
		{rules: OnlyOneOf(
			AtLeastOneOf(Key("3"), Key("4")),
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), in: []string{"3"}},

		{rules: AllOf(
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), missing: []string{"5", "1", "3"}},
		{rules: AllOf(
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), in: []string{"5"}, missing: []string{"1", "3"}},
		{rules: AllOf(
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), in: []string{"5", "1"}, missing: []string{"3"}},
		{rules: AllOf(
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), in: []string{"5", "3"}, missing: []string{"1"}},
	}

	for i, tcase := range tcases {
		if got, want := tcase.rules.Missing(tcase.in), tcase.missing; !reflect.DeepEqual(got, want) {
			t.Fatalf("missing: %d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestValidateRule(t *testing.T) {
	tcases := []struct {
		rules       Rule
		in          []string
		expectErr   bool
		errContains []string
	}{
		{rules: AllOf()},
		{rules: OnlyOneOf()},
		{rules: AtLeastOneOf()},

		{rules: AllOf(Key("1")), in: []string{"1", "2"}},
		{rules: OnlyOneOf(Key("1")), in: []string{"1"}},
		{rules: AtLeastOneOf(Key("1")), in: []string{"1"}},

		{rules: AllOf(Key("1"), Key("2")), in: []string{"1", "2"}},
		{rules: OnlyOneOf(Key("1"), Key("2")), in: []string{"1"}},
		{rules: OnlyOneOf(Key("1"), Key("2")), in: []string{"2"}},
		{rules: AtLeastOneOf(Key("1"), Key("2")), in: []string{"1"}},
		{rules: AtLeastOneOf(Key("1"), Key("2")), in: []string{"2"}},

		{rules: AllOf(Key("1"), Key("2")), in: []string{"1"}, expectErr: true, errContains: []string{"2"}},
		{rules: AllOf(Key("1")), in: []string{"2"}, expectErr: true, errContains: []string{"1"}},
		{rules: OnlyOneOf(Key("1")), in: []string{"2"}, expectErr: true, errContains: []string{"1"}},
		{rules: OnlyOneOf(Key("1"), Key("2")), in: []string{"1", "2"}, expectErr: true, errContains: []string{"1", "2"}},
		{rules: AtLeastOneOf(Key("1"), Key("2")), in: []string{"0"}, expectErr: true, errContains: []string{"1", "2"}},
		{rules: AtLeastOneOf(Key("1")), in: []string{"2"}, expectErr: true, errContains: []string{"1"}},

		{rules: AllOf(
			OnlyOneOf(Key("1"), Key("2")),
			Key("3")),
			in: []string{"4"}, expectErr: true, errContains: []string{"1", "2", "3"}},

		{rules: AllOf(
			AtLeastOneOf(Key("1"), Key("2")),
			Key("3")),
			in: []string{"4"}, expectErr: true, errContains: []string{"1", "3"}},
	}

	for i, tcase := range tcases {
		err := tcase.rules.Validate(tcase.in)
		if !tcase.expectErr && err != nil {
			t.Fatalf("%d: expected no error got %s", i+1, err)
		}

		if tcase.expectErr {
			if err == nil {
				t.Fatalf("%d: expected error got none", i+1)
			}
			msg := err.Error()
			for _, e := range tcase.errContains {
				if !strings.Contains(msg, e) {
					t.Fatalf("%d: expected %s to contain %s", i+1, msg, e)
				}
			}
		}
	}
}

func TestRulesRequired(t *testing.T) {
	tcases := []struct {
		rules    Rule
		required []string
	}{
		{rules: AllOf()},
		{rules: OnlyOneOf()},
		{rules: AtLeastOneOf()},

		{rules: AllOf(Key("1")), required: []string{"1"}},
		{rules: OnlyOneOf(Key("1")), required: []string{"1"}},
		{rules: AtLeastOneOf(Key("1")), required: []string{"1"}},

		{rules: AllOf(Key("1"), Key("2")), required: []string{"1", "2"}},
		{rules: OnlyOneOf(Key("1"), Key("2")), required: []string{"1"}},
		{rules: AtLeastOneOf(Key("2"), Key("1")), required: []string{"2"}},

		{rules: AllOf(
			Key("5"),
			OnlyOneOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), required: []string{"5", "1", "3"}},
		{rules: OnlyOneOf(
			AllOf(Key("1"), Key("2")),
			AtLeastOneOf(Key("3"), Key("4"))), required: []string{"1", "2"}},
		{rules: OnlyOneOf(
			AtLeastOneOf(Key("3"), Key("4")),
			AllOf(Key("1"), Key("2"))), required: []string{"3"}},
	}

	for _, tcase := range tcases {
		if got, want := tcase.rules.Required(), tcase.required; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
