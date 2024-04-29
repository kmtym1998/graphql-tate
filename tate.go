package tate

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Tate struct {
	messageBuilder    func(ctx context.Context, fieldName string) string
	extensionsBuilder func(ctx context.Context, fieldName string) map[string]interface{}
	permission        RootFieldPermission
}

func New(perm RootFieldPermission) (*Tate, error) {
	for op, permitter := range perm {
		childFieldPermission, ok := permitter.(ChildFieldPermission)
		if !ok {
			continue
		}

		if err := childFieldPermission.validate(); err != nil {
			return nil, fmt.Errorf("invalid permission for %s: %w", op, err)
		}
	}

	return &Tate{
		messageBuilder: func(ctx context.Context, fieldName string) string {
			return fmt.Sprintf("permission denied for %s", fieldName)
		},
		extensionsBuilder: func(ctx context.Context, fieldName string) map[string]interface{} {
			return map[string]interface{}{
				"fieldName": fieldName,
			}
		},
		permission: perm,
	}, nil
}

func (t *Tate) SetErrorMessageBuilder(f func(ctx context.Context, fieldName string) string) {
	t.messageBuilder = f
}

func (t *Tate) AroundFields(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	operationCtx := graphql.GetOperationContext(ctx)
	fieldCtx := graphql.GetFieldContext(ctx)

	fieldName := fieldCtx.Field.Name
	operationName := operationCtx.Operation.Operation
	variables := operationCtx.Variables

	// TODO: fetch fields recursively
	rule, ok := t.permission[operationName][fieldName]
	if !ok {
		return next(ctx)
	}

	if err := rule(ctx, fieldCtx.Field.Arguments, variables); err != nil {
		return nil, &gqlerror.Error{
			Message: fmt.Sprintf("%s: %s", t.messageBuilder(ctx, fieldName), err.Error()),
			Path:    fieldCtx.Path(),
			Locations: []gqlerror.Location{
				{
					Line:   fieldCtx.Field.Position.Line,
					Column: fieldCtx.Field.Position.Column,
				},
			},
			Extensions: t.extensionsBuilder(ctx, fieldName),
		}
	}

	return next(ctx)
}
