//go:generate go run github.com/99designs/gqlgen
package graphql

import (
	"http-api/app/http/graphql/model"
)

type Resolver struct{
	todos []*model.Todo
}
