package fcm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
)

// ValidateFCMData checks that the supplied FCM data does not use re
func ValidateFCMData(data map[string]string) error {
	if data != nil {
		fcmReservedWords := []string{"from", "notification", "message_type"}
		for _, reservedWord := range fcmReservedWords {
			_, present := data[reservedWord]
			if present {
				return fmt.Errorf("invalid use of FCM reserved word: %s", reservedWord)
			}
		}
		fcmReservedPrefixes := []string{"gcm", "google"}
		for _, reservedPrefix := range fcmReservedPrefixes {
			for k := range data {
				if strings.HasPrefix(k, reservedPrefix) {
					return fmt.Errorf("illegal FCM prefix: %s", reservedPrefix)
				}
			}
		}
	}
	return nil
}

// ValidateSendNotificationPayload checks that the request payload supplied in the indicated request are valid
func ValidateSendNotificationPayload(w http.ResponseWriter, r *http.Request) (*firebasetools.SendNotificationPayload, error) {
	payload := &firebasetools.SendNotificationPayload{}
	serverutils.DecodeJSONToTargetStruct(w, r, payload)

	if payload.RegistrationTokens == nil {
		err := fmt.Errorf("can't send FCM notifications to nil registration tokens")
		errorcodeutil.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}

	err := ValidateFCMData(payload.Data)
	if err != nil {
		errorcodeutil.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}

	return payload, nil
}
