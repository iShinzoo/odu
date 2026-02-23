package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iShinzoo/odu/internal/config"
	"github.com/iShinzoo/odu/internal/db"
	"github.com/iShinzoo/odu/internal/order"
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

	//logic part
	repo := order.NewPostgresOrderRepository(database)
	service := order.NewOrderService(repo)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := uuid.New().String()

	_, err = database.ExecContext(ctx,
		`INSERT INTO users (id, name, email)
	VALUES ($1, $2, $3)
	`, userID, "krishna", "krsna@example.com")

	if err != nil {
		logger.Log.Fatal("Failed to create user", zap.Error(err))
	}

	newOrder, err := service.CreateOrder(ctx, userID, 250.50)
	if err != nil {
		logger.Log.Fatal("Failed to create order", zap.Error(err))
	}

	logger.Log.Info("Order Created", zap.String("OrderID", newOrder.ID))

}

// helper for structured error logging
func zapError(err error) zap.Field {
	return zap.Error(err)
}
