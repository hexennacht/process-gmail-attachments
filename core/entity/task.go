package entity

import "google.golang.org/api/gmail/v1"

type TaskFetchMessageFromListRequest struct {
	*ReadFromRequest
	Messages *gmail.ListMessagesResponse
}

type TaskFetchAttachmentFromMessageRequest struct {
	MessageID    string
	AttachmentID string
}
