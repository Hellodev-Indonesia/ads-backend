package config

import "log"

var CentrifugoConfig struct {
	APIKey      string
	TokenSecret string
	URL         string
}

func InitCentrifugo() {
	CentrifugoConfig.APIKey = GetEnv("CENTRIFUGO_API_KEY", "")
	CentrifugoConfig.TokenSecret = GetEnv("CENTRIFUGO_TOKEN_SECRET", "")
	CentrifugoConfig.URL = GetEnv("CENTRIFUGO_URL", "http://localhost:8000")
	log.Printf(
		"Centrifugo config loaded: url=%s, api_key_set=%t, token_secret_set=%t",
		CentrifugoConfig.URL,
		CentrifugoConfig.APIKey != "",
		CentrifugoConfig.TokenSecret != "",
	)
}
