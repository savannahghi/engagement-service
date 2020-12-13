package scheduling

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// SecretPassPhraseEnvVarName is the name of the environment variable that
// holds the decryption key for embedded GPG encrypted credentials
const SecretPassPhraseEnvVarName = "SECRET_PASSPHRASE"

func getJSONGoogleApplicationCredentials() ([]byte, error) {
	// a future iteration of this needs to decrypt GPG encoded creds using the
	// secret pass phrase
	return base.GPGEncryptedJSONGoogleApplicationCredentials, nil
}

// GetTokenSource gets a token source to be used in Google Cloud APIs that
// require impersonation of a user e.g Google Calendar
func GetTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	jsonCredentials, err := getJSONGoogleApplicationCredentials()
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(
		jsonCredentials,
		calendar.CalendarScope,
	)
	if err != nil {
		return nil, fmt.Errorf("JWTConfigFromJSON: %v", err)
	}
	config.Subject = DefaultCalendarEmail

	ts := config.TokenSource(ctx)
	return ts, nil
}
