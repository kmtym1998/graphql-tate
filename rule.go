package tate

import (
	"context"
	"errors"
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

// Permitter is an interface that all permission types must implement
// It can be RuleFunc or ChildFieldPermission
type Permitter interface {
	isPermitter()
}

type RootFieldPermission map[ast.Operation]Permitter

// ChildFieldPermission is a map of field name to RuleFunc
type ChildFieldPermission map[string]Permitter

var _ Permitter = ChildFieldPermission{}

func (ChildFieldPermission) isPermitter() {}

// RuleFunc is a function that checks if the user has permission to access the field
// It returns nil if the user has permission, otherwise it returns an error
type RuleFunc func(ctx context.Context, args ast.ArgumentList, variable interface{}) error

var _ Permitter = RuleFunc(nil)

func (RuleFunc) isPermitter() {}

func (p ChildFieldPermission) validate() error {
	for _, permitter := range p {
		switch v := permitter.(type) {
		case RuleFunc:
			if v == nil {
				return errors.New("invalid permitter. RuleFunc cannot be nil")
			}
		case ChildFieldPermission:
			if err := v.validate(); err != nil {
				return fmt.Errorf("invalid permitter: %w", err)
			}
		default:
			return errors.New("invalid permitter")
		}
	}

	return nil
}

func OR(
	rules ...RuleFunc,
) RuleFunc {
	return func(ctx context.Context, args ast.ArgumentList, variable interface{}) error {
		errs := make([]error, 0, len(rules))
		for _, rule := range rules {
			if err := rule(ctx, args, variable); err != nil {
				errs = append(errs, err)

				continue
			}

			return nil
		}

		return errors.Join(errs...)
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

func NOT(
	rule RuleFunc,
	msg string,
) RuleFunc {
	return func(ctx context.Context, args ast.ArgumentList, variable interface{}) error {
		if rule == nil {
			return nil
		}

		if err := rule(ctx, args, variable); err != nil {
			return nil
		}

		return fmt.Errorf(msg)
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
