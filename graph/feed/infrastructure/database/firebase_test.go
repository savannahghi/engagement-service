package db_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	db "gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/database"
)

const (
	intMax          = 9223372036854775807
	base64PNGSample = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAFNeavDAAAACklEQVQIHWNgAAAAAgABz8g15QAAAABJRU5ErkJggg=="
	base64PDFSample = "JVBERi0xLjUKJbXtrvsKNCAwIG9iago8PCAvTGVuZ3RoIDUgMCBSCiAgIC9GaWx0ZXIgL0ZsYXRlRGVjb2RlCj4+CnN0cmVhbQp4nDNUMABCXUMQpWdkopCcy1XIFcgFADCwBFQKZW5kc3RyZWFtCmVuZG9iago1IDAgb2JqCiAgIDI3CmVuZG9iagozIDAgb2JqCjw8Cj4+CmVuZG9iagoyIDAgb2JqCjw8IC9UeXBlIC9QYWdlICUgMQogICAvUGFyZW50IDEgMCBSCiAgIC9NZWRpYUJveCBbIDAgMCAwLjI0IDAuMjQgXQogICAvQ29udGVudHMgNCAwIFIKICAgL0dyb3VwIDw8CiAgICAgIC9UeXBlIC9Hcm91cAogICAgICAvUyAvVHJhbnNwYXJlbmN5CiAgICAgIC9JIHRydWUKICAgICAgL0NTIC9EZXZpY2VSR0IKICAgPj4KICAgL1Jlc291cmNlcyAzIDAgUgo+PgplbmRvYmoKMSAwIG9iago8PCAvVHlwZSAvUGFnZXMKICAgL0tpZHMgWyAyIDAgUiBdCiAgIC9Db3VudCAxCj4+CmVuZG9iago2IDAgb2JqCjw8IC9Qcm9kdWNlciAoY2Fpcm8gMS4xNi4wIChodHRwczovL2NhaXJvZ3JhcGhpY3Mub3JnKSkKICAgL0NyZWF0aW9uRGF0ZSAoRDoyMDIwMTAzMDA4MDkwOCswMycwMCkKPj4KZW5kb2JqCjcgMCBvYmoKPDwgL1R5cGUgL0NhdGFsb2cKICAgL1BhZ2VzIDEgMCBSCj4+CmVuZG9iagp4cmVmCjAgOAowMDAwMDAwMDAwIDY1NTM1IGYgCjAwMDAwMDAzODEgMDAwMDAgbiAKMDAwMDAwMDE2MSAwMDAwMCBuIAowMDAwMDAwMTQwIDAwMDAwIG4gCjAwMDAwMDAwMTUgMDAwMDAgbiAKMDAwMDAwMDExOSAwMDAwMCBuIAowMDAwMDAwNDQ2IDAwMDAwIG4gCjAwMDAwMDA1NjIgMDAwMDAgbiAKdHJhaWxlcgo8PCAvU2l6ZSA4CiAgIC9Sb290IDcgMCBSCiAgIC9JbmZvIDYgMCBSCj4+CnN0YXJ0eHJlZgo2MTQKJSVFT0YK"
	sampleVideoURL  = "https://www.youtube.com/watch?v=bPiofmZGb8o"
)

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getTextExpiry() time.Time {
	return time.Now().Add(time.Hour * 24000)
}

func getTestImage() feed.Image {
	return feed.Image{
		ID:     uuid.New().String(),
		Base64: base64PNGSample,
	}
}

func getTestVideo() feed.Video {
	return feed.Video{
		ID:  uuid.New().String(),
		URL: sampleVideoURL,
	}
}

func getTestMessage() feed.Message {
	return feed.Message{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           uuid.New().String(),
		ReplyTo:        uuid.New().String(),
		PostedByUID:    uuid.New().String(),
		PostedByName:   uuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestEvent() feed.Event {
	return feed.Event{
		ID:   uuid.New().String(),
		Name: "TEST_EVENT",
		Context: feed.Context{
			UserID:         uuid.New().String(),
			Flavour:        feed.FlavourConsumer,
			OrganizationID: uuid.New().String(),
			LocationID:     uuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() feed.Action {
	return feed.Action{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
		Event:          getTestEvent(),
	}
}

func getTestDocument() feed.Document {
	return feed.Document{
		ID:     uuid.New().String(),
		Base64: base64PDFSample,
	}
}

func testItem() *feed.Item {
	return &feed.Item{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         getTextExpiry(),
		Persistent:     true,
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Icon:           getTestImage(),
		Author:         uuid.New().String(),
		Tagline:        uuid.New().String(),
		Label:          uuid.New().String(),
		Timestamp:      time.Now(),
		Summary:        uuid.New().String(),
		Text:           uuid.New().String(),
		Images: []feed.Image{
			getTestImage(),
		},
		Videos: []feed.Video{
			getTestVideo(),
		},
		Actions: []feed.Action{
			getTestAction(),
		},
		Conversations: []feed.Message{
			getTestMessage(),
		},
		Users: []string{
			uuid.New().String(),
		},
		Groups: []string{
			uuid.New().String(),
		},
		Documents: []feed.Document{
			getTestDocument(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelEmail,
			feed.ChannelFcm,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func testNudge() *feed.Nudge {
	return &feed.Nudge{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Title:          uuid.New().String(),
		Image:          getTestImage(),
		Text:           uuid.New().String(),
		Actions: []feed.Action{
			getTestAction(),
		},
		Users: []string{
			uuid.New().String(),
		},
		Groups: []string{
			uuid.New().String(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelEmail,
			feed.ChannelFcm,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func TestNewFirebaseRepository(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case - should succeed",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.NewFirebaseRepository(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewFirebaseRepository() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
func TestFirebaseRepository_GetFeed(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourConsumer
	persistent := feed.BooleanFilterBoth
	status := feed.StatusPending
	visibility := feed.VisibilityHide
	expired := feed.BooleanFilterFalse

	type args struct {
		uid          string
		flavour      feed.Flavour
		persistent   feed.BooleanFilter
		status       feed.Status
		visibility   feed.Visibility
		expired      feed.BooleanFilter
		filterParams *feed.FilterParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no filter params",
			args: args{
				uid:          uid,
				flavour:      flavour,
				persistent:   persistent,
				status:       status,
				visibility:   visibility,
				expired:      expired,
				filterParams: nil,
			},
			wantErr: false,
		},
		{
			name: "with filter params",
			args: args{
				uid:        uid,
				flavour:    flavour,
				persistent: persistent,
				status:     status,
				visibility: visibility,
				expired:    expired,
				filterParams: &feed.FilterParams{
					Labels: []string{uuid.New().String()},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed, err := fr.GetFeed(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.persistent,
				&tt.args.status,
				&tt.args.visibility,
				&tt.args.expired,
				tt.args.filterParams,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetFeed() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, feed)
			}
		})
	}
}

func TestFirebaseRepository_GetFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()
	uid := uuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feed.Flavour
		itemID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "get back saved feed items",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
			},
			wantErr: false,
		},
		{
			name: "non existent feed item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				itemID:  uuid.New().String(),
			},
			wantErr: false,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetFeedItem(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetFeedItem() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantNil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.itemID, got.ID)
			}
		})
	}
}

func TestFirebaseRepository_SaveFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	proItem := testItem()
	consumerItem := testItem()

	type args struct {
		uid     string
		flavour feed.Flavour
		item    *feed.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pro item",
			args: args{
				uid:     uuid.New().String(),
				flavour: feed.FlavourConsumer,
				item:    proItem,
			},
			wantErr: false,
		},
		{
			name: "consumer item",
			args: args{
				uid:     uuid.New().String(),
				flavour: feed.FlavourPro,
				item:    consumerItem,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SaveFeedItem(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.item,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.SaveFeedItem() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)

				bs, err := got.ValidateAndMarshal()
				assert.Nil(t, err)
				assert.NotNil(t, bs)
			}
		})
	}
}

func TestFirebaseRepository_DeleteFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()
	uid := uuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		uid     string
		flavour feed.Flavour
		itemID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing feed item",
			args: args{
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
			},
			wantErr: false,
		},
		{
			name: "non existing feed item",
			args: args{
				uid:     uid,
				flavour: flavour,
				itemID:  uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteFeedItem(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.DeleteFeedItem() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFirebaseRepository_GetNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		uid     string
		flavour feed.Flavour
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - retrieves successfully",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudgeID: savedNudge.ID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetNudge(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudgeID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetNudge() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_SaveNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	type args struct {
		uid     string
		flavour feed.Flavour
		nudge   *feed.Nudge
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - save nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudge:   nudge,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SaveNudge(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudge,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.SaveNudge() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_DeleteNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		uid     string
		flavour feed.Flavour
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
		{
			name: "non existing nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudgeID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteNudge(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudgeID,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.DeleteNudge() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFirebaseRepository_GetAction(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourPro

	action := getTestAction()
	savedAction, err := fr.SaveAction(ctx, uid, flavour, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		uid      string
		flavour  feed.Flavour
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - will save",
			args: args{
				uid:      uid,
				flavour:  flavour,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetAction(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.actionID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetAction() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_SaveAction(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourPro
	action := getTestAction()

	type args struct {
		uid     string
		flavour feed.Flavour
		action  *feed.Action
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case, should save",
			args: args{
				uid:     uid,
				flavour: flavour,
				action:  &action,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SaveAction(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.action,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.SaveAction() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_DeleteAction(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourConsumer
	action := getTestAction()

	savedAction, err := fr.SaveAction(ctx, uid, flavour, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		uid      string
		flavour  feed.Flavour
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing action",
			args: args{
				uid:      uid,
				flavour:  flavour,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
		{
			name: "non existing action",
			args: args{
				uid:      uid,
				flavour:  flavour,
				actionID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteAction(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.actionID,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.DeleteAction() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFirebaseRepository_PostMessage(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()
	uid := uuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()

	type args struct {
		uid     string
		flavour feed.Flavour
		itemID  string
		message *feed.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - should save",
			args: args{
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
				message: &message,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.PostMessage(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
				tt.args.message,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.PostMessage() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_UpdateFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()
	uid := uuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	item.Text = "updated"

	type args struct {
		ctx     context.Context
		uid     string
		flavour feed.Flavour
		item    *feed.Item
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		{
			name: "valid case - will update",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				item:    item,
			},
			wantText: "updated",
			wantErr:  false,
		},
		{
			name: "error case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				item:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.UpdateFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.UpdateFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantText, got.Text)
			}
		})
	}
}

func TestFirebaseRepository_UpdateNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := uuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	savedNudge.Text = "updated"

	type args struct {
		ctx     context.Context
		uid     string
		flavour feed.Flavour
		nudge   *feed.Nudge
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		{
			name: "valid case - update an existing nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				nudge:   savedNudge,
			},
			wantText: "updated",
			wantErr:  false,
		},
		{
			name: "nil nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				nudge:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.UpdateNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudge)
			if (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.UpdateNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantText, got.Text)
			}
		})
	}
}

func TestFirebaseRepository_DeleteMessage(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()
	uid := uuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()
	savedMessage, err := fr.PostMessage(ctx, uid, flavour, item.ID, &message)
	assert.Nil(t, err)
	assert.NotNil(t, savedMessage)

	type args struct {
		ctx       context.Context
		uid       string
		flavour   feed.Flavour
		itemID    string
		messageID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing message, should delete",
			args: args{
				ctx:       ctx,
				uid:       uid,
				flavour:   flavour,
				itemID:    item.ID,
				messageID: savedMessage.ID,
			},
			wantErr: false,
		},
		{
			name: "non existent message, should not error",
			args: args{
				ctx:       ctx,
				uid:       uid,
				flavour:   flavour,
				itemID:    uuid.New().String(),
				messageID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteMessage(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
				tt.args.messageID,
			); (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFirebaseRepository_SaveIncomingEvent(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	event := getTestEvent()

	type args struct {
		ctx   context.Context
		event *feed.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx:   ctx,
				event: &event,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx:   ctx,
				event: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.SaveIncomingEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.SaveIncomingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFirebaseRepository_SaveOutgoingEvent(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	event := getTestEvent()

	type args struct {
		ctx   context.Context
		event *feed.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx:   ctx,
				event: &event,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx:   ctx,
				event: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.SaveOutgoingEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.SaveOutgoingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
