package tate

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

type Tate struct {
	permission PermissionDef
}

func NewTate(perm PermissionDef) (*Tate, error) {
	if err := perm.validate(); err != nil {
		return nil, fmt.Errorf("invalid permission: %w", err)
	}

	return &Tate{
		permission: perm,
	}, nil
}

// TODO: ほんとうにエラーを返すで良いかどうかは検討が必要
func (t *Tate) Check(ctx context.Context) error {
	operationCtx := graphql.GetOperationContext(ctx)
	if operationCtx == nil {
		return fmt.Errorf("operation context not found")
	}

	op := operationCtx.Operation.Operation
	variables := operationCtx.Variables

	var fieldNames []string
	var args ast.ArgumentList
	for _, sel := range operationCtx.Operation.SelectionSet {
		f, ok := sel.(*ast.Field)
		if !ok {
			continue
		}

		fieldNames = append(fieldNames, f.Name)
		args = f.Arguments
	}

	for _, fieldName := range fieldNames {
		rule, ok := t.permission[op][fieldName]
		if !ok {
			return fmt.Errorf("permission not found for %s.%s", op, fieldName)
		}

		if err := rule(ctx, args, variables); err != nil {
			return fmt.Errorf("permission denied for %s.%s", op, fieldName)
		}
	}

	return nil
}
