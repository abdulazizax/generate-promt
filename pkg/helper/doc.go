package helper

import (
	"context"
	"fmt"

	"google.golang.org/api/docs/v1"
)

// CreateNewDoc
func CreateNewDoc(docsService *docs.Service, title string, tabNames []string) (*docs.Document, error) {
	doc, err := docsService.Documents.Create(&docs.Document{
		Title: title,
	}).Context(context.Background()).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create doc: %v", err)
	}

	var requests []*docs.Request

	for _, name := range tabNames {
		requests = append(requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				EndOfSegmentLocation: &docs.EndOfSegmentLocation{
					SegmentId: "",
				},
				Text: fmt.Sprintf("%s\n\n", name),
			},
		})
	}

	_, err = docsService.Documents.BatchUpdate(doc.DocumentId, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Context(context.Background()).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to batch-update doc: %v", err)
	}

	return doc, nil
}
