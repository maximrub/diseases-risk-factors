package graph

import "github.com/maximrub/thesis-diseases-risk-factors-server/internal/dal"

//go:generate go run gen/generator.go
//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DAL *dal.DAL
}
