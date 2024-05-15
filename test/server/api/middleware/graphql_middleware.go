package middleware

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/vektah/gqlparser/v2/ast"
)

var IsAdmin tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ map[string]interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "admin" {
		return nil
	}

	return fmt.Errorf("role is not admin")
}

var IsEditor tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ map[string]interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "editor" {
		return nil
	}

	return fmt.Errorf("role is not editor")
}

var IsViewer tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ map[string]interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "viewer" {
		return nil
	}

	return fmt.Errorf("role is not viewer")
}

var IsAnonymous tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ map[string]interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "anonymous" {
		return nil
	}

	return fmt.Errorf("role is not anonymous")
}

var OnlyAnonymousMustHaveLimit tate.RuleFunc = func(ctx context.Context, args ast.ArgumentList, variables map[string]interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "" {
		roleName = "anonymous"
	}

	if roleName != "anonymous" {
		return nil
	}

	var varName string
	for _, arg := range args {
		if arg.Name != "limit" {
			continue
		}

		if arg.Value.Kind == ast.Variable {
			if limitVar, ok := variables[varName]; ok {
				if limitVal, ok := limitVar.(int64); ok {
					if limitVal > 50 {
						return errors.New("limit is too large (from variable)")
					}

					return nil
				}

				return errors.New("limit is invalid (from variable)")
			}

			return errors.New("limit arg exists but variable is not set")
		} else {
			limitVal, err := strconv.Atoi(arg.Value.String())
			if err != nil {
				return errors.New("limit exists but is invalid (from arg)")
			}

			if limitVal > 50 {
				return errors.New("limit is too large (from variable)")
			}

			return nil
		}
	}

	return errors.New("limit not set")
}

func NewPermission() tate.RootFieldPermission {
	return tate.RootFieldPermission{
		ast.Query: tate.ChildFieldPermission{
			"user": tate.ChildFieldPermission{
				"id":   tate.OR(IsEditor, IsAdmin),
				"name": tate.OR(IsViewer, IsEditor, IsAdmin),
			},
			"todos": tate.OR(IsViewer, IsEditor, IsAdmin),
			"users": OnlyAnonymousMustHaveLimit,
		},
		ast.Mutation: tate.ChildFieldPermission{
			"createTodo": IsAnonymous,
		},
	}
}
