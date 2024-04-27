package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/kmtym1998/graphql-tate/example/generated"
	"github.com/kmtym1998/graphql-tate/example/model"
	"github.com/kmtym1998/graphql-tate/example/resolver"
	"github.com/lmittmann/tint"
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
	es := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			TodoList: []*model.Todo{
				{ID: "1", Text: "todo1", Done: false},
				{ID: "2", Text: "todo2", Done: true},
			},
		},
	})

	srv := handler.NewDefaultServer(es)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	})
}
