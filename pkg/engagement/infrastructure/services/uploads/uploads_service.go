package uploads

import (
	"bytes"
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"encoding/base64"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/rs/xid"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement-service/pkg/engagement/services/uploads")

// Constants used to save and retrieve upload content
const (
	BucketNameBase                = "healthcloud_bewell_api_gateway_uploads"
	BucketLocation                = "EU"
	UploadFirestoreCollectionName = "uploads"
)

// ServiceUploads ...
type ServiceUploads interface {
	Upload(
		ctx context.Context,
		inp profileutils.UploadInput,
	) (*profileutils.Upload, error)

	FindUploadByID(
		ctx context.Context,
		id string,
	) (*profileutils.Upload, error)
}

func getBucketName() string {
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	projectSlug := strings.ReplaceAll(projectID, "-", "_")
	bucketName := fmt.Sprintf("%s_%s", projectSlug, BucketNameBase)
	return bucketName
}

// GCSClient initializes a new Google Cloud Storage client
// and ensures that the bucket(s) we need exists
func GCSClient(ctx context.Context) (*storage.Client, error) {
	ctx, span := tracer.Start(ctx, "GCSClient")
	defer span.End()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't initialize storage client: %w", err)
	}

	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	bucket := storageClient.Bucket(getBucketName())
	_, err = bucket.Attrs(ctx)
	if err != nil {
		bucketAccess := storage.UniformBucketLevelAccess{
			Enabled:    true,
			LockedTime: time.Now(),
		}
		bucketAttrs := &storage.BucketAttrs{
			Location:                 BucketLocation,
			UniformBucketLevelAccess: bucketAccess,
			VersioningEnabled:        true,
			StorageClass:             "STANDARD",
		}
		err := bucket.Create(ctx, projectID, bucketAttrs)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("can't create bucket: %w", err)
		}
	}
	return storageClient, nil
}

// NewUploadsService initializes an upload service
func NewUploadsService() *Service {
	ctx := context.Background()
	client, err := GCSClient(ctx)
	if err != nil {
		log.Panicf(
			"unable to initialize GCS client for upload service: %s", err)
	}

	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Panicf(
			"unable to initialize Firebase  for upload service: %s", err)
	}
	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Panicf(
			"unable to initialize Firestore client for upload service: %s", err)
	}

	return &Service{
		storageClient: client,
		contentTypeMap: map[string]string{
			"PNG": "image/png",
			"JPG": "image/jpeg",
			"PDF": "application/pdf",
		},
		firestoreClient: firestoreClient,
	}
}

// Service is an upload service
type Service struct {
	storageClient   *storage.Client
	firestoreClient *firestore.Client
	contentTypeMap  map[string]string
}

func (s Service) enforcePreconditions() {
	if s.storageClient == nil {
		log.Panicf("uploads.Service *storage.Client is nil")
	}

	if s.firestoreClient == nil {
		log.Panicf("uploads.Service *firestore.Client is nil")
	}

	if s.contentTypeMap == nil {
		log.Panicf("uploads.Service contentTypeMap is nil")
	}
}

// StorageClient returns the upload service's Google Cloud Storage Client
func (s Service) StorageClient() *storage.Client {
	s.enforcePreconditions()
	return s.storageClient
}

// FirestoreClient returns the upload service's Firebase Firestore client
func (s Service) FirestoreClient() *firestore.Client {
	s.enforcePreconditions()
	return s.firestoreClient
}

// Upload uploads the file to cloud storage
func (s Service) Upload(
	ctx context.Context,
	inp profileutils.UploadInput,
) (*profileutils.Upload, error) {
	ctx, span := tracer.Start(ctx, "Upload")
	defer span.End()
	s.enforcePreconditions()

	// data decoding
	data, err := base64.StdEncoding.DecodeString(inp.Base64data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("upload base64 decode error: %w", err)
	}

	// sha hash
	h := sha512.New()
	_, err = h.Write(data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to calculate upload hash: %w", err)
	}
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// upload to GCS
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	objectName := fmt.Sprintf("%s-%s", xid.New(), inp.Filename)
	dest := s.storageClient.Bucket(
		getBucketName()).Object(objectName).NewWriter(ctx)
	source := bytes.NewReader(data)
	if _, err = io.Copy(dest, source); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}
	if err := dest.Close(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}
	url := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		BucketNameBase,
		inp.Filename,
	)

	id := ksuid.New().String()
	u := &profileutils.Upload{
		ID:          id,
		Title:       inp.Title,
		Creation:    time.Now(),
		ContentType: inp.ContentType,
		Language:    inp.Language,
		Size:        len(data),
		Hash:        hash,
		URL:         url,
		Base64data:  inp.Base64data,
	}

	_, _, err = firebasetools.CreateNode(ctx, u)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't save upload: %w", err)
	}
	return u, nil
}

// FindUploadByID retrieves an upload by it's ID
func (s Service) FindUploadByID(
	ctx context.Context,
	id string,
) (*profileutils.Upload, error) {
	ctx, span := tracer.Start(ctx, "FindUploadByID")
	defer span.End()
	node, err := firebasetools.RetrieveNode(ctx, id, &profileutils.Upload{})
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to retrieve upload: %w", err)
	}
	upload, ok := node.(*profileutils.Upload)
	if !ok {
		return nil, fmt.Errorf("unable to cast node to upload")
	}
	return upload, nil
}
