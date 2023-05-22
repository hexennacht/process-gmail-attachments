package queue

import (
	"context"
	"github.com/hexennacht/process-gmail-attachments/core/entity"
	jsoniter "github.com/json-iterator/go"

	"github.com/hibiken/asynq"
	"golang.org/x/oauth2"

	messageModule "github.com/hexennacht/process-gmail-attachments/core/module/message"
	"github.com/hexennacht/process-gmail-attachments/pkg"
)

type messageHandler struct {
	app            *asynq.ServeMux
	oauthConfig    *oauth2.Config
	messageService messageModule.Module

	tokenFilePath string
	state         string
}

func RegisterMessageHandler(app *asynq.ServeMux, oauthConfig *oauth2.Config, messageService messageModule.Module) {
	handler := newMessageHandler(app, oauthConfig, messageService)
	handler.app.HandleFunc(pkg.TaskProcessMessageList, handler.ProcessMessagesList)
	handler.app.HandleFunc(pkg.TaskProcessMessageAttachment, handler.ProcessMessageAttachment)
}

func newMessageHandler(app *asynq.ServeMux, oauthConfig *oauth2.Config, messageService messageModule.Module) *messageHandler {
	return &messageHandler{app: app, oauthConfig: oauthConfig, messageService: messageService}
}

func (m *messageHandler) ProcessMessagesList(ctx context.Context, t *asynq.Task) error {
	var req *entity.TaskFetchMessageFromListRequest

	if err := jsoniter.Unmarshal(t.Payload(), &req); err != nil {
		return err
	}

	if err := m.messageService.FetchMessageFromList(ctx, req); err != nil {
		return err
	}

	return nil
}

func (m *messageHandler) ProcessMessageAttachment(ctx context.Context, t *asynq.Task) error {
	var req *entity.TaskFetchMessageFromListRequest

	if err := jsoniter.Unmarshal(t.Payload(), &req); err != nil {
		return err
	}

	if err := m.messageService.FetchMessageFromList(ctx, req); err != nil {
		return err
	}

	return nil
}
