package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.slade360emr.com/go/base"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/exceptions"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
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
func (p PresentationHandlersImpl) GoogleCloudPubSubHandler(w http.ResponseWriter, r *http.Request) {
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
	case helpers.AddPubSubNamespace(common.ItemPublishTopic):
		err = p.interactor.Notification.HandleItemPublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemDeleteTopic):
		err = p.interactor.Notification.HandleItemDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemResolveTopic):
		err = p.interactor.Notification.HandleItemResolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemUnresolveTopic):
		err = p.interactor.Notification.HandleItemUnresolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemHideTopic):
		err = p.interactor.Notification.HandleItemHide(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemShowTopic):
		err = p.interactor.Notification.HandleItemShow(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemPinTopic):
		err = p.interactor.Notification.HandleItemPin(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemUnpinTopic):
		err = p.interactor.Notification.HandleItemUnpin(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgePublishTopic):
		err = p.interactor.Notification.HandleNudgePublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeDeleteTopic):
		err = p.interactor.Notification.HandleNudgeDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeResolveTopic):
		err = p.interactor.Notification.HandleNudgeResolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeUnresolveTopic):
		err = p.interactor.Notification.HandleNudgeUnresolve(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeHideTopic):
		err = p.interactor.Notification.HandleNudgeHide(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeShowTopic):
		err = p.interactor.Notification.HandleNudgeShow(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ActionPublishTopic):
		err = p.interactor.Notification.HandleActionPublish(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.ActionDeleteTopic):
		err = p.interactor.Notification.HandleActionDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.MessagePostTopic):
		err = p.interactor.Notification.HandleMessagePost(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.MessageDeleteTopic):
		err = p.interactor.Notification.HandleMessageDelete(ctx, m)
		if err != nil {
			base.WriteJSONResponse(w, base.ErrorMap(err), http.StatusBadRequest)
			return
		}
	case helpers.AddPubSubNamespace(common.IncomingEventTopic):
		err = p.interactor.Notification.HandleIncomingEvent(ctx, m)
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

		filterParams, err := getOptionalFilterParamsQueryParam(r, "filterParams")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		feed, err := p.interactor.Feed.GetFeed(
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

		item, err := p.interactor.Feed.GetFeedItem(ctx, *uid, *flavour, itemID)
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

		nudge, err := p.interactor.Feed.GetNudge(ctx, *uid, *flavour, nudgeID)
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

		action, err := p.interactor.Feed.GetAction(ctx, *uid, *flavour, actionID)
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

		publishedItem, err := p.interactor.Feed.PublishFeedItem(ctx, *uid, *flavour, item)
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

		err = p.interactor.Feed.DeleteFeedItem(ctx, *uid, *flavour, itemID)
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

		publishedNudge, err := p.interactor.Feed.PublishNudge(ctx, *uid, *flavour, nudge)
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

		nudge, err := p.interactor.Feed.GetDefaultNudgeByTitle(ctx, *uid, *flavour, title)
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

		_, err = p.interactor.Feed.ResolveNudge(ctx, *uid, *flavour, nudge.ID)
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

		err = p.interactor.Feed.DeleteNudge(ctx, *uid, *flavour, nudgeID)
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

		publishedAction, err := p.interactor.Feed.PublishAction(ctx, *uid, *flavour, action)
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

		err = p.interactor.Feed.DeleteAction(ctx, *uid, *flavour, actionID)
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

		postedMessage, err := p.interactor.Feed.PostMessage(ctx, *uid, *flavour, itemID, message)
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

		err = p.interactor.Feed.DeleteMessage(ctx, *uid, *flavour, itemID, messageID)
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

		err = p.interactor.Feed.ProcessEvent(ctx, *uid, *flavour, event)
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
func (p PresentationHandlersImpl) Upload(ctx context.Context) http.HandlerFunc {
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
func (p PresentationHandlersImpl) FindUpload(ctx context.Context) http.HandlerFunc {
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
func (p PresentationHandlersImpl) SendEmail(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &domain.EMailMessage{}
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
