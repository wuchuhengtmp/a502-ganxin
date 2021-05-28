package schema

// This files will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"http-api/app/http/graph/generated"
	"http-api/app/http/graph/model"
	mutation_resolver "http-api/app/http/graph/schema/mutation-resolver"
)
//type LoginMutationResolver =


func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.todos, nil
}

func (r *queryResolver) Hello(ctx context.Context) (*model.User, error) {
	return &model.User {
		ID: "hello",
		Name: "1123",
	}, nil
}

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{
		ID:   obj.UserID,
		Name: fmt.Sprintf("user %s", obj.UserID),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{
	*Resolver
	*mutation_resolver.MutationResolver
}

type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
