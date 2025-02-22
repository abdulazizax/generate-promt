package helper

import (
	"context"
	"fmt"

	"google.golang.org/api/docs/v1"
)

// CreateNewDoc creates a new Google Doc with the given title.
func CreateNewDoc(docsService *docs.Service, title string) (*docs.Document, error) {
	doc, err := docsService.Documents.Create(&docs.Document{
		Title: title,
	}).Context(context.Background()).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create doc: %v", err)
	}

	err = makeFilePublic(doc.DocumentId)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
