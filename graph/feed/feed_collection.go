package feed

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/base"
)

// NewCollection initializes a user feed that is backed by an in-memory map
func NewCollection(
	repository Repository,
	notificationService NotificationService,
) (*Collection, error) {
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		return nil, fmt.Errorf("can't load inter-service config from YAML: %w", err)
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
	if err != nil {
		return nil, fmt.Errorf(
			"can't set up profile interservice client: %w", err)
	}

	feedCollection := &Collection{
		repository:          repository,
		notificationService: notificationService,
		profileClient:       profileClient,
	}

	err = feedCollection.checkPreconditions()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize feeds: %w", err)
	}
	return feedCollection, nil
}

// Collection organizes the top level methods for interaction with feeds
type Collection struct {
	repository          Repository
	notificationService NotificationService
	profileClient       *base.InterServiceClient
}

func (agg Collection) checkPreconditions() error {
	if agg.repository == nil {
		return fmt.Errorf(
			"incorrectly initialized feed aggregate: nil repository")
	}

	if agg.notificationService == nil {
		return fmt.Errorf(
			"incorrectly initialized feed aggregate: nil notification service")
	}

	if agg.profileClient == nil {
		return fmt.Errorf(
			"incorrectly initialized feed aggregate: nil profile client")
	}

	return nil
}

// GetFeed retrieves a feed
func (agg Collection) GetFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour base.Flavour,
	persistent base.BooleanFilter,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
	filterParams *FilterParams,
) (*Feed, error) {
	if err := agg.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("precondition check failed: %w", err)
	}
	feed, err := agg.repository.GetFeed(
		ctx,
		uid,
		isAnonymous,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		return nil, fmt.Errorf("feed retrieval error: %w", err)
	}

	// set the ID (computed, not stored)
	feed.ID = feed.getID()

	// inject the repository and notification service into the returned feed
	feed.repository = agg.repository
	feed.notificationService = agg.notificationService
	feed.SequenceNumber = int(time.Now().Unix())

	if err := feed.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"the retrieved feed failed precondition checks: %w", err)
	}

	return feed, nil
}

// GetThinFeed gets a feed with only the UID, flavour and dependencies
// filled in.
//
// It is used for efficient instantiation of feeds by code that does not need
// the full detail.
func (agg Collection) GetThinFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour base.Flavour,
) (*Feed, error) {
	if err := agg.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("precondition check failed: %w", err)
	}
	feed := &Feed{
		UID:         *uid,
		Flavour:     flavour,
		Actions:     []base.Action{},
		Items:       []base.Item{},
		Nudges:      []base.Nudge{},
		IsAnonymous: isAnonymous,
	}

	// set the ID (computed, not stored)
	feed.ID = feed.getID()

	// inject the repository and notification service into the returned feed
	feed.repository = agg.repository
	feed.notificationService = agg.notificationService
	feed.SequenceNumber = int(time.Now().Unix())

	if err := feed.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"the retrieved feed failed precondition checks: %w", err)
	}

	return feed, nil
}
