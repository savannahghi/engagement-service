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

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/generated"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	db "gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging/pubsubhandlers"
	"gitlab.slade360emr.com/go/engagement/graph/uploads"
)

const (
	// StaticDir is the directory that contains schemata, default images etc
	StaticDir = "gitlab.slade360emr.com/go/engagement:/static/"

	mbBytes              = 1048576
	serverTimeoutSeconds = 120
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

// these are expensive initializations that should be done only once
var ns feed.NotificationService
var fr feed.Repository
var firebaseApp base.IFirebaseApp
var err error

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {

	if fr == nil {
		fr, err = db.NewFirebaseRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
		}
	}

	if ns == nil {
		projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
		if err != nil {
			return nil, fmt.Errorf(
				"can't get projectID from env var `%s`: %w", base.GoogleCloudProjectIDEnvVarName, err)
		}

		ns, err = messaging.NewPubSubNotificationService(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("can't instantiate notification service in resolver: %w", err)
		}
	}

	if firebaseApp == nil {
		fc := &base.FirebaseClient{}
		firebaseApp, err = fc.InitFirebase()
		if err != nil {
			return nil, err
		}
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
	r.Path(base.PubSubHandlerPath).Methods(
		http.MethodPost).HandlerFunc(GoogleCloudPubSubHandler)

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

	// Bulk routes
	bulk := r.PathPrefix("/bulk/").Subrouter()
	bulk.Use(base.InterServiceAuthenticationMiddleware())

	// Interservice Authenticated routes
	feedISC := r.PathPrefix("/feed/{uid}/{flavour}/{isAnonymous}/").Subrouter()
	feedISC.Use(base.InterServiceAuthenticationMiddleware())

	// retrieval
	feedISC.Methods(
		http.MethodGet,
	).Path("/").HandlerFunc(
		GetFeed(ctx, fr, ns),
	).Name("getFeed")

	feedISC.Methods(
		http.MethodGet,
	).Path("/items/{itemID}/").HandlerFunc(
		GetFeedItem(ctx, fr, ns),
	).Name("getFeedItem")

	feedISC.Methods(
		http.MethodGet,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		GetNudge(ctx, fr, ns),
	).Name("getNudge")

	feedISC.Methods(
		http.MethodGet,
	).Path("/actions/{actionID}/").HandlerFunc(
		GetAction(ctx, fr, ns),
	).Name("getAction")

	// creation
	feedISC.Methods(
		http.MethodPost,
	).Path("/items/").HandlerFunc(
		PublishFeedItem(ctx, fr, ns),
	).Name("publishFeedItem")

	feedISC.Methods(
		http.MethodPost,
	).Path("/nudges/").HandlerFunc(
		PublishNudge(ctx, fr, ns),
	).Name("publishNudge")

	feedISC.Methods(
		http.MethodPost,
	).Path("/actions/").HandlerFunc(
		PublishAction(ctx, fr, ns),
	).Name("publishAction")

	feedISC.Methods(
		http.MethodPost,
	).Path("/{itemID}/messages/").HandlerFunc(
		PostMessage(ctx, fr, ns),
	).Name("postMessage")

	feedISC.Methods(
		http.MethodPost,
	).Path("/events/").HandlerFunc(
		ProcessEvent(ctx, fr, ns),
	).Name("postEvent")

	// deleting
	feedISC.Methods(
		http.MethodDelete,
	).Path("/items/{itemID}/").HandlerFunc(
		DeleteFeedItem(ctx, fr, ns),
	).Name("deleteFeedItem")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		DeleteNudge(ctx, fr, ns),
	).Name("deleteNudge")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/actions/{actionID}/").HandlerFunc(
		DeleteAction(ctx, fr, ns),
	).Name("deleteAction")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/{itemID}/messages/{messageID}/").HandlerFunc(
		DeleteMessage(ctx, fr, ns),
	).Name("deleteMessage")

	// modifying (patching)
	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/resolve/").HandlerFunc(
		ResolveFeedItem(ctx, fr, ns),
	).Name("resolveFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unresolve/").HandlerFunc(
		UnresolveFeedItem(ctx, fr, ns),
	).Name("unresolveFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/pin/").HandlerFunc(
		PinFeedItem(ctx, fr, ns),
	).Name("pinFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unpin/").HandlerFunc(
		UnpinFeedItem(ctx, fr, ns),
	).Name("unpinFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/hide/").HandlerFunc(
		HideFeedItem(ctx, fr, ns),
	).Name("hideFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/show/").HandlerFunc(
		ShowFeedItem(ctx, fr, ns),
	).Name("showFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/resolve/").HandlerFunc(
		ResolveNudge(ctx, fr, ns),
	).Name("resolveNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/defaultnudges/{title}/resolve/").HandlerFunc(
		ResolveDefaultNudge(ctx, fr, ns),
	).Name("resolveDefaultNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/unresolve/").HandlerFunc(
		UnresolveNudge(ctx, fr, ns),
	).Name("unresolveNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/show/").HandlerFunc(
		ShowNudge(ctx, fr, ns),
	).Name("showNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/hide/").HandlerFunc(
		HideNudge(ctx, fr, ns),
	).Name("hideNudge")

	isc := r.PathPrefix("/internal/").Subrouter()
	isc.Use(base.InterServiceAuthenticationMiddleware())

	isc.Methods(
		http.MethodGet,
	).Path("/upload/{uploadID}/").HandlerFunc(
		FindUpload(ctx),
	).Name("getUpload")

	isc.Methods(
		http.MethodPost,
	).Path("/upload/").HandlerFunc(
		Upload(ctx),
	).Name("upload")

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

// GoogleCloudPubSubHandler receives push messages from Google Cloud Pub-Sub
func GoogleCloudPubSubHandler(w http.ResponseWriter, r *http.Request) {
	m, err := base.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
		return
	}

	topicID, err := base.GetPubSubTopic(m)
	if err != nil {
		base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	switch topicID {
	case feed.ItemPublishTopic:
		err = pubsubhandlers.HandleItemPublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemDeleteTopic:
		err = pubsubhandlers.HandleItemDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemResolveTopic:
		err = pubsubhandlers.HandleItemResolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemUnresolveTopic:
		err = pubsubhandlers.HandleItemUnresolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemHideTopic:
		err = pubsubhandlers.HandleItemHide(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemShowTopic:
		err = pubsubhandlers.HandleItemShow(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemPinTopic:
		err = pubsubhandlers.HandleItemPin(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ItemUnpinTopic:
		err = pubsubhandlers.HandleItemUnpin(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.NudgePublishTopic:
		err = pubsubhandlers.HandleNudgePublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.NudgeDeleteTopic:
		err = pubsubhandlers.HandleNudgeDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.NudgeResolveTopic:
		err = pubsubhandlers.HandleNudgeResolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.NudgeUnresolveTopic:
		err = pubsubhandlers.HandleNudgeUnresolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.NudgeHideTopic:
		err = pubsubhandlers.HandleNudgeHide(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.NudgeShowTopic:
		err = pubsubhandlers.HandleNudgeShow(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ActionPublishTopic:
		err = pubsubhandlers.HandleActionPublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.ActionDeleteTopic:
		err = pubsubhandlers.HandleActionDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.MessagePostTopic:
		err = pubsubhandlers.HandleMessagePost(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.MessageDeleteTopic:
		err = pubsubhandlers.HandleMessageDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case feed.IncomingEventTopic:
		err = pubsubhandlers.HandleIncomingEvent(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	default:
		// the topic should be anticipated/handled here
		errMsg := fmt.Sprintf(
			"pub sub handler error: unknown topic `%s`",
			topicID,
		)
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	resp := map[string]string{"status": "success"}
	marshalledSuccessMsg, err := json.Marshal(resp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write(marshalledSuccessMsg)
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
		uid, flavour, anonymous, err := getUIDFlavourAndIsAnonymous(r)
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
			anonymous,
			*flavour,
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

		item := &base.Item{}
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

		nudge := &base.Nudge{}
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

// ResolveDefaultNudge marks a default nudges as resolved
func ResolveDefaultNudge(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getStringVar(r, "title")

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		thinFeed, err := getThinFeed(ctx, fr, ns, r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudge, err := thinFeed.GetDefaultNudgeByTitle(ctx, title)
		if err != nil {
			if errors.Is(err, feed.ErrNilNudge) {
				respondWithError(w, http.StatusNotFound, err)
				return
			}
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		_, err = thinFeed.ResolveNudge(ctx, nudge.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, nil)

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

		action := &base.Action{}
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

		message := &base.Message{}
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

		event := &base.Event{}
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

// Upload saves an upload in cloud storage
func Upload(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uploadInput := base.UploadInput{}
		err = json.Unmarshal(data, &uploadInput)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		if uploadInput.Base64data == "" {
			err := fmt.Errorf("blank upload base64 data")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		if uploadInput.Filename == "" {
			err := fmt.Errorf("blank upload filename")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		if uploadInput.Title == "" {
			err := fmt.Errorf("blank upload title")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uploadService := uploads.NewService()
		upload, err := uploadService.Upload(ctx, uploadInput)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if upload == nil {
			err := fmt.Errorf("nil upload in response from upload service")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		marshalled, err := json.Marshal(upload)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// FindUpload retrieves an upload by it's ID
func FindUpload(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadID, err := getStringVar(r, "uploadID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uploadService := uploads.NewService()
		upload, err := uploadService.FindUploadByID(ctx, uploadID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if upload == nil {
			err := fmt.Errorf("nil upload in response from upload service")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		marshalled, err := json.Marshal(upload)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}
