//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

var ErrStopped = errors.New("limiter stopped")

type RateLimiter struct {
	timeWindow  time.Duration
	timerPool   chan []*time.Timer
	stopChannel chan struct{}
	maxRequests int
}

func NewLimiter(maxRequests int, timeWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		timeWindow:  timeWindow,
		timerPool:   createTimerPool(maxRequests),
		stopChannel: make(chan struct{}),
		maxRequests: maxRequests,
	}
}

func createTimerPool(maxRequests int) chan []*time.Timer {
	timers := make([]*time.Timer, maxRequests)
	for i := range timers {
		timers[i] = time.NewTimer(0)
	}
	timerChannel := make(chan []*time.Timer, 1)
	timerChannel <- timers
	return timerChannel
}

func (rl *RateLimiter) Acquire(ctx context.Context) error {
	if err := rl.checkIfStopped(); err != nil {
		return err
	}

	if err := rl.checkContextCancellation(ctx); err != nil {
		return err
	}

	return rl.waitForAvailableSlot(ctx)
}

func (rl *RateLimiter) checkIfStopped() error {
	select {
	case <-rl.stopChannel:
		return ErrStopped
	default:
		return nil
	}
}

func (rl *RateLimiter) checkContextCancellation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (rl *RateLimiter) waitForAvailableSlot(ctx context.Context) error {
	for {
		select {
		case <-rl.stopChannel:
			return ErrStopped
		case <-ctx.Done():
			return ctx.Err()
		case timers := <-rl.timerPool:
			if rl.tryAcquireSlot(timers) {
				return nil
			}
		}
	}
}

func (rl *RateLimiter) tryAcquireSlot(timers []*time.Timer) bool {
	defer func() {
		rl.timerPool <- timers
	}()

	for i, timer := range timers {
		if isSlotReady(timer) {
			timers[i] = time.NewTimer(rl.timeWindow)
			return true
		}
	}
	return false
}

func isSlotReady(timer *time.Timer) bool {
	select {
	case <-timer.C:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChannel)
}
