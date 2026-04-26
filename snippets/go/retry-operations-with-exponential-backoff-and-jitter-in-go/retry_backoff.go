package main

import (
	"context"
	"errors"
	"log"
	"math"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts    int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	JitterFraction float64
}

var randomFloat64 = rand.Float64

var waitForDelay = func(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Flow:
//
//	run the operation
//	   |
//	   +-> success -> return nil
//	   +-> non-retryable or last-attempt error -> return the error
//	   `-> retryable error -> wait for jittered backoff or stop early on context cancellation
func RetryOperation(ctx context.Context, operation func(context.Context) error, isRetryable func(error) bool, config RetryConfig) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if operation == nil {
		return errors.New("operation is required")
	}
	if isRetryable == nil {
		return errors.New("retryable predicate is required")
	}
	if config.MaxAttempts <= 0 {
		return errors.New("max attempts must be greater than 0")
	}
	if config.BaseDelay <= 0 {
		return errors.New("base delay must be greater than 0")
	}
	if config.MaxDelay < config.BaseDelay {
		return errors.New("max delay must be greater than or equal to base delay")
	}
	if config.JitterFraction < 0 || config.JitterFraction > 1 {
		return errors.New("jitter fraction must be between 0 and 1")
	}

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := operation(ctx)
		if err == nil {
			return nil
		}
		if !isRetryable(err) || attempt == config.MaxAttempts {
			return err
		}

		delay := backoffDelay(attempt, config.BaseDelay, config.MaxDelay, config.JitterFraction)
		if err := waitForDelay(ctx, delay); err != nil {
			return err
		}
	}

	return nil
}

func backoffDelay(attempt int, baseDelay, maxDelay time.Duration, jitterFraction float64) time.Duration {
	maxDelayFloat := float64(maxDelay)
	delayFloat := float64(baseDelay) * math.Pow(2, float64(attempt-1))
	if delayFloat > maxDelayFloat {
		delayFloat = maxDelayFloat
	}

	jitterMultiplier := 1 + ((randomFloat64()*2)-1)*jitterFraction
	jitteredDelayFloat := delayFloat * jitterMultiplier
	if jitteredDelayFloat < 0 {
		jitteredDelayFloat = 0
	}
	if jitteredDelayFloat > maxDelayFloat {
		jitteredDelayFloat = maxDelayFloat
	}
	return time.Duration(jitteredDelayFloat)
}

var errTemporary = errors.New("temporary failure")

func main() {
	attempts := 0
	err := RetryOperation(
		context.Background(),
		func(context.Context) error {
			attempts++
			if attempts < 3 {
				return errTemporary
			}
			return nil
		},
		func(err error) bool {
			return errors.Is(err, errTemporary)
		},
		RetryConfig{
			MaxAttempts:    4,
			BaseDelay:      20 * time.Millisecond,
			MaxDelay:       80 * time.Millisecond,
			JitterFraction: 0.25,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("attempts=%d", attempts)
}
