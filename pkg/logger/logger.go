package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init() error {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	_ = Log.Sync()
}

// helper for structured error logging
func ZapError(err error) zap.Field {
	return zap.Error(err)
}
