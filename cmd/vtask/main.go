package main

import (
	"fmt"
	"log"
	"os"

	"github.com/syed.fazil/vtask/internal/config"
	"github.com/syed.fazil/vtask/internal/db"
	"github.com/syed.fazil/vtask/internal/server"
)

func main() {
	// load config
	config.Init()
	if err := os.MkdirAll(config.App.UploadPath, 0755); err != nil {
		log.Fatalf("failed to create upload dir %s: %v", config.App.UploadPath, err)
	}
	cfg := config.App
	// initialize db
	dbConn, err := db.Open(db.Config{
		Path:  cfg.DatabasePath,
		Debug: true,
	})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	// auto migrate task model
	err = db.Migrate(dbConn)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	// setup server
	router := server.New(server.Config{
		DB:            dbConn,
		UploadPath:    cfg.UploadPath,
		MaxUploadSize: cfg.MaxUploadSize,
	})
	// start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	log.Fatal(router.Run(addr))
}
