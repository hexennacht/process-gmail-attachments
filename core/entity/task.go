package entity

import "google.golang.org/api/gmail/v1"

type TaskProcessMessageListRequest struct {
	*ReadFromRequest
	Messages *gmail.ListMessagesResponse
}
