package helper

import (
	"fmt"

	"google.golang.org/api/docs/v1"
)

func ExportDataToDoc(service *docs.Service, documentID string, resp string) error {
	_, err := service.Documents.Get(documentID).Do()
	if err != nil {
		return fmt.Errorf("failed to read doc %s: %w", documentID, err)
	}

	tab := "############################################################################"

	requests := []*docs.Request{
		{
			InsertText: &docs.InsertTextRequest{
				EndOfSegmentLocation: &docs.EndOfSegmentLocation{},
				Text:                 resp + "\n\n" + tab + "\n\n",
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
