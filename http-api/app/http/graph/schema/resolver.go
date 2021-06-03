//go:generate go run github.com/99designs/gqlgen
package schema

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"http-api/app/http/graph/directives"
	"http-api/app/http/graph/generated"
	"http-api/app/http/graph/model"
)

type Resolver struct{
	todos []*model.Todo
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
