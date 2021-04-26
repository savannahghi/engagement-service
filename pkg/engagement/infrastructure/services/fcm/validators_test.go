package fcm_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"
)

func Test_validateFCMData(t *testing.T) {
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil data",
			args: args{
				data: nil,
			},
			wantErr: false,
		},
		{
			name: "good case",
			args: args{
				data: map[string]string{
					"a": "1",
					"b": "2",
				},
			},
			wantErr: false,
		},
		{
			name: "reserved words",
			args: args{
				data: map[string]string{
					"a":    "1",
					"b":    "2",
					"from": "should not be used",
				},
			},
			wantErr: true,
		},
		{
			name: "illegal prefix",
			args: args{
				data: map[string]string{
					"a":            "1",
					"b":            "2",
					"gcmgibberish": "gcm is an illegal prefix",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fcm.ValidateFCMData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("fcm.ValidateFCMData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func GetSendNotificationPayload() *base.SendNotificationPayload {

	img := "https://www.wxpr.org/sites/wxpr/files/styles/medium/public/202007/chipmunk-5401165_1920.jpg"
	key := uuid.New().String()
	fakeToken := uuid.New().String()
	pckg := "video"

	return &base.SendNotificationPayload{
		RegistrationTokens: []string{fakeToken},
		Data: map[string]string{
			"some": "data",
		},
		Notification: &base.FirebaseSimpleNotificationInput{
			Title:    "Test Notification",
			Body:     "From Integration Tests",
			ImageURL: &img,
			Data: map[string]interface{}{
				"more": "data",
			},
		},
		Android: &base.FirebaseAndroidConfigInput{
			Priority:              "high",
			CollapseKey:           &key,
			RestrictedPackageName: &pckg,
		},
	}
}

func TestValidateSendNotificationPayload(t *testing.T) {

	goodData := GetSendNotificationPayload()
	goodDataJSONBytes, err := json.Marshal(goodData)
	if err != nil {
		t.Errorf("struct could not be marshalled: %v", err)
		return
	}
	if goodDataJSONBytes == nil {
		t.Errorf("nil json Bytes")
		return
	}

	validRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validRequest.Body = ioutil.NopCloser(bytes.NewReader(goodDataJSONBytes))

	emptyDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	emptyDataRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *base.SendNotificationPayload
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want:    goodData,
			wantErr: false,
		},
		{
			name: "invalid data",
			args: args{
				w: httptest.NewRecorder(),
				r: emptyDataRequest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fcm.ValidateSendNotificationPayload(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSendNotificationPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSendNotificationPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
