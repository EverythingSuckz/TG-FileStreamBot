package utils

import (
	"time"

	"go.uber.org/zap"
)

// TimeFuncWithResult wraps a function call and logs its execution time, returning a result
func TimeFuncWithResult[T any](log *zap.Logger, funcName string, fn func() (T, error)) (T, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		log.Info(funcName+" executed", zap.Int64("duration_ms", duration))
	}()
	return fn()
}
