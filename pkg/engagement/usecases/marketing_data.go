package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/profileutils"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/apiclient"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/authorization"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/authorization/permission"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

const (
	campaignDataFileName = "campaign.dataset.csv"
)

const dataLoadingTemaplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Template</title>
</head>

<body>
    <h2>Loading Error : {{.LoadingError}}</h2>

	<h2>Entries Unique Loaded On Firebase : {{.EntriesUniqueLoadedOnFirebase}}</h2>

	<h2>Entries Unique Loaded On CRM : {{.EntriesUniqueLoadedOnCRM}}</h2>

	<h2>Total Entries Found On File : {{.TotalEntriesFoundOnFile}}</h2>

	<h3>Entries</h3>

	<table style="width:100%">
	<tr>
		<th>Identifier</th>
		<th>HasLoadedToFirebase</th>
		<th>HasBeenRollBackFromFirebase</th>
		<th>HasLoadedToCRM</th>
		<th>FirebaseLoadError</th>
		<th>CRMLoadError</th>
	</tr>
	{{range .Entries}}
		<tr>
			<td>{{.Identifier}}</td>
			<td>{{.HasLoadedToFirebase}}</td>
			<td>{{.HasBeenRollBackFromFirebase}}</td>
			<td>{{.HasLoadedToCRM}}</td>
			<td>{{.FirebaseLoadError}}</td>
			<td>{{.CRMLoadError}}</td>			
		</tr>
    {{end}}	
	</table>   

</body>

</html>
`

// MarketingDataUseCases represents all the marketing data business logic
type MarketingDataUseCases interface {
	GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*apiclient.Segment, error)
	UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error
	BeWellAware(ctx context.Context, email string) error
	LoadCampaignDataset(ctx context.Context, phone string, emails []string)
	GetUserMarketingData(ctx context.Context, phonenumber string) (*apiclient.Segment, error)
}

// MarketingDataImpl represents the marketing usecase implementation
type MarketingDataImpl struct {
	repository repository.Repository
	hubspot    hubspot.ServiceHubSpotInterface
	mail       *mail.Service
}

// NewMarketing initialises a marketing usecase
func NewMarketing(
	repository repository.Repository, hubspot hubspot.ServiceHubSpotInterface, mail *mail.Service,
) *MarketingDataImpl {
	return &MarketingDataImpl{
		repository: repository,
		hubspot:    hubspot,
		mail:       mail,
	}
}

// GetMarketingData fetches all the marketing data from a collection
func (m MarketingDataImpl) GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*apiclient.Segment, error) {
	ctx, span := tracer.Start(ctx, "GetMarketingData")
	defer span.End()
	segmentData, err := m.repository.RetrieveMarketingData(ctx, data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to retrieve the marketing data")
	}

	return segmentData, nil
}

// UpdateUserCRMEmail updates a user CRM contact with the supplied email
func (m MarketingDataImpl) UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error {
	ctx, span := tracer.Start(ctx, "UpdateUserCRMEmail")
	defer span.End()
	CRMContactProperties := domain.ContactProperties{
		Email: email,
	}

	if err := m.repository.UpdateUserCRMEmail(ctx, phonenumber, &dto.UpdateContactPSMessage{
		Properties: CRMContactProperties,
		Phone:      phonenumber,
	}); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("failed to create CRM staging payload %v", err)
	}
	return nil

}

//BeWellAware toggles the user identified by the provided email= as bewell-aware on the CRM
func (m MarketingDataImpl) BeWellAware(ctx context.Context, email string) error {
	ctx, span := tracer.Start(ctx, "BeWellAware")
	defer span.End()
	CRMContactProperties := domain.ContactProperties{
		BeWellAware: domain.GeneralOptionTypeYes,
	}

	if err := m.repository.UpdateUserCRMBewellAware(ctx, email, &dto.UpdateContactPSMessage{
		Properties: CRMContactProperties,
	}); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("failed to set user as BeWell Aware %v", err)
	}
	return nil
}

// LoadCampaignDataset publishes the campaign dataset into firestore and CRM
// this call is idempotent,hence its safeguards against duplicates on both firstore and CRM
// this function expects that the CSV containing the dataset has been preprocessed and contains
// attributes expected by firestore and CRM
// column order mattters. Should be as;
// be_well_enrolled
// opt_out,
// be_well_aware
//be_well_persona
// has_wellness_card
// has_cover
// payor
// first_channel_of_contact
// initial_segment
// has_virtual_card
// email
// phone_number
// firstname
// lastname
// wing
// message_sent
// payer_slade_code
// member_number
func (m MarketingDataImpl) LoadCampaignDataset(ctx context.Context, phone string, emails []string) {
	ctx, span := tracer.Start(ctx, "LoadCampaignDataset")
	defer span.End()
	logrus.Info("loading campaign dataset started")
	res := dto.MarketingDataLoadOutput{}
	res.StartedAt = time.Now()

	sendMail := func(data dto.MarketingDataLoadOutput) {
		for _, email := range emails {
			logrus.Infof("sending load campaign dataset processing output to  %v", email)
			data.StoppedAt = time.Now()
			data.HoursTaken = data.StartedAt.Sub(data.StoppedAt).Hours()

			t := template.Must(template.New("LoadCampaignDataset").Parse(dataLoadingTemaplate))
			buf := new(bytes.Buffer)
			_ = t.Execute(buf, res)
			content := buf.String()
			if _, _, err := m.mail.SendEmail(
				ctx,
				"Load Campaign Dataset Processing",
				content,
				nil,
				email,
			); err != nil {
				helpers.RecordSpanError(span, err)
				logrus.Errorf("failed to send Load Campaign Dataset Processing email: %v", err)
			}
		}
	}

	if p := converterandformatter.StringSliceContains(base.AuthorizedPhones, phone); !p {
		res.LoadingError = fmt.Errorf("not authorized to access this resource")
		sendMail(res)
		return
	}
	isAuthorized, err := authorization.IsAuthorized(&profileutils.UserInfo{
		PhoneNumber: phone,
	}, permission.LoadMarketingData)
	if err != nil {
		helpers.RecordSpanError(span, err)
		res.LoadingError = err
		sendMail(res)
		return
	}
	if !isAuthorized {
		res.LoadingError = fmt.Errorf("user not authorized to access this resource")
		sendMail(res)
		return
	}

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, campaignDataFileName)

	csvFile, err := os.Open(path)
	if err != nil {
		helpers.RecordSpanError(span, err)
		res.LoadingError = fmt.Errorf("error opening the CSV file: %w", err)
		sendMail(res)
		return
	}
	defer csvFile.Close()

	csvContent, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		res.LoadingError = fmt.Errorf("failed to read data from the CSV file :%w", err)
		sendMail(res)
		return
	}

	csvPayload := csvContent[1:]

	for idx, line := range csvPayload {
		logrus.Infof("processing entry of identifier %v : Index %v", line[10], idx)

		entry := dto.MarketingDataLoadEntriesOutput{}

		data := apiclient.Segment{
			Properties: domain.ContactProperties{
				BeWellEnrolled:        domain.GeneralOptionType(line[0]),
				OptOut:                domain.GeneralOptionType(line[1]),
				BeWellAware:           domain.GeneralOptionType(line[2]),
				BeWellPersona:         domain.Persona(line[3]),
				HasWellnessCard:       domain.GeneralOptionType(line[4]),
				HasCover:              domain.GeneralOptionType(line[5]),
				Payor:                 domain.Payor(line[6]),
				FirstChannelOfContact: domain.ChannelOfContact(line[7]),
				InitialSegment:        line[8],
				HasVirtualCard:        domain.GeneralOptionType(line[9]),
				Email:                 line[10],
				Phone:                 line[11],
				FirstName:             line[12],
				LastName:              line[13],
			},
			Wing:           line[14],
			MessageSent:    line[15],
			PayerSladeCode: line[16],
			MemberNumber:   line[17],
		}

		// publish to firestore
		i, err := m.repository.LoadMarketingData(ctx, data)
		if i == 0 && err != nil {
			helpers.RecordSpanError(span, err)
			entry.FirebaseLoadError = fmt.Errorf("%v", err)
			entry.HasLoadedToFirebase = false
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)
			continue
		}

		// record has been created on firestore or it exists. proceeed to crm now
		if (i == -1 || i == 1) && err == nil {
			entry.HasLoadedToFirebase = true
			entry.Identifier = line[10]
		}

		resH, err := m.hubspot.SearchContactByPhone(line[11])
		if err != nil {
			helpers.RecordSpanError(span, err)
			// Do not roll back
			// We shall sync this data further
			// _ = m.repository.RollBackMarketingData(ctx, data)
			entry.HasBeenRollBackFromFirebase = true
			entry.HasLoadedToCRM = false
			entry.CRMLoadError = err
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)
			continue
		}

		// contact already exists. Nothing to do here
		if len(resH.Results) >= 1 {
			entry.HasLoadedToCRM = false
			entry.CRMLoadError = fmt.Errorf("contact already exists on the CRM")
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)
			continue
		}

		resS, err := m.hubspot.CreateContact(domain.CRMContact{
			Properties: domain.ContactProperties{
				BeWellEnrolled:        domain.GeneralOptionType(line[0]),
				OptOut:                domain.GeneralOptionType(line[1]),
				BeWellAware:           domain.GeneralOptionType(line[2]),
				BeWellPersona:         domain.Persona(line[3]),
				HasWellnessCard:       domain.GeneralOptionType(line[4]),
				HasCover:              domain.GeneralOptionType(line[5]),
				Payor:                 domain.Payor(line[6]),
				FirstChannelOfContact: domain.ChannelOfContact(line[7]),
				HasVirtualCard:        domain.GeneralOptionType(line[9]),
				Email:                 line[10],
				Phone:                 line[11],
				FirstName:             line[12],
				LastName:              line[13],
			},
		})

		if err != nil {
			helpers.RecordSpanError(span, err)
			// Do not roll back
			// We shall sync this data further
			// _ = m.repository.RollBackMarketingData(ctx, data)
			entry.HasBeenRollBackFromFirebase = true
			entry.HasLoadedToCRM = false
			entry.CRMLoadError = err
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)

			// protection from hubspot rate limiting
			time.Sleep(time.Millisecond * 500)
			continue
		}

		// means the record exists on crm
		if resS == nil && err == nil {
			entry.HasLoadedToCRM = true
			entry.CRMLoadError = nil
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)

			// protection from hubspot rate limiting
			time.Sleep(time.Millisecond * 500)
			continue
		}

		// means a new record has been created on the crm
		if resS != nil && err == nil {
			// append entries that have loaded on both firebase and crm
			entry.HasLoadedToCRM = true
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)

			// protection from hubspot rate limiting
			time.Sleep(time.Millisecond * 500)
			continue
		}

	}

	res.TotalEntriesFoundOnFile = len(csvPayload)
	res.EntriesUniqueLoadedOnFirebase = func(ents []dto.MarketingDataLoadEntriesOutput) int {
		count := 0
		for _, e := range ents {
			if e.HasBeenRollBackFromFirebase && !e.HasBeenRollBackFromFirebase {
				count++
			}
		}
		return 0
	}(res.Entries)

	res.EntriesUniqueLoadedOnCRM = func(ents []dto.MarketingDataLoadEntriesOutput) int {
		count := 0
		for _, e := range ents {
			if e.HasLoadedToCRM {
				count++
			}
		}
		return 0
	}(res.Entries)

	sendMail(res)
}

// GetUserMarketingData is used to retrieve the data of a targeted slader using their phonenumber
func (m MarketingDataImpl) GetUserMarketingData(ctx context.Context, phonenumber string) (*apiclient.Segment, error) {
	sladerData, err := m.repository.GetSladerDataByPhone(ctx, phonenumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the user's marketing data: %w", err)
	}

	return sladerData, nil
}
