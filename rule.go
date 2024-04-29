package tate

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/vektah/gqlparser/v2/ast"
)

type RuleFunc func(ctx context.Context, args ast.ArgumentList, variable interface{}) error

type PermissionDef map[ast.Operation]map[string]RuleFunc

func (p PermissionDef) validate() error {
	operations := []ast.Operation{ast.Query, ast.Mutation, ast.Subscription}
	for op, rules := range p {
		if !slices.Contains(operations, op) {
			return errors.New("invalid operation")
		}

		for fieldName, rule := range rules {
			if rule == nil {
				return fmt.Errorf("RuleFunc for %s is nil", fieldName)
			}
		}
	}

	return nil
}

func OR(
	rules ...RuleFunc,
) RuleFunc {
	return func(ctx context.Context, args ast.ArgumentList, variable interface{}) error {
		for _, rule := range rules {
			if err := rule(ctx, args, variable); err != nil {
				return err
			}
		}

		return nil
	}
}

func AND(
	rules ...RuleFunc,
) RuleFunc {
	return func(ctx context.Context, args ast.ArgumentList, variable interface{}) error {
		for _, rule := range rules {
			if err := rule(ctx, args, variable); err != nil {
				return err
			}
		}

		return nil
	}
}

func Any() RuleFunc {
	return func(ctx context.Context, args ast.ArgumentList, variable interface{}) error {
		return nil
	}
}

func None() RuleFunc {
	return func(ctx context.Context, args ast.ArgumentList, variable interface{}) error {
		return fmt.Errorf("no permission allowed for this field")
	}
}
