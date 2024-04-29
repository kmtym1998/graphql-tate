package tate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
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

	slog.Debug(
		"checking permission",
		"fieldName", fieldName,
		"operationName", operationName,
		"variables", variables,
	)

	fieldNames := getCurrentAndParentFieldNames(ctx)

	slog.Info(
		"fields",
		"fieldNames", fieldNames,
		"len(fieldNames)", len(fieldNames),
	)

	return next(ctx)
}

func getCurrentAndParentFieldNames(ctx context.Context) []string {
	fieldCtx := graphql.GetFieldContext(ctx)

	fieldNames := []string{fieldCtx.Field.Name}
	fieldCtxItr := fieldCtx
	for {
		fieldCtxItr = fieldCtxItr.Parent

		if fieldCtxItr == nil {
			break
		}

		var parentFieldName string
		if fieldCtxItr.Field.Field != nil {
			parentFieldName = fieldCtxItr.Field.Field.Name
		} else if fieldCtxItr.Parent != nil {
			parentFieldName = fieldCtxItr.Parent.Field.Name
			fieldCtxItr = fieldCtxItr.Parent
		} else {
			continue
		}

		fieldNames = append(fieldNames, parentFieldName)
	}

	return reverse(fieldNames)
}

func reverse[T any](s []T) []T {
	n := len(s)
	for i := 0; i < n/2; i++ {
		s[i], s[n-1-i] = s[n-1-i], s[i]
	}
	return s
}
