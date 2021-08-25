package graph

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/firebasetools"

	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/interactor"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen

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

func (r Resolver) checkPreconditions() {
	if r.interactor.Feed == nil {
		log.Panicf("expected feed usecases to be defined resolver")
	}

	if r.interactor.Notification == nil {
		log.Panicf("expected notification usecases to be define in resolver ")
	}

	if r.interactor.Uploads == nil {
		log.Panicf("expected uploads usecases to be define in resolver ")
	}

	if r.interactor.Mail == nil {
		log.Panicf("expected mail usecases to be define in resolver ")
	}
}

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
