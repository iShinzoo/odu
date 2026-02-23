package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/iShinzoo/odu/internal/config"
	"github.com/iShinzoo/odu/internal/db"
	"github.com/iShinzoo/odu/internal/order"
	"github.com/iShinzoo/odu/internal/worker"
	"github.com/iShinzoo/odu/pkg/logger"
	"github.com/lib/pq"
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

	// create repository and service instances
	repo := order.NewPostgresOrderRepository(database)
	service := order.NewOrderService(repo)

	// create context for workers
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	// start worker pool
	pool := worker.NewPool(service)
	pool.Start(workerCtx, 5)

	logger.Log.Info("Worker pool started")

	// simulate order creation
	orderCtx, orderCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer orderCancel()

	userID := uuid.New().String()

	_, err = database.ExecContext(orderCtx,
		`INSERT INTO users (id, name, email)
	VALUES ($1, $2, $3)
	`, userID, "krsnaAA", "krsna@example.ZOHOIN")

	if err != nil {
		if isDuplicateError(err) {
			logger.Log.Warn("user already exists", zapError(err))
		} else {
			logger.Log.Fatal("Failed to create user", zapError(err))
		}
	}

	newOrder, err := service.CreateOrder(orderCtx, userID, 250.50)
	if err != nil {
		logger.Log.Fatal("Failed to create order", zapError(err))
	}

	logger.Log.Info("Order Created", zap.String("OrderID", newOrder.ID))

	pool.Submit(worker.Job{
		OrderID: newOrder.ID,
	})

	waitForShutdown()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

// helper for structured error logging
func zapError(err error) zap.Field {
	return zap.Error(err)
}

func isDuplicateError(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505"
	}
	return false
}
