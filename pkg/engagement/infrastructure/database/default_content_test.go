package database_test

import (
	"context"
	"testing"

	"github.com/savannahghi/engagement-service/pkg/engagement/repository"

	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"

	db "github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/database"
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
		flavour    feedlib.Flavour
		repository repository.Repository
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
				flavour:    feedlib.FlavourConsumer,
				repository: fr,
			},
		},
		{
			name: "new user - generate default pro actions",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feedlib.FlavourPro,
				repository: fr,
			},
		},
		{
			name: "fail: missing UID in PRO",
			args: args{
				ctx:        ctx,
				uid:        "",
				flavour:    feedlib.FlavourPro,
				repository: fr,
			},
			wantErr: true,
		},
		{
			name: "fail: missing UID in CONSUMER",
			args: args{
				ctx:        ctx,
				uid:        "",
				flavour:    feedlib.FlavourConsumer,
				repository: fr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.SetDefaultActions(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected non nil result")
					return
				}
				if len(got) < 1 {
					t.Errorf("got <1 action after defaults initialization")
					return
				}

				// refetch actions
				actions, err := fr.GetActions(
					ctx,
					tt.args.uid,
					tt.args.flavour,
				)
				if err != nil {
					t.Errorf("unable to re-fetch actions: %s", err)
					return
				}
				if len(actions) < 1 {
					t.Errorf("nil actions after re-fetching newly initialized actions")
					return
				}
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
		flavour    feedlib.Flavour
		repository repository.Repository
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
				flavour:    feedlib.FlavourConsumer,
				repository: fr,
			},
		},
		{
			name: "new user - generate default pro nudges",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feedlib.FlavourPro,
				repository: fr,
			},
		},
		{
			name: "fail: missing UID in PRO",
			args: args{
				ctx:        ctx,
				uid:        "",
				flavour:    feedlib.FlavourPro,
				repository: fr,
			},
			wantErr: true,
		},
		{
			name: "fail: missing UID in CONSUMER",
			args: args{
				ctx:        ctx,
				uid:        "",
				flavour:    feedlib.FlavourConsumer,
				repository: fr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.SetDefaultNudges(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultNudges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected non nil result")
					return
				}
				if len(got) < 1 {
					t.Errorf("less than 1 nudge initialized")
					return
				}

				// refetch nudges
				pending := feedlib.StatusPending
				show := feedlib.VisibilityShow
				expired := feedlib.BooleanFilterFalse
				nudges, err := fr.GetNudges(
					ctx,
					tt.args.uid,
					tt.args.flavour,
					&pending,
					&show,
					&expired,
				)
				if err != nil {
					t.Errorf("unable to fetch nudges after default initialiation: %s", err)
					return
				}
				if len(nudges) < 1 {
					t.Errorf("zero nudges after re-fetching newly initialized nudges")
					return
				}
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
		flavour    feedlib.Flavour
		repository repository.Repository
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
				flavour:    feedlib.FlavourConsumer,
				repository: fr,
			},
		},
		{
			name: "new user - generate default pro nudges",
			args: args{
				ctx:        ctx,
				uid:        ksuid.New().String(),
				flavour:    feedlib.FlavourPro,
				repository: fr,
			},
		},
		{
			name: "fail: missing UID in PRO",
			args: args{
				ctx:        ctx,
				uid:        "",
				flavour:    feedlib.FlavourPro,
				repository: fr,
			},
			wantErr: true,
		},
		{
			name: "fail: missing UID in CONSUMER",
			args: args{
				ctx:        ctx,
				uid:        "",
				flavour:    feedlib.FlavourConsumer,
				repository: fr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.SetDefaultItems(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected non nil result")
					return
				}

				if len(got) < 1 {
					t.Errorf("got < 1 item after defaults initialization")
					return
				}

				// refetch items
				items, err := fr.GetItems(
					ctx,
					tt.args.uid,
					tt.args.flavour,
					feedlib.BooleanFilterBoth,
					nil,
					nil,
					nil,
					nil,
				)
				if err != nil {
					t.Errorf("unable to re-fetch items: %s", err)
					return
				}
				if len(items) < 1 {
					t.Errorf("nil items after re-fetching newly initialized items")
					return
				}
			}
		})
	}
}

func TestTruncateStringWithEllipses(t *testing.T) {
	type args struct {
		str    string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "short string",
			args: args{
				str:    "drink",
				length: 3,
			},
			want: "dri",
		},
		{
			name: "empty string",
			args: args{
				str:    "something",
				length: 0,
			},
			want: "",
		},
		{
			name: "long string",
			args: args{
				str:    "This is an epic tale that is intended to exceed 140 characters. At that point, it will be truncated to the indicated target length, including getting some ellipses added at the end.",
				length: 140,
			},
			want: "This is an epic tale that is intended to exceed 140 characters. At that point, it will be truncated to the indicated target length, incl...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.TruncateStringWithEllipses(tt.args.str, tt.args.length); got != tt.want {
				t.Errorf("TruncateStringWithEllipses() = %v, want %v", got, tt.want)
			}
		})
	}
}
