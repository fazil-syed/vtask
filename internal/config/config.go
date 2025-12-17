package config

import (
	"os"
	"strconv"
)
type Config struct {
	UploadPath    string
	DatabasePath  string
	MaxUploadSize int64
	STTProvider   string
	ServerPort   int
}
func LoadConfig() (*Config, error) {
	maxUploadSizeStr := os.Getenv("MAX_UPLOAD_SIZE")
	maxUploadSize, err := strconv.ParseInt(maxUploadSizeStr, 10, 64)
	if err != nil {
		maxUploadSize = 10 * 1024 * 1024 // default to 10 MB
	}
	UploadPath := os.Getenv("UPLOAD_PATH")
	if UploadPath == "" {
		UploadPath = "./uploads"
	}
	DatabasePath := os.Getenv("DATABASE_PATH")
	if DatabasePath == "" {
		DatabasePath = "app.db"
	}
	STTProvider := os.Getenv("STT_PROVIDER")
	if STTProvider == "" {
		STTProvider = "stub"
	}
	ServerPortStr := os.Getenv("SERVER_PORT")
	ServerPort, err := strconv.Atoi(ServerPortStr)
	if err != nil {
		ServerPort = 8080 // default to port 8080
	}
	config := &Config{
		UploadPath:    UploadPath,
		DatabasePath:  DatabasePath,
		MaxUploadSize: maxUploadSize,
		STTProvider:   STTProvider,
		ServerPort:   ServerPort,
	}
	return config, nil
}
	