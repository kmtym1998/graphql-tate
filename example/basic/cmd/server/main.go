package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
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

	mux := http.NewServeMux()

	mux.Handle("POST /v1/graphql", v1postGraphQLHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info(
		"Server is running",
		slog.String("port", port),
		"address", "http://localhost:"+port+"/v1/graphql",
	)
	if err := http.ListenAndServe(
		":"+port,
		mux,
	); err != nil {
		panic(err)
	}
}

func v1postGraphQLHandler() http.HandlerFunc {
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

	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		operationCtx := graphql.GetOperationContext(ctx)

		slog.Info(
			"OperationContext",
			"operationCtx.Doc.Operations", fmt.Sprintf("%#v", operationCtx.Doc.Operations),
		)

		var rootFieldNames []string
		for _, sel := range operationCtx.Doc.Operations[0].SelectionSet {
			slog.Info(
				"SelectionSet",
				"sel", fmt.Sprintf("%#v", sel),
			)

			f, ok := sel.(*ast.Field)
			if !ok {
				continue
			}

			rootFieldNames = append(rootFieldNames, f.Name)
		}

		slog.Info(
			"Responses",
			"operationCtx", fmt.Sprintf("%#v", operationCtx),
		)

		return next(ctx)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	})
}
