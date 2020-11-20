package inbox

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("%v_%v", "testing", time.Now().UnixNano()))
	os.Exit(m.Run())
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
