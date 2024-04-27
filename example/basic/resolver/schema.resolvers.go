package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/kmtym1998/graphql-tate/example/generated"
	"github.com/kmtym1998/graphql-tate/example/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	r.TodoList = append(r.TodoList, &model.Todo{
		ID:   fmt.Sprintf("T%d", len(r.TodoList)+1),
		Text: input.Text,
		Done: false,
	})

	return r.TodoList[len(r.TodoList)-1], nil
}

// UpdateTodoDone is the resolver for the updateTodoDone field.
func (r *mutationResolver) UpdateTodoDone(ctx context.Context, id string, done bool) (*model.Todo, error) {
	for _, todo := range r.TodoList {
		if todo.ID == id {
			todo.Done = done
			return todo, nil
		}
	}

	return nil, fmt.Errorf("Todo not found")
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, name string) (*model.User, error) {
	if name == "" {
		return nil, fmt.Errorf("Name is empty")
	}

	r.UserList = append(r.UserList, &model.User{
		ID:   fmt.Sprintf("U%d", len(r.UserList)+1),
		Name: name,
	})

	return r.UserList[len(r.UserList)-1], nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.TodoList, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.UserList, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	for _, user := range r.UserList {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User not found")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
