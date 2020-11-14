package feed_test

import (
	"bytes"
	"testing"

	"gitlab.slade360emr.com/go/feed/graph/feed"
)

func getBlankActionType() *feed.ActionType {
	at := feed.ActionType("")
	return &at
}

func TestActionType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.ActionType
		want bool
	}{
		{
			name: "valid case",
			e:    feed.ActionTypeFloating,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.ActionType("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ActionType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionType_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.ActionType
		want string
	}{
		{
			name: "primary",
			e:    feed.ActionTypePrimary,
			want: "PRIMARY",
		},
		{
			name: "secondary",
			e:    feed.ActionTypeSecondary,
			want: "SECONDARY",
		},
		{
			name: "overflow",
			e:    feed.ActionTypeOverflow,
			want: "OVERFLOW",
		},
		{
			name: "floating",
			e:    feed.ActionTypeFloating,
			want: "FLOATING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ActionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionType_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.ActionType
		args    args
		wantErr bool
	}{
		{
			name: "primary",
			e:    getBlankActionType(),
			args: args{
				v: "PRIMARY",
			},
			wantErr: false,
		},
		{
			name: "secondary",
			e:    getBlankActionType(),
			args: args{
				v: "SECONDARY",
			},
			wantErr: false,
		},
		{
			name: "overflow",
			e:    getBlankActionType(),
			args: args{
				v: "OVERFLOW",
			},
			wantErr: false,
		},
		{
			name: "floating",
			e:    getBlankActionType(),
			args: args{
				v: "FLOATING",
			},
			wantErr: false,
		},
		{
			name: "invalid - should error",
			e:    getBlankActionType(),
			args: args{
				v: "bogus bonoko",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"ActionType.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestActionType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.ActionType
		wantW string
	}{
		{
			name:  "floating",
			e:     feed.ActionTypeFloating,
			wantW: `"FLOATING"`,
		},
		{
			name:  "primary",
			e:     feed.ActionTypePrimary,
			wantW: `"PRIMARY"`,
		},
		{
			name:  "secondary",
			e:     feed.ActionTypeSecondary,
			wantW: `"SECONDARY"`,
		},
		{
			name:  "overflow",
			e:     feed.ActionTypeOverflow,
			wantW: `"OVERFLOW"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"ActionType.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestHandling_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Handling
		want bool
	}{
		{
			name: "valid case",
			e:    feed.HandlingFullPage,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.Handling("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Handling.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandling_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Handling
		want string
	}{
		{
			name: "simple case",
			e:    feed.HandlingInline,
			want: "INLINE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Handling.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandling_UnmarshalGQL(t *testing.T) {
	target := feed.Handling("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.Handling
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "INLINE",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Handling.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestHandling_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.Handling
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.HandlingFullPage,
			wantW: `"FULL_PAGE"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"Handling.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Status
		want bool
	}{
		{
			name: "valid case",
			e:    feed.StatusDone,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.Status("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Status.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Status
		want string
	}{
		{
			name: "simple case",
			e:    feed.StatusDone,
			want: "DONE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_UnmarshalGQL(t *testing.T) {
	target := feed.Status("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.Status
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "DONE",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Status.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.Status
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.StatusDone,
			wantW: `"DONE"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Status.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestVisibility_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Visibility
		want bool
	}{
		{
			name: "valid case",
			e:    feed.VisibilityHide,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.Visibility("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Visibility.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibility_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Visibility
		want string
	}{

		{
			name: "simple case",
			e:    feed.VisibilityShow,
			want: "SHOW",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Visibility.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibility_UnmarshalGQL(t *testing.T) {
	target := feed.Visibility("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.Visibility
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "SHOW",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Visibility.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestVisibility_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.Visibility
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.VisibilityHide,
			wantW: `"HIDE"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"Visibility.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestChannel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Channel
		want bool
	}{
		{
			name: "valid case",
			e:    feed.ChannelEmail,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.Channel("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Channel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Channel
		want string
	}{
		{
			name: "simple case",
			e:    feed.ChannelEmail,
			want: "EMAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Channel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_UnmarshalGQL(t *testing.T) {
	target := feed.Channel("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.Channel
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "EMAIL",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Channel.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestChannel_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.Channel
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.ChannelEmail,
			wantW: `"EMAIL"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Channel.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFlavour_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Flavour
		want bool
	}{
		{
			name: "valid case",
			e:    feed.FlavourConsumer,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.Flavour("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Flavour.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlavour_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Flavour
		want string
	}{
		{
			name: "simple case",
			e:    feed.FlavourConsumer,
			want: "CONSUMER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Flavour.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlavour_UnmarshalGQL(t *testing.T) {
	target := feed.Flavour("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.Flavour
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "PRO",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Flavour.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFlavour_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.Flavour
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.FlavourPro,
			wantW: `"PRO"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Flavour.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestKeys_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Keys
		want bool
	}{
		{
			name: "valid case",
			e:    feed.KeysActions,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.Keys("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Keys.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeys_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.Keys
		want string
	}{
		{
			name: "simple case",
			e:    feed.KeysActions,
			want: "actions",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Keys.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeys_UnmarshalGQL(t *testing.T) {
	target := feed.Keys("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.Keys
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "actions",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(
				tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Keys.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestKeys_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.Keys
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.KeysActions,
			wantW: `"actions"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Keys.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestBooleanFilter_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.BooleanFilter
		want bool
	}{
		{
			name: "valid case",
			e:    feed.BooleanFilterBoth,
			want: true,
		},
		{
			name: "invalid case",
			e:    feed.BooleanFilter("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("BooleanFilter.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooleanFilter_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.BooleanFilter
		want string
	}{
		{
			name: "simple case",
			e:    feed.BooleanFilterFalse,
			want: "FALSE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("BooleanFilter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooleanFilter_UnmarshalGQL(t *testing.T) {
	target := feed.BooleanFilter("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.BooleanFilter
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "BOTH",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"BooleanFilter.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestBooleanFilter_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.BooleanFilter
		wantW string
	}{
		{
			name:  "simple case",
			e:     feed.BooleanFilterBoth,
			wantW: `"BOTH"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"BooleanFilter.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestLinkType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     feed.LinkType
		wantW string
	}{
		{
			name:  "PDF document",
			e:     feed.LinkTypePdfDocument,
			wantW: `"PDF_DOCUMENT"`,
		},
		{
			name:  "PNG Image",
			e:     feed.LinkTypePngImage,
			wantW: `"PNG_IMAGE"`,
		},
		{
			name:  "YouTube Video",
			e:     feed.LinkTypeYoutubeVideo,
			wantW: `"YOUTUBE_VIDEO"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("LinkType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestLinkType_UnmarshalGQL(t *testing.T) {
	l := feed.LinkType("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *feed.LinkType
		args    args
		wantErr bool
	}{
		{
			name: "invalid link type",
			e:    &l,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
		{
			name: "valid - pdf",
			e:    &l,
			args: args{
				v: "PDF_DOCUMENT",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("LinkType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    feed.LinkType
		want bool
	}{
		{
			name: "PDF document",
			e:    feed.LinkTypePdfDocument,
			want: true,
		},
		{
			name: "PNG Image",
			e:    feed.LinkTypePngImage,
			want: true,
		},
		{
			name: "YouTube Video",
			e:    feed.LinkTypeYoutubeVideo,
			want: true,
		},
		{
			name: "invalid link type",
			e:    feed.LinkType("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("LinkType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkType_String(t *testing.T) {
	tests := []struct {
		name string
		e    feed.LinkType
		want string
	}{
		{
			name: "YouTube video",
			e:    feed.LinkTypeYoutubeVideo,
			want: "YOUTUBE_VIDEO",
		},
		{
			name: "PDF document",
			e:    feed.LinkTypePdfDocument,
			want: "PDF_DOCUMENT",
		},
		{
			name: "PNG image",
			e:    feed.LinkTypePngImage,
			want: "PNG_IMAGE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("LinkType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
