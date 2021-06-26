package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/otp"

	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.slade360emr.com/go/base"

	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/exceptions"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/interactor"
)

const (
	// StaticDir is the directory that contains schemata, default images etc
	StaticDir = "gitlab.slade360emr.com/go/engagement:/static/"

	mbBytes              = 1048576
	serverTimeoutSeconds = 120
)

var errNotFound = fmt.Errorf("not found")

// PresentationHandlers represents all the REST API logic
type PresentationHandlers interface {
	GoogleCloudPubSubHandler(w http.ResponseWriter, r *http.Request)
	GetFeed(
		ctx context.Context,
	) http.HandlerFunc

	GetFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	GetNudge(
		ctx context.Context,
	) http.HandlerFunc

	GetAction(
		ctx context.Context,
	) http.HandlerFunc

	PublishFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	DeleteFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	ResolveFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	PinFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	UnpinFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	HideFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	ShowFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	UnresolveFeedItem(
		ctx context.Context,
	) http.HandlerFunc

	PublishNudge(
		ctx context.Context,
	) http.HandlerFunc

	ResolveNudge(
		ctx context.Context,
	) http.HandlerFunc

	ResolveDefaultNudge(
		ctx context.Context,
	) http.HandlerFunc

	UnresolveNudge(
		ctx context.Context,
	) http.HandlerFunc

	HideNudge(
		ctx context.Context,
	) http.HandlerFunc

	ShowNudge(
		ctx context.Context,
	) http.HandlerFunc

	DeleteNudge(
		ctx context.Context,
	) http.HandlerFunc

	PublishAction(
		ctx context.Context,
	) http.HandlerFunc

	DeleteAction(
		ctx context.Context,
	) http.HandlerFunc

	PostMessage(
		ctx context.Context,
	) http.HandlerFunc

	DeleteMessage(
		ctx context.Context,
	) http.HandlerFunc

	ProcessEvent(
		ctx context.Context,
	) http.HandlerFunc

	Upload(ctx context.Context) http.HandlerFunc

	FindUpload(ctx context.Context) http.HandlerFunc

	SendEmail(ctx context.Context) http.HandlerFunc

	SendToMany(ctx context.Context) http.HandlerFunc

	GetAITSMSDeliveryCallback(ctx context.Context) http.HandlerFunc

	GetNotificationHandler(ctx context.Context) http.HandlerFunc

	GetIncomingMessageHandler(ctx context.Context) http.HandlerFunc

	GetFallbackHandler(ctx context.Context) http.HandlerFunc

	PhoneNumberVerificationCodeHandler(ctx context.Context) http.HandlerFunc

	SendOTPHandler() http.HandlerFunc

	SendRetryOTPHandler(ctx context.Context) http.HandlerFunc

	VerifyRetryOTPHandler(ctx context.Context) http.HandlerFunc

	VerifyRetryEmailOTPHandler(ctx context.Context) http.HandlerFunc

	SendNotificationHandler(ctx context.Context) http.HandlerFunc

	GetContactLists() http.HandlerFunc
	GetContactListByID() http.HandlerFunc
	GetContactsInAList() http.HandlerFunc

	SetBewellAware() http.HandlerFunc
}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	interactor *interactor.Interactor
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(i *interactor.Interactor) PresentationHandlers {
	return &PresentationHandlersImpl{i}
}

// GoogleCloudPubSubHandler receives push messages from Google Cloud Pub-Sub
func (p PresentationHandlersImpl) GoogleCloudPubSubHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	// get the UID frrom the payload
	var envelope dto.NotificationEnvelope
	err = json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
		return
	}
	ctx := addUIDToContext(envelope.UID)

	switch topicID {
	case helpers.AddPubSubNamespace(common.ItemPublishTopic):
		err = p.interactor.Notification.HandleItemPublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemDeleteTopic):
		err = p.interactor.Notification.HandleItemDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemResolveTopic):
		err = p.interactor.Notification.HandleItemResolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemUnresolveTopic):
		err = p.interactor.Notification.HandleItemUnresolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemHideTopic):
		err = p.interactor.Notification.HandleItemHide(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemShowTopic):
		err = p.interactor.Notification.HandleItemShow(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemPinTopic):
		err = p.interactor.Notification.HandleItemPin(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemUnpinTopic):
		err = p.interactor.Notification.HandleItemUnpin(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgePublishTopic):
		err = p.interactor.Notification.HandleNudgePublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeDeleteTopic):
		err = p.interactor.Notification.HandleNudgeDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeResolveTopic):
		err = p.interactor.Notification.HandleNudgeResolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeUnresolveTopic):
		err = p.interactor.Notification.HandleNudgeUnresolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeHideTopic):
		err = p.interactor.Notification.HandleNudgeHide(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeShowTopic):
		err = p.interactor.Notification.HandleNudgeShow(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ActionPublishTopic):
		err = p.interactor.Notification.HandleActionPublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ActionDeleteTopic):
		err = p.interactor.Notification.HandleActionDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.MessagePostTopic):
		err = p.interactor.Notification.HandleMessagePost(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.MessageDeleteTopic):
		err = p.interactor.Notification.HandleMessageDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.IncomingEventTopic):
		err = p.interactor.Notification.HandleIncomingEvent(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.FcmPublishTopic):
		err = p.interactor.Notification.HandleSendNotification(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.SentEmailTopic):
		err = p.interactor.Notification.SendEmail(ctx, m)
		if err != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusBadRequest,
			)
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

// GetFeed retrieves and serves a feed
func (p PresentationHandlersImpl) GetFeed(
	ctx context.Context,
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

		filterParams, err := getOptionalFilterParamsQueryParam(
			r,
			"filterParams",
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		feed, err := p.interactor.Feed.GetFeed(
			addUIDToContext(*uid),
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
func (p PresentationHandlersImpl) GetFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		item, err := p.interactor.Feed.GetFeedItem(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			itemID,
		)
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
func (p PresentationHandlersImpl) GetNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudgeID, err := getStringVar(r, "nudgeID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		ctx = addUIDToContext(*uid)
		nudge, err := p.interactor.Feed.GetNudge(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			nudgeID,
		)
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
func (p PresentationHandlersImpl) GetAction(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		actionID, err := getStringVar(r, "actionID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		action, err := p.interactor.Feed.GetAction(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			actionID,
		)
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
func (p PresentationHandlersImpl) PublishFeedItem(
	ctx context.Context,
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

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedItem, err := p.interactor.Feed.PublishFeedItem(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			item,
		)
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
func (p PresentationHandlersImpl) DeleteFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.interactor.Feed.DeleteFeedItem(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			itemID,
		)
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
func (p PresentationHandlersImpl) ResolveFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchItem(ctx, p.interactor.Feed.ResolveFeedItem, w, r)
	}
}

// PinFeedItem marks a feed item as done
func (p PresentationHandlersImpl) PinFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchItem(ctx, p.interactor.Feed.PinFeedItem, w, r)
	}
}

// UnpinFeedItem marks a feed item as done
func (p PresentationHandlersImpl) UnpinFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchItem(ctx, p.interactor.Feed.UnpinFeedItem, w, r)
	}
}

// HideFeedItem marks a feed item as done
func (p PresentationHandlersImpl) HideFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchItem(ctx, p.interactor.Feed.HideFeedItem, w, r)
	}
}

// ShowFeedItem marks a feed item as done
func (p PresentationHandlersImpl) ShowFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchItem(ctx, p.interactor.Feed.ShowFeedItem, w, r)
	}
}

// UnresolveFeedItem marks a feed item as not resolved
func (p PresentationHandlersImpl) UnresolveFeedItem(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchItem(ctx, p.interactor.Feed.UnresolveFeedItem, w, r)
	}
}

// PublishNudge posts a new nudge
func (p PresentationHandlersImpl) PublishNudge(
	ctx context.Context,
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

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedNudge, err := p.interactor.Feed.PublishNudge(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			nudge,
		)
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
func (p PresentationHandlersImpl) ResolveNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchNudge(ctx, p.interactor.Feed.ResolveNudge, w, r)
	}
}

// ResolveDefaultNudge marks a default nudges as resolved
func (p PresentationHandlersImpl) ResolveDefaultNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getStringVar(r, "title")

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudge, err := p.interactor.Feed.GetDefaultNudgeByTitle(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			title,
		)
		if err != nil {
			if errors.Is(err, exceptions.ErrNilNudge) {
				respondWithError(w, http.StatusNotFound, err)
				return
			}
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if nudge.Status == base.StatusDone {
			respondWithJSON(w, http.StatusOK, marshalled)
		}

		_, err = p.interactor.Feed.ResolveNudge(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			nudge.ID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// UnresolveNudge marks a nudge as not resolved
func (p PresentationHandlersImpl) UnresolveNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchNudge(ctx, p.interactor.Feed.UnresolveNudge, w, r)
	}
}

// HideNudge marks a nudge as not resolved
func (p PresentationHandlersImpl) HideNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchNudge(ctx, p.interactor.Feed.HideNudge, w, r)
	}
}

// ShowNudge marks a nudge as not resolved
func (p PresentationHandlersImpl) ShowNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patchNudge(ctx, p.interactor.Feed.ShowNudge, w, r)
	}
}

// DeleteNudge permanently deletes a nudge
func (p PresentationHandlersImpl) DeleteNudge(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nudgeID, err := getStringVar(r, "nudgeID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.interactor.Feed.DeleteNudge(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			nudgeID,
		)
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
func (p PresentationHandlersImpl) PublishAction(
	ctx context.Context,
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

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedAction, err := p.interactor.Feed.PublishAction(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			action,
		)
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
func (p PresentationHandlersImpl) DeleteAction(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionID, err := getStringVar(r, "actionID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.interactor.Feed.DeleteAction(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			actionID,
		)
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
func (p PresentationHandlersImpl) PostMessage(
	ctx context.Context,
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

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		postedMessage, err := p.interactor.Feed.PostMessage(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			itemID,
			message,
		)
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
func (p PresentationHandlersImpl) DeleteMessage(
	ctx context.Context,
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

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.interactor.Feed.DeleteMessage(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			itemID,
			messageID,
		)
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
func (p PresentationHandlersImpl) ProcessEvent(
	ctx context.Context,
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

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.interactor.Feed.ProcessEvent(
			addUIDToContext(*uid),
			*uid,
			*flavour,
			event,
		)
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
func (p PresentationHandlersImpl) Upload(
	ctx context.Context,
) http.HandlerFunc {
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

		upload, err := p.interactor.Uploads.Upload(ctx, uploadInput)
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
func (p PresentationHandlersImpl) FindUpload(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadID, err := getStringVar(r, "uploadID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		upload, err := p.interactor.Uploads.FindUploadByID(ctx, uploadID)
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

// SendEmail sends the specified email to the recipient(s) specified in `to`
// and returns the status
func (p PresentationHandlersImpl) SendEmail(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.EMailMessage{}
		base.DecodeJSONToTargetStruct(w, r, payload)
		if payload.Subject == "" {
			err := fmt.Errorf("blank email subject")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if payload.Text == "" {
			err := fmt.Errorf("blank email text")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if len(payload.To) == 0 {
			err := fmt.Errorf("no destination email supplied")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		resp, _, err := p.interactor.Mail.SendEmail(
			payload.Subject,
			payload.Text,
			payload.To...,
		)
		if err != nil {
			err := fmt.Errorf("email not sent: %s", err)
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// SendToMany sends a data message to the specified recipient
func (p PresentationHandlersImpl) SendToMany(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.SendSMSPayload{}
		base.DecodeJSONToTargetStruct(w, r, payload)

		for _, phoneNo := range payload.To {
			_, err := base.NormalizeMSISDN(phoneNo)
			if err != nil {
				err := fmt.Errorf(
					"can't send sms, expected a valid phone number",
				)
				respondWithError(w, http.StatusBadRequest, err)
				return
			}
		}

		if payload.Message == "" {
			err := fmt.Errorf("can't send sms, expected a message")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		_, err := p.interactor.SMS.SendToMany(
			payload.Message,
			payload.To,
			payload.Sender,
		)
		if err != nil {
			err := fmt.Errorf("sms not sent: %s", err)

			isBadReq := strings.Contains(err.Error(), "http error status: 400")

			if isBadReq {
				respondWithError(w, http.StatusBadRequest, err)
				return
			}

			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		type okResp struct {
			Status string `json:"status"`
		}

		marshalled, err := json.Marshal(okResp{Status: "ok"})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetAITSMSDeliveryCallback generates an SMS Delivery Report by saving the callback data for future analysis.
func (p PresentationHandlersImpl) GetAITSMSDeliveryCallback(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parses the request body
		err := r.ParseForm()
		if err != nil {
			log.Printf("unable to parse request data %v", err)
			return
		}
		if r.Form == nil || len(r.Form) == 0 {
			return
		}

		err = p.interactor.SMS.SaveAITCallbackResponse(
			ctx,
			dto.CallbackData{Values: r.Form},
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		marshalled, err := json.Marshal(dto.CallbackData{Values: r.Form})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetNotificationHandler returns a handler that processes an Africa's Talking payment notification
func (p PresentationHandlersImpl) GetNotificationHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.Message{}
		base.DecodeJSONToTargetStruct(w, r, payload)
		if payload.AccountSID == "" {
			err := fmt.Errorf(
				"twilio notification payload not parsed correctly",
			)
			log.Printf("Twilio callback error: %s", err)
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusInternalServerError,
			)
		}

		// save Twilio response for audit purposes
		err := p.interactor.Whatsapp.SaveTwilioCallbackResponse(ctx, *payload)
		if err != nil {
			err := fmt.Errorf("twilio notification payload not saved")
			log.Printf("Twilio callback error: %s", err)
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusInternalServerError,
			)
		}
		// TODO Common pathway for saving, returning OK etc

		type okResp struct {
			Status string `json:"status"`
		}
		base.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// GetIncomingMessageHandler returns a handler that processes an Africa's Talking payment notification
func (p PresentationHandlersImpl) GetIncomingMessageHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.Message{}
		base.DecodeJSONToTargetStruct(w, r, payload)
		if payload.AccountSID == "" {
			err := fmt.Errorf(
				"twilio notification payload not parsed correctly",
			)
			log.Printf("Twilio callback error: %s", err)
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusInternalServerError,
			)
		}

		// save Twilio response for audit purposes
		err := p.interactor.Whatsapp.SaveTwilioCallbackResponse(ctx, *payload)
		if err != nil {
			err := fmt.Errorf("twilio notification payload not saved")
			log.Printf("Twilio callback error: %s", err)
			base.WriteJSONResponse(
				w,
				base.ErrorMap(err),
				http.StatusInternalServerError,
			)
		}
		// TODO Common pathway for saving, returning OK etc

		type okResp struct {
			Status string `json:"status"`
		}
		base.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// GetFallbackHandler returns a handler that processes an Africa's Talking payment notification
func (p PresentationHandlersImpl) GetFallbackHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO ErrorCode and ErrorURL sent here as params
		// TODO Implement WhatsAPP fallback handler: base.DecodeJSONToTargetStruct(w, r, notificationPayload)
		// base.ErrorMap(fmt.Errorf("unbound mandatory notification payload fields")),
		// base.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// PhoneNumberVerificationCodeHandler process ISC request to PhoneNumberVerificationCode
func (p PresentationHandlersImpl) PhoneNumberVerificationCodeHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		type PayloadRequest struct {
			To               string `json:"to"`
			Code             string `json:"code"`
			MarketingMessage string `json:"marketingMessage"`
		}

		payloadRequest := &PayloadRequest{}

		base.DecodeJSONToTargetStruct(rw, r, payloadRequest)

		ok, err := p.interactor.Whatsapp.PhoneNumberVerificationCode(
			ctx,
			payloadRequest.To,
			payloadRequest.Code,
			payloadRequest.MarketingMessage,
		)
		if err != nil {
			base.RespondWithError(rw, http.StatusInternalServerError, err)
			return
		}

		type PayloadResponse struct {
			Status bool `json:"status"`
		}

		response := &PayloadResponse{Status: ok}
		base.WriteJSONResponse(rw, response, http.StatusOK)
	}
}

// SendOTPHandler is an isc api that generates and sends an otp to an msisdn
func (p PresentationHandlersImpl) SendOTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := otp.NewService()
		msisdn, err := otp.ValidateSendOTPPayload(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		code, codeErr := s.GenerateAndSendOTP(msisdn)
		if codeErr != nil {
			base.WriteJSONResponse(
				w,
				base.ErrorMap(
					fmt.Errorf("unable to generate and send otp: %v", codeErr),
				),
				http.StatusInternalServerError,
			)
		}

		base.WriteJSONResponse(w, code, http.StatusOK)
	}
}

// SendRetryOTPHandler is an isc api that generates
// fallback OTPs when Africa is talking sms fails
func (p PresentationHandlersImpl) SendRetryOTPHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := otp.NewService()
		payload, err := otp.ValidateGenerateRetryOTPPayload(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		code, codeErr := s.GenerateRetryOTP(
			ctx,
			payload.Msisdn,
			payload.RetryStep,
		)
		if codeErr != nil {
			err := base.ErrorMap(
				fmt.Errorf(
					"unable to generate and send a fallback OTP: %v",
					codeErr,
				),
			)
			base.WriteJSONResponse(w, err, http.StatusInternalServerError)
		}

		base.WriteJSONResponse(w, code, http.StatusOK)
	}
}

// VerifyRetryOTPHandler is an isc api that confirms OTPs earlier sent
func (p PresentationHandlersImpl) VerifyRetryOTPHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := otp.NewService()
		payload, err := otp.ValidateVerifyOTPPayload(w, r, false)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		isVerified, err := s.VerifyOtp(
			ctx,
			payload.Msisdn,
			payload.VerificationCode,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		type otpResponse struct {
			IsVerified bool `json:"IsVerified"`
		}

		base.WriteJSONResponse(
			w,
			otpResponse{IsVerified: isVerified},
			http.StatusOK,
		)
	}
}

// VerifyRetryEmailOTPHandler is an isc api that confirms OTPs earlier sent via email.
func (p PresentationHandlersImpl) VerifyRetryEmailOTPHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := otp.NewService()
		payload, err := otp.ValidateVerifyOTPPayload(w, r, true)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		isVerified, err := s.VerifyEmailOtp(
			ctx,
			payload.Email,
			payload.VerificationCode,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		type otpResponse struct {
			IsVerified bool `json:"IsVerified"`
		}

		base.WriteJSONResponse(
			w,
			otpResponse{IsVerified: isVerified},
			http.StatusOK,
		)
	}
}

// SendNotificationHandler sends a data message to the specified registration tokens.
func (p PresentationHandlersImpl) SendNotificationHandler(
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, payloadErr := fcm.ValidateSendNotificationPayload(w, r)
		if payloadErr != nil {
			base.ReportErr(w, payloadErr, http.StatusBadRequest)
			return
		}

		_, err := p.interactor.FCM.SendNotification(
			ctx,
			payload.RegistrationTokens,
			payload.Data,
			payload.Notification,
			payload.Android,
			payload.Ios,
			payload.Web,
		)
		if err != nil {
			err := fmt.Errorf("notification not sent: %s", err)

			isBadReq := strings.Contains(err.Error(), "http error status: 400")

			if isBadReq {
				base.ReportErr(w, err, http.StatusBadRequest)
				return
			}

			base.ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		type okResp struct {
			Status string `json:"status"`
		}
		base.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// GetContactLists fetches all the Contact Lists on hubspot
// todo write automated tests for this (it has already been hand-tested to work)
func (p PresentationHandlersImpl) GetContactLists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactLists, err := p.interactor.CRM.GetContactLists()
		if err != nil {
			base.RespondWithError(w, http.StatusBadRequest, err)
			return
		}
		base.WriteJSONResponse(w, contactLists, http.StatusOK)
	}
}

// GetContactListByID fetches a specific Contact List given its listId
// todo write automated tests for this (it has already been hand-tested to work)
func (p PresentationHandlersImpl) GetContactListByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.ListID{}
		base.DecodeJSONToTargetStruct(w, r, payload)
		contactList, err := p.interactor.CRM.GetContactListByID(payload.ListID)
		if err != nil {
			base.RespondWithError(w, http.StatusBadRequest, err)
			return
		}
		base.WriteJSONResponse(w, contactList, http.StatusOK)
	}
}

// GetContactsInAList fetches all the contacts segmented in a Contact List
// todo write automated tests for this (it has already been hand-tested to work)
func (p PresentationHandlersImpl) GetContactsInAList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.ListID{}
		base.DecodeJSONToTargetStruct(w, r, payload)
		contactList, err := p.interactor.CRM.GetContactsInAList(payload.ListID)
		if err != nil {
			base.RespondWithError(w, http.StatusBadRequest, err)
			return
		}
		base.WriteJSONResponse(w, contactList, http.StatusOK)
	}
}

//SetBewellAware the user identified by the provided email= as bewell-aware on the CRM
// todo write automated tests for this (it has already been hand-tested to work)
func (p PresentationHandlersImpl) SetBewellAware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.SetBewellAwareInput{}
		base.DecodeJSONToTargetStruct(w, r, payload)

		filters := CRMDomain.Filters{
			Value:        payload.EmailAddress,
			PropertyName: "email",
			Operator:     "EQ",
		}
		filtergroup := CRMDomain.FilterGroups{Filters: []CRMDomain.Filters{filters}}
		searchParams := CRMDomain.SearchParams{
			FilterGroups: []CRMDomain.FilterGroups{filtergroup},
			Properties:   []string{"email", "phone", "firstname", "lastname"},
		}

		usercontacts, err := p.interactor.CRM.SearchContact(searchParams)
		if err != nil {
			base.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("failed to search contact %v", err))
			return
		}

		crmContactProperties := CRMDomain.ContactProperties{
			BeWellAware: CRMDomain.GeneralOptionTypeYes,
			Phone:       usercontacts.Results[0].Properties.Phone,
			Email:       usercontacts.Results[0].Properties.Email,
		}

		if _, err := p.interactor.CRM.UpdateContact(usercontacts.Results[0].Properties.Phone, crmContactProperties); err != nil {
			base.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("failed to update contatct %v", err))
			return
		}
		base.WriteJSONResponse(w, dto.OKResp{Status: "SUCCESS"}, http.StatusOK)
	}
}
