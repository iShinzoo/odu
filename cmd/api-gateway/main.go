package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iShinzoo/odu/pkg/logger"
	orderpb "github.com/iShinzoo/odu/proto"
	"google.golang.org/grpc"
)

type Gateway struct {
	client orderpb.OrderServiceClient
}

type CreateOrderRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
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

func (g *Gateway) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	var req CreateOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.client.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId: req.UserID,
		Amount: req.Amount,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

