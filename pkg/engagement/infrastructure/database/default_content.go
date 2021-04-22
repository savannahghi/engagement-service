package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/library"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"

	"github.com/markbates/pkger"
	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
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
	addInsuranceActionName        = "ADD_INSURANCE"
	addNHIFActionName             = "ADD_NHIF"
	partnerAccountSetupActionName = "PARTNER_ACCOUNT_SETUP"
	verifyEmailActionName         = "VERIFY_EMAIL"

	defaultOrg        = "default-org-id-please-change"
	defaultLocation   = "default-location-id-please-change"
	defaultContentDir = "/static/"
	defaultAuthor     = "Be.Well Team"
)

// embed default content assets (e.g images and documents) in the binary
var _ = pkger.Dir(defaultContentDir)

type actionGenerator func(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error)

type nudgeGenerator func(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Nudge, error)

type itemGenerator func(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Item, error)

// SetDefaultActions ensures that a feed has default actions
func SetDefaultActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Action, error) {
	actions := []base.Action{}

	switch flavour {
	case base.FlavourConsumer:
		consumerActions, err := defaultConsumerActions(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer actions: %w", err)
		}
		actions = consumerActions
	case base.FlavourPro:
		proActions, err := defaultProActions(ctx, uid, flavour, repository)
		if err != nil {
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
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Nudge, error) {
	var nudges []base.Nudge

	switch flavour {
	case base.FlavourConsumer:
		consumerNudges, err := defaultConsumerNudges(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer nudges: %w", err)
		}
		nudges = consumerNudges
	case base.FlavourPro:
		proNudges, err := defaultProNudges(ctx, uid, flavour, repository)
		if err != nil {
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
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Item, error) {
	var items []base.Item

	switch flavour {
	case base.FlavourConsumer:
		consumerItems, err := defaultConsumerItems(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer items: %w", err)
		}
		items = consumerItems
	case base.FlavourPro:
		proItems, err := defaultProItems(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default pro items: %w", err)
		}
		items = proItems
	}

	// fetch CMS items from the CMS feed tag
	cmsItems := feedItemsFromCMSFeedTag(ctx)
	for _, cmsItem := range cmsItems {
		_, err := repository.SaveFeedItem(ctx, uid, flavour, &cmsItem)
		if err != nil {
			return nil, fmt.Errorf("unable to CMS save item: %w", err)
		}

	}

	items = append(items, cmsItems...)

	return items, nil
}

func defaultConsumerNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Nudge, error) {
	var nudges []base.Nudge
	fns := []nudgeGenerator{
		addInsuranceNudge,
		verifyEmailNudge,
	}
	// TODO: return the descoped NHIF nudge
	for _, fn := range fns {
		nudge, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating nudge: %w", err)
		}
		nudges = append(nudges, *nudge)
	}
	return nudges, nil
}

func defaultProNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Nudge, error) {
	var nudges []base.Nudge
	fns := []nudgeGenerator{
		partnerAccountSetupNudge,
		verifyEmailNudge,
	}
	for _, fn := range fns {
		nudge, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating nudge: %w", err)
		}
		nudges = append(nudges, *nudge)
	}
	return nudges, nil
}

func defaultConsumerActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Action, error) {
	var actions []base.Action
	fns := []actionGenerator{
		defaultGetInsuranceAction,
		defaultGetTestAction,
		defaultBuyMedicineAction,
		defaultSeeDoctorAction,
	}
	for _, fn := range fns {
		action, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating action: %w", err)
		}
		actions = append(actions, *action)
	}
	return actions, nil
}

func defaultProActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Action, error) {
	var actions []base.Action
	fns := []actionGenerator{
		defaultAddPatientAction,
		defaultSearchPatientAction,
	}
	for _, fn := range fns {
		action, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating action: %w", err)
		}
		actions = append(actions, *action)
	}
	return actions, nil
}

func defaultSeeDoctorAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getConsultationActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		common.StaticBase+"/actions/svg/see_doctor.svg",
		"See Doctor",
		"See a doctor",
		repository,
	)
}

func defaultBuyMedicineAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getMedicineActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		common.StaticBase+"/actions/svg/medicine.svg",
		"Get Medicine",
		"Get medicines",
		repository,
	)
}

func defaultGetTestAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getTestActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		common.StaticBase+"/actions/svg/get_tested.svg",
		"Get tests",
		"Get diagnostic tests",
		repository,
	)
}

func defaultGetInsuranceAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getInsuranceActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		common.StaticBase+"/actions/svg/buy_cover.svg",
		"Buy Cover",
		"Buy medical insurance",
		repository,
	)
}

func defaultSearchPatientAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		searchPatientActionName,
		base.ActionTypeSecondary,
		base.HandlingFullPage,
		common.StaticBase+"/actions/svg/search_user.svg",
		"Search a patient",
		"Search for a patient",
		repository,
	)
}

func defaultAddPatientAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		addPatientActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		common.StaticBase+"/actions/svg/add_user.svg",
		"Register patient",
		"Register a patient",
		repository,
	)
}

func addInsuranceNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Nudge, error) {
	title := common.AddInsuranceNudgeTitle
	text := "Link your existing medical cover"
	imgURL := common.StaticBase + "/nudges/add_insurance.png"
	addInsuranceAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		addInsuranceActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", addInsuranceActionName, err)
	}
	actions := []base.Action{
		*addInsuranceAction,
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
	)
}

func partnerAccountSetupNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Nudge, error) {
	title := common.PartnerAccountSetupNudgeTitle
	text := "Create a partner account to begin transacting on Be.Well"
	imgURL := common.StaticBase + "/nudges/complete_profile.png"
	partnerAccountSetupAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		partnerAccountSetupActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", partnerAccountSetupActionName, err)
	}
	actions := []base.Action{
		*partnerAccountSetupAction,
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
	)
}

func verifyEmailNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Nudge, error) {
	title := common.AddPrimaryEmailNudgeTitle
	text := "Please add and verify your primary email address"
	imgURL := common.StaticBase + "/nudges/verify_email.png"
	verifyEmailAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		verifyEmailActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", verifyEmailActionName, err)
	}
	actions := []base.Action{
		*verifyEmailAction,
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
	)
}

func createNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	title string,
	text string,
	imageURL string,
	imageTitle string,
	imageDescription string,
	actions []base.Action,
	repository repository.Repository,
) (*base.Nudge, error) {
	future := time.Now().Add(time.Hour * futureHours)
	nudge := &base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Visibility:     base.VisibilityShow,
		Status:         base.StatusPending,
		Expiry:         future,
		Title:          title,
		Text:           text,
		Links: []base.Link{
			base.GetPNGImageLink(
				imageURL, imageTitle, imageDescription, imageURL),
		},
		Actions:              actions,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []base.Channel{},
	}
	_, err := nudge.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("nudge validation error: %w", err)
	}

	nudge, err = repository.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}
	return nudge, nil
}

func createGlobalAction(
	ctx context.Context,
	uid string,
	allowAnonymous bool,
	flavour base.Flavour,
	name string,
	actionType base.ActionType,
	handling base.Handling,
	iconLink string,
	iconTitle string,
	iconDescription string,
	repository repository.Repository,
) (*base.Action, error) {
	action := &base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		Icon: base.GetSVGImageLink(
			iconLink, iconTitle, iconDescription, iconLink),
		ActionType:     actionType,
		Handling:       handling,
		AllowAnonymous: allowAnonymous,
	}
	_, err := action.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("action validation error: %w", err)
	}

	action, err = repository.SaveAction(ctx, uid, flavour, action)
	if err != nil {
		return nil, fmt.Errorf("unable to save action: %w", err)
	}
	return action, nil
}

func createLocalAction(
	ctx context.Context,
	uid string,
	allowAnonymous bool,
	flavour base.Flavour,
	name string,
	actionType base.ActionType,
	handling base.Handling,
	repository repository.Repository,
) (*base.Action, error) {
	action := &base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		Icon: base.GetPNGImageLink(
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
		return nil, fmt.Errorf("action validation error: %w", err)
	}
	// not saved...intentionally
	// it will save embedded in a nudge or feed item

	return action, nil
}

func createFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	author string,
	tagline string,
	label string,
	iconImageURL string,
	iconTitle string,
	iconDescription string,
	summary string,
	text string,
	links []base.Link,
	actions []base.Action,
	conversations []base.Message,
	persistent bool,
	repository repository.Repository,
) (*base.Item, error) {
	future := time.Now().Add(time.Hour * futureHours)
	item := &base.Item{
		ID:             itemID,
		SequenceNumber: defaultSequenceNumber,
		Expiry:         future,
		Persistent:     persistent,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon: base.GetPNGImageLink(
			iconImageURL, iconTitle, iconDescription, iconImageURL),
		Author:               author,
		Tagline:              tagline,
		Label:                label,
		Timestamp:            time.Now(),
		Summary:              summary,
		Text:                 text,
		TextType:             base.TextTypeMarkdown,
		Links:                links,
		Actions:              actions,
		Conversations:        conversations,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []base.Channel{},
	}
	_, err := item.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("item validation error: %w", err)
	}
	item, err = repository.SaveFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to save item: %w", err)
	}
	return item, nil
}

func defaultConsumerItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Item, error) {
	var items []base.Item
	fns := []itemGenerator{
		simpleConsumerWelcome,
	}
	for _, fn := range fns {
		item, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func defaultProItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Item, error) {
	var items []base.Item
	fns := []itemGenerator{
		simpleProWelcome,
	}
	for _, fn := range fns {
		item, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func simpleConsumerWelcome(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Item, error) {
	persistent := true // at least one persistent message in welcome data
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to access affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos()
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	itemID := ksuid.New().String()
	conversations, err := getConsumerWelcomeThread(ctx, uid, flavour, itemID, repository)
	if err != nil {
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
	flavour base.Flavour,
	repository repository.Repository,
) (*base.Item, error) {
	persistent := true // at least one persistent message in welcome data
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to provide affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos()
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	itemID := ksuid.New().String()
	conversations, err := getProWelcomeThread(ctx, uid, flavour, itemID, repository)
	if err != nil {
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
	flavour base.Flavour,
	itemID string,
	text string,
	replyTo *base.Message,
	postedByName string,
	repository repository.Repository,
) (*base.Message, error) {
	msg := &base.Message{
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
	flavour base.Flavour,
	itemID string,
	repository repository.Repository,
) ([]base.Message, error) {
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
		return nil, err
	}

	return []base.Message{
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
	flavour base.Flavour,
	itemID string,
	repository repository.Repository,
) ([]base.Message, error) {
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
		return nil, err
	}
	if err != nil {
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
		return nil, err
	}

	return []base.Message{
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

func getFeedWelcomeVideos() []base.Link {
	return []base.Link{
		base.GetYoutubeVideoLink(
			"https://youtu.be/mKnlXcS3_Z0",
			"Slade 360",
			"Slade 360. HealthCare. Simplified.",
			common.StaticBase+"/items/videos/thumbs/04_slade.png",
		),
	}
}

func feedItemsFromCMSFeedTag(ctx context.Context) []base.Item {
	libraryService := library.NewLibraryService()
	items := []base.Item{}
	feedPosts, err := libraryService.GetFeedContent(ctx)
	if err != nil {
		//  non-fatal, intentionally
		log.Printf("ERROR: unable to fetch welcome feed posts from CMS: %s", err)
	}
	for _, post := range feedPosts {
		if post == nil {
			// non fatal, intentionally
			log.Printf("ERROR: nil CMS post when adding welcome posts to feed")
			continue
		}
		items = append(items, feedItemFromCMSPost(*post))
	}
	return items
}

func feedItemFromCMSPost(post library.GhostCMSPost) base.Item {
	future := time.Now().Add(time.Hour * futureHours)
	return base.Item{
		ID:                   post.UUID,
		SequenceNumber:       int(post.UpdatedAt.Unix()),
		Expiry:               future,
		Persistent:           false,
		Status:               base.StatusPending,
		Visibility:           base.VisibilityShow,
		Icon:                 base.GetPNGImageLink(common.DefaultIconPath, "Icon", "Feed Item Icon", common.DefaultIconPath),
		Author:               defaultAuthor,
		Tagline:              post.Slug,
		Label:                common.DefaultLabel,
		Summary:              TruncateStringWithEllipses(post.Excerpt, 140),
		Timestamp:            post.UpdatedAt,
		Text:                 post.HTML,
		TextType:             base.TextTypeHTML,
		Links:                getLinks(post),
		Actions:              []base.Action{},
		Conversations:        []base.Message{},
		Users:                []string{},
		Groups:               []string{},
		NotificationChannels: []base.Channel{},
	}
}

func getLinks(post library.GhostCMSPost) []base.Link {
	featureImageLink := post.FeatureImage
	defaultLinkTitle := "CMS Item default Icon"
	if strings.HasSuffix(featureImageLink, ".png") {
		return []base.Link{
			{
				ID:          ksuid.New().String(),
				URL:         featureImageLink,
				LinkType:    base.LinkTypePngImage,
				Title:       defaultLinkTitle,
				Description: defaultLinkTitle,
				Thumbnail:   featureImageLink,
			},
		}
	}
	return []base.Link{
		{
			ID:          ksuid.New().String(),
			URL:         common.DefaultIconPath,
			LinkType:    base.LinkTypeDefault,
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
	flavour base.Flavour,
	repository repository.Repository,
) ([]base.Action, error) {
	resolveAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		common.ResolveItemActionName,
		base.ActionTypePrimary,
		base.HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create resolve action: %w", err)
	}

	pinAction, err := createLocalAction(
		ctx,
		uid,
		true,
		flavour,
		common.PinItemActionName,
		base.ActionTypePrimary,
		base.HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create pin action: %w", err)
	}

	hideAction, err := createLocalAction(
		ctx,
		uid,
		true,
		flavour,
		common.HideItemActionName,
		base.ActionTypePrimary,
		base.HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create hide action: %w", err)
	}
	actions := []base.Action{
		*resolveAction,
		*pinAction,
		*hideAction,
	}

	return actions, nil
}
