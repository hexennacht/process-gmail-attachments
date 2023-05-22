package message

import (
	"context"
	"fmt"
	"github.com/hexennacht/process-gmail-attachments/core/entity"
	"github.com/hexennacht/process-gmail-attachments/pkg"
	"github.com/hibiken/asynq"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/api/gmail/v1"
)

type Module interface {
	ReadFrom(ctx context.Context, req *entity.ReadFromRequest) error

	ProcessMessageListToFetchAttachments(ctx context.Context, req *entity.TaskProcessMessageListRequest) error
}

type module struct {
	gmailService *gmail.Service
	task         *asynq.Client
	userID       string
}

func NewModule(gmailService *gmail.Service, task *asynq.Client, userID string) Module {
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

	err = m.createTaskDefaultPriority(ctx, pkg.TaskProcessMessageList, &entity.TaskProcessMessageListRequest{
		ReadFromRequest: req,
		Messages:        listMessage,
	})
	if err != nil {
		return err
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

func (m *module) ProcessMessageListToFetchAttachments(ctx context.Context, req *entity.TaskProcessMessageListRequest) error {
	//TODO implement me
	panic("implement me")
}
