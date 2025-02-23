package client

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryConfig defines the configuration for retry behavior
type RetryConfig struct {
	MaxAttempts      int           // Maximum number of retry attempts
	InitialBackoff   time.Duration // Initial backoff duration
	MaxBackoff       time.Duration // Maximum backoff duration
	BackoffMultipler float64       // Multiplier for exponential backoff
}

// DefaultRetryConfig provides reasonable default values for retry behavior
var DefaultRetryConfig = RetryConfig{
	MaxAttempts:      5,
	InitialBackoff:   100 * time.Millisecond,
	MaxBackoff:       5 * time.Second,
	BackoffMultipler: 2.0,
}

// shouldRetry determines if an error should trigger a retry
func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() { // Check for gRPC status codes to determine whether to retry 
	case codes.DeadlineExceeded,
		codes.Unavailable,
		codes.ResourceExhausted,
		codes.Aborted:
		return true
	default:
		return false
	}
}

// withRetry executes the given operation with retry logic
func withRetry[T any](ctx context.Context, config RetryConfig, operation func(context.Context) (T, error)) (T, error) {
	var lastErr error
	var result T
	backoff := config.InitialBackoff

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return result, ctx.Err()
			case <-timer.C:
			}
			backoff = time.Duration(float64(backoff) * config.BackoffMultipler)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}

		result, lastErr = operation(ctx)
		if lastErr == nil {
			return result, nil
		}

		if !shouldRetry(lastErr) {
			return result, lastErr
		}
	}

	return result, lastErr
}

// withRetryNoReturn executes void operations with retry logic
func withRetryNoReturn(ctx context.Context, config RetryConfig, operation func(context.Context) error) error {
	var lastErr error
	backoff := config.InitialBackoff

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			backoff = time.Duration(float64(backoff) * config.BackoffMultipler)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}

		lastErr = operation(ctx)
		if lastErr == nil {
			return nil
		}

		if !shouldRetry(lastErr) {
			return lastErr
		}
	}

	return lastErr
}
