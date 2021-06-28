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
func (*Resolver) Mutation() generated.MutationResolver { return &mutationResolver{} }

func (*Resolver) Query() generated.QueryResolver {
	return &queryResolver{ }
}
type mutationResolver struct{
	*Resolver
	*mutation_resolver.MutationResolver
}

type queryResolver struct{
	*Resolver
	query_resolver.QueryResolver
}

func (*Resolver) RepositoryItem () generated.RepositoryItemResolver {
	return query_resolver.RepositoryItemResolver{}
}

func (*Resolver) SpecificationItem() generated.SpecificationItemResolver {
	return query_resolver.SpecificationItemResolver{}
}
func (*Resolver)DeviceItem() generated.DeviceItemResolver {
	return query_resolver.DeviceItemResolver{}
}

func (*Resolver) UserItem() generated.UserItemResolver {
	return query_resolver.UserItemResolver{}
}
func (*Resolver)CompanyItem() generated.CompanyItemResolver {
	return query_resolver.CompanyItemResolver{}
}

func (*Resolver)SteelItem() generated.SteelItemResolver {
	return query_resolver.SteelItemResolver{}
}
func (*Resolver)ProjectItem() generated.ProjectItemResolver {
	return query_resolver.ProjectItemResolver{}
}

func (*Resolver)OrderItem() generated.OrderItemResolver {
	return query_resolver.OrderItemResolver{}
}
func (*Resolver)OrderExpressItem() generated.OrderExpressItemResolver {
	return query_resolver.OrderExpressItemResolver{}
}

func (*Resolver) OrderSpecificationItem() generated.OrderSpecificationItemResolver {
	return query_resolver.OrderSpecificationItemResolver{}
}
func (*Resolver)MaintenanceRecordItem() generated.MaintenanceRecordItemResolver {
	return query_resolver.MaintenanceRecordItemResolver{}
}
func(*Resolver) SteelInProject() generated.SteelInProjectResolver {
	return query_resolver.SteelInProjectResolver{}
}
func (*Resolver) OrderSpecificationSteelItem() generated.OrderSpecificationSteelItemResolver{
	return query_resolver.OrderSpecificationSteelItemResolver{}
}

/**
 * graphQL 解析器入口
 */
func Handler() *handler.Server {
	conf := generated.Config{Resolvers: &Resolver{}}
	// 加入角色指令解析器
	conf.Directives.HasRole = directives.HasRole
	// 必须是设备
	conf.Directives.MustBeDevice = directives.MustBeDevice
	conf.Directives.HasRole = directives.HasRole

	return handler.NewDefaultServer( generated.NewExecutableSchema(conf))
}
