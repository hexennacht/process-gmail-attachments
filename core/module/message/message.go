package message

import (
	"context"
	"fmt"
	"github.com/hexennacht/process-gmail-attachments/core/entity"
	"github.com/hexennacht/process-gmail-attachments/pkg"
	"github.com/hibiken/asynq"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/api/gmail/v1"
	"log"
	"strings"
)

type Module interface {
	ReadFrom(ctx context.Context, req *entity.ReadFromRequest) error

	FetchMessageFromList(ctx context.Context, req *entity.TaskFetchMessageFromListRequest) error
	FetchAttachmentFromMessage(ctx context.Context, req *entity.TaskFetchAttachmentFromMessageRequest) error
}

type module struct {
	gmailService    *gmail.Service
	task            *asynq.Client
	userID          string
	allowedMimeType map[string]bool
}

func NewModule(gmailService *gmail.Service, task *asynq.Client, userID, mimeType string) Module {
	allowedMimeType := make(map[string]bool)
	for _, mime := range strings.Split(mimeType, "|") {
		allowedMimeType[mime] = true
	}

	return &module{gmailService: gmailService, userID: userID, task: task, allowedMimeType: allowedMimeType}
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

	err = m.createTaskDefaultPriority(ctx, pkg.TaskProcessMessageList, &entity.TaskFetchMessageFromListRequest{
		ReadFromRequest: req,
		Messages:        listMessage,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *module) FetchMessageFromList(ctx context.Context, req *entity.TaskFetchMessageFromListRequest) error {
	for _, msg := range req.Messages.Messages {
		message, err := m.gmailService.Users.Messages.Get(m.userID, msg.Id).Do()
		if err != nil {
			log.Printf("cannot fetch message with id: %s with error: %+v\n", msg.Id, err)
			continue
		}

		if !m.isMessageHaveHeaderWithNameAndValue(message, "From", req.EmailFrom) {
			log.Printf("message with id: %s is not from %s\n", msg.Id, req.EmailFrom)
			continue
		}

		if len(message.Payload.Parts) <= 1 {
			log.Printf("message with id: %s didn't have attachments\n", msg.Id)
			continue
		}

		if err = m.queueMessageAttachments(ctx, message); err != nil {
			log.Printf("failed to queue attachment with message id: %s\n", msg.Id)
			continue
		}

		_, err = m.gmailService.Users.Messages.
			Modify(
				m.userID,
				msg.Id,
				&gmail.ModifyMessageRequest{
					RemoveLabelIds: []string{"UNREAD"},
				}).
			Do()
		if err != nil {
			log.Printf("failed to update message with id %s status to read\n", msg.Id)
		}
	}

	return nil
}

func (m *module) FetchAttachmentFromMessage(ctx context.Context, req *entity.TaskFetchAttachmentFromMessageRequest) error {
	//TODO implement me
	panic("implement me")
}

func (m *module) queueMessageAttachments(ctx context.Context, message *gmail.Message) error {
	for _, attachment := range message.Payload.Parts {
		if !m.allowedMimeType[attachment.MimeType] {
			log.Printf("message with id: %s didn't have attachment with allowed mime type please check on config env\n", message.Id)
			continue
		}

		err := m.createTaskDefaultPriority(ctx, pkg.TaskProcessMessageAttachment, &entity.TaskFetchAttachmentFromMessageRequest{
			MessageID:    message.Id,
			AttachmentID: attachment.Body.AttachmentId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *module) createTaskDefaultPriority(ctx context.Context, taskName string, taskRequest interface{}) error {
	taskBody, err := jsoniter.Marshal(taskRequest)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskName, taskBody, asynq.Queue(pkg.TaskQueueDefault))

	_, err = m.task.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (m *module) isMessageHaveHeaderWithNameAndValue(message *gmail.Message, headerName, value string) bool {
	var isHeaderExists = make(map[string]*gmail.MessagePartHeader)
	for _, header := range message.Payload.Headers {
		isHeaderExists[header.Name] = header
	}

	if isHeaderExists[headerName] == nil {
		return false
	}

	return strings.EqualFold(isHeaderExists[headerName].Value, value)
}
