package otp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
)

func TestValidateSendOTPPayload(t *testing.T) {
	goodData := &resources.Msisdn{
		Msisdn: "+254723002959",
	}
	goodDataJSONBytes, err := json.Marshal(goodData)
	assert.Nil(t, err)
	assert.NotNil(t, goodDataJSONBytes)

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
		want    string
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "invalid data",
			args: args{
				w: httptest.NewRecorder(),
				r: emptyDataRequest,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSendOTPPayload(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSendOTPPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateSendOTPPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateGenerateRetryOTPPayload(t *testing.T) {
	phoneNumber := base.TestUserPhoneNumber
	goodData := &resources.GenerateRetryOTP{
		Msisdn:    &phoneNumber,
		RetryStep: 2,
	}
	goodDataJSONBytes, err := json.Marshal(goodData)
	assert.Nil(t, err)
	assert.NotNil(t, goodDataJSONBytes)

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
		want    *resources.GenerateRetryOTP
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want: &resources.GenerateRetryOTP{
				Msisdn:    &phoneNumber,
				RetryStep: 2,
			},
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
			got, err := ValidateGenerateRetryOTPPayload(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGenerateRetryOTPPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateGenerateRetryOTPPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateVerifyOTPPayload_Phone(t *testing.T) {
	phoneNumber := base.TestUserPhoneNumber
	verificationCode := "45225"

	goodData := &resources.VerifyOTP{
		Msisdn:           &phoneNumber,
		VerificationCode: &verificationCode,
	}
	goodDataJSONBytes, err := json.Marshal(goodData)
	assert.Nil(t, err)
	assert.NotNil(t, goodDataJSONBytes)

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
		want    *resources.VerifyOTP
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want: &resources.VerifyOTP{
				Msisdn:           &phoneNumber,
				VerificationCode: &verificationCode,
			},
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
			got, err := ValidateVerifyOTPPayload(tt.args.w, tt.args.r, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVerifyOTPPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateVerifyOTPPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateVerifyOTPPayload_Email(t *testing.T) {
	phoneNumber := base.TestUserPhoneNumber
	email := base.TestUserEmail
	verificationCode := "45225"

	goodData := &resources.VerifyOTP{
		Email:            &email,
		VerificationCode: &verificationCode,
	}
	goodDataJSONBytes, err := json.Marshal(goodData)
	assert.Nil(t, err)
	assert.NotNil(t, goodDataJSONBytes)

	validRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validRequest.Body = ioutil.NopCloser(bytes.NewReader(goodDataJSONBytes))

	emptyDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	emptyDataRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	invalidData := &resources.VerifyOTP{
		Msisdn:           &phoneNumber,
		VerificationCode: &verificationCode,
	}
	invalidDataJSONBytes, err := json.Marshal(invalidData)
	assert.Nil(t, err)
	assert.NotNil(t, invalidDataJSONBytes)

	invalidRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidRequest.Body = ioutil.NopCloser(bytes.NewReader(invalidDataJSONBytes))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *resources.VerifyOTP
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want: &resources.VerifyOTP{
				Email:            &email,
				VerificationCode: &verificationCode,
			},
			wantErr: false,
		},
		{
			name: "invalid : phone instead of email",
			args: args{
				w: httptest.NewRecorder(),
				r: invalidRequest,
			},
			want:    nil,
			wantErr: true,
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
			got, err := ValidateVerifyOTPPayload(tt.args.w, tt.args.r, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVerifyOTPPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateVerifyOTPPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
