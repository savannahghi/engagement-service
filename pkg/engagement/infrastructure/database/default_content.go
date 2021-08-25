package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"

	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/library"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/engagement-service/pkg/engagement/repository"

	"github.com/markbates/pkger"
	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"
)

const (
	defaultSequenceNumber = 1
	defaultPostedByUID    = "hOcaUv8dqqgmWYf9HEhjdudgf0b2"
	futureHours           = 878400 // hours in a century of leap years...

	getConsultationActionName     = "GET_CONSULTATION"
	getMedicineActionName         = "GET_MEDICINE"
	getTestActionName             = "GET_TEST"
	getInsuranceActionName        = "GET_INSURANCE"
	addPatientActionName          = "ADD_PATIENT"
	searchPatientActionName       = "SEARCH_PATIENT"
	addNHIFActionName             = "ADD_NHIF"
	partnerAccountSetupActionName = "PARTNER_ACCOUNT_SETUP"
	verifyEmailActionName         = "VERIFY_EMAIL"

	defaultOrg           = "default-org-id-please-change"
	defaultLocation      = "default-location-id-please-change"
	defaultContentDir    = "/static/"
	defaultAuthor        = "Be.Well Team"
	defaultInsuranceText = "Insurance Simplified"
	onboardingService    = "profile"
)

// embed default content assets (e.g images and documents) in the binary
var _ = pkger.Dir(defaultContentDir)

type actionGenerator func(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error)

type nudgeGenerator func(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Nudge, error)

type itemGenerator func(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Item, error)

// SetDefaultActions ensures that a feed has default actions
func SetDefaultActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "SetDefaultActions")
	defer span.End()
	actions := []feedlib.Action{}

	switch flavour {
	case feedlib.FlavourConsumer:
		consumerActions, err := defaultConsumerActions(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to initialize default consumer actions: %w", err)
		}
		actions = consumerActions
	case feedlib.FlavourPro:
		proActions, err := defaultProActions(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to initialize default pro actions: %w", err)
		}
		actions = proActions
	}

	return actions, nil
}

// SetDefaultNudges ensures that a feed has default nudges
func SetDefaultNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "SetDefaultNudges")
	defer span.End()
	var nudges []feedlib.Nudge

	switch flavour {
	case feedlib.FlavourConsumer:
		consumerNudges, err := defaultConsumerNudges(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to initialize default consumer nudges: %w", err)
		}
		nudges = consumerNudges
	case feedlib.FlavourPro:
		proNudges, err := defaultProNudges(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to initialize default pro nudges: %w", err)
		}
		nudges = proNudges
	}

	return nudges, nil
}

// SetDefaultItems ensures that a feed has default feed items
func SetDefaultItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "SetDefaultItems")
	defer span.End()
	var items []feedlib.Item

	switch flavour {
	case feedlib.FlavourConsumer:
		consumerItems, err := defaultConsumerItems(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to initialize default consumer items: %w", err)
		}
		items = consumerItems
	case feedlib.FlavourPro:
		proItems, err := defaultProItems(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to initialize default pro items: %w", err)
		}
		items = proItems
	}

	return items, nil
}

func defaultConsumerNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "defaultConsumerNudges")
	defer span.End()
	var nudges []feedlib.Nudge
	fns := []nudgeGenerator{
		verifyEmailNudge,
	}
	for _, fn := range fns {
		nudge, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("error when generating nudge: %w", err)
		}
		nudges = append(nudges, *nudge)
	}
	return nudges, nil
}

func defaultProNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "defaultProNudges")
	defer span.End()
	var nudges []feedlib.Nudge
	fns := []nudgeGenerator{
		partnerAccountSetupNudge,
		verifyEmailNudge,
	}
	for _, fn := range fns {
		nudge, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("error when generating nudge: %w", err)
		}
		nudges = append(nudges, *nudge)
	}
	return nudges, nil
}

func defaultConsumerActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "defaultConsumerActions")
	defer span.End()
	var actions []feedlib.Action
	fns := []actionGenerator{
		defaultGetInsuranceAction,
		defaultGetTestAction,
		defaultBuyMedicineAction,
		defaultSeeDoctorAction,
	}
	for _, fn := range fns {
		action, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("error when generating action: %w", err)
		}
		actions = append(actions, *action)
	}
	return actions, nil
}

func defaultProActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "defaultProActions")
	defer span.End()
	var actions []feedlib.Action
	fns := []actionGenerator{
		defaultAddPatientAction,
		defaultSearchPatientAction,
	}
	for _, fn := range fns {
		action, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("error when generating action: %w", err)
		}
		actions = append(actions, *action)
	}
	return actions, nil
}

func defaultSeeDoctorAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getConsultationActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		common.StaticBase+"/actions/svg/see_doctor.svg",
		"See Doctor",
		"See a doctor",
		repository,
	)
}

func defaultBuyMedicineAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getMedicineActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		common.StaticBase+"/actions/svg/medicine.svg",
		"Get Medicine",
		"Get medicines",
		repository,
	)
}

func defaultGetTestAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getTestActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		common.StaticBase+"/actions/svg/get_tested.svg",
		"Get tests",
		"Get diagnostic tests",
		repository,
	)
}

func defaultGetInsuranceAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getInsuranceActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		common.StaticBase+"/actions/svg/buy_cover.svg",
		"Buy Cover",
		"Buy medical insurance",
		repository,
	)
}

func defaultSearchPatientAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		searchPatientActionName,
		feedlib.ActionTypeSecondary,
		feedlib.HandlingFullPage,
		common.StaticBase+"/actions/svg/search_user.svg",
		"Search a patient",
		"Search for a patient",
		repository,
	)
}

func defaultAddPatientAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		addPatientActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		common.StaticBase+"/actions/svg/add_user.svg",
		"Register patient",
		"Register a patient",
		repository,
	)
}

func partnerAccountSetupNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "partnerAccountSetupNudge")
	defer span.End()
	title := common.PartnerAccountSetupNudgeTitle
	text := "Create a partner account to begin transacting on Be.Well"
	imgURL := common.StaticBase + "/nudges/complete_profile.png"
	partnerAccountSetupAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		partnerAccountSetupActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"can't create %s action: %w", partnerAccountSetupActionName, err)
	}
	actions := []feedlib.Action{
		*partnerAccountSetupAction,
	}
	notificationBody := feedlib.NotificationBody{
		ResolveMessage: "Thank you for setting up your partner set up account.",
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
		notificationBody,
	)
}

func verifyEmailNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "verifyEmailNudge")
	defer span.End()
	title := common.AddPrimaryEmailNudgeTitle
	text := "Please add and verify your primary email address"
	imgURL := common.StaticBase + "/nudges/verify_email.png"
	verifyEmailAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		verifyEmailActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingFullPage,
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"can't create %s action: %w", verifyEmailActionName, err)
	}
	actions := []feedlib.Action{
		*verifyEmailAction,
	}
	notificationBody := feedlib.NotificationBody{
		ResolveMessage: "Thank you for adding your primary email address.",
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
		notificationBody,
	)
}

func createNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
	text string,
	imageURL string,
	imageTitle string,
	imageDescription string,
	actions []feedlib.Action,
	repository repository.Repository,
	notificationBody feedlib.NotificationBody,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "createNudge")
	defer span.End()
	future := time.Now().Add(time.Hour * futureHours)
	nudge := &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Visibility:     feedlib.VisibilityShow,
		Status:         feedlib.StatusPending,
		Expiry:         future,
		Title:          title,
		Text:           text,
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(
				imageURL, imageTitle, imageDescription, imageURL),
		},
		Actions:              actions,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []feedlib.Channel{},
		NotificationBody:     notificationBody,
	}
	_, err := nudge.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("nudge validation error: %w", err)
	}

	nudge, err = repository.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}
	return nudge, nil
}

func createGlobalAction(
	ctx context.Context,
	uid string,
	allowAnonymous bool,
	flavour feedlib.Flavour,
	name string,
	actionType feedlib.ActionType,
	handling feedlib.Handling,
	iconLink string,
	iconTitle string,
	iconDescription string,
	repository repository.Repository,
) (*feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "createGlobalAction")
	defer span.End()
	action := &feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		Icon: feedlib.GetSVGImageLink(
			iconLink, iconTitle, iconDescription, iconLink),
		ActionType:     actionType,
		Handling:       handling,
		AllowAnonymous: allowAnonymous,
	}
	_, err := action.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("action validation error: %w", err)
	}

	action, err = repository.SaveAction(ctx, uid, flavour, action)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save action: %w", err)
	}
	return action, nil
}

func createLocalAction(
	ctx context.Context,
	uid string,
	allowAnonymous bool,
	flavour feedlib.Flavour,
	name string,
	actionType feedlib.ActionType,
	handling feedlib.Handling,
	repository repository.Repository,
) (*feedlib.Action, error) {
	_, span := tracer.Start(ctx, "createLocalAction")
	defer span.End()
	action := &feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		Icon: feedlib.GetPNGImageLink(
			common.StaticBase+"/1px.png",
			"Blank Image",
			"Default Blank Image",
			common.StaticBase+"/1px.png",
		),
		ActionType:     actionType,
		Handling:       handling,
		AllowAnonymous: allowAnonymous,
	}
	_, err := action.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("action validation error: %w", err)
	}
	// not saved...intentionally
	// it will save embedded in a nudge or feed item

	return action, nil
}

func createFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	author string,
	tagline string,
	label string,
	iconImageURL string,
	iconTitle string,
	iconDescription string,
	summary string,
	text string,
	links []feedlib.Link,
	actions []feedlib.Action,
	conversations []feedlib.Message,
	persistent bool,
	repository repository.Repository,
) (*feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "createFeedItem")
	defer span.End()
	future := time.Now().Add(time.Hour * futureHours)
	item := &feedlib.Item{
		ID:             itemID,
		SequenceNumber: defaultSequenceNumber,
		Expiry:         future,
		Persistent:     persistent,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon: feedlib.GetPNGImageLink(
			iconImageURL, iconTitle, iconDescription, iconImageURL),
		Author:               author,
		Tagline:              tagline,
		Label:                label,
		Timestamp:            time.Now(),
		Summary:              summary,
		Text:                 text,
		TextType:             feedlib.TextTypeMarkdown,
		Links:                links,
		Actions:              actions,
		Conversations:        conversations,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []feedlib.Channel{},
	}
	_, err := item.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("item validation error: %w", err)
	}
	item, err = repository.SaveFeedItem(ctx, uid, flavour, item)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save item: %w", err)
	}
	return item, nil
}

func defaultConsumerItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "defaultConsumerItems")
	defer span.End()
	var items []feedlib.Item
	fns := []itemGenerator{
		simpleConsumerWelcome,
	}
	for _, fn := range fns {
		item, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("error when generating item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func defaultProItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "defaultProItems")
	defer span.End()
	var items []feedlib.Item
	fns := []itemGenerator{
		simpleProWelcome,
	}
	for _, fn := range fns {
		item, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("error when generating item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func simpleConsumerWelcome(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "simpleConsumerWelcome")
	defer span.End()
	persistent := true // at least one persistent message in welcome data
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to access affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos(flavour)
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	itemID := ksuid.New().String()
	conversations, err := getConsumerWelcomeThread(ctx, uid, flavour, itemID, repository)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to initialize welcome message thread: %w", err)
	}

	return createFeedItem(
		ctx,
		uid,
		flavour,
		itemID,
		defaultAuthor,
		tagline,
		common.DefaultLabel,
		common.DefaultIconPath,
		"Feed Item Icon",
		"Feed Item Icon",
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func simpleProWelcome(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) (*feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "simpleProWelcome")
	defer span.End()
	persistent := true // at least one persistent message in welcome data
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to provide affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos(flavour)
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	itemID := ksuid.New().String()
	conversations, err := getProWelcomeThread(ctx, uid, flavour, itemID, repository)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to initialize welcome message thread: %w", err)
	}

	return createFeedItem(
		ctx,
		uid,
		flavour,
		itemID,
		defaultAuthor,
		tagline,
		common.DefaultLabel,
		common.DefaultIconPath,
		"Feed Item Icon",
		"Feed Item Icon",
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func getMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	text string,
	replyTo *feedlib.Message,
	postedByName string,
	repository repository.Repository,
) (*feedlib.Message, error) {
	ctx, span := tracer.Start(ctx, "getMessage")
	defer span.End()
	msg := &feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Text:           text,
		PostedByUID:    defaultPostedByUID,
		PostedByName:   postedByName,
		Timestamp:      time.Now(),
	}
	if replyTo != nil {
		msg.ReplyTo = replyTo.ID
	}

	savedMsg, err := repository.PostMessage(ctx, uid, flavour, itemID, msg)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't save message for default welcome thread(s): %w", err)
	}

	if savedMsg == nil {
		return nil, fmt.Errorf("nil saved message")
	}

	return savedMsg, nil
}

func getConsumerWelcomeThread(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	repository repository.Repository,
) ([]feedlib.Message, error) {
	ctx, span := tracer.Start(ctx, "getConsumerWelcomeThread")
	defer span.End()
	welcome, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"Welcome to Be.Well. We are glad to meet you!",
		nil,
		"Be.Well",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	pharmacyReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the medications service. I'll ensure that you get quality and affordable medications, on time. 👋!",
		welcome,
		"Medications Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	deliveryAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the delivery assistant. I help the medications service get medicines to you on time. 👋!",
		pharmacyReply,
		"Delivery Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	dispensingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the dispensing assistant. I help your preferred pharmacy prepare your order before you go for it. 👋!",
		pharmacyReply,
		"Dispensing Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	testsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the tests service. I'll ensure that you get quality and affordable diagnostic tests. 👋!",
		welcome,
		"Tests Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	consultationsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the consultations service. I'll ensure that you can get in-person or remote(tele) advice from qualified medical professionals. 👋!",
		welcome,
		"Consultations Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	teleconsultAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the teleconsultations assistant. I'll ensure that you can reach a qualified medical professional via video or audio conference, whenever you need to. If you have an emergency, I'll help you find the nearest hospital for emergencies. 👋!",
		consultationsReply,
		"Teleconsultations Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	bookingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the booking assistant. I'll help you book appointments for your care and remind you when it's time. 👋!",
		consultationsReply,
		"Booking Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	insuranceReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the insurance service. I'll get you great quotes for medical cover and assist you when you need to use your insurance. 👋!",
		welcome,
		"Insurance Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	remindersReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the reminders service. I'll help you remember things related to your health. It could be an appointment or when you need to take some medication etc. Try me 👋!",
		welcome,
		"Reminders Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	return []feedlib.Message{
		*welcome,
		*pharmacyReply,
		*deliveryAssistant,
		*dispensingAssistant,
		*testsReply,
		*consultationsReply,
		*teleconsultAssistant,
		*bookingAssistant,
		*insuranceReply,
		*remindersReply,
	}, nil
}
func getProWelcomeThread(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	repository repository.Repository,
) ([]feedlib.Message, error) {
	ctx, span := tracer.Start(ctx, "getProWelcomeThread")
	defer span.End()
	welcome, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"Welcome to Be.Well. We are glad to meet you!",
		nil,
		"Be.Well",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	pharmacyReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the medications service. I'll help you deliver quality and affordable medications, on time. 👋!",
		welcome,
		"Medications Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	deliveryAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the delivery assistant. I help the medications service deliver medicines on time. 👋!",
		pharmacyReply,
		"Delivery Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	dispensingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the dispensing assistant. I help you prepare your orders. 👋!",
		pharmacyReply,
		"Dispensing Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	testsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the tests service. I'll help you deliver quality and affordable diagnostic tests. 👋!",
		welcome,
		"Tests Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	consultationsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the consultations service. I'll set up in-person and remote consultations for you. 👋!",
		welcome,
		"Consultations Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	teleconsultAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the teleconsultations assistant. I'll ensure that you can conduct consultations via video or audio conference, whenever you need to. If you have an emergency, I'll help you find the nearest hospital for emergencies. 👋!",
		consultationsReply,
		"Teleconsultations Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	bookingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the booking assistant. I'll help you book appointments and remind you when it's time. 👋!",
		consultationsReply,
		"Booking Assistant",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	remindersReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the reminders service. I'll help you remember things that you need to do. 👋!",
		welcome,
		"Reminders Service",
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	return []feedlib.Message{
		*welcome,
		*pharmacyReply,
		*deliveryAssistant,
		*dispensingAssistant,
		*testsReply,
		*consultationsReply,
		*teleconsultAssistant,
		*bookingAssistant,
		*remindersReply,
	}, nil
}

func getFeedWelcomeVideos(flavour feedlib.Flavour) []feedlib.Link {
	videos := []feedlib.Link{}

	switch flavour {
	case feedlib.FlavourConsumer:
		// Videos for Consumer
		consumerVideos := []feedlib.Link{
			feedlib.GetYoutubeVideoLink(
				"https://youtu.be/-mlr9rjRXmc",
				"Slade 360",
				" View your health insurance cover benefits on your Be.Well app.",
				common.StaticBase+"/items/videos/thumbs/01_lead.png",
			),
			feedlib.GetYoutubeVideoLink(
				"https://youtu.be/-iSB8yrSIps",
				"Slade 360",
				"How to add your health insurance cover to your Be.Well app.",
				common.StaticBase+"/items/videos/thumbs/01_lead.png",
			),
		}
		videos = append(videos, consumerVideos...)

	case feedlib.FlavourPro:
		// Videos for PRO

	}

	return videos
}

func feedItemsFromCMSFeedTag(ctx context.Context, flavour feedlib.Flavour) []feedlib.Item {
	ctx, span := tracer.Start(ctx, "feedItemsFromCMSFeedTag")
	defer span.End()
	// Initialize ISC clients
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)

	// Initialize new instances of the infrastructure services
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	libraryService := library.NewLibraryService(onboarding)

	items := []feedlib.Item{}

	feedPosts, err := libraryService.GetFeedContent(ctx, flavour)
	if err != nil {
		helpers.RecordSpanError(span, err)
		//  non-fatal, intentionally
		log.Printf("ERROR: unable to fetch welcome feed posts from CMS: %s", err)
	}

	// retrieve video posts from datastore and extend to the list of posts
	videosLinks := getFeedWelcomeVideos(flavour)

	future := time.Now().Add(time.Hour * futureHours)
	for _, videoLink := range videosLinks {
		tagline := ""
		summary := ""
		text := ""
		sequenceNumber := int(time.Now().Unix())

		if videoLink.URL == "https://youtu.be/-mlr9rjRXmc" {
			tagline = "See what you can do on your Be.Well app"
			summary = "See what you can do on your Be.Well app"
			text = "How to add your health insurance cover to your Be.Well app."
		}
		if videoLink.URL == "https://youtu.be/-iSB8yrSIps" {
			tagline = "Learn how to add your cover in 3 easy steps"
			summary = "Learn how to add your cover in 3 easy steps"
			text = "View your health insurance cover benefits on your Be.Well app."
			sequenceNumber = int(time.Now().Unix()) + 1
		}

		items = append(items, feedlib.Item{
			ID:             ksuid.New().String(),
			SequenceNumber: sequenceNumber,
			Expiry:         future,
			Persistent:     false,
			Status:         feedlib.StatusPending,
			Visibility:     feedlib.VisibilityShow,
			Icon:           feedlib.GetPNGImageLink(common.DefaultIconPath, "Icon", "Feed Item Icon", common.DefaultIconPath),
			Author:         defaultAuthor,
			Tagline:        tagline,
			Label:          common.DefaultLabel,
			Summary:        summary,
			Timestamp:      time.Now(),
			Text:           text,
			TextType:       feedlib.TextTypeHTML,
			Links: []feedlib.Link{
				feedlib.GetYoutubeVideoLink(
					videoLink.URL,
					videoLink.Title,
					videoLink.Description,
					videoLink.Thumbnail,
				),
			},
			Actions:              []feedlib.Action{},
			Conversations:        []feedlib.Message{},
			Users:                []string{},
			Groups:               []string{},
			NotificationChannels: []feedlib.Channel{},
		})
	}

	for _, post := range feedPosts {
		if post == nil {
			// non fatal, intentionally
			log.Printf("ERROR: nil CMS post when adding welcome posts to feed")
			continue
		}
		items = append(items, feedItemFromCMSPost(*post))
	}

	// add the slade 360 video last
	items = append(items, feedlib.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: int(time.Now().Unix()),
		Expiry:         future,
		Persistent:     false,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon:           feedlib.GetPNGImageLink(common.DefaultIconPath, "Icon", "Feed Item Icon", common.DefaultIconPath),
		Author:         defaultAuthor,
		Tagline:        "Learn what is Be.Well and how you can benefit from using it",
		Label:          common.DefaultLabel,
		Summary:        "Be.Well is a virtual and physical healthcare community.",
		Timestamp:      time.Now(),
		Text:           "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to access affordable high-quality healthcare - whether online or in person.",
		TextType:       feedlib.TextTypeHTML,
		Links: []feedlib.Link{
			feedlib.GetYoutubeVideoLink(
				"https://youtu.be/mKnlXcS3_Z0",
				"Slade 360",
				"Slade 360. HealthCare. Simplified.",
				common.StaticBase+"/items/videos/thumbs/04_slade.png",
			),
		},
		Actions:              []feedlib.Action{},
		Conversations:        []feedlib.Message{},
		Users:                []string{},
		Groups:               []string{},
		NotificationChannels: []feedlib.Channel{},
	})

	return items
}

func feedItemFromCMSPost(post library.GhostCMSPost) feedlib.Item {
	future := time.Now().Add(time.Hour * futureHours)
	return feedlib.Item{
		ID:                   post.UUID,
		SequenceNumber:       int(post.PublishedAt.Unix()),
		Expiry:               future,
		Persistent:           false,
		Status:               feedlib.StatusPending,
		Visibility:           feedlib.VisibilityShow,
		Icon:                 feedlib.GetPNGImageLink(common.DefaultIconPath, "Icon", "Feed Item Icon", common.DefaultIconPath),
		Author:               defaultAuthor,
		Tagline:              post.Slug,
		Label:                common.DefaultLabel,
		Summary:              TruncateStringWithEllipses(post.Excerpt, 140),
		Timestamp:            post.UpdatedAt,
		Text:                 post.HTML,
		TextType:             feedlib.TextTypeHTML,
		Links:                getLinks(post),
		Actions:              []feedlib.Action{},
		Conversations:        []feedlib.Message{},
		Users:                []string{},
		Groups:               []string{},
		NotificationChannels: []feedlib.Channel{},
	}
}

func getLinks(post library.GhostCMSPost) []feedlib.Link {
	featureImageLink := post.FeatureImage
	defaultLinkTitle := "CMS Item default Icon"
	if strings.HasSuffix(featureImageLink, ".png") {
		return []feedlib.Link{
			{
				ID:          ksuid.New().String(),
				URL:         featureImageLink,
				LinkType:    feedlib.LinkTypePngImage,
				Title:       defaultLinkTitle,
				Description: defaultLinkTitle,
				Thumbnail:   featureImageLink,
			},
		}
	}
	return []feedlib.Link{
		{
			ID:          ksuid.New().String(),
			URL:         common.DefaultIconPath,
			LinkType:    feedlib.LinkTypeDefault,
			Title:       defaultLinkTitle,
			Description: defaultLinkTitle,
			Thumbnail:   common.DefaultIconPath,
		},
	}
}

// TruncateStringWithEllipses truncates a string at the indicated length and adds trailing ellipses
func TruncateStringWithEllipses(str string, length int) string {
	if length <= 0 {
		return ""
	}

	targetLength := length
	addEllipses := false
	if length >= 140 {
		targetLength = length - 4 // room for ellipses for longer strings
		addEllipses = true
	}

	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= targetLength {
			break
		}
	}
	if addEllipses {
		return truncated + "..."
	}
	return truncated
}

func defaultActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	repository repository.Repository,
) ([]feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "defaultActions")
	defer span.End()
	resolveAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		common.ResolveItemActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingInline,
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to create resolve action: %w", err)
	}

	pinAction, err := createLocalAction(
		ctx,
		uid,
		true,
		flavour,
		common.PinItemActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingInline,
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to create pin action: %w", err)
	}

	hideAction, err := createLocalAction(
		ctx,
		uid,
		true,
		flavour,
		common.HideItemActionName,
		feedlib.ActionTypePrimary,
		feedlib.HandlingInline,
		repository,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to create hide action: %w", err)
	}
	actions := []feedlib.Action{
		*resolveAction,
		*pinAction,
		*hideAction,
	}

	return actions, nil
}
