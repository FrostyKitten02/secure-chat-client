package env

import (
	"github.com/joho/godotenv"
	"os"
)

type Options struct {
	WebServiceUrl string
	WebSocketUrl  string
}

var Opts Options

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env file not found")
	}
	webServiceURL := os.Getenv("WEB_SERVICE_URL")
	Opts.WebServiceUrl = webServiceURL

	webSocketURL := os.Getenv("WEB_SOCKET_URL")
	Opts.WebSocketUrl = webSocketURL
}
