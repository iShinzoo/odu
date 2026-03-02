package interceptor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iShinzoo/odu/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ctxKey string

const RequestIDKey ctxKey = "requestID"

func UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	start := time.Now()

	resp, err := handler(ctx, req)

	logger.Log.Info("grpc request",
		zap.String("method", info.FullMethod),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return resp, err
}

func UnaryRecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Panic recovered in gRPC",
				zap.Any("panic", r),
			)
			err = grpc.Errorf(grpc.Code(err), "Internal Server Error")
		}
	}()

	return handler(ctx, req)
}

func UnaryRequestIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, RequestIDKey, requestID)

	return handler(ctx, req)
}
