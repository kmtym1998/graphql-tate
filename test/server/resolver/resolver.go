package resolver

import "github.com/kmtym1998/graphql-tate/test/server/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoList []*model.Todo
	UserList []*model.User
}
