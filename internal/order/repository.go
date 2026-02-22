package order

import "context"

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetById(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
