package fcm_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"

	"github.com/rs/xid"
	"gitlab.slade360emr.com/go/base"
)

func getNotificationPayload(t *testing.T) base.SendNotificationPayload {
	imgURL := "https://example.com/img.png"
	return base.SendNotificationPayload{
		RegistrationTokens: []string{xid.New().String(), xid.New().String()},
		Data: map[string]string{
			xid.New().String(): xid.New().String(),
			xid.New().String(): xid.New().String(),
		},
		Notification: &base.FirebaseSimpleNotificationInput{
			Title:    xid.New().String(),
			Body:     xid.New().String(),
			ImageURL: &imgURL,
			Data: map[string]interface{}{
				xid.New().String(): xid.New().String(),
				xid.New().String(): xid.New().String(),
			},
		},
	}
}

func TestNewRemotePushService(t *testing.T) {
	ctx := context.Background()
	rfs, err := fcm.NewRemotePushService(ctx)
	if err != nil {
		t.Errorf("error setting up remote FCM push service: %v", err)
		return
	}
	if rfs == nil {
		t.Errorf("nil remote FCM push service")
		return
	}
}

func TestRemotePushService_Push(t *testing.T) {
	ctx := context.Background()
	rfs, err := fcm.NewRemotePushService(ctx)
	if err != nil {
		t.Errorf("error setting up remote FCM push service: %v", err)
		return
	}
	if rfs == nil {
		t.Errorf("nil remote FCM push service")
		return
	}

	type args struct {
		ctx                 context.Context
		sender              string
		notificationPayload base.SendNotificationPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid send - won't actually push but won't error",
			args: args{
				ctx:                 ctx,
				sender:              "test",
				notificationPayload: getNotificationPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := rfs.Push(tt.args.ctx, tt.args.sender, tt.args.notificationPayload); (err != nil) != tt.wantErr {
				t.Errorf("RemotePushService.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
