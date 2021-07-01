package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/authorization"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/authorization/permission"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
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
    <h2>Loading Error, {{.LoadingError}}</h2>

	<h2>Entries Unique Loaded On Firebase, {{.EntriesUniqueLoadedOnFirebase}}</h2>

	<h2>Entries Unique Loaded On CRM, {{.EntriesUniqueLoadedOnCRM}}</h2>

	<h2>Total Entries Found On File, {{.TotalEntriesFoundOnFile}}</h2>

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

type MarketingDataUseCases interface {
	GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*dto.Segment, error)
	UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error
	BeWellAware(ctx context.Context, email string) error
	LoadCampaignDataset(ctx context.Context, phone string, email string)
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

func (m MarketingDataImpl) GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*dto.Segment, error) {
	segmentData, err := m.repository.RetrieveMarketingData(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the marketing data")
	}

	return segmentData, nil
}

// UpdateUserCRMEmail updates a user CRM contact with the supplied email
func (m MarketingDataImpl) UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error {
	CRMContactProperties := CRMDomain.ContactProperties{
		Email: email,
	}

	if err := m.repository.UpdateUserCRMEmail(ctx, phonenumber, &dto.UpdateContactPSMessage{
		Properties: CRMContactProperties,
		Phone:      phonenumber,
	}); err != nil {
		return fmt.Errorf("failed to create CRM staging payload %v", err)
	}
	return nil

}

//BeWellAware toggles the user identified by the provided email= as bewell-aware on the CRM
func (m MarketingDataImpl) BeWellAware(ctx context.Context, email string) error {
	CRMContactProperties := CRMDomain.ContactProperties{
		BeWellAware: CRMDomain.GeneralOptionTypeYes,
	}

	logrus.Printf("bewellAware payload: %v with email: %v", CRMContactProperties, email)

	if err := m.repository.UpdateUserCRMBewellAware(ctx, email, &dto.UpdateContactPSMessage{
		Properties: CRMContactProperties,
	}); err != nil {
		return fmt.Errorf("failed to create CRM staging payload %v", err)
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
func (m MarketingDataImpl) LoadCampaignDataset(ctx context.Context, phone string, email string) {
	logrus.Info("loading campaign dataset started")
	res := dto.MarketingDataLoadOutput{}

	sendMail := func(data dto.MarketingDataLoadOutput) {
		logrus.Infof("sending load campaign dataset processing output to  %v", email)

		t := template.Must(template.New("LoadCampaignDataset").Parse(dataLoadingTemaplate))
		buf := new(bytes.Buffer)
		_ = t.Execute(buf, res)
		content := buf.String()
		if _, _, err := m.mail.SendEmail(
			"Load Campaign Dataset Processing",
			content,
			nil,
			email,
		); err != nil {
			logrus.Errorf("failed to send Load Campaign Dataset Processing email: %v", err)
		}
	}

	if p := base.StringSliceContains(base.AuthorizedPhones, phone); !p {
		res.LoadingError = fmt.Errorf("not authorized to access this resource")
		sendMail(res)
		return
	}
	isAuthorized, err := authorization.IsAuthorized(&base.UserInfo{
		PhoneNumber: phone,
	}, permission.LoadMarketingData)
	if err != nil {
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
		res.LoadingError = fmt.Errorf("error opening the CSV file: %w", err)
		sendMail(res)
		return
	}
	defer csvFile.Close()

	csvContent, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		res.LoadingError = fmt.Errorf("failed to read data from the CSV file :%w", err)
		sendMail(res)
		return
	}

	csvPayload := csvContent[1:]

	for _, line := range csvPayload {

		entry := dto.MarketingDataLoadEntriesOutput{}

		data := dto.Segment{
			BeWellEnrolled:        line[0],
			OptOut:                line[1],
			BeWellAware:           line[2],
			BeWellPersona:         line[3],
			HasWellnessCard:       line[4],
			HasCover:              line[5],
			Payor:                 line[6],
			FirstChannelOfContact: line[7],
			InitialSegment:        line[8],
			HasVirtualCard:        line[9],
			Email:                 line[10],
			PhoneNumber:           line[11],
			FirstName:             line[12],
			LastName:              line[13],
			Wing:                  line[14],
			MessageSent:           line[15],
		}

		logrus.Infof("processing entry of identifier %v", line[10])

		// publish to firestore
		i, err := m.repository.LoadMarketingData(ctx, data)
		if i == 0 && err != nil {
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
			_ = m.repository.RollBackMarketingData(ctx, data)
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

		convertor := func(d string) CRMDomain.GeneralOptionType {
			if d == "YES" {
				return CRMDomain.GeneralOptionTypeYes
			}
			return CRMDomain.GeneralOptionTypeNo
		}

		convertorPersona := func(d string) CRMDomain.Persona {
			if d == "ALICE" {
				return CRMDomain.PersonaAlice
			}

			if d == "JUMA" {
				return CRMDomain.PersonaJuma
			}

			if d == "BOB" {
				return CRMDomain.PersonaBob
			}

			if d == "ANDREW" {
				return CRMDomain.PersonaAndrew
			}
			return CRMDomain.PersonaSlader
		}

		convertorPayor := func(d string) CRMDomain.Payor {
			if d == "RESOLUTION" || d == "RESOLUTION INSURANCE" {
				return CRMDomain.PayorResolution
			}

			if d == "APA" || d == "APA INSURANCE" {
				return CRMDomain.PayorApa
			}

			if d == "JUBILEE" || d == "JUBILEE INSURANCE" {
				return CRMDomain.PayorJubilee
			}

			if d == "BRITAM" || d == "BRITAM INSURANCE" {
				return CRMDomain.PayorBritam
			}

			if d == "MADISON" || d == "MADISON INSURANCE" {
				return CRMDomain.PayorMadison
			}

			return CRMDomain.PayorJubilee
		}

		convertorChannel := func(d string) CRMDomain.ChannelOfContact {
			if d == "APP" {
				return CRMDomain.ChannelOfContactApp
			}
			if d == "USSD" {
				return CRMDomain.ChannelOfContactUssd
			}
			return CRMDomain.ChannelOfContactShortcode
		}

		if _, err := m.hubspot.CreateContact(CRMDomain.CRMContact{
			Properties: CRMDomain.ContactProperties{
				BeWellEnrolled:        convertor(line[0]),
				OptOut:                convertor(line[1]),
				BeWellAware:           convertor(line[2]),
				BeWellPersona:         convertorPersona(line[3]),
				HasSladeID:            convertor(line[4]),
				HasCover:              convertor(line[5]),
				Payor:                 convertorPayor(line[6]),
				FirstChannelOfContact: convertorChannel(line[7]),
				HasVirtualCard:        convertor(line[9]),
				Email:                 line[10],
				Phone:                 line[11],
				FirstName:             line[12],
				LastName:              line[13],
			},
		}); err != nil {
			_ = m.repository.RollBackMarketingData(ctx, data)
			entry.HasBeenRollBackFromFirebase = true
			entry.HasLoadedToCRM = false
			entry.CRMLoadError = err
			entry.Identifier = line[10]
			res.Entries = append(res.Entries, entry)

			// protection from hubspot rate limiting
			time.Sleep(time.Second * 2)
			continue
		}

		// append entries that have loaded on both firebase and crm
		entry.HasLoadedToCRM = true
		entry.Identifier = line[10]
		res.Entries = append(res.Entries, entry)

		// protection from hubspot rate limiting
		time.Sleep(time.Second * 2)
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
