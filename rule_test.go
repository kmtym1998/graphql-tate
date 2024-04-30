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

var ctx = context.Background()
var ruleWithNoErr RuleFunc = func(_ context.Context, _ ast.ArgumentList, _ interface{}) error {
	return nil
}
var ruleWithErr RuleFunc = func(_ context.Context, _ ast.ArgumentList, _ interface{}) error {
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
			err := OR(tc.rules...)(ctx, nil, nil)

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
			err := AND(tc.rules...)(ctx, nil, nil)

			if err != nil && (tc.expected == nil || err.Error() != tc.expected.Error()) {
				t.Errorf("expected: `%#v`, got: `%#v`", tc.expected, err)
			}
		})
	}
}
