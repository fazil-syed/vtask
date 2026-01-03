package main

import (
	"fmt"
	"log"
	"os"

	"github.com/syed.fazil/vtask/internal/config"
	"github.com/syed.fazil/vtask/internal/db"
	"github.com/syed.fazil/vtask/internal/handlers"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/server"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
		log.Fatalf("failed to create upload dir %s: %v", cfg.UploadPath, err)
	}
	// initialize db
	db, err := db.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// auto migrate task model
	err = db.AutoMigrate(&models.Task{}, &models.VoiceNote{}, &models.User{}, &models.Identity{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS identities_password_email_unique
		ON identities (email)
		WHERE issuer = 'password';
	`)
	// initialize stt provider
	var sttProvider handlers.STTService
	switch cfg.STTProvider {
	case "stub":
		sttProvider = &handlers.StubSTTService{Transcription: "remind me to finish this project tommorow"}
	default:
		log.Fatalf("Unsupported STT provider: %s", cfg.STTProvider)
	}
	// setup server
	router := server.SetupServer(db, sttProvider, cfg.UploadPath, cfg.MaxUploadSize)
	// start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
