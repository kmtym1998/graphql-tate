package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/example/generated"
	"github.com/kmtym1998/graphql-tate/example/model"
	"github.com/kmtym1998/graphql-tate/example/resolver"
)

func PostV1GraphQLHandler(tate *tate.Tate) http.HandlerFunc {
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
