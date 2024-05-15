package tate

import (
	"context"
	"errors"
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

type ruleFuncTestCase struct {
	name     string
	rules    []RuleFunc
	expected error
}

var ruleWithNoErr RuleFunc = func(_ context.Context, _ ast.ArgumentList, _ map[string]interface{}) error {
	return nil
}
var ruleWithErr RuleFunc = func(_ context.Context, _ ast.ArgumentList, _ map[string]interface{}) error {
	return errors.New("error")
}

func TestOR(t *testing.T) {
	for _, tc := range []ruleFuncTestCase{
		{
			name:     "no rules",
			rules:    nil,
			expected: nil,
		}, {
			name:     "one rule with no error",
			rules:    []RuleFunc{ruleWithNoErr},
			expected: nil,
		}, {
			name:     "one rule with error",
			rules:    []RuleFunc{ruleWithErr},
			expected: errors.New("error"),
		}, {
			name:     "two rules with no error",
			rules:    []RuleFunc{ruleWithNoErr, ruleWithNoErr},
			expected: nil,
		}, {
			name:     "two rules with error",
			rules:    []RuleFunc{ruleWithErr, ruleWithErr},
			expected: errors.Join([]error{errors.New("error"), errors.New("error")}...),
		}, {
			name:     "one rule with error, one rule with no error",
			rules:    []RuleFunc{ruleWithErr, ruleWithNoErr},
			expected: nil,
		}, {
			name:     "two rule with no error, one rule with error",
			rules:    []RuleFunc{ruleWithNoErr, ruleWithNoErr, ruleWithErr},
			expected: nil,
		}, {
			name:     "two rule with error, one rule with no error",
			rules:    []RuleFunc{ruleWithErr, ruleWithErr, ruleWithNoErr},
			expected: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := OR(tc.rules...)(nil, nil, nil)

			if err != nil && (tc.expected == nil || err.Error() != tc.expected.Error()) {
				t.Errorf("expected: `%#v`, got: `%#v`", tc.expected, err)
			}
		})
	}
}

func TestAND(t *testing.T) {
	for _, tc := range []ruleFuncTestCase{
		{
			name:     "no rules",
			rules:    nil,
			expected: nil,
		}, {
			name:     "one rule with no error",
			rules:    []RuleFunc{ruleWithNoErr},
			expected: nil,
		}, {
			name:     "one rule with error",
			rules:    []RuleFunc{ruleWithErr},
			expected: errors.New("error"),
		}, {
			name:     "two rules with no error",
			rules:    []RuleFunc{ruleWithNoErr, ruleWithNoErr},
			expected: nil,
		}, {
			name:     "two rules with error",
			rules:    []RuleFunc{ruleWithErr, ruleWithErr},
			expected: errors.New("error"),
		}, {
			name:     "one rule with error, one rule with no error",
			rules:    []RuleFunc{ruleWithErr, ruleWithNoErr},
			expected: errors.New("error"),
		}, {
			name:     "two rule with no error, one rule with error",
			rules:    []RuleFunc{ruleWithNoErr, ruleWithNoErr, ruleWithErr},
			expected: errors.New("error"),
		}, {
			name:     "two rule with error, one rule with no error",
			rules:    []RuleFunc{ruleWithErr, ruleWithErr, ruleWithNoErr},
			expected: errors.New("error"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := AND(tc.rules...)(nil, nil, nil)

			if err != nil && (tc.expected == nil || err.Error() != tc.expected.Error()) {
				t.Errorf("expected: `%#v`, got: `%#v`", tc.expected, err)
			}
		})
	}
}

func TestNOT(t *testing.T) {
	for _, tc := range []ruleFuncTestCase{
		{
			name:     "no rules",
			rules:    nil,
			expected: nil,
		}, {
			name:     "one rule with no error",
			rules:    []RuleFunc{ruleWithNoErr},
			expected: errors.New("error"),
		}, {
			name:     "one rule with error",
			rules:    []RuleFunc{ruleWithErr},
			expected: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := func() RuleFunc {
				if len(tc.rules) > 0 {
					return tc.rules[0]
				}

				return nil
			}()

			err := NOT(rule, "error")(nil, nil, nil)

			if err != nil && (tc.expected == nil || err.Error() != tc.expected.Error()) {
				t.Errorf("expected: `%#v`, got: `%#v`", tc.expected, err)
			}
		})
	}
}

type validateTestCase struct {
	name      string
	permitter ChildFieldPermission
	expected  error
}

func TestValidate(t *testing.T) {
	for _, tc := range []validateTestCase{
		{
			name:      "no rules",
			permitter: ChildFieldPermission{},
			expected:  nil,
		}, {
			name: "with field permission",
			permitter: ChildFieldPermission{
				"field1": Any(),
				"field2": Any(),
				"field3": None(),
			},
			expected: nil,
		}, {
			name: "with field permission but no rules",
			permitter: ChildFieldPermission{
				"field1": OR(),
				"field2": nil,
				"field3": AND(),
			},
			expected: errors.New("invalid permitter type: <nil> field: field2"),
		}, {
			name: "with field permission but invalid rule",
			permitter: ChildFieldPermission{
				"field1": RuleFunc(nil),
				"field2": Any(),
				"field3": AND(),
			},
			expected: errors.New("invalid permitter. RuleFunc cannot be nil for field: field1"),
		}, {
			name: "with child field permission",
			permitter: ChildFieldPermission{
				"field1": ChildFieldPermission{
					"field1-1": Any(),
					"field1-2": Any(),
					"field1-3": None(),
				},
				"field2": Any(),
				"field3": None(),
			},
			expected: nil,
		}, {
			name: "with child field permission but invalid rule in child field permission",
			permitter: ChildFieldPermission{
				"field1": ChildFieldPermission{
					"field1-1": RuleFunc(nil),
					"field1-2": Any(),
					"field1-3": AND(),
				},
				"field2": Any(),
				"field3": None(),
			},
			expected: errors.New("child field has invalid permitter field: field1: invalid permitter. RuleFunc cannot be nil for field: field1-1"),
		}, {
			name: "with child field permission but invalid rule in parent field permission",
			permitter: ChildFieldPermission{
				"field1": ChildFieldPermission{
					"field1-1": Any(),
					"field1-2": Any(),
					"field1-3": None(),
				},
				"field2": RuleFunc(nil),
				"field3": AND(),
			},
			expected: errors.New("invalid permitter. RuleFunc cannot be nil for field: field2"),
		}, {
			name: "with child field permission but invalid rule in grandchild field permission",
			permitter: ChildFieldPermission{
				"field1": ChildFieldPermission{
					"field1-1": Any(),
					"field1-2": Any(),
					"field1-3": AND(),
				},
				"field2": ChildFieldPermission{
					"field2-1": Any(),
					"field2-2": ChildFieldPermission{
						"field2-2-1": RuleFunc(nil),
						"field2-2-2": Any(),
						"field2-2-3": Any(),
					},
					"field2-3": Any(),
				},
				"field3": AND(),
			},
			expected: errors.New("child field has invalid permitter field: field2: child field has invalid permitter field: field2-2: invalid permitter. RuleFunc cannot be nil for field: field2-2-1"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.permitter.validate()

			if err != nil && (tc.expected == nil || err.Error() != tc.expected.Error()) {
				t.Errorf("expected: `%#v`, got: `%#v`", tc.expected, err)
			}
		})
	}
}
