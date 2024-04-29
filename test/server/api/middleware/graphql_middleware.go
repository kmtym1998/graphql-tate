package middleware

import (
	"context"
	"fmt"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/vektah/gqlparser/v2/ast"
)

var isAnonymous tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "anonymous" {
		return nil
	}

	return fmt.Errorf("role is not anonymous")
}

var isAdmin tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "admin" {
		return nil
	}

	return fmt.Errorf("role is not admin")
}

var isEditor tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "editor" {
		return nil
	}

	return fmt.Errorf("role is not editor")
}

var isViewer tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "viewer" {
		return nil
	}

	return fmt.Errorf("role is not viewer")
}

func NewPermission() tate.RootFieldPermission {
	return tate.RootFieldPermission{
		ast.Query: tate.ChildFieldPermission{
			"user": tate.ChildFieldPermission{
				"id":   tate.OR(isEditor, isAdmin),
				"name": tate.OR(isViewer, isEditor, isAdmin),
			},
			"todos": tate.OR(isViewer, isEditor, isAdmin),
		},
		ast.Mutation: tate.ChildFieldPermission{
			"createTodo": isAnonymous,
		},
	}
}
