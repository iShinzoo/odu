package main

import (
	"github.com/iShinzoo/odu/internal/config"
	"github.com/iShinzoo/odu/internal/db"
	"github.com/iShinzoo/odu/pkg/logger"
	"go.uber.org/zap"
)

func main() {

	// Initialize the logger
	err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting the Application...")

	// Load Config
	cfg := config.LoadConfig()

	// connect to the database
	database, err := db.NewPostgresConnection(cfg.DBUrl)
	if err != nil {
		logger.Log.Fatal("DB Connection failed", zapError(err))
	}
	defer database.Close()

	logger.Log.Info("Database started successfully")

	logger.Log.Info("Application started")
}

// helper for structured error logging
func zapError(err error) zap.Field {
	return zap.Error(err)
}
