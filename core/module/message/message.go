package message

import (
	"context"
	"fmt"
	"github.com/hexennacht/process-gmail-attachments/core/entity"
	"google.golang.org/api/gmail/v1"
)

type Module interface {
	ReadFrom(ctx context.Context, req *entity.ReadFromRequest) error
}

type module struct {
	gmailService *gmail.Service

	userID string
}

func NewModule(gmailService *gmail.Service, userID string) Module {
	return &module{gmailService: gmailService, userID: userID}
}

func (m *module) ReadFrom(ctx context.Context, req *entity.ReadFromRequest) error {
	listMessage, err := m.gmailService.Users.Messages.
		List(m.userID).
		LabelIds("INBOX").
		Q(fmt.Sprintf("is:unread has:attachment from:%s subject:%s", req.EmailFrom, req.EmailSubject)).Context(ctx).
		Do()
	if err != nil {
		return err
	}

	_ = listMessage

	return nil
}
