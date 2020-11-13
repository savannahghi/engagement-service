package graph

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"net/http"

	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	db "gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/database"
	"gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/messaging"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

const (
	mbBytes              = 1048576
	serverTimeoutSeconds = 120
	schemaDir            = "gitlab.slade360emr.com/go/feed:/graph/feed/schema/"
)

var allowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
	"https://feed-staging.healthcloud.co.ke",
	"https://feed-testing.healthcloud.co.ke",
	"https://feed-prod.healthcloud.co.ke",
}
var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type",
}
var errNotFound = fmt.Errorf("not found")

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w", base.GoogleCloudProjectIDEnvVarName, err)
	}

	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate notification service in resolver: %w", err)
	}

	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(base.RequestDebugMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	// static files
	schemaFileHandler, err := schemaHandler()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate schema file handler: %w", err)
	}
	r.PathPrefix("/schema/").Handler(schemaFileHandler)

	// Authenticated routes
	authR := r.Path("/graphql").Subrouter()
	authR.Use(base.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, fr, ns))

	// REST routes
	// Interservice Authenticated routes
	isc := r.PathPrefix("/feed/{uid}/{flavour}/").Subrouter()
	isc.Use(base.InterServiceAuthenticationMiddleware())

	// retrieval
	isc.Methods(
		http.MethodGet,
	).Path("/").HandlerFunc(
		GetFeed(ctx, fr, ns),
	).Name("getFeed")

	isc.Methods(
		http.MethodGet,
	).Path("/items/{itemID}/").HandlerFunc(
		GetFeedItem(ctx, fr, ns),
	).Name("getFeedItem")

	isc.Methods(
		http.MethodGet,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		GetNudge(ctx, fr, ns),
	).Name("getNudge")

	isc.Methods(
		http.MethodGet,
	).Path("/actions/{actionID}/").HandlerFunc(
		GetAction(ctx, fr, ns),
	).Name("getAction")

	// creation
	isc.Methods(
		http.MethodPost,
	).Path("/items/").HandlerFunc(
		PublishFeedItem(ctx, fr, ns),
	).Name("publishFeedItem")

	isc.Methods(
		http.MethodPost,
	).Path("/nudges/").HandlerFunc(
		PublishNudge(ctx, fr, ns),
	).Name("publishNudge")

	isc.Methods(
		http.MethodPost,
	).Path("/actions/").HandlerFunc(
		PublishAction(ctx, fr, ns),
	).Name("publishAction")

	isc.Methods(
		http.MethodPost,
	).Path("/{itemID}/messages/").HandlerFunc(
		PostMessage(ctx, fr, ns),
	).Name("postMessage")

	isc.Methods(
		http.MethodPost,
	).Path("/events/").HandlerFunc(
		ProcessEvent(ctx, fr, ns),
	).Name("postEvent")

	// deleting
	isc.Methods(
		http.MethodDelete,
	).Path("/items/{itemID}/").HandlerFunc(
		DeleteFeedItem(ctx, fr, ns),
	).Name("deleteFeedItem")

	isc.Methods(
		http.MethodDelete,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		DeleteNudge(ctx, fr, ns),
	).Name("deleteNudge")

	isc.Methods(
		http.MethodDelete,
	).Path("/actions/{actionID}/").HandlerFunc(
		DeleteAction(ctx, fr, ns),
	).Name("deleteAction")

	isc.Methods(
		http.MethodDelete,
	).Path("/{itemID}/messages/{messageID}/").HandlerFunc(
		DeleteMessage(ctx, fr, ns),
	).Name("deleteMessage")

	// modifying (patching)
	isc.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/resolve/").HandlerFunc(
		ResolveFeedItem(ctx, fr, ns),
	).Name("resolveFeedItem")

	isc.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unresolve/").HandlerFunc(
		UnresolveFeedItem(ctx, fr, ns),
	).Name("unresolveFeedItem")

	isc.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/pin/").HandlerFunc(
		PinFeedItem(ctx, fr, ns),
	).Name("pinFeedItem")

	isc.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unpin/").HandlerFunc(
		UnpinFeedItem(ctx, fr, ns),
	).Name("unpinFeedItem")

	isc.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/hide/").HandlerFunc(
		HideFeedItem(ctx, fr, ns),
	).Name("hideFeedItem")

	isc.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/show/").HandlerFunc(
		ShowFeedItem(ctx, fr, ns),
	).Name("showFeedItem")

	isc.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/resolve/").HandlerFunc(
		ResolveNudge(ctx, fr, ns),
	).Name("resolveNudge")

	isc.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/unresolve/").HandlerFunc(
		UnresolveNudge(ctx, fr, ns),
	).Name("unresolveNudge")

	isc.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/show/").HandlerFunc(
		ShowNudge(ctx, fr, ns),
	).Name("showNudge")

	isc.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/hide/").HandlerFunc(
		HideNudge(ctx, fr, ns),
	).Name("hideNudge")

	// return the combined router
	return r, nil
}

// PrepareServer starts up a server
func PrepareServer(ctx context.Context, port int) *http.Server {
	// start up the router
	r, err := Router(ctx)
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// start the server
	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)
	h = handlers.CORS(
		handlers.AllowedHeaders(allowedHeaders),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)
	h = handlers.CombinedLoggingHandler(os.Stdout, h)
	h = handlers.ContentTypeHandler(h, "application/json")
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: serverTimeoutSeconds * time.Second,
		ReadTimeout:  serverTimeoutSeconds * time.Second,
	}
	log.Infof("Server running at port %v", addr)
	return srv
}

//HealthStatusCheck endpoint to check if the server is working.
func HealthStatusCheck(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	resolver, err := NewResolver(ctx, fr, ns)
	if err != nil {
		base.LogStartupError(ctx, err)
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

// GetFeed retrieves and serves a feed
func GetFeed(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, flavour, err := getUIDAndFlavour(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		persistent, err := getRequiredBooleanFilterQueryParam(r, "persistent")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		status, err := getOptionalStatusQueryParam(r, "status")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		visibility, err := getOptionalVisibilityQueryParam(r, "visibility")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		expired, err := getOptionalBooleanFilterQueryParam(r, "expired")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		filterParams, err := getOptionalFilterParamsQueryParam(r, "filterParams")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		agg, err := feed.NewCollection(fr, ns)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		feed, err := agg.GetFeed(
			ctx,
			uid,
			flavour,
			persistent,
			status,
			visibility,
			expired,
			filterParams,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := feed.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetFeedItem retrieves a single feed item
func GetFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		item, err := thinFeed.GetFeedItem(ctx, itemID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if item == nil {
			respondWithError(w, http.StatusNotFound, errNotFound)
		}

		marshalled, err := item.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetNudge retrieves a single nudge
func GetNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nudgeID, err := getStringVar(r, "nudgeID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudge, err := thinFeed.GetNudge(ctx, nudgeID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if nudge == nil {
			respondWithError(w, http.StatusNotFound, errNotFound)
		}

		marshalled, err := nudge.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetAction retrieves a single action
func GetAction(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionID, err := getStringVar(r, "actionID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		action, err := thinFeed.GetAction(ctx, actionID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if action == nil {
			respondWithError(w, http.StatusNotFound, errNotFound)
		}

		marshalled, err := action.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

func readBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, mbBytes))
	if err != nil {
		return nil, fmt.Errorf("can't read request body: %w", err)
	}
	return body, nil
}

// PublishFeedItem posts a feed item
func PublishFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		item := &feed.Item{}
		err = item.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedItem, err := thinFeed.PublishFeedItem(ctx, item)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := publishedItem.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// DeleteFeedItem removes a feed item
func DeleteFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = thinFeed.DeleteFeedItem(ctx, itemID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// ResolveFeedItem marks a feed item as done
func ResolveFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchItem(ctx, fr, ns, thinFeed.ResolveFeedItem, w, r)
	}
}

// PinFeedItem marks a feed item as done
func PinFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchItem(ctx, fr, ns, thinFeed.PinFeedItem, w, r)
	}
}

// UnpinFeedItem marks a feed item as done
func UnpinFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchItem(ctx, fr, ns, thinFeed.UnpinFeedItem, w, r)
	}
}

// HideFeedItem marks a feed item as done
func HideFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchItem(ctx, fr, ns, thinFeed.HideFeedItem, w, r)
	}
}

// ShowFeedItem marks a feed item as done
func ShowFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchItem(ctx, fr, ns, thinFeed.ShowFeedItem, w, r)
	}
}

// UnresolveFeedItem marks a feed item as not resolved
func UnresolveFeedItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchItem(ctx, fr, ns, thinFeed.UnresolveFeedItem, w, r)
	}
}

// PublishNudge posts a new nudge
func PublishNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudge := &feed.Nudge{}
		err = nudge.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedNudge, err := thinFeed.PublishNudge(ctx, nudge)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := publishedNudge.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// ResolveNudge marks a nudge as resolved
func ResolveNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchNudge(ctx, fr, ns, thinFeed.ResolveNudge, w, r)
	}
}

// UnresolveNudge marks a nudge as not resolved
func UnresolveNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchNudge(ctx, fr, ns, thinFeed.UnresolveNudge, w, r)
	}
}

// HideNudge marks a nudge as not resolved
func HideNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchNudge(ctx, fr, ns, thinFeed.HideNudge, w, r)
	}
}

// ShowNudge marks a nudge as not resolved
func ShowNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		patchNudge(ctx, fr, ns, thinFeed.ShowNudge, w, r)
	}
}

// DeleteNudge permanently deletes a nudge
func DeleteNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nudgeID, err := getStringVar(r, "nudgeID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = thinFeed.DeleteNudge(ctx, nudgeID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// PublishAction posts a new action to a user's feed
func PublishAction(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		action := &feed.Action{}
		err = action.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedAction, err := thinFeed.PublishAction(ctx, action)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := publishedAction.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// DeleteAction permanently removes an action from a user's feed
func DeleteAction(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionID, err := getStringVar(r, "actionID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = thinFeed.DeleteAction(ctx, actionID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// PostMessage adds a message to a thread
func PostMessage(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		message := &feed.Message{}
		err = message.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		postedMessage, err := thinFeed.PostMessage(ctx, itemID, message)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(postedMessage)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// DeleteMessage removes a message from a thread
func DeleteMessage(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		messageID, err := getStringVar(r, "messageID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = thinFeed.DeleteMessage(ctx, itemID, messageID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// ProcessEvent saves an event
func ProcessEvent(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		event := &feed.Event{}
		err = event.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = thinFeed.ProcessEvent(ctx, event)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	errMap := base.ErrorMap(err)
	errBytes, err := json.Marshal(errMap)
	if err != nil {
		errBytes = []byte(fmt.Sprintf("error: %s", err))
	}
	respondWithJSON(w, code, errBytes)
}

func respondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(payload)
	if err != nil {
		log.Printf(
			"unable to write payload `%s` to the http.ResponseWriter: %s",
			string(payload),
			err,
		)
	}
}

func getThinFeed(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
	r *http.Request,
) (*feed.Feed, error) {
	uid, flavour, err := getUIDAndFlavour(r)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate thin feed: %w", err)
	}

	agg, err := feed.NewCollection(fr, ns)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate thin feed: %w", err)
	}

	thinFeed, err := agg.GetThinFeed(ctx, uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate thin feed: %w", err)
	}

	return thinFeed, nil
}

func getUIDAndFlavour(r *http.Request) (string, feed.Flavour, error) {
	if r == nil {
		return "", "", fmt.Errorf("nil request")
	}

	uid, err := getStringVar(r, "uid")
	if err != nil {
		return "", "", fmt.Errorf("can't get `uid` path var")
	}

	flavourStr, err := getStringVar(r, "flavour")
	if err != nil {
		return "", "", fmt.Errorf("can't get `flavour` path var: %w", err)
	}

	flavour := feed.Flavour(flavourStr)
	if !flavour.IsValid() {
		return "", "", fmt.Errorf("`%s` is not a valid feed flavour", err)
	}

	return uid, flavour, nil
}

type patchItemFunc func(ctx context.Context, itemID string) (*feed.Item, error)

func patchItem(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
	patchFunc patchItemFunc,
	w http.ResponseWriter,
	r *http.Request,
) {
	itemID, err := getStringVar(r, "itemID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	item, err := patchFunc(ctx, itemID)
	if err != nil {
		if errors.Is(err, feed.ErrNilFeedItem) {
			respondWithError(w, http.StatusNotFound, err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(item)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, marshalled)
}

type patchNudgeFunc func(ctx context.Context, nudgeID string) (*feed.Nudge, error)

func patchNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
	patchFunc patchNudgeFunc,
	w http.ResponseWriter,
	r *http.Request,
) {
	nudgeID, err := getStringVar(r, "nudgeID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	element, err := patchFunc(ctx, nudgeID)
	if err != nil {
		if errors.Is(err, feed.ErrNilNudge) {
			respondWithError(w, http.StatusNotFound, err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(element)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, marshalled)
}

func getOptionalBooleanFilterQueryParam(r *http.Request, paramName string) (*feed.BooleanFilter, error) {
	val := r.FormValue(paramName)
	if val == "" {
		return nil, nil // optional
	}

	boolFilter := feed.BooleanFilter(val)
	if !boolFilter.IsValid() {
		return nil, fmt.Errorf("optional bool: `%s` is not a valid boolean filter value", val)
	}

	return &boolFilter, nil
}

func getRequiredBooleanFilterQueryParam(r *http.Request, paramName string) (feed.BooleanFilter, error) {
	val := r.FormValue(paramName)
	if val == "" {
		return "", fmt.Errorf("required BooleanFilter `%s` not set", paramName)
	}

	boolFilter := feed.BooleanFilter(val)
	if !boolFilter.IsValid() {
		return "", fmt.Errorf("required bool: `%s` is not a valid boolean filter value", val)
	}

	return boolFilter, nil
}

func getOptionalStatusQueryParam(
	r *http.Request,
	paramName string,
) (*feed.Status, error) {
	val, err := getStringVar(r, paramName)
	if err != nil {
		return nil, nil // this is an optional param
	}

	status := feed.Status(val)
	if !status.IsValid() {
		return nil, fmt.Errorf("`%s` is not a valid status", val)
	}

	return &status, nil
}

func getOptionalVisibilityQueryParam(
	r *http.Request,
	paramName string,
) (*feed.Visibility, error) {
	val, err := getStringVar(r, paramName)
	if err != nil {
		return nil, nil // this is an optional param
	}

	visibility := feed.Visibility(val)
	if !visibility.IsValid() {
		return nil, fmt.Errorf("`%s` is not a valid visibility value", val)
	}

	return &visibility, nil
}

func getOptionalFilterParamsQueryParam(
	r *http.Request,
	paramName string,
) (*feed.FilterParams, error) {
	// expect the filter params value to be JSON encoded
	val, err := getStringVar(r, paramName)
	if err != nil {
		return nil, nil // this is an optional param
	}

	filterParams := &feed.FilterParams{}
	err = json.Unmarshal([]byte(val), filterParams)
	if err != nil {
		return nil, fmt.Errorf(
			"filter params should be a valid JSON representation of `feed.FilterParams`. `%s` is not", val)
	}

	return filterParams, nil
}

func getStringVar(r *http.Request, varName string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("can't get string var from a nil request")
	}
	pathVars := mux.Vars(r)
	pathVar, found := pathVars[varName]
	if !found {
		return "", fmt.Errorf("the request does not have a path var named `%s`", varName)
	}
	return pathVar, nil
}

func schemaHandler() (http.Handler, error) {
	f, err := pkger.Open(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("can't open pkger schema dir: %w", err)
	}
	defer f.Close()

	return http.StripPrefix("/schema", http.FileServer(f)), nil
}
