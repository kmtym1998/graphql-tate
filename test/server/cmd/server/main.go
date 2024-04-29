package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/example/api/middleware"
	"github.com/kmtym1998/graphql-tate/example/generated"
	"github.com/kmtym1998/graphql-tate/example/model"
	"github.com/kmtym1998/graphql-tate/example/resolver"
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
	r.Post("/v1/graphql", v1postGraphQLHandler(tate))

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

func v1postGraphQLHandler(tate *tate.Tate) http.HandlerFunc {
	user1 := &model.User{ID: "U1", Name: "user1"}
	user2 := &model.User{ID: "U2", Name: "user2"}
	todo1 := &model.Todo{ID: "1", Text: "todo1", Done: false, User: user1}
	todo2 := &model.Todo{ID: "2", Text: "todo2", Done: true, User: user2}
	todo3 := &model.Todo{ID: "3", Text: "todo2", Done: true, User: user2}
	user1.Todos = []*model.Todo{todo1}
	user2.Todos = []*model.Todo{todo2, todo3}

	es := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			UserList: []*model.User{
				user1,
				user2,
			},
			TodoList: []*model.Todo{
				todo1,
				todo2,
				todo3,
			},
		},
	})

	srv := handler.NewDefaultServer(es)

	srv.AroundFields(tate.AroundFields)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	})
}
