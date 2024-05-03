package middleware

import (
	"context"
	"fmt"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/vektah/gqlparser/v2/ast"
)

var IsAnonymous tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := RoleFrom(ctx)
	if roleName == "anonymous" {
		return nil
	}

	return fmt.Errorf("role is not anonymous")
}

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

func NewPermission() tate.RootFieldPermission {
	return tate.RootFieldPermission{
		ast.Query: tate.ChildFieldPermission{
			"user": tate.ChildFieldPermission{
				"id":   tate.OR(IsEditor, IsAdmin),
				"name": tate.OR(IsViewer, IsEditor, IsAdmin),
			},
			"todos": tate.OR(IsViewer, IsEditor, IsAdmin),
		},
		ast.Mutation: tate.ChildFieldPermission{
			"createTodo": IsAnonymous,
		},
	}
}
