package http

import (
	"github.com/hexennacht/process-gmail-attachments/core/entity"
	messageModule "github.com/hexennacht/process-gmail-attachments/core/module/message"

	"github.com/gofiber/fiber/v2"
)

type messageHandler struct {
	app            *fiber.App
	messageService messageModule.Module

	tokenFilePath string
	state         string
}

func RegisterMessageHandler(app *fiber.App, messageService messageModule.Module) {
	handler := newMessageHandler(app, messageService)

	message := app.Group("/message")
	message.Post("/read-from", handler.ReadFrom)

	attachment := message.Group("attachments")
	attachment.Get("/", handler.Attachments)
}

func newMessageHandler(app *fiber.App, messageService messageModule.Module) *messageHandler {
	return &messageHandler{app: app, messageService: messageService}
}

func (h *messageHandler) ReadFrom(ctx *fiber.Ctx) error {
	var req entity.ReadFromRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.DefaultErrorHandler(ctx.Status(fiber.StatusUnprocessableEntity), err)
	}

	if err := h.messageService.ReadFrom(ctx.Context(), &req); err != nil {
		return fiber.DefaultErrorHandler(ctx.Status(fiber.StatusInternalServerError), err)
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (h *messageHandler) Attachments(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
