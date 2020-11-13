package feed_test

import (
	"context"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	db "gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/database"
)

func TestSetDefaultActions(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
	}

	type args struct {
		ctx        context.Context
		uid        string
		flavour    feed.Flavour
		repository feed.Repository
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new user - generate default consumer actions",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feed.FlavourConsumer,
				repository: fr,
			},
		},
		{
			name: "new user - generate default pro actions",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feed.FlavourPro,
				repository: fr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := feed.SetDefaultActions(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected non nil result")
					return
				}
				assert.GreaterOrEqual(t, len(got), 1)
			}
		})
	}
}

func TestSetDefaultNudges(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
	}

	type args struct {
		ctx        context.Context
		uid        string
		flavour    feed.Flavour
		repository feed.Repository
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new user - generate default consumer nudges",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feed.FlavourConsumer,
				repository: fr,
			},
		},
		{
			name: "new user - generate default pro nudges",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feed.FlavourPro,
				repository: fr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := feed.SetDefaultNudges(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultNudges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected non nil result")
					return
				}
				assert.GreaterOrEqual(t, len(got), 1)
			}
		})
	}
}

func TestSetDefaultItems(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
	}

	type args struct {
		ctx        context.Context
		uid        string
		flavour    feed.Flavour
		repository feed.Repository
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new user - generate default consumer nudges",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feed.FlavourConsumer,
				repository: fr,
			},
		},
		{
			name: "new user - generate default pro nudges",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feed.FlavourPro,
				repository: fr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := feed.SetDefaultItems(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected non nil result")
					return
				}
				assert.GreaterOrEqual(t, len(got), 1)
			}
		})
	}
}
