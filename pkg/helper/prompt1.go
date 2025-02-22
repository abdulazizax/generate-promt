package helper

import (
	"context"
	"fmt"
	"generate-promt-v1/api/models"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/sheets/v4"
)

func ExportFunctionalRequirementsToSheet(srv *sheets.Service, spreadsheetID, sheetName string, parsedResponse *models.ProjectResponse) error {
	if err := WriteAllRows(srv, spreadsheetID, sheetName, parsedResponse); err != nil {
		return fmt.Errorf("failed to append rows: %w", err)
	}

	if err := ApplyStylingAndMerges(srv, spreadsheetID, sheetName, parsedResponse); err != nil {
		return fmt.Errorf("failed to apply styling: %w", err)
	}

	return nil
}
func WriteAllRows(srv *sheets.Service, spreadsheetID, sheetName string, parsedResponse *models.ProjectResponse) error {
	var rows [][]interface{}
	rows = append(rows, []interface{}{"#", "Epic", "Story", "Task"})

	rowNumber := 1
	for _, epic := range parsedResponse.FunctionalRequirements {
		for _, story := range epic.Stories {
			for _, task := range story.Tasks {
				rows = append(rows, []interface{}{
					rowNumber,
					epic.Epic,
					story.Story,
					task,
				})
				rowNumber++
			}
		}
	}

	targetRange := fmt.Sprintf("%s!A1", sheetName)
	valueRange := &sheets.ValueRange{
		Values: rows,
	}

	_, err := srv.Spreadsheets.Values.Append(spreadsheetID, targetRange, valueRange).
		ValueInputOption("RAW").
		Context(context.Background()).
		Do()
	return err
}
func ApplyStylingAndMerges(srv *sheets.Service, spreadsheetID, sheetName string, parsedResponse *models.ProjectResponse) error {
	sheetID, err := findSheetID(srv, spreadsheetID, sheetName)
	if err != nil {
		return fmt.Errorf("findSheetID: %w", err)
	}
	if sheetID < 0 {
		return fmt.Errorf("sheet %q not found", sheetName)
	}

	var requests []*sheets.Request

	currentRow := int64(1)

	for _, epic := range parsedResponse.FunctionalRequirements {
		epicStart := currentRow

		var epicRowCount int64
		for _, story := range epic.Stories {
			epicRowCount += int64(len(story.Tasks))
		}

		if epicRowCount > 0 {
			requests = append(requests, buildMergeRequest(sheetID, epicStart, epicRowCount, 1, 2, "MERGE_ALL"))

			requests = append(requests, buildUpdateCellFormatRequest(sheetID, epicStart, epicRowCount, 1, 2, &sheets.CellFormat{
				BackgroundColor: &sheets.Color{
					Red:   0.80,
					Green: 0.90,
					Blue:  0.90,
				},
				TextFormat: &sheets.TextFormat{
					Bold: true,
				},
				HorizontalAlignment: "CENTER",
				VerticalAlignment:   "MIDDLE",
			}))
		}

		for _, story := range epic.Stories {
			storyRowCount := int64(len(story.Tasks))
			if storyRowCount == 0 {
				continue
			}

			// The story starts at currentRow
			storyStart := currentRow

			// Merge the story cells in column C (index 2) across storyRowCount rows
			requests = append(requests, buildMergeRequest(sheetID, storyStart, storyRowCount, 2, 3, "MERGE_ALL"))

			// Style the story cells (slightly different color)
			requests = append(requests, buildUpdateCellFormatRequest(sheetID, storyStart, storyRowCount, 2, 3, &sheets.CellFormat{
				BackgroundColor: &sheets.Color{
					Red:   0.85,
					Green: 0.85,
					Blue:  0.95, // Example: a light purple tint
				},
				TextFormat: &sheets.TextFormat{
					Bold: true,
				},
				HorizontalAlignment: "CENTER",
				VerticalAlignment:   "MIDDLE",
			}))

			// Tasks are in column D (index 3). We won't merge them, but we can style them if you like.
			// We'll skip styling tasks for now.

			// Move currentRow pointer
			currentRow += storyRowCount
		}
	}

	// Also style the header row (row 0, columns 0..4)
	headerRequest := buildUpdateCellFormatRequest(sheetID, 0, 1, 0, 4, &sheets.CellFormat{
		BackgroundColor: &sheets.Color{
			Red:   0.2,
			Green: 0.5,
			Blue:  0.8, // a deeper blue
		},
		TextFormat: &sheets.TextFormat{
			Bold: true,
			ForegroundColor: &sheets.Color{
				Red:   1.0,
				Green: 1.0,
				Blue:  1.0, // White text
			},
		},
		HorizontalAlignment: "CENTER",
		VerticalAlignment:   "MIDDLE",
	})
	requests = append(requests, headerRequest)

	// Optionally auto-resize columns A-D
	autoResizeCols := &sheets.Request{
		AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
			Dimensions: &sheets.DimensionRange{
				SheetId:    sheetID,
				Dimension:  "COLUMNS",
				StartIndex: 0,
				EndIndex:   4, // columns A-D
			},
		},
	}
	requests = append(requests, autoResizeCols)

	batchReq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}
	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, batchReq).Do()
	return err
}
func buildMergeRequest(sheetID, startRow, rowCount, startCol, endCol int64, mergeType string) *sheets.Request {
	return &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    startRow,
				EndRowIndex:      startRow + rowCount,
				StartColumnIndex: startCol,
				EndColumnIndex:   endCol,
			},
			MergeType: mergeType,
		},
	}
}
func findSheetID(srv *sheets.Service, spreadsheetID, sheetName string) (int64, error) {
	ss, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return -1, err
	}
	for _, sheet := range ss.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return -1, nil
}
func buildUpdateCellFormatRequest(sheetID, startRow, rowCount, startCol, endCol int64, format *sheets.CellFormat) *sheets.Request {
	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    startRow,
				EndRowIndex:      startRow + rowCount,
				StartColumnIndex: startCol,
				EndColumnIndex:   endCol,
			},
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredFormat: format,
						},
					},
				},
			},
			Fields: "userEnteredFormat(backgroundColor,textFormat,horizontalAlignment,verticalAlignment)",
		},
	}
}

func ExportProjectDataToDoc(service *docs.Service, documentID string, projectBrief *struct {
	ProjectGoal       string `json:"project_goal"`
	PrimaryObjectives string `json:"primary_objectives"`
	ExpectedOutcomes  string `json:"expected_outcomes"`
	SuccessMetrics    string `json:"success_metrics"`
}) error {
	_, err := service.Documents.Get(documentID).Do()
	if err != nil {
		return fmt.Errorf("failed to read doc %s: %w", documentID, err)
	}

	textToInsert := fmt.Sprintf(
		"\nThe first promt\n\nProject Goal: %s\nPrimary Objectives: %s\nExpected Outcomes: %s\nSuccess Metrics: %s\n\n\n\n",
		projectBrief.ProjectGoal,
		projectBrief.PrimaryObjectives,
		projectBrief.ExpectedOutcomes,
		projectBrief.SuccessMetrics,
	)

	requests := []*docs.Request{
		{
			InsertText: &docs.InsertTextRequest{
				EndOfSegmentLocation: &docs.EndOfSegmentLocation{},
				Text:                 textToInsert,
			},
		},
	}

	_, err = service.Documents.
		BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
			Requests: requests,
		}).
		Do()

	if err != nil {
		return fmt.Errorf("failed to update doc %s: %w", documentID, err)
	}

	return nil
}
