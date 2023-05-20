package server

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/hexennacht/process-gmail-attachments/config"
	"github.com/hexennacht/process-gmail-attachments/handler"
	"github.com/hexennacht/process-gmail-attachments/pkg"
)

func Serve(conf *config.Configuration) {
	app := newServer(conf)

	app.Get("/_health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "success",
		})
	})

	oauthConfig := pkg.NewOauth2Config(conf.GoogleClientID, conf.GoogleClientSecret, conf.GoogleOAuth2RedirectURL)

	handler.RegisterLoginHandler(app, oauthConfig, pkg.RandomString(pkg.DefaultRandomStringLength), conf.GoogleOAuth2TokenFile)

	log.Fatalln(app.Listen(conf.BaseURL))
}

func newServer(conf *config.Configuration) *fiber.App {
	return fiber.New(fiber.Config{
		Prefork:           true,
		AppName:           conf.AppName,
		JSONEncoder:       jsoniter.Marshal,
		JSONDecoder:       jsoniter.Unmarshal,
		EnablePrintRoutes: true,
		ColorScheme:       fiber.DefaultColors,
	})
}

func newGmailService(conf *config.Configuration) (*gmail.Service, error) {
	httpClient := option.WithHTTPClient(pkg.NewOauth2Client(conf.GoogleOAuth2TokenFile))

	return gmail.NewService(context.Background(), httpClient)
}
