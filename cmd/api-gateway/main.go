package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iShinzoo/odu/pkg/logger"
	orderpb "github.com/iShinzoo/odu/proto"
	"google.golang.org/grpc"
)

type Gateway struct {
	client orderpb.OrderServiceClient
}

func main() {

	// connect to grpc server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	client := orderpb.NewOrderServiceClient(conn)

	gateway := &Gateway{
		client: client,
	}

	r := chi.NewRouter()

	// define routes
	r.Post("/orders", gateway.CreateOrderHandler)
	r.Get("/orders/{id}", gateway.GetOrderHandler)

	logger.Log.Info("API Gateway is running on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
