package pkg

import (
	"bufio"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
	"os"
)

func NewOauth2Config(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       []string{gmail.MailGoogleComScope},
	}
}

func NewOauth2Client(tokenFilePath string) *http.Client {
	token, err := readOauth2TokenFromFile(tokenFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	source := oauth2.StaticTokenSource(token)

	return oauth2.NewClient(context.Background(), oauth2.ReuseTokenSource(token, source))
}

func readOauth2TokenFromFile(tokenFilePath string) (*oauth2.Token, error) {
	tokenByte, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := jsoniter.Unmarshal(tokenByte, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func WriteTokenToFile(tokenFilePath string, token *oauth2.Token) (err error) {
	var f *os.File

	if !checkFileExists(tokenFilePath) {
		f, err = os.Create(tokenFilePath)
		if err != nil {
			return err
		}
	}

	if checkFileExists(tokenFilePath) {
		f, err = os.OpenFile(tokenFilePath, os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}
	}

	defer f.Close()

	writer := bufio.NewWriter(f)

	fileToken, err := jsoniter.Marshal(&token)
	if err != nil {
		return err
	}

	_, err = writer.Write(fileToken)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func checkFileExists(tokenFilePath string) bool {
	if _, err := os.Stat(tokenFilePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
