package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var App *Config

func Init() {
	_ = godotenv.Load()

	viper.AutomaticEnv()
	// Defaults
	viper.SetDefault("UPLOAD_PATH", "./uploads")
	viper.SetDefault("DATABASE_PATH", "app.db")
	viper.SetDefault("MAX_UPLOAD_SIZE", int64(10*1024*1024))
	viper.SetDefault("STT_PROVIDER", "stub")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("APP_BASE_URL", "http://localhost:8080")
	viper.SetDefault("JWT_SECRET", "JWT_SECRET")
	viper.SetDefault("NLP_SERVER", "localhost:50051")
	viper.SetDefault("APP_UI_BASE_URL", "http://localhost:5173")

	cfg := &Config{
		UploadPath:              viper.GetString("UPLOAD_PATH"),
		DatabasePath:            viper.GetString("DATABASE_PATH"),
		MaxUploadSize:           viper.GetInt64("MAX_UPLOAD_SIZE"),
		STTProvider:             viper.GetString("STT_PROVIDER"),
		ServerPort:              viper.GetInt("SERVER_PORT"),
		OAuthGoogleClientID:     viper.GetString("OAUTH_GOOGLE_CLIENT_ID"),
		OAuthGoogleClientSecret: viper.GetString("OAUTH_GOOGLE_CLIENT_SECRET"),
		AppBaseURL:              viper.GetString("APP_BASE_URL"),
		JWTSecret:               viper.GetString("JWT_SECRET"),
		NLPServer:               viper.GetString("NLP_SERVER"),
		APPUIBaseURL:            viper.GetString("APP_UI_BASE_URL"),
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	App = cfg
}
