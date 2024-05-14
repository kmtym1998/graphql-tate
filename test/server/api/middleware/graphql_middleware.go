package middleware

import (
	"context"
	"fmt"
	"strconv"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/vektah/gqlparser/v2/ast"
)

var IsAdmin tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "admin" {
		return nil
	}

	return fmt.Errorf("role is not admin")
}

var IsEditor tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "editor" {
		return nil
	}

	return fmt.Errorf("role is not editor")
}

var IsViewer tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "viewer" {
		return nil
	}

	return fmt.Errorf("role is not viewer")
}

var IsAnonymous tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "anonymous" {
		return nil
	}

	return fmt.Errorf("role is not anonymous")
}

var OnlyAnonymousMustHaveLimit tate.RuleFunc = func(ctx context.Context, args ast.ArgumentList, vars interface{}) error {
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

		limitValOrVarName, err := strconv.Atoi(arg.Value.Raw)
		if err == nil {
			if limitValOrVarName > 50 {
				return fmt.Errorf("limit is too large (from arg)")
			}
		}

		varName = arg.Value.Raw

		break
	}

	varsMap, ok := vars.(map[string]interface{})
	if !ok {
		return nil
	}

	if limitVar, ok := varsMap[varName]; ok {
		limitVal, ok := limitVar.(int64)
		if ok {
			if limitVal > 50 {
				return fmt.Errorf("limit is too large (from variable)")
			}

			return nil
		}
	}

	return fmt.Errorf("limit is not set")
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
