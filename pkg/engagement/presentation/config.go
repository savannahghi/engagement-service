package presentation

import (
	"context"

	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagement-service/pkg/engagement/usecases"
	engLibPresentation "github.com/savannahghi/engagement/pkg/engagement/presentation"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph"
	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"

	"net/http"

	"github.com/gorilla/mux"
)

// AllowedOrigins is list of CORS origins allowed to interact with
// this service
var AllowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"https://a.bewell.co.ke",
	"https://b.bewell.co.ke",
	"http://localhost:5000",
	"https://europe-west3-bewell-app.cloudfunctions.net",
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	// Initialize new instances of the infrastructure services
	infrastructure := infrastructure.NewInfrastructureInteractor()
	usecases := usecases.NewUsecasesInteractor(infrastructure)

	r := mux.NewRouter() // gorilla mux
	engLibPresentation.SharedUnauthenticatedRoutes(ctx, r)

	// Authenticated routes
	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, usecases))

	engLibPresentation.SharedAuthenticatedISCRoutes(ctx, r)
	return r, nil
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	usecases usecases.Usecases,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, usecases)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
