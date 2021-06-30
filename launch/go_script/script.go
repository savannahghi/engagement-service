package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
)

// This script is used to read data from a CSV and push it to a collection

// ReadCSVFile ..
func ReadCSVFile(path string) ([]dto.Segment, error) {
	var data []dto.Segment
	// firstRow := true
	// var column map[string]int

	csvFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the CSV file: %w", err)
	}
	defer csvFile.Close()

	csvContent, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read data from the CSV file :%w", err)
	}

	for _, line := range csvContent {
		segmentData := dto.Segment{
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
			IsSynced:              line[16],
			TimeSynced:            line[17],
		}
		data = append(data, segmentData)
	}

	return data, nil
}

func main() {
	ReadCSVFile("/home/sala/Documents/sil_dry_run_data.csv")
}
