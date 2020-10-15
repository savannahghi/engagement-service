//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"gitlab.slade360emr.com/go/feed/graph/feed"
)

// NewResolver sets up the dependencies needed for WhatsApp query and mutation resolvers to work
func NewResolver() *Resolver {
	return &Resolver{
		feedService: feed.NewService(),
	}
}

// Resolver serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	feedService *feed.Service
}
