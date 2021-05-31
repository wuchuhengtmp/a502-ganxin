package schema

// This files will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"http-api/app/http/graph/generated"
	"http-api/app/http/graph/schema/mutation_resolver"
	"http-api/app/http/graph/schema/query_resolver"
)

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{} }

type mutationResolver struct{
	*Resolver
	*mutation_resolver.MutationResolver
}

type queryResolver struct{
	*Resolver
	query_resolver.QueryResolver
}
