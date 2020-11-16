package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/markbates/pkger"
	"github.com/segmentio/ksuid"
)

const (
	defaultSequenceNumber = 1
	defaultPostedByUID    = "hOcaUv8dqqgmWYf9HEhjdudgf0b2"
	futureHours           = 878400 // hours in a century of leap years...

	getConsultationActionName = "GET_CONSULTATION"
	getMedicineActionName     = "GET_MEDICINE"
	getTestActionName         = "GET_TEST"
	getInsuranceActionName    = "GET_INSURANCE"
	getCoachingActionName     = "GET_COACHING"
	findPatientActionName     = "FIND_PATIENT"
	addInsuranceActionName    = "ADD_INSURANCE"
	addNHIFActionName         = "ADD_NHIF"
	completeProfileActionName = "COMPLETE_PROFILE"
	completeKYCActionName     = "COMPLETE_KYC"
	hideItemActionName        = "HIDE_ITEM"
	pinItemActionName         = "PIN_ITEM"
	resolveItemActionName     = "RESOLVE_ITEM"

	defaultOrg        = "default-org-id-please-change"
	defaultLocation   = "default-location-id-please-change"
	defaultContentDir = "/graph/feed/static"
	defaultAuthor     = "Be.Well Team"
	defaultLabel      = "WELCOME"
	staticBase        = "https://assets.healthcloud.co.ke"
	defaultIconPath   = staticBase + "/bewell_logo.png"
)

// embed default content assets (e.g images and documents) in the binary
var _ = pkger.Dir(defaultContentDir)

type actionGenerator func(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error)

type nudgeGenerator func(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Nudge, error)

type itemGenerator func(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Item, error)

// SetDefaultActions ensures that a feed has default actions
func SetDefaultActions(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) ([]Action, error) {
	actions := []Action{}

	switch flavour {
	case FlavourConsumer:
		consumerActions, err := defaultConsumerActions(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer actions: %w", err)
		}
		actions = consumerActions
	case FlavourPro:
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
	flavour Flavour,
	repository Repository,
) ([]Nudge, error) {
	var nudges []Nudge

	switch flavour {
	case FlavourConsumer:
		consumerNudges, err := defaultConsumerNudges(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer nudges: %w", err)
		}
		nudges = consumerNudges
	case FlavourPro:
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
	flavour Flavour,
	repository Repository,
) ([]Item, error) {
	var items []Item

	switch flavour {
	case FlavourConsumer:
		consumerItems, err := defaultConsumerItems(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer items: %w", err)
		}
		items = consumerItems
	case FlavourPro:
		proItems, err := defaultProItems(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default pro items: %w", err)
		}
		items = proItems
	}

	return items, nil
}

func defaultConsumerNudges(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) ([]Nudge, error) {
	var nudges []Nudge
	fns := []nudgeGenerator{
		addInsuranceNudge,
		addNHIFNudge,
		completeProfileNudge,
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

func defaultProNudges(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) ([]Nudge, error) {
	var nudges []Nudge
	fns := []nudgeGenerator{
		completeKYCNudge,
		completeProfileNudge,
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
	flavour Flavour,
	repository Repository,
) ([]Action, error) {
	var actions []Action
	fns := []actionGenerator{
		defaultSeeDoctorAction,
		defaultBuyMedicineAction,
		defaultGetTestAction,
		defaultGetInsuranceAction,
		defaultCoachingAction,
		defaultHelpAction,
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
	flavour Flavour,
	repository Repository,
) ([]Action, error) {
	var actions []Action
	fns := []actionGenerator{
		defaultFindPatientAction,
		defaultHelpAction,
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
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		getConsultationActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
}

func defaultBuyMedicineAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		getMedicineActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
}

func defaultGetTestAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		getMedicineActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
}

func defaultGetInsuranceAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		getInsuranceActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
}

func defaultCoachingAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		getCoachingActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
}

func defaultHelpAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		getCoachingActionName,
		ActionTypeFloating,
		HandlingFullPage,
		repository,
	)
}

func defaultFindPatientAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		flavour,
		findPatientActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
}

func addInsuranceNudge(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Nudge, error) {
	title := "Add Insurance"
	text := "Link your existing medical cover"
	imgURL := staticBase + "/nudges/add_insurance.png"
	addInsuranceAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		addInsuranceActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", addInsuranceActionName, err)
	}
	actions := []Action{
		*addInsuranceAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		actions,
		repository,
	)
}

func addNHIFNudge(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Nudge, error) {
	title := "Add NHIF"
	text := "Link your NHIF cover"
	imgURL := staticBase + "/nudges/add_insurance.png"
	addNHIFAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		addNHIFActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", addNHIFActionName, err)
	}
	actions := []Action{
		*addNHIFAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		actions,
		repository,
	)
}

func completeProfileNudge(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Nudge, error) {
	title := "Complete your profile"
	text := "Fill in your Be.Well profile to unlock more rewards"
	imgURL := staticBase + "/nudges/complete_profile.png"
	completeProfileAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		completeProfileActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", completeProfileActionName, err)
	}
	actions := []Action{
		*completeProfileAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		actions,
		repository,
	)
}

func completeKYCNudge(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Nudge, error) {
	title := "Complete your business profile"
	text := "Fill in your Be.Well usiness profile in order to start transacting"
	imgURL := staticBase + "/nudges/complete_kyc.png"
	completeKYCAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		completeKYCActionName,
		ActionTypePrimary,
		HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", completeKYCActionName, err)
	}
	actions := []Action{
		*completeKYCAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		actions,
		repository,
	)
}

func createNudge(
	ctx context.Context,
	uid string,
	flavour Flavour,
	title string,
	text string,
	imageURL string,
	actions []Action,
	repository Repository,
) (*Nudge, error) {
	future := time.Now().Add(time.Hour * futureHours)
	nudge := &Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Visibility:     VisibilityShow,
		Status:         StatusPending,
		Expiry:         future,
		Title:          title,
		Text:           text,
		Links: []Link{
			GetPNGImageLink(imageURL),
		},
		Actions:              actions,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []Channel{},
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
	flavour Flavour,
	name string,
	actionType ActionType,
	handling Handling,
	repository Repository,
) (*Action, error) {
	action := &Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		ActionType:     actionType,
		Handling:       handling,
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
	flavour Flavour,
	name string,
	actionType ActionType,
	handling Handling,
	repository Repository,
) (*Action, error) {
	action := &Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		ActionType:     actionType,
		Handling:       handling,
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
	flavour Flavour,
	author string,
	tagline string,
	label string,
	iconImageURL string,
	summary string,
	text string,
	links []Link,
	actions []Action,
	conversations []Message,
	persistent bool,
	repository Repository,
) (*Item, error) {
	future := time.Now().Add(time.Hour * futureHours)
	item := &Item{
		ID:                   ksuid.New().String(),
		SequenceNumber:       defaultSequenceNumber,
		Expiry:               future,
		Persistent:           persistent,
		Status:               StatusPending,
		Visibility:           VisibilityShow,
		Icon:                 GetPNGImageLink(iconImageURL),
		Author:               author,
		Tagline:              tagline,
		Label:                label,
		Timestamp:            time.Now(),
		Summary:              summary,
		Text:                 text,
		Links:                links,
		Actions:              actions,
		Conversations:        conversations,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []Channel{},
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
	flavour Flavour,
	repository Repository,
) ([]Item, error) {
	var items []Item
	fns := []itemGenerator{
		ultimateComposite,
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
	flavour Flavour,
	repository Repository,
) ([]Item, error) {
	var items []Item
	fns := []itemGenerator{
		ultimateComposite,
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
	flavour Flavour,
	repository Repository,
) (*Item, error) {
	persistent := false
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to access affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos()
	actions := []Action{}
	conversations := getConsumerWelcomeThread()
	return createFeedItem(
		ctx,
		uid,
		flavour,
		defaultAuthor,
		tagline,
		defaultLabel,
		defaultIconPath,
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
	flavour Flavour,
	repository Repository,
) (*Item, error) {
	persistent := false
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to provide affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos()
	actions := []Action{}
	conversations := getProWelcomeThread()
	return createFeedItem(
		ctx,
		uid,
		flavour,
		defaultAuthor,
		tagline,
		defaultLabel,
		defaultIconPath,
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func ultimateComposite(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Item, error) {
	// here's what Be.Well can do for you... / help you do for your patients
	persistent := false
	tagline := "This is Be.Well..."
	summary := "This is Be.Well..."
	text := "This is Be.Well..."
	links := []Link{
		GetPNGImageLink(staticBase + "/items/images/bewell_banner01.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner02.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner03.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner04.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner05.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner06.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner07.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner08.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner09.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner10.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner11.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner12.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner13.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner14.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner15.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner16.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner17.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner18.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner19.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner20.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner21.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner22.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner23.png"),
		GetPNGImageLink(staticBase + "/items/images/bewell_banner24.png"),

		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_25.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_26.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_27.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_28.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_29.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_30.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_31.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_32.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_33.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_34.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_35.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_36.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_37.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_38.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_39.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_40.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_41.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_42.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_43.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_44.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_45.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_46.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_47.pdf"),
		GetPDFDocumentLink(staticBase + "/items/documents/bewell_banner_48.pdf"),
	}
	resolveAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		resolveItemActionName,
		ActionTypePrimary,
		HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create resolve action: %w", err)
	}

	pinAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		pinItemActionName,
		ActionTypePrimary,
		HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create pin action: %w", err)
	}

	hideAction, err := createLocalAction(
		ctx,
		uid,
		flavour,
		hideItemActionName,
		ActionTypePrimary,
		HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create hide action: %w", err)
	}

	actions := []Action{
		*resolveAction,
		*pinAction,
		*hideAction,
	}
	conversations := []Message{}
	return createFeedItem(
		ctx,
		uid,
		flavour,
		defaultAuthor,
		tagline,
		defaultLabel,
		defaultIconPath,
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
	text string,
	replyTo *Message,
	postedByName string,
) Message {
	msg := Message{
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
	return msg
}

func getConsumerWelcomeThread() []Message {
	welcome := getMessage(
		"Welcome to Be.Well. We are glad to meet you!",
		nil,
		"Be.Well",
	)
	pharmacyReply := getMessage(
		"I'm the medications service. I'll ensure that you get quality and affordable medications, on time. 👋!",
		&welcome,
		"Medications Service",
	)

	deliveryAssistant := getMessage(
		"I'm the delivery assistant. I help the medications service get medicines to you on time. 👋!",
		&pharmacyReply,
		"Delivery Assistant",
	)

	dispensingAssistant := getMessage(
		"I'm the dispensing assistant. I help your preferred pharmacy prepare your order before you go for it. 👋!",
		&pharmacyReply,
		"Dispensing Assistant",
	)

	testsReply := getMessage(
		"I'm the tests service. I'll ensure that you get quality and affordable diagnostic tests. 👋!",
		&welcome,
		"Tests Service",
	)

	consultationsReply := getMessage(
		"I'm the consultations service. I'll ensure that you can get in-person or remote(tele) advice from qualified medical professionals. 👋!",
		&welcome,
		"Consultations Service",
	)

	teleconsultAssistant := getMessage(
		"I'm the teleconsultations assistant. I'll ensure that you can reach a qualified medical professional via video or audio conference, whenever you need to. If you have an emergency, I'll help you find the nearest hospital for emergencies. 👋!",
		&consultationsReply,
		"Teleconsultations Assistant",
	)

	bookingAssistant := getMessage(
		"I'm the booking assistant. I'll help you book appointments for your care and remind you when it's time. 👋!",
		&consultationsReply,
		"Booking Assistant",
	)

	coachingReply := getMessage(
		"I'm the coaching service. I'll link you up to *awesome* wellness and fitness coaches. 👋!",
		&welcome,
		"Coaching Service",
	)

	insuranceReply := getMessage(
		"I'm the insurance service. I'll get you great quotes for medical cover and assist you when you need to use your insurance. 👋!",
		&welcome,
		"Coaching Service",
	)

	remindersReply := getMessage(
		"I'm the reminders service. I'll help you remember things related to your health. It could be an appointment or when you need to take some medication etc. Try me 👋!",
		&welcome,
		"Reminders Service",
	)

	return []Message{
		welcome,
		pharmacyReply,
		deliveryAssistant,
		dispensingAssistant,
		testsReply,
		consultationsReply,
		teleconsultAssistant,
		bookingAssistant,
		coachingReply,
		insuranceReply,
		remindersReply,
	}
}
func getProWelcomeThread() []Message {
	welcome := getMessage(
		"Welcome to Be.Well. We are glad to meet you!",
		nil,
		"Be.Well",
	)
	pharmacyReply := getMessage(
		"I'm the medications service. I'll help you deliver quality and affordable medications, on time. 👋!",
		&welcome,
		"Medications Service",
	)

	deliveryAssistant := getMessage(
		"I'm the delivery assistant. I help the medications service deliver medicines on time. 👋!",
		&pharmacyReply,
		"Delivery Assistant",
	)

	dispensingAssistant := getMessage(
		"I'm the dispensing assistant. I help you prepare your orders. 👋!",
		&pharmacyReply,
		"Dispensing Assistant",
	)

	testsReply := getMessage(
		"I'm the tests service. I'll help you deliver quality and affordable diagnostic tests. 👋!",
		&welcome,
		"Tests Service",
	)

	consultationsReply := getMessage(
		"I'm the consultations service. I'll set up in-person and remote consultations for you. 👋!",
		&welcome,
		"Consultations Service",
	)

	teleconsultAssistant := getMessage(
		"I'm the teleconsultations assistant. I'll ensure that you can conduct consultations via video or audio conference, whenever you need to. If you have an emergency, I'll help you find the nearest hospital for emergencies. 👋!",
		&consultationsReply,
		"Teleconsultations Assistant",
	)

	bookingAssistant := getMessage(
		"I'm the booking assistant. I'll help you book appointments and remind you when it's time. 👋!",
		&consultationsReply,
		"Booking Assistant",
	)

	coachingReply := getMessage(
		"I'm the coaching service. I'll help you deliver your *awesome* coaching services to clients. 👋!",
		&welcome,
		"Coaching Service",
	)

	remindersReply := getMessage(
		"I'm the reminders service. I'll help you remember things that you need to do. 👋!",
		&welcome,
		"Reminders Service",
	)

	return []Message{
		welcome,
		pharmacyReply,
		deliveryAssistant,
		dispensingAssistant,
		testsReply,
		consultationsReply,
		teleconsultAssistant,
		bookingAssistant,
		coachingReply,
		remindersReply,
	}
}

func getFeedWelcomeVideos() []Link {
	return []Link{
		GetYoutubeVideoLink("https://www.youtube.com/watch?v=gcv2Z2AdpjM"),
		GetYoutubeVideoLink("https://www.youtube.com/watch?v=W_daZjDET9Q"),
		GetYoutubeVideoLink("https://www.youtube.com/watch?v=IbtVBXNvpSA"),
		GetYoutubeVideoLink("https://www.youtube.com/watch?v=mKnlXcS3_Z0"),
	}
}
