package pubsubhandlers_test

import (
	"strconv"
	"testing"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging/pubsubhandlers"
)

func getTestPubsubPayload(t *testing.T) *base.PubSubPayload {
	return &base.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: base.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      []byte(strconv.Quote(ksuid.New().String())),
			Attributes: map[string]string{
				"topicID": ksuid.New().String(),
			},
		},
	}
}

func TestHandleFeedRetrieval(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleFeedRetrieval(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleFeedRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleThinFeedRetrieval(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleThinFeedRetrieval(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleThinFeedRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemRetrieval(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemRetrieval(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemPublish(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemPublish(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemDelete(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemDelete(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemResolve(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemResolve(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnresolve(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemUnresolve(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemHide(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemHide(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemShow(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemShow(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemPin(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemPin(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnpin(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemUnpin(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnpin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeRetrieval(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeRetrieval(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgePublish(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgePublish(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgePublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeDelete(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeDelete(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeResolve(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeResolve(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeUnresolve(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeUnresolve(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeHide(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeHide(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeShow(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeShow(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionRetrieval(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleActionRetrieval(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionPublish(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleActionPublish(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionDelete(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleActionDelete(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessagePost(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleMessagePost(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessagePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessageDelete(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleMessageDelete(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleIncomingEvent(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleIncomingEvent(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleIncomingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
