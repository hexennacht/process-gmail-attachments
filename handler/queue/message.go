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
}

func newMessageHandler(app *asynq.ServeMux, oauthConfig *oauth2.Config, messageService messageModule.Module) *messageHandler {
	return &messageHandler{app: app, oauthConfig: oauthConfig, messageService: messageService}
}

func (m *messageHandler) ProcessMessagesList(ctx context.Context, t *asynq.Task) error {
	var req *entity.TaskProcessMessageListRequest

	if err := jsoniter.Unmarshal(t.Payload(), &req); err != nil {
		return err
	}

	if err := m.messageService.ProcessMessageListToFetchAttachments(ctx, req); err != nil {
		return err
	}

	return nil
}
