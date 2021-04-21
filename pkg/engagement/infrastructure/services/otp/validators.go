package otp

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
)

// ValidateSendOTPPayload checks the validity of the request payload
func ValidateSendOTPPayload(w http.ResponseWriter, r *http.Request) (string, error) {
	payload := &resources.Msisdn{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if payload.Msisdn == "" {
		err := fmt.Errorf("invalid generate and send otp payload")
		base.ReportErr(w, err, http.StatusBadRequest)
		return "", err
	}
	return payload.Msisdn, nil
}

// ValidateGenerateRetryOTPPayload checks the validity of the request payload
func ValidateGenerateRetryOTPPayload(w http.ResponseWriter, r *http.Request) (*resources.GenerateRetryOTP, error) {
	payload := &resources.GenerateRetryOTP{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if payload.Msisdn == nil || payload.RetryStep == 0 {
		err := fmt.Errorf("invalid generate retry and fallback otp payload")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return payload, nil
}

// ValidateVerifyOTPPayload checks the validity of the request payload
func ValidateVerifyOTPPayload(w http.ResponseWriter, r *http.Request, isEmail bool) (*resources.VerifyOTP, error) {
	payload := &resources.VerifyOTP{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if isEmail {
		if payload.Email == nil || payload.VerificationCode == nil {
			err := fmt.Errorf("invalid verify otp payload")
			base.ReportErr(w, err, http.StatusBadRequest)
			return nil, err
		}
	} else {
		if payload.Msisdn == nil || payload.VerificationCode == nil {
			err := fmt.Errorf("invalid verify otp payload")
			base.ReportErr(w, err, http.StatusBadRequest)
			return nil, err
		}
	}
	return payload, nil
}
