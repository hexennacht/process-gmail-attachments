package config

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	AppName string `envconfig:"APP_NAME" default:"process-gmail-attachment"`
	Mode    string `envconfig:"MODE" default:"development"`
	BaseURL string `envconfig:"BASE_URL" default:"127.0.0.1:8080"`

	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgresql://user:password@127.0.0.1:5432/gmail_attachments"`
	RedisURL    string `envconfig:"REDIS_URL" default:"localhost:6379"`

	GoogleClientID          string `envconfig:"GOOGLE_CLIENT_ID" default:""`
	GoogleClientSecret      string `envconfig:"GOOGLE_CLIENT_SECRET" default:""`
	GoogleOAuth2RedirectURL string `envconfig:"GOOGLE_OAUTH2_REDIRECT_URL" default:""`
	GoogleOAuth2TokenFile   string `envconfig:"GOOGLE_OAUTH2_TOKEN_FILE" default:""`
	GoogleEmail             string `envconfig:"GOOGLE_EMAIL" default:"me"`

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	SupportedFileMimeType string `envconfig:"SUPPORTED_FILE_MIME_TYPE" default:"application/gzip|application/zip|text/csv|application/pdf"`
}

func Read() *Configuration {
	var configuration Configuration
	envconfig.MustProcess("", &configuration)

	return &configuration
}
