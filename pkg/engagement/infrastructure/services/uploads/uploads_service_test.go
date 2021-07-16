package uploads

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/profileutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	m.Run()
}

func TestUpload(t *testing.T) {
	ctx := context.Background()

	bs, err := ioutil.ReadFile("testdata/gandalf.jpg")
	assert.Nil(t, err)
	sEnc := base64.StdEncoding.EncodeToString(bs)

	service := NewUploadsService()
	tests := map[string]struct {
		inp                  profileutils.UploadInput
		expectError          bool
		expectedErrorMessage string
	}{
		"simple_case": {
			inp: profileutils.UploadInput{
				Title:       "Test file from automated tests",
				ContentType: "JPG",
				Language:    "en",
				Filename:    fmt.Sprintf("%s.jpg", uuid.New().String()),
				Base64data:  sEnc,
			},
			expectError: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			upload, err := service.Upload(ctx, tc.inp)
			if tc.expectError {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			}
			if !tc.expectError {
				assert.NotNil(t, upload)
				assert.Nil(t, err)

				assert.NotZero(t, upload.ID)
				assert.NotZero(t, upload.URL)
				assert.NotZero(t, upload.Size)
				assert.NotZero(t, upload.Hash)
				assert.NotZero(t, upload.Creation)
				assert.NotZero(t, upload.Title)
				assert.NotZero(t, upload.ContentType)
				assert.NotZero(t, upload.Language)
				assert.NotZero(t, upload.Base64data)
			}
		})
	}
}
