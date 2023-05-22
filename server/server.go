package server

import (
	"context"
	"github.com/hexennacht/process-gmail-attachments/handler/queue"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/hexennacht/process-gmail-attachments/config"
	"github.com/hexennacht/process-gmail-attachments/core/module/message"
	"github.com/hexennacht/process-gmail-attachments/handler/http"
	"github.com/hexennacht/process-gmail-attachments/pkg"
)

func Serve(conf *config.Configuration) {
	httpServer := newHttpServer(conf)
	queueServer := newQueueServer(conf.RedisURL)
	queueClient := asynq.NewClient(&asynq.RedisClientOpt{Addr: conf.RedisURL})

	httpServer.Get("/_health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "success",
		})
	})

	oauthConfig := pkg.NewOauth2Config(conf.GoogleClientID, conf.GoogleClientSecret, conf.GoogleOAuth2RedirectURL)
	gmailService, err := newGmailService(conf)
	if err != nil {
		log.Fatalln(err)
	}

	messageSvc := message.NewModule(gmailService, queueClient, conf.GoogleEmail)

	http.RegisterLoginHandler(httpServer, oauthConfig, pkg.RandomString(pkg.DefaultRandomStringLength), conf.GoogleOAuth2TokenFile)
	http.RegisterMessageHandler(httpServer, messageSvc)

	go func() {
		mux := asynq.NewServeMux()

		queue.RegisterMessageHandler(mux, oauthConfig, messageSvc)

		if err := queueServer.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	log.Fatalln(httpServer.Listen(conf.BaseURL))
}

func newHttpServer(conf *config.Configuration) *fiber.App {
	return fiber.New(fiber.Config{
		Prefork:           true,
		AppName:           conf.AppName,
		JSONEncoder:       jsoniter.Marshal,
		JSONDecoder:       jsoniter.Unmarshal,
		EnablePrintRoutes: true,
		ColorScheme:       fiber.DefaultColors,
	})
}

func newQueueServer(redisAddr string) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				pkg.TaskQueueCritical: 6,
				pkg.TaskQueueDefault:  3,
				pkg.TaskQueueLow:      1,
			},
			// See the godoc for other configuration options
			LogLevel: asynq.ErrorLevel,
		},
	)
}

func newGmailService(conf *config.Configuration) (*gmail.Service, error) {
	httpClient := option.WithHTTPClient(pkg.NewOauth2Client(conf.GoogleOAuth2TokenFile))

	return gmail.NewService(context.Background(), httpClient)
}
