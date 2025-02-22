package helper

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"google.golang.org/api/sheets/v4"
)

func CreateNewSpreadsheet(srv *sheets.Service, title string, sheetNames []string) (*sheets.Spreadsheet, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}

	createdSheet, err := srv.Spreadsheets.Create(spreadsheet).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	requests := []*sheets.Request{}

	for _, name := range sheetNames {
		requests = append(requests, &sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: name,
				},
			},
		})
	}
	requests = append(requests, &sheets.Request{
		DeleteSheet: &sheets.DeleteSheetRequest{
			SheetId: createdSheet.Sheets[0].Properties.SheetId,
		},
	})
	if len(requests) > 0 {
		rb := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}
		_, err = srv.Spreadsheets.BatchUpdate(createdSheet.SpreadsheetId, rb).
			Context(context.Background()).Do()
		if err != nil {
			return nil, err
		}
	}

	updatedSpreadsheet, err := srv.Spreadsheets.Get(createdSheet.SpreadsheetId).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}
	err = makeFilePublic(updatedSpreadsheet.SpreadsheetId)
	if err != nil {
		return nil, err
	}
	return updatedSpreadsheet, nil
}

func makeFilePublic(fileID string) error {
	driveService, err := drive.NewService(
		context.Background(),
		option.WithCredentialsFile("service_account.json"),
	)
	if err != nil {
		return fmt.Errorf("unable to create Drive client: %v", err)
	}

	perm := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err = driveService.Permissions.Create(fileID, perm).Do()
	if err != nil {
		return fmt.Errorf("unable to set file as public: %v", err)
	}

	return nil
}
