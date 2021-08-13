package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/edi"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/library"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/surveys"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/savannahghi/engagement/pkg/engagement/presentation/graph"
	"github.com/savannahghi/engagement/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/serverutils"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/fcm"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/uploads"
	"github.com/savannahghi/engagement/pkg/engagement/presentation/rest"
	"github.com/savannahghi/engagement/pkg/engagement/usecases"

	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	crmExt "github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement/pkg/engagement/presentation/interactor"
	hubspotHandlers "gitlab.slade360emr.com/go/commontools/crm/pkg/presentation/rest"
)

const (
	mbBytes              = 1048576
	serverTimeoutSeconds = 120
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
var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type", " X-Authorization", " Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers",
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate firebase repository in resolver: %w",
			err,
		)
	}
	fcmNotification, err := fcm.NewRemotePushService(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate push notification service : %w",
			err,
		)
	}

	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}

	// Initialize new instances of the infrastructure services
	onboarding := onboarding.NewRemoteProfileService(onboarding.NewOnboardingClient())
	fcm := fcm.NewService(fr, onboarding)
	mail := mail.NewService(fr)
	edi := edi.NewEdiService(edi.NewEDIClient())

	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(ctx, hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr, hubspotService)

	notification := usecases.NewNotification(
		fr,
		fcmNotification,
		onboarding,
		fcm,
		mail,
		hubspotService,
	)
	uploads := uploads.NewUploadsService()
	library := library.NewLibraryService(onboarding)
	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate notification service in resolver: %w",
			err,
		)
	}

	crmExt := crmExt.NewCrmService(hubspotUsecases, mail)
	sms := sms.NewService(fr, crmExt, ns, edi)
	feed := usecases.NewFeed(fr, ns)
	whatsapp := whatsapp.NewService()
	twilio := twilio.NewService(sms)
	otp := otp.NewService(whatsapp, mail, sms, twilio)
	surveys := surveys.NewService(fr)

	// Initialize the interactor
	i, err := interactor.NewEngagementInteractor(
		feed,
		notification,
		uploads,
		library,
		sms,
		*mail,
		whatsapp,
		otp,
		twilio,
		fcm,
		surveys,
		hubspotService,
		crmExt,
		onboarding,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}

	husbspotHandlers := hubspotHandlers.NewHandlers(hubspotUsecases)
	h := rest.NewPresentationHandlers(i, husbspotHandlers)

	r := mux.NewRouter() // gorilla mux
	r.Use(otelmux.Middleware(serverutils.MetricsCollectorService("engagement")))
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(serverutils.RequestDebugMiddleware())

	// Add Middleware that records the metrics for our HTTP routes
	r.Use(serverutils.CustomHTTPRequestMetricsMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	// static files
	schemaFileHandler, err := rest.SchemaHandler()
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate schema file handler: %w",
			err,
		)
	}
	r.PathPrefix("/schema/").Handler(schemaFileHandler)

	// Upload route.
	// The reason for the below endpoint is to help upload base64 data.
	// It is solving a problem ("error": "Unexpected token u in JSON at position 0")
	// that occurs in https://graph-test.bewell.co.ke/ while trying to upload large sized photos
	// This patch allows for the upload of a photo of any size.
	r.Path("/upload").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(h.Upload())

	r.Path("/send_marketing_sms").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(h.SendMarketingSMS())

	r.Path("/set_bewell_aware").Methods(
		http.MethodPost,
	).HandlerFunc(h.SetBewellAware())

	//Authenticated engagement routes

	// Expose a bulk SMS sending endpoint
	engagementAuthenticatedRoutes := r.PathPrefix("").Subrouter()
	engagementAuthenticatedRoutes.Use(interserviceclient.InterServiceAuthenticationMiddleware())
	// HubSpot CRM specific endpoints
	engagementAuthenticatedRoutes.Methods(
		http.MethodPost,
	).Path("/sync_contacts").HandlerFunc(h.HubSpotFirestoreSync())

	engagementAuthenticatedRoutes.Methods(
		http.MethodPost,
	).Path("/collect_email_address").HandlerFunc(h.CollectEmailAddress())

	// Callbacks
	r.Path("/ait_callback").
		Methods(http.MethodPost).
		HandlerFunc(h.GetAITSMSDeliveryCallback())
	r.Path("/twilio_notification").
		Methods(http.MethodPost).
		HandlerFunc(h.GetNotificationHandler())
	r.Path("/twilio_incoming_message").
		Methods(http.MethodPost).
		HandlerFunc(h.GetIncomingMessageHandler())
	r.Path("/twilio_fallback").
		Methods(http.MethodPost).
		HandlerFunc(h.GetFallbackHandler())
	r.Path("/facebook_data_deletion_callback").Methods(
		http.MethodPost,
	).HandlerFunc(h.DataDeletionRequestCallback())

	// Authenticated routes
	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, i))

	// REST routes

	// Interservice Authenticated routes
	feedISC := r.PathPrefix("/feed/{uid}/{flavour}/{isAnonymous}/").Subrouter()
	feedISC.Use(interserviceclient.InterServiceAuthenticationMiddleware())

	// retrieval
	feedISC.Methods(
		http.MethodGet,
	).Path("/").HandlerFunc(
		h.GetFeed(),
	).Name("getFeed")

	feedISC.Methods(
		http.MethodGet,
	).Path("/items/{itemID}/").HandlerFunc(
		h.GetFeedItem(),
	).Name("getFeedItem")

	feedISC.Methods(
		http.MethodGet,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		h.GetNudge(),
	).Name("getNudge")

	feedISC.Methods(
		http.MethodGet,
	).Path("/actions/{actionID}/").HandlerFunc(
		h.GetAction(),
	).Name("getAction")

	// creation
	feedISC.Methods(
		http.MethodPost,
	).Path("/items/").HandlerFunc(
		h.PublishFeedItem(),
	).Name("publishFeedItem")

	feedISC.Methods(
		http.MethodPost,
	).Path("/nudges/").HandlerFunc(
		h.PublishNudge(),
	).Name("publishNudge")

	feedISC.Methods(
		http.MethodPost,
	).Path("/actions/").HandlerFunc(
		h.PublishAction(),
	).Name("publishAction")

	feedISC.Methods(
		http.MethodPost,
	).Path("/{itemID}/messages/").HandlerFunc(
		h.PostMessage(),
	).Name("postMessage")

	feedISC.Methods(
		http.MethodPost,
	).Path("/events/").HandlerFunc(
		h.ProcessEvent(),
	).Name("postEvent")

	// deleting
	feedISC.Methods(
		http.MethodDelete,
	).Path("/items/{itemID}/").HandlerFunc(
		h.DeleteFeedItem(),
	).Name("deleteFeedItem")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		h.DeleteNudge(),
	).Name("deleteNudge")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/actions/{actionID}/").HandlerFunc(
		h.DeleteAction(),
	).Name("deleteAction")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/{itemID}/messages/{messageID}/").HandlerFunc(
		h.DeleteMessage(),
	).Name("deleteMessage")

	// modifying (patching)
	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/resolve/").HandlerFunc(
		h.ResolveFeedItem(),
	).Name("resolveFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unresolve/").HandlerFunc(
		h.UnresolveFeedItem(),
	).Name("unresolveFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/pin/").HandlerFunc(
		h.PinFeedItem(),
	).Name("pinFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unpin/").HandlerFunc(
		h.UnpinFeedItem(),
	).Name("unpinFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/hide/").HandlerFunc(
		h.HideFeedItem(),
	).Name("hideFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/show/").HandlerFunc(
		h.ShowFeedItem(),
	).Name("showFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/resolve/").HandlerFunc(
		h.ResolveNudge(),
	).Name("resolveNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/defaultnudges/{title}/resolve/").HandlerFunc(
		h.ResolveDefaultNudge(),
	).Name("resolveDefaultNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/unresolve/").HandlerFunc(
		h.UnresolveNudge(),
	).Name("unresolveNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/show/").HandlerFunc(
		h.ShowNudge(),
	).Name("showNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/hide/").HandlerFunc(
		h.HideNudge(),
	).Name("hideNudge")

	isc := r.PathPrefix("/internal/").Subrouter()
	isc.Use(interserviceclient.InterServiceAuthenticationMiddleware())

	isc.Methods(
		http.MethodGet,
	).Path("/upload/{uploadID}/").HandlerFunc(
		h.FindUpload(),
	).Name("getUpload")

	isc.Methods(
		http.MethodPost,
	).Path("/upload/").HandlerFunc(
		h.Upload(),
	).Name("upload")

	isc.Methods(
		http.MethodPost,
	).Path("/send_email").HandlerFunc(
		h.SendEmail(),
	).Name("sendEmail")

	isc.Methods(
		http.MethodPost,
	).Path("/mailgun_delivery_webhook").HandlerFunc(
		h.UpdateMailgunDeliveryStatus(),
	).Name("mailgun_delivery_webhook")

	isc.Methods(
		http.MethodPost,
	).Path("/send_sms").HandlerFunc(
		h.SendToMany(),
	).Name("sendToMany")

	isc.Path("/verify_phonenumber").Methods(http.MethodPost).HandlerFunc(
		h.PhoneNumberVerificationCodeHandler(),
	)

	isc.Path("/send_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendOTPHandler())

	isc.Path("/send_retry_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendRetryOTPHandler())

	isc.Path("/verify_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.VerifyRetryOTPHandler())

	isc.Path("/verify_email_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.VerifyRetryEmailOTPHandler())

	isc.Path("/send_notification").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendNotificationHandler())

	// return the combined router
	return r, nil
}

// PrepareServer starts up a server
func PrepareServer(
	ctx context.Context,
	port int,
	allowedOrigins []string,
) *http.Server {
	// start up the router
	r, err := Router(ctx)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
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
	h = handlers.ContentTypeHandler(
		h,
		"application/json",
		"application/x-www-form-urlencoded",
	)
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
	service *interactor.Interactor,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, service)
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
