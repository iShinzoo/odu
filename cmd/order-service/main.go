package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iShinzoo/odu/internal/config"
	"github.com/iShinzoo/odu/internal/db"
	"github.com/iShinzoo/odu/internal/order"
	"github.com/iShinzoo/odu/internal/worker"
	"github.com/iShinzoo/odu/pkg/logger"
	orderpb "github.com/iShinzoo/odu/proto"
	"google.golang.org/grpc"
)

type grpcServer struct {
	orderpb.UnimplementedOrderServiceServer
	service *order.OrderService
	pool    *worker.Pool
}

func (s *grpcServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {

	order, err := s.service.CreateOrder(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	// async processing
	s.pool.Submit(worker.Job{
		OrderID: order.ID,
	})

	return &orderpb.CreateOrderResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}

func (s *grpcServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {

	order, err := s.service.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrderResponse{
		OrderId: order.ID,
		UserId:  order.UserID,
		Amount:  order.Amount,
		Status:  order.Status,
	}, nil
}

func main() {
	// Initialize logger
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Order gRPC service")

	// Load configuration
	cfg := config.LoadConfig()

	// connect to database
	database, err := db.NewPostgresConnection(cfg.DBUrl)
	if err != nil {
		logger.Log.Fatal("Database connection failed", logger.ZapError(err))
	}
	defer database.Close()

	// setup repository and service
	repo := order.NewPostgresOrderRepository(database)
	service := order.NewOrderService(repo)

	// setup worker pool
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	pool := worker.NewPool(service)
	pool.Start(workerCtx, 5)

	// start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("Failed to listen", logger.ZapError(err))
	}

	gServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(gServer, &grpcServer{
		service: service,
		pool:    pool,
	})

	logger.Log.Info("Order gRPC service is running on port 50051")

	// Graceful shutdown handling
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		logger.Log.Info("Shutting down Order gRPC service")
		cancelWorkers()
		gServer.GracefulStop()
	}()

	if err := gServer.Serve(lis); err != nil {
		logger.Log.Fatal("Failed to serve", logger.ZapError(err))
	}
}
