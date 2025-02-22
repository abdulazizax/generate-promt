package helper

import (
	"context"

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

<<<<<<< HEAD
=======
	requests = append(requests, &sheets.Request{
		DeleteSheet: &sheets.DeleteSheetRequest{
			SheetId: createdSheet.Sheets[0].Properties.SheetId,
		},
	})

>>>>>>> 049446241bd00e2e2ad3459a6a6b1b83753e4d88
	for _, name := range sheetNames {
		requests = append(requests, &sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: name,
				},
			},
		})
	}
<<<<<<< HEAD
	requests = append(requests, &sheets.Request{
		DeleteSheet: &sheets.DeleteSheetRequest{
			SheetId: createdSheet.Sheets[0].Properties.SheetId,
		},
	})
=======

>>>>>>> 049446241bd00e2e2ad3459a6a6b1b83753e4d88
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

	return updatedSpreadsheet, nil
}
