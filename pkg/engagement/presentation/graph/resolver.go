package graph

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/firebasetools"

	"github.com/savannahghi/engagement-service/pkg/engagement/usecases"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	usecases usecases.Usecases
}

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	usecases usecases.Usecases,
) (*Resolver, error) {
	return &Resolver{
		usecases: usecases,
	}, nil
}

func (r Resolver) checkPreconditions() {}

func (r Resolver) getLoggedInUserUID(ctx context.Context) (string, error) {
	authToken, err := firebasetools.GetUserTokenFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("auth token not found in context: %w", err)
	}
	return authToken.UID, nil
}

// CheckUserTokenInContext ensures that the context has a valid Firebase auth token
func (r *Resolver) CheckUserTokenInContext(ctx context.Context) *auth.Token {
	token, err := firebasetools.GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("graph.Resolver: context user token is nil")
	}
	return token
}
