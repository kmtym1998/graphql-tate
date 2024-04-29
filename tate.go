package tate

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Tate struct {
	messageBuilder func(ctx context.Context, fieldName string) string
	permission     PermissionDef
}

func NewTate(perm PermissionDef) (*Tate, error) {
	if err := perm.validate(); err != nil {
		return nil, fmt.Errorf("invalid permission: %w", err)
	}

	return &Tate{
		messageBuilder: func(ctx context.Context, fieldName string) string {
			return fmt.Sprintf("permission denied for %s", fieldName)
		},
		permission: perm,
	}, nil
}

func (t *Tate) SetErrorMessageBuilder(f func(ctx context.Context, fieldName string) string) {
	t.messageBuilder = f
}

func (t *Tate) AroundResponses(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	operationCtx := graphql.GetOperationContext(ctx)
	fieldCtx := graphql.GetFieldContext(ctx)

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

	operationName := operationCtx.Operation.Operation
	variables := operationCtx.Variables

	var graphqlErrors gqlerror.List
	for _, fieldName := range fieldNames {
		rule, ok := t.permission[operationName][fieldName]
		if !ok {
			graphqlErrors = append(graphqlErrors, &gqlerror.Error{
				Message: t.messageBuilder(ctx, fieldName),
			})
		}

		if err := rule(ctx, args, variables); err != nil {
			graphqlErrors = append(graphqlErrors, &gqlerror.Error{
				Message: fmt.Sprintf("%s: %s", t.messageBuilder(ctx, fieldName), err.Error()),
			})
		}
	}

	if len(graphqlErrors) > 0 {
		return &graphql.Response{
			Errors: graphqlErrors,
			Path:   fieldCtx.Path(),
		}
	}

	return next(ctx)
}
