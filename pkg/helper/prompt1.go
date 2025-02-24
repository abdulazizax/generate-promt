package helper

import (
	"context"
	"fmt"
	"generate-promt-v1/api/models"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/sheets/v4"
)

// ExportFunctionalRequirementsToSheet writes a two-row header, a blank row,
// then your epic/story/task data starting from row 4. Epics are numbered in column B.
func ExportFunctionalRequirementsToSheet(
	srv *sheets.Service,
	spreadsheetID, sheetName string,
	parsedResponse *models.ProjectResponse,
) error {

	if err := WriteAllRows(srv, spreadsheetID, sheetName, parsedResponse); err != nil {
		return fmt.Errorf("failed to append rows: %w", err)
	}

	if err := ApplyStylingAndMerges(srv, spreadsheetID, sheetName, parsedResponse); err != nil {
		return fmt.Errorf("failed to apply styling: %w", err)
	}

	return nil
}

// WriteAllRows:
//
//	Row 1:  A="Version", B=">1m",   C="",       D=">2w",  E=">1w"
//	Row 2:  A="Месяц",   B="#",     C="Epic",   D="Story",E="Task"
//	Row 3:  (blank)
//	Row 4+: data
//
// We will store an “epicNumber” in column B for each Epic, rather than
// numbering each Task.
func WriteAllRows(
	srv *sheets.Service,
	spreadsheetID, sheetName string,
	parsedResponse *models.ProjectResponse,
) error {

	var rows [][]interface{}

	// ---------------------------
	// Header Row 1 (A..E)
	// ---------------------------
	// A: "Version"
	// B..C: merged in styling => ">1m"
	// D: ">2w"
	// E: ">1w"
	rows = append(rows, []interface{}{
		"Version", // A1
		">1m",     // B1
		"",        // C1 (will merge with B1 in styling)
		">2w",     // D1
		">1w",     // E1
	})

	// ---------------------------
	// Header Row 2 (A..E)
	// ---------------------------
	// A: "Месяц"
	// B: "#"
	// C: "Epic"
	// D: "Story"
	// E: "Task"
	rows = append(rows, []interface{}{
		"Месяц",
		"#",
		"Epic",
		"Story",
		"Task",
	})

	// ---------------------------
	// Blank Row 3 (A..E)
	// ---------------------------
	rows = append(rows, []interface{}{"", "", "", "", ""})

	// ---------------------------
	// Data rows start at row 4
	// ---------------------------
	// We'll store epicNumber in column B, epic text in C, story in D, tasks in E.
	epicNumber := 1

	for _, epic := range parsedResponse.FunctionalRequirements {
		// For each Epic, we append as many rows as the total number of tasks in that epic.
		// We'll fill column B with the same epicNumber, column C with the same epic text,
		// but actually they will be merged in styling.
		for _, story := range epic.Stories {
			for _, task := range story.Tasks {
				rows = append(rows, []interface{}{
					"",          // A: (Месяц) blank or fill with e.g. "Feb 24"
					epicNumber,  // B: Epic # (same for all tasks under this epic)
					epic.Epic,   // C: Epic text
					story.Story, // D: Story text
					task,        // E: Task
				})
			}
		}
		// Increment epicNumber for the next epic
		epicNumber++
	}

	// Now we use the Sheets API to append these rows to the target sheet.
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

func ApplyStylingAndMerges(
	srv *sheets.Service,
	spreadsheetID, sheetName string,
	parsedResponse *models.ProjectResponse,
) error {

	sheetID, err := findSheetID(srv, spreadsheetID, sheetName)
	if err != nil {
		return fmt.Errorf("findSheetID: %w", err)
	}
	if sheetID < 0 {
		return fmt.Errorf("sheet %q not found", sheetName)
	}

	var requests []*sheets.Request

	// ------------------------------------------------------------
	// 1) Merge header row 1 (B1..C1) => ">1m"
	// ------------------------------------------------------------
	requests = append(requests, buildMergeRequest(
		sheetID,
		0, 1, // row=0, rowCount=1
		1, 3, // col=1..3 merges B1..C1
		"MERGE_ALL",
	))

	// ------------------------------------------------------------
	// 2) Style row 1 (A1..E1) - dark background, white bold text
	// ------------------------------------------------------------
	headerRow1Format := &sheets.CellFormat{
		BackgroundColor: &sheets.Color{Red: 0.05, Green: 0.20, Blue: 0.30}, // dark
		TextFormat: &sheets.TextFormat{
			Bold:            true,
			ForegroundColor: &sheets.Color{Red: 1, Green: 1, Blue: 1}, // white
		},
		HorizontalAlignment: "CENTER",
		VerticalAlignment:   "MIDDLE",
	}
	requests = append(requests, buildUpdateCellFormatRequest(
		sheetID,
		0, 1, // row 0..1
		0, 5, // col 0..5 => A..E
		headerRow1Format,
	))

	// ------------------------------------------------------------
	// 3) Style row 2 (A2..E2) - medium background, white bold text
	// ------------------------------------------------------------
	headerRow2Format := &sheets.CellFormat{
		BackgroundColor: &sheets.Color{Red: 0.2, Green: 0.5, Blue: 0.8}, // teal-ish
		TextFormat: &sheets.TextFormat{
			Bold:            true,
			ForegroundColor: &sheets.Color{Red: 1, Green: 1, Blue: 1}, // white
		},
		HorizontalAlignment: "CENTER",
		VerticalAlignment:   "MIDDLE",
	}
	requests = append(requests, buildUpdateCellFormatRequest(
		sheetID,
		1, 2, // row 1..2
		0, 5, // col 0..5 => A..E
		headerRow2Format,
	))

	// Row 3 is blank; we do not style it specifically (you can if you wish).

	// ------------------------------------------------------------
	// Data starts at row=3 in zero-based index (that's row 4 visually).
	// We'll merge column B (the epic number) and column C (the epic text)
	// for all tasks belonging to each epic. We'll also merge column D
	// for each story. Column E is for tasks (no merge).
	// ------------------------------------------------------------
	currentRow := int64(3) // zero-based => the first data row is row 3

	// We need to track the same epicNumber logic used in WriteAllRows.
	epicNumber := 1

	for _, epic := range parsedResponse.FunctionalRequirements {
		// Count how many total tasks for this epic
		var epicRowCount int64
		for _, story := range epic.Stories {
			epicRowCount += int64(len(story.Tasks))
		}

		if epicRowCount == 0 {
			// If an epic has no tasks, skip merges.
			continue
		}

		epicStart := currentRow

		// Merge column B for epic number (row epicStart..epicStart+epicRowCount)
		requests = append(requests, buildMergeRequest(
			sheetID,
			epicStart, epicRowCount,
			1, 2, // column B => col index=1
			"MERGE_ALL",
		))
		// Merge column C for epic text
		requests = append(requests, buildMergeRequest(
			sheetID,
			epicStart, epicRowCount,
			2, 3, // column C => col index=2
			"MERGE_ALL",
		))
		// Style epic columns B & C
		epicFormat := &sheets.CellFormat{
			BackgroundColor: &sheets.Color{
				Red:   0.74, // approximate teal color
				Green: 0.82,
				Blue:  0.91,
			},
			TextFormat: &sheets.TextFormat{
				Bold: true,
			},
			HorizontalAlignment: "CENTER",
			VerticalAlignment:   "MIDDLE",
		}
		requests = append(requests, buildUpdateCellFormatRequest(
			sheetID,
			epicStart, epicRowCount,
			1, 3,
			epicFormat,
		))

		// Now do merges for each story in column D
		storyStart := epicStart
		for _, story := range epic.Stories {
			storyRowCount := int64(len(story.Tasks))
			if storyRowCount == 0 {
				continue
			}
			// Merge column D for these tasks
			requests = append(requests, buildMergeRequest(
				sheetID,
				storyStart, storyRowCount,
				3, 4, // col D => index=3
				"MERGE_ALL",
			))
			// Style story
			storyFormat := &sheets.CellFormat{
				BackgroundColor: &sheets.Color{
					Red:   0.87,
					Green: 0.81,
					Blue:  0.94, // lighter purple/pink
				},
				TextFormat: &sheets.TextFormat{
					Bold: true,
				},
				HorizontalAlignment: "CENTER",
				VerticalAlignment:   "MIDDLE",
			}
			requests = append(requests, buildUpdateCellFormatRequest(
				sheetID,
				storyStart, storyRowCount,
				3, 4,
				storyFormat,
			))

			storyStart += storyRowCount
		}

		currentRow += epicRowCount
		epicNumber++
	}

	// ------------------------------------------------------------
	// Finally, auto-resize columns A..E
	// ------------------------------------------------------------
	requests = append(requests, &sheets.Request{
		AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
			Dimensions: &sheets.DimensionRange{
				SheetId:    sheetID,
				Dimension:  "COLUMNS",
				StartIndex: 0,
				EndIndex:   5, // columns A=0 .. E=4
			},
		},
	})

	// Execute batch update
	batchReq := &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}
	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, batchReq).Do()
	return err
}

// findSheetID returns the sheet ID for the named sheet
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

// buildMergeRequest merges the cell range [startRow..startRow+rowCount)
// in columns [startCol..endCol).
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

// buildUpdateCellFormatRequest applies the given CellFormat to all cells
// in [startRow..startRow+rowCount) x [startCol..endCol).
func buildUpdateCellFormatRequest(
	sheetID, startRow, rowCount, startCol, endCol int64,
	format *sheets.CellFormat,
) *sheets.Request {
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

// ExportProjectDataToDoc remains the same as your original:
func ExportProjectDataToDoc(
	service *docs.Service,
	documentID string,
	projectBrief *struct {
		ProjectGoal       string `json:"project_goal"`
		PrimaryObjectives string `json:"primary_objectives"`
		ExpectedOutcomes  string `json:"expected_outcomes"`
		SuccessMetrics    string `json:"success_metrics"`
	},
) error {
	// ...
	return nil
}
