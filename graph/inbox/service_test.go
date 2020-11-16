package inbox

import (
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("%v_%v", "testing", time.Now().UnixNano()))
	os.Exit(m.Run())
}

func GetFirestoreClient(t *testing.T) *firestore.Client {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, firestoreClient)
	return firestoreClient
}

func GetFirebaseAuthClient(t *testing.T) (*auth.Client, error) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Firebase app: %w", err)
	}
	ctx := base.GetAuthenticatedContext(t)
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Firebase auth client: %w", err)
	}
	return client, nil
}

func TestNewService(t *testing.T) {
	s := NewService()
	s.checkPreconditions() // should not panic
}

func TestCollectionSuffix(t *testing.T) {
	col := base.SuffixCollection(userInboxCollectionName)
	assert.Contains(t, col, "testing")
}

func TestCheckIfUserIsLoggedIn(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	s := NewService()
	uid, _ := s.getLoggedInUserUID(ctx)
	assert.NotNil(t, uid)
}

func TestSendWelcomeMessage(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	s := NewService()

	uid, _ := s.getLoggedInUserUID(ctx)

	_, err := s.SendWelcomeMessageToUser(ctx)
	if err != nil {
		t.Fatalf("unable send welcome message: %v", err)
	}

	// retrieve the message
	msgs, err := s.GetUserMessages(ctx)
	if err != nil {
		t.Fatalf("unable to get user message: %v", err)
	}

	thisUserMsgs := []bool{}
	for _, msg := range msgs {
		if msg.RecipientUID == *uid {
			thisUserMsgs = append(thisUserMsgs, true)
		}
	}

	assert.Equal(t, len(msgs), len(thisUserMsgs))

	wlMsg := []*Message{}

	for _, msg := range msgs {
	INNER:
		for _, tag := range msg.Tags {
			if tag.Slug == welcomeMessageTagSlug {
				wlMsg = append(wlMsg, msg)
				break INNER
			}
		}
	}

	assert.Equal(t, len(wlMsg), 1)

	assert.Equal(t, wlMsg[0].Body, welcomeMessageBody)

}

func TestFetchmessage(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	s := NewService()
	en, _ := base.EncryptMessage("test body")

	uid, _ := s.getLoggedInUserUID(ctx)

	UUID := uuid.New()
	_uuid := UUID.String()
	newMessage := &Message{
		ID:            _uuid,
		SenderName:    "sytem",
		SenderUID:     "systemuuid",
		CreatedAt:     time.Now(),
		RecipientName: "test user",
		RecipientUID:  *uid,
		Body:          *en,
		Channel: MessageChannel{
			ID:   _uuid,
			Name: "test",
		},
		Tags: []MessageTag{
			{
				ID:   _uuid,
				Name: "test",
			},
		},
	}

	_, err := base.SaveDataToFirestore(s.firestoreClient, base.SuffixCollection(userInboxCollectionName), newMessage)
	if err != nil {
		t.Fatalf("unable to save user message: %v", err)
	}

	// retrieve
	msgs, err := s.GetUserMessages(ctx)
	if err != nil {
		t.Fatalf("unable to get user message: %v", err)
	}

	assert.GreaterOrEqual(t, len(msgs), 1)

}

func TestEncryptDecrypt(t *testing.T) {

	message1 := "test message 1"

	enc, err := base.EncryptMessage(message1)
	if err != nil {
		t.Fatalf("unable to encrypt message: %v", err)
	}

	assert.NotEqual(t, message1, *enc)

	dec, err := base.DecrypMessage(*enc)
	if err != nil {
		t.Fatalf("unable to decrypt message: %v", err)
	}

	assert.NotEqual(t, *enc, *dec)
	assert.Equal(t, message1, *dec)
}
