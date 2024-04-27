package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/kmtym1998/graphql-tate/example/generated"
	"github.com/kmtym1998/graphql-tate/example/model"
	"github.com/kmtym1998/graphql-tate/example/resolver"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("POST /v1/graphql", v1postGraphQLHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
			Todos: []*model.Todo{
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
