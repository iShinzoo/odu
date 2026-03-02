package worker

import (
	"context"
	"time"

	"github.com/iShinzoo/odu/internal/order"
	"github.com/iShinzoo/odu/internal/ws"
	"github.com/iShinzoo/odu/pkg/logger"
	"go.uber.org/zap"
)

type Pool struct {
	jobQueue chan Job
	service  *order.OrderService
	hub      *ws.Hub
}

func NewPool(service *order.OrderService, hub *ws.Hub) *Pool {
	return &Pool{
		jobQueue: make(chan Job, 100), // buffered channel to hold jobs
		service:  service,
		hub:      hub,
	}
}

func (p *Pool) Start(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go p.worker(ctx, i)
	}
}

func (p *Pool) worker(ctx context.Context, id int) {
	logger.Log.Info("worker started", zap.Int("workerId", id))

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Worker Shutting down", zap.Int("WorkerID", id))
			return
		case job := <-p.jobQueue:
			logger.Log.Info("Processing Job",
				zap.Int("WorkerID", id),
				zap.String("orderID", job.OrderID),
			)

			time.Sleep(3 * time.Second)

			err := p.service.UpdateOrderStatus(ctx, job.OrderID, "PROCESSED")
			if err == nil {
				logger.Log.Info("Worker notifying websocket",
					zap.String("orderID", job.OrderID),
				)
				p.hub.Notify(job.OrderID, "PROCESSED")
			} else {
				logger.Log.Error("Failed to update order status", zap.Error(err))
			}
		}
	}
}

func (p *Pool) Submit(job Job) {
	p.jobQueue <- job
}
