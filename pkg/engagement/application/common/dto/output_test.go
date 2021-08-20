package dto_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
)

func Test_Message_IsNode(t *testing.T) {
	type fields struct {
		ID                  string
		AccountSID          string
		APIVersion          string
		Body                string
		DateCreated         string
		DateSent            string
		DateUpdated         string
		Direction           string
		ErrorCode           *string
		ErrorMessage        *string
		From                string
		MessagingServiceSID string
		NumMedia            string
		NumSegments         string
		Price               *string
		PriceUnit           *string
		SID                 string
		Status              string
		SubresourceURLs     map[string]string
		To                  string
		URI                 string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &dto.Message{
				ID:                  tt.fields.ID,
				AccountSID:          tt.fields.AccountSID,
				APIVersion:          tt.fields.APIVersion,
				Body:                tt.fields.Body,
				DateCreated:         tt.fields.DateCreated,
				DateSent:            tt.fields.DateSent,
				DateUpdated:         tt.fields.DateUpdated,
				Direction:           tt.fields.Direction,
				ErrorCode:           tt.fields.ErrorCode,
				ErrorMessage:        tt.fields.ErrorMessage,
				From:                tt.fields.From,
				MessagingServiceSID: tt.fields.MessagingServiceSID,
				NumMedia:            tt.fields.NumMedia,
				NumSegments:         tt.fields.NumSegments,
				Price:               tt.fields.Price,
				PriceUnit:           tt.fields.PriceUnit,
				SID:                 tt.fields.SID,
				Status:              tt.fields.Status,
				SubresourceURLs:     tt.fields.SubresourceURLs,
				To:                  tt.fields.To,
				URI:                 tt.fields.URI,
			}
			m.IsNode()
		})
	}
}

func Test_Message_GetID(t *testing.T) {
	type fields struct {
		ID                  string
		AccountSID          string
		APIVersion          string
		Body                string
		DateCreated         string
		DateSent            string
		DateUpdated         string
		Direction           string
		ErrorCode           *string
		ErrorMessage        *string
		From                string
		MessagingServiceSID string
		NumMedia            string
		NumSegments         string
		Price               *string
		PriceUnit           *string
		SID                 string
		Status              string
		SubresourceURLs     map[string]string
		To                  string
		URI                 string
	}
	tests := []struct {
		name   string
		fields fields
		want   firebasetools.ID
	}{
		{
			name: "good case",
			fields: fields{
				ID: "an ID",
			},
			want: firebasetools.IDValue("an ID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &dto.Message{
				ID:                  tt.fields.ID,
				AccountSID:          tt.fields.AccountSID,
				APIVersion:          tt.fields.APIVersion,
				Body:                tt.fields.Body,
				DateCreated:         tt.fields.DateCreated,
				DateSent:            tt.fields.DateSent,
				DateUpdated:         tt.fields.DateUpdated,
				Direction:           tt.fields.Direction,
				ErrorCode:           tt.fields.ErrorCode,
				ErrorMessage:        tt.fields.ErrorMessage,
				From:                tt.fields.From,
				MessagingServiceSID: tt.fields.MessagingServiceSID,
				NumMedia:            tt.fields.NumMedia,
				NumSegments:         tt.fields.NumSegments,
				Price:               tt.fields.Price,
				PriceUnit:           tt.fields.PriceUnit,
				SID:                 tt.fields.SID,
				Status:              tt.fields.Status,
				SubresourceURLs:     tt.fields.SubresourceURLs,
				To:                  tt.fields.To,
				URI:                 tt.fields.URI,
			}
			if got := m.GetID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Message_SetID(t *testing.T) {
	type fields struct {
		ID                  string
		AccountSID          string
		APIVersion          string
		Body                string
		DateCreated         string
		DateSent            string
		DateUpdated         string
		Direction           string
		ErrorCode           *string
		ErrorMessage        *string
		From                string
		MessagingServiceSID string
		NumMedia            string
		NumSegments         string
		Price               *string
		PriceUnit           *string
		SID                 string
		Status              string
		SubresourceURLs     map[string]string
		To                  string
		URI                 string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "good case",
			args: args{
				id: "an ID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &dto.Message{
				ID:                  tt.fields.ID,
				AccountSID:          tt.fields.AccountSID,
				APIVersion:          tt.fields.APIVersion,
				Body:                tt.fields.Body,
				DateCreated:         tt.fields.DateCreated,
				DateSent:            tt.fields.DateSent,
				DateUpdated:         tt.fields.DateUpdated,
				Direction:           tt.fields.Direction,
				ErrorCode:           tt.fields.ErrorCode,
				ErrorMessage:        tt.fields.ErrorMessage,
				From:                tt.fields.From,
				MessagingServiceSID: tt.fields.MessagingServiceSID,
				NumMedia:            tt.fields.NumMedia,
				NumSegments:         tt.fields.NumSegments,
				Price:               tt.fields.Price,
				PriceUnit:           tt.fields.PriceUnit,
				SID:                 tt.fields.SID,
				Status:              tt.fields.Status,
				SubresourceURLs:     tt.fields.SubresourceURLs,
				To:                  tt.fields.To,
				URI:                 tt.fields.URI,
			}
			m.SetID(tt.args.id)
			assert.Equal(t, m.GetID(), firebasetools.IDValue(tt.args.id))
		})
	}
}

func TestNPSResponse_IsNode(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Score     int
		SladeCode string
		Email     *string
		MSISDN    *string
		Feedback  []dto.Feedback
		Timestamp time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &dto.NPSResponse{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Score:     tt.fields.Score,
				SladeCode: tt.fields.SladeCode,
				Email:     tt.fields.Email,
				MSISDN:    tt.fields.MSISDN,
				Feedback:  tt.fields.Feedback,
				Timestamp: tt.fields.Timestamp,
			}
			e.IsNode()
		})
	}
}

func TestNPSResponse_GetID(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Score     int
		SladeCode string
		Email     *string
		MSISDN    *string
		Feedback  []dto.Feedback
		Timestamp time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   firebasetools.ID
	}{
		{
			name: "good case",
			fields: fields{
				ID: "an ID",
			},
			want: firebasetools.IDValue("an ID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &dto.NPSResponse{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Score:     tt.fields.Score,
				SladeCode: tt.fields.SladeCode,
				Email:     tt.fields.Email,
				MSISDN:    tt.fields.MSISDN,
				Feedback:  tt.fields.Feedback,
				Timestamp: tt.fields.Timestamp,
			}
			if got := e.GetID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NPSResponse.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNPSResponse_SetID(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Score     int
		SladeCode string
		Email     *string
		MSISDN    *string
		Feedback  []dto.Feedback
		Timestamp time.Time
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "good case",
			args: args{
				id: "an ID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &dto.NPSResponse{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Score:     tt.fields.Score,
				SladeCode: tt.fields.SladeCode,
				Email:     tt.fields.Email,
				MSISDN:    tt.fields.MSISDN,
				Feedback:  tt.fields.Feedback,
				Timestamp: tt.fields.Timestamp,
			}
			e.SetID(tt.args.id)
		})
	}
}

func TestAccessToken_IsEntity(t *testing.T) {
	type fields struct {
		JWT             string
		UniqueName      string
		SID             string
		DateUpdated     time.Time
		Status          string
		Type            string
		MaxParticipants int
		Duration        *int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := dto.AccessToken{
				JWT:             tt.fields.JWT,
				UniqueName:      tt.fields.UniqueName,
				SID:             tt.fields.SID,
				DateUpdated:     tt.fields.DateUpdated,
				Status:          tt.fields.Status,
				Type:            tt.fields.Type,
				MaxParticipants: tt.fields.MaxParticipants,
				Duration:        tt.fields.Duration,
			}
			a.IsEntity()
		})
	}
}

func TestSavedNotification_IsEntity(t *testing.T) {
	type fields struct {
		ID                string
		RegistrationToken string
		MessageID         string
		Timestamp         time.Time
		Data              map[string]interface{}
		Notification      *dto.FirebaseSimpleNotification
		AndroidConfig     *dto.FirebaseAndroidConfig
		WebpushConfig     *dto.FirebaseWebpushConfig
		APNSConfig        *dto.FirebaseAPNSConfig
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := dto.SavedNotification{
				ID:                tt.fields.ID,
				RegistrationToken: tt.fields.RegistrationToken,
				MessageID:         tt.fields.MessageID,
				Timestamp:         tt.fields.Timestamp,
				Data:              tt.fields.Data,
				Notification:      tt.fields.Notification,
				AndroidConfig:     tt.fields.AndroidConfig,
				WebpushConfig:     tt.fields.WebpushConfig,
				APNSConfig:        tt.fields.APNSConfig,
			}
			u.IsEntity()
		})
	}
}

func TestNewOKResp(t *testing.T) {
	type args struct {
		rawResponse interface{}
	}

	tests := []struct {
		name string
		args args
		want *dto.OKResp
	}{
		{
			name: "happy case",
			args: args{
				rawResponse: "some raw response",
			},
			want: &dto.OKResp{Status: "OK", Response: "some raw response"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dto.NewOKResp(tt.args.rawResponse); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOKResp() = %v, want %v", got, tt.want)
			}
		})
	}
}
