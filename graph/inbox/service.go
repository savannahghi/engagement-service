package inbox

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// inbox service constants
const (
	userInboxCollectionName = "user_inbox"
	welcomeMessageBody      = `
	Hello there, welcome to Be.Well. We trust you will find something useful here.
	Happy staying healthy.
	`
	welcomeMessageTagSlug = "welcome-message"
	welcomeMessageTagName = "Welcome Message"

	defaultChannelSlug = "direct-messages"
	defaultChannelName = "Direct Messages"

	systemSenderName = "Bewell Team"
	systemUID        = "e47c6e5b-ea35-4b4e-a266-c5e91294d070"
)

// NewService returns a new authentication service
func NewService() *Service {
	fc := &base.FirebaseClient{}
	ctx := context.Background()

	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("can't initialize Firebase app when setting up inbox service: %s", err)
	}

	auth, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up inbox service: %s", err)
	}

	firestore, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("can't initialize Firestore client when setting up inbox service: %s", err)
	}

	return &Service{
		firestoreClient: firestore,
		firebaseAuth:    auth,
		collectionName:  base.SuffixCollection(userInboxCollectionName),
	}
}

// Service organizes inbox functionality
type Service struct {
	firestoreClient *firestore.Client
	firebaseAuth    *auth.Client
	collectionName  string
}

func (s Service) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("inbox service does not have an initialized firestoreClient")
	}

	if s.firebaseAuth == nil {
		log.Panicf("inbox service does not have an initialized firebaseAuth")
	}

	if s.collectionName == "" {
		log.Panicf("collection name not set on inbox service")
	}
}

func (s Service) getLoggedInUserUID(ctx context.Context) (*string, error) {
	authToken, err := base.GetUserTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth token not found in context: %w", err)
	}
	return &authToken.UID, nil
}

// GetUserMessages fetches the logged in user's messages
func (s Service) GetUserMessages(ctx context.Context) ([]*Message, error) {
	s.checkPreconditions()

	uid, err := s.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	collection := s.firestoreClient.Collection(s.collectionName)
	query := collection.Where("recipientUID", "==", *uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get inbox documents: %w", err)
	}

	var messages []*Message
	for _, dsnap := range docs {
		m := &Message{}
		err = dsnap.DataTo(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read user messages: %w", err)
		}
		messages = append(messages, m)
	}

	if len(messages) == 0 {
		success, err := s.createLoggedInUserMessage(ctx, systemSenderName, systemUID,
			welcomeMessageBody,
			MessageChannel{
				ID:   ksuid.New().String(),
				Name: defaultChannelName,
				Slug: defaultChannelSlug,
			},
			[]MessageTag{
				{
					ID:   ksuid.New().String(),
					Name: welcomeMessageTagName,
					Slug: welcomeMessageTagSlug,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("can't set up welcome messages: %w", err)
		}

		if !success {
			return nil, fmt.Errorf("not successful in setting up welcome messages: %w", err)
		}

		// potentially dangerous recursion
		return s.GetUserMessages(ctx)
	}

	return messages, nil
}

func (s Service) createLoggedInUserMessage(
	ctx context.Context,
	senderName string,
	senderUID string,
	body string,
	channel MessageChannel,
	tags []MessageTag,
) (bool, error) {
	s.checkPreconditions()

	uid, err := s.getLoggedInUserUID(ctx)
	if err != nil {
		return false, err
	}

	newMessage := &Message{
		ID:            ksuid.New().String(),
		SenderName:    senderName,
		SenderUID:     senderUID,
		CreatedAt:     time.Now(),
		RecipientName: "Me", // once IPC comes, replace this
		RecipientUID:  *uid,
		Body:          body,
		Channel:       channel,
		Tags:          tags,
	}

	_, err = base.SaveDataToFirestore(s.firestoreClient, s.collectionName, newMessage)
	if err != nil {
		return false, fmt.Errorf("unable to save user welcome message: %v", err)
	}

	return true, nil
}
