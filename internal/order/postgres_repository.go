package order

import (
	"context"
	"database/sql"
	"time"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *Order) error {

	query := `INSERT INTO orders (id, user_id, amount, status, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		order.ID,
		order.UserID,
		order.Amount,
		order.Status,
		time.Now(),
	)
	return err
}

func (r *PostgresOrderRepository) GetById(ctx context.Context, id string) (*Order, error) {

	query := `SELECT id, user_id, amount, status, created_at FROM orders WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var order Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {

	query := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
