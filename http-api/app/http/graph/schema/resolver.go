//go:generate go run github.com/99designs/gqlgen
package schema

import (
	"http-api/app/http/graph/model"
)

type Resolver struct{
	todos []*model.Todo
}
