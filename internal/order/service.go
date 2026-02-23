package order

import (
	"context"

	"github.com/google/uuid"
)

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, amount float64) (*Order, error) {
	order := &Order{
		ID:     uuid.New().String(),
		UserID: userID,
		Amount: amount,
		Status: "CREATED",
	}

	err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*Order, error) {
	return s.repo.GetById(ctx, id)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
