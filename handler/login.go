package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hexennacht/process-gmail-attachments/pkg"
	"golang.org/x/oauth2"
	"net/http"
)

type loginHandler struct {
	app         *fiber.App
	oauthConfig *oauth2.Config

	tokenFilePath string
	state         string
}

func RegisterLoginHandler(app *fiber.App, oauthConfig *oauth2.Config, state, tokenFilePath string) {
	handler := newLoginHandler(app, oauthConfig, state, tokenFilePath)

	oauth := app.Group("/oauth2")

	oauth.Get("/login", handler.Login)
	oauth.Get("/redirect", handler.Callback)
}

func newLoginHandler(app *fiber.App, oauthConfig *oauth2.Config, state string, tokenFilePath string) *loginHandler {
	return &loginHandler{app: app, oauthConfig: oauthConfig, state: state, tokenFilePath: tokenFilePath}
}

func (h *loginHandler) Login(ctx *fiber.Ctx) error {
	url := h.oauthConfig.AuthCodeURL(
		h.state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("response", "ACCEPT"),
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (h *loginHandler) Callback(ctx *fiber.Ctx) error {
	code := ctx.Get("code")

	token, err := h.oauthConfig.Exchange(ctx.Context(), code)
	if err != nil {
		return fiber.DefaultErrorHandler(ctx.Status(fiber.StatusUnprocessableEntity), err)
	}

	if err := pkg.WriteTokenToFile(h.tokenFilePath, token); err != nil {
		return fiber.DefaultErrorHandler(ctx.Status(fiber.StatusInternalServerError), err)
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
