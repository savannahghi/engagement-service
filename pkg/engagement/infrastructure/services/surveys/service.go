package surveys

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/surveys")

// NPSResponseCollectionName firestore collection name where nps responses are stored
const NPSResponseCollectionName = "nps_response"

// ServiceSurveys defines the interactions with the surveys service
type ServiceSurveys interface {
	RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error)
}

// NewService initializes a surveys service
func NewService(repository repository.Repository) *Service {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app for Surveys service: %s", err)
	}
	ctx := context.Background()

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore client: %s", err)
	}

	srv := &Service{
		firestoreClient: firestoreClient,
		Repository:      repository,
	}
	srv.checkPreconditions()
	return srv
}

func (s Service) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("surveys service has a nil firestore client")
	}
}

// Service is an surveys service
type Service struct {
	firestoreClient *firestore.Client
	Repository      repository.Repository
}

// RecordNPSResponse ...
func (s Service) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	ctx, span := tracer.Start(ctx, "RecordNPSResponse")
	defer span.End()
	s.checkPreconditions()
	response := &dto.NPSResponse{
		Name:      input.Name,
		Score:     input.Score,
		SladeCode: input.SladeCode,
	}

	response.SetID(uuid.New().String())

	if input.Email != nil {
		response.Email = input.Email
	}

	if input.PhoneNumber != nil {
		response.MSISDN = input.PhoneNumber
	}

	feedbacks := []dto.Feedback{}
	if input.Feedback != nil {

		for _, input := range input.Feedback {
			feedback := dto.Feedback{
				Question: input.Question,
				Answer:   input.Answer,
			}
			feedbacks = append(feedbacks, feedback)
		}

		response.Feedback = feedbacks
	}

	response.Timestamp = time.Now()

	err := s.Repository.SaveNPSResponse(ctx, response)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("cannot save nps response to firestore: %w", err)
	}

	return true, nil
}
