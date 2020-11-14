package feed

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/markbates/pkger"
	"github.com/segmentio/ksuid"
)

const (
	defaultSequenceNumber = 1
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

	defaultOrg        = "default-org-id-please-change"
	defaultLocation   = "default-location-id-please-change"
	defaultContentDir = "/graph/feed/content"
	defaultIconPath   = "/graph/feed/content/bewell_logo.png"
	defaultAuthor     = "Be.Well Team"
	defaultLabel      = "WELCOME"
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

func defaultConsumerItems(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) ([]Item, error) {
	var items []Item
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
	flavour Flavour,
	repository Repository,
) ([]Item, error) {
	var items []Item
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
	flavour Flavour,
	repository Repository,
) (*Item, error) {
	persistent := false
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to access affordable high-quality healthcare - whether online or in person."
	images := []Image{}
	documents := []Document{}
	videos := []Video{}
	actions := []Action{}
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
		images,
		documents,
		videos,
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
	images := []Image{}
	documents := []Document{}
	videos := []Video{}
	actions := []Action{}
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
		images,
		documents,
		videos,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func defaultSeeDoctorAction(
	ctx context.Context,
	uid string,
	flavour Flavour,
	repository Repository,
) (*Action, error) {
	return createAction(
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
	return createAction(
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
	return createAction(
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
	return createAction(
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
	return createAction(
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
	return createAction(
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
	return createAction(
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
	pkgerImgPath := "/graph/feed/content/nudges/add_insurance.png"
	addInsuranceAction, err := createAction(
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
		pkgerImgPath,
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
	pkgerImgPath := "/graph/feed/content/nudges/add_insurance.png"
	addNHIFAction, err := createAction(
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
		pkgerImgPath,
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
	pkgerImgPath := "/graph/feed/content/nudges/complete_profile.png"
	completeProfileAction, err := createAction(
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
		pkgerImgPath,
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
	pkgerImgPath := "/graph/feed/content/nudges/complete_kyc.png"
	completeKYCAction, err := createAction(
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
		pkgerImgPath,
		actions,
		repository,
	)
}

func imagePathToBase64(path string) (string, error) {
	img, err := pkger.Open(path)
	if err != nil {
		return "", fmt.Errorf("can't open pkger image path: %w", err)
	}
	defer img.Close()

	imgBytes, err := ioutil.ReadAll(img)
	if err != nil {
		return "", fmt.Errorf("can't read image: %w", err)
	}
	return base64.StdEncoding.EncodeToString(imgBytes), nil
}

func createNudge(
	ctx context.Context,
	uid string,
	flavour Flavour,
	title string,
	text string,
	pkgerImagePath string,
	actions []Action,
	repository Repository,
) (*Nudge, error) {
	b64, err := imagePathToBase64(pkgerImagePath)
	if err != nil {
		return nil, fmt.Errorf("unable to load nudge image: %w", err)
	}

	future := time.Now().Add(time.Hour * futureHours)
	nudge := &Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Visibility:     VisibilityShow,
		Status:         StatusPending,
		Expiry:         future,
		Title:          title,
		Text:           text,
		Image: Image{
			ID:     ksuid.New().String(),
			Base64: b64,
		},
		Actions:              actions,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []Channel{},
	}
	_, err = nudge.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("nudge validation error: %w", err)
	}

	nudge, err = repository.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}
	return nudge, nil
}

func createAction(
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

func createFeedItem(
	ctx context.Context,
	uid string,
	flavour Flavour,
	author string,
	tagline string,
	label string,
	iconImagePath string,
	summary string,
	text string,
	images []Image,
	documents []Document,
	videos []Video,
	actions []Action,
	conversations []Message,
	persistent bool,
	repository Repository,
) (*Item, error) {
	b64, err := imagePathToBase64(iconImagePath)
	if err != nil {
		return nil, fmt.Errorf("unable to load nudge image: %w", err)
	}

	future := time.Now().Add(time.Hour * futureHours)
	item := &Item{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Expiry:         future,
		Persistent:     persistent,
		Status:         StatusPending,
		Visibility:     VisibilityShow,
		Icon: Image{
			ID:     ksuid.New().String(),
			Base64: b64,
		},
		Author:               author,
		Tagline:              tagline,
		Label:                label,
		Timestamp:            time.Now(),
		Summary:              summary,
		Text:                 text,
		Images:               images,
		Documents:            documents,
		Videos:               videos,
		Actions:              actions,
		Conversations:        conversations,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []Channel{},
	}
	_, err = item.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("item validation error: %w", err)
	}
	item, err = repository.SaveFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to save item: %w", err)
	}
	return item, nil
}
