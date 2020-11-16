package inbox

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
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

// GetUserMessages ...
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
		return nil, err
	}

	var messages []*Message

	for _, dsnap := range docs {
		m := &Message{}
		err = dsnap.DataTo(m)
		if err != nil {
			log.Errorf("unable to read user messges: %v", err)
			break
		}

		d, err := base.DecrypMessage(m.Body)
		if err != nil {
			log.Errorf("unable to read user messges: %v", err)
			break
		}

		m.Body = *d

		messages = append(messages, m)
	}

	log.Println(messages)

	return messages, nil
}

// SendWelcomeMessageToUser adds a welcome message to the users inbox.
// if a similar exists a success is returned
func (s Service) SendWelcomeMessageToUser(ctx context.Context) (*bool, error) {
	s.checkPreconditions()

	resp := true

	userMessages, err := s.GetUserMessages(ctx)
	if err != nil {
		return nil, err
	}

	hasWelcomeMessage := false

	for _, msg := range userMessages {
		for _, tag := range msg.Tags {
			if tag.Slug == welcomeMessageTagSlug {
				hasWelcomeMessage = true
				break
			}
		}
	}

	if !hasWelcomeMessage {
		// create welcome message
		UUID := uuid.New()
		_uuid := UUID.String()
		return s.createLoggedInUserMessage(ctx, systemSenderName, systemUID,
			welcomeMessageBody,
			MessageChannel{
				ID:   _uuid,
				Name: defaultChannelName,
				Slug: defaultChannelSlug,
			},
			[]MessageTag{
				{
					ID:   _uuid,
					Name: welcomeMessageTagName,
					Slug: welcomeMessageTagSlug,
				},
			},
		)
	}

	return &resp, nil
}

func (s Service) createLoggedInUserMessage(ctx context.Context, senderName string, senderUID string,
	body string, channel MessageChannel, tags []MessageTag) (*bool, error) {

	resp := false
	s.checkPreconditions()

	uid, err := s.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	en, err := base.EncryptMessage(body)
	if err != nil {
		return nil, err
	}

	UUID := uuid.New()
	_uuid := UUID.String()
	newMessage := &Message{
		ID:            _uuid,
		SenderName:    senderName,
		SenderUID:     senderUID,
		CreatedAt:     time.Now(),
		RecipientName: "Me", // once IPc comes, replace this
		RecipientUID:  *uid,
		Body:          *en,
		Channel:       channel,
		Tags:          tags,
	}

	_, errS := base.SaveDataToFirestore(s.firestoreClient, s.collectionName, newMessage)
	if errS != nil {
		return &resp, fmt.Errorf("unable to save user welcome message: %v", err)
	}

	resp = true
	return &resp, nil
}
