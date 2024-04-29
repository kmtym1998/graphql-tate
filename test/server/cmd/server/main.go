package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/example/api/handler"
	"github.com/kmtym1998/graphql-tate/example/api/middleware"
	"github.com/lmittmann/tint"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(
		os.Stdout,
		&tint.Options{
			AddSource: true,
		},
	)))

	r := chi.NewRouter()
	r.Use(middleware.InjectRole)

	tate, err := tate.New(permission)
	if err != nil {
		panic(err)
	}
	r.Post("/v1/graphql", handler.PostV1GraphQLHandler(tate))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info(
		"Server is running",
		slog.String("port", port),
		"address", "http://localhost:"+port+"/v1/graphql",
	)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

var isAnonymous tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := middleware.RoleFrom(ctx)
	if roleName == "anonymous" {
		return nil
	}

	return fmt.Errorf("role is not anonymous")
}

var isAdmin tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := middleware.RoleFrom(ctx)
	if roleName == "admin" {
		return nil
	}

	return fmt.Errorf("role is not admin")
}

var isEditor tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := middleware.RoleFrom(ctx)
	if roleName == "editor" {
		return nil
	}

	return fmt.Errorf("role is not editor")
}

var isViewer tate.RuleFunc = func(ctx context.Context, _ ast.ArgumentList, _ interface{}) error {
	roleName := middleware.RoleFrom(ctx)
	if roleName == "viewer" {
		return nil
	}

	return fmt.Errorf("role is not viewer")
}

var permission = tate.RootFieldPermission{
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
