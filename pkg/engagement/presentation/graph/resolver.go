package graph

import (
	"context"
	"fmt"

	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/interactor"
	"github.com/savannahghi/firebasetools"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/savannahghi/engagement-service/cmd/generator

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	interactor *interactor.Interactor
}

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	interactor *interactor.Interactor,
) (*Resolver, error) {
	return &Resolver{
		interactor: interactor,
	}, nil
}

func (r Resolver) getLoggedInUserUID(ctx context.Context) (string, error) {
	authToken, err := firebasetools.GetUserTokenFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("auth token not found in context: %w", err)
	}
	return authToken.UID, nil
}
