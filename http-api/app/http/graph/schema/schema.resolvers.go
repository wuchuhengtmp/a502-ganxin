package schema

// This files will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"http-api/app/http/graph/directives"
	"http-api/app/http/graph/generated"
	"http-api/app/http/graph/schema/mutation_resolver"
	"http-api/app/http/graph/schema/query_resolver"
)

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{ }
}

func (r *Resolver) RepositoryItem () generated.RepositoryItemResolver {
	return query_resolver.RepositoryItemResolver{}
}

func (r *Resolver) SpecificationItem() generated.SpecificationItemResolver {
	return query_resolver.SpecificationItemResolver{}
}
func (r *Resolver)DeviceItem() generated.DeviceItemResolver {
	return query_resolver.DeviceItemResolver{}
}

func (r *Resolver) UserItem() generated.UserItemResolver {
	return query_resolver.UserItemResolver{}
}
func (c *Resolver)CompanyItem() generated.CompanyItemResolver {
	return query_resolver.CompanyItemResolver{}
}

func (c *Resolver)SteelItem() generated.SteelItemResolver {
	return query_resolver.SteelItemResolver{}
}
func (r *Resolver)ProjectItem() generated.ProjectItemResolver {
	return query_resolver.ProjectItemResolver{}
}


type mutationResolver struct{
	*Resolver
	*mutation_resolver.MutationResolver
}

type queryResolver struct{
	*Resolver
	query_resolver.QueryResolver
}

/**
 * graphQL 解析器入口
 */
func Handler() *handler.Server {
	conf := generated.Config{Resolvers: &Resolver{}}
	// 加入角色指令解析器
	conf.Directives.HasRole = directives.HasRole
	return handler.NewDefaultServer( generated.NewExecutableSchema(conf))
}
