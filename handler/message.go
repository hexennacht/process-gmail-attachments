package handler

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/gmail/v1"
)

type messageHandler struct {
	app          *fiber.App
	gmailService *gmail.Service

	tokenFilePath string
	state         string
}

func RegisterMessageHandler(app *fiber.App, gmailService *gmail.Service) {
	handler := newMessageHandler(app, gmailService)

	oauth := app.Group("/message")

	oauth.Get("/attachments", handler.Attachments)
}

func newMessageHandler(app *fiber.App, gmailService *gmail.Service) *messageHandler {
	return &messageHandler{app: app, gmailService: gmailService}
}

func (h *messageHandler) Attachments(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
