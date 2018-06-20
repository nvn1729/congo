// Package latch provides a CountDownLatch concurrency primitive.
// A CountDownLatch is used to signal the completion of a specified number of events.
//
// In typical usage, a main goroutine creates the latch with a provided count and passes the latch to a number of goroutines.
// The goroutines invoke CountDown (or WeightedCountDown) on the latch to signal completion and reduce the remaining count on the latch.
// The main goroutine invokes Wait or WaitTimeout on the latch to wait for all events to complete.
//
// CountDownLatch is similar to sync.WaitGroup but differs in a few different ways:
// - It has the ability to wait with a timeout (WaitTimeout) for latch to complete.
// - It has ability to do a WeightedCountDown, i.e. reduce the remaining count by more than 1.
// - The count is set one time at latch creation instead of through a separate Add call, which may prevent misuse of Add in goroutines.
package latch

import (
	"time"
	"sync"
	"errors"
)

type CountDownLatch struct {
	m sync.Mutex
	remainingCount uint
	countDownCompleteCh chan struct{}
}

var ErrLatchCompleted = errors.New("Latch count down already complete")

// New creates a new CountDownLatch with the provided count.
// If the count is 0, the latch will be set to immediately signal count down completion for any goroutines that subsequently call Wait or WaitTimeout.
func New(count uint) *CountDownLatch {
	latch := &CountDownLatch{
		remainingCount: count,
		countDownCompleteCh: make(chan struct{}),
	}
	if latch.remainingCount == 0 {
		close(latch.countDownCompleteCh)
	}
	return latch
}

// Count returns the remaining count.
func (latch *CountDownLatch) Count() uint {
	latch.m.Lock()
	defer latch.m.Unlock()
	return latch.remainingCount
}

// CountDown invokes WeightedCountDown with a weight of 1.
func (latch *CountDownLatch) CountDown() error {
	return latch.WeightedCountDown(1)
}

// WeightedCountDown reduces the count on the CountDownLatch by the given weight.
// If given weight exceeds the current latch's count, the latch's count is set to 0.
// If the count hits 0, WeightedCountDown signals latch count down completion, waking up any goroutines waiting with Wait or WaitTimeout.
// An error, ErrLatchCompleted, is returned if count down has already been completed (Count() returns 0).
func (latch *CountDownLatch) WeightedCountDown(weight uint) error {
	latch.m.Lock()
	defer latch.m.Unlock()

	select {
	case <-latch.countDownCompleteCh:
		return ErrLatchCompleted
	default:
		if latch.remainingCount > weight {
			latch.remainingCount -= weight
		} else {
			latch.remainingCount = 0
			close(latch.countDownCompleteCh)
		}
		return nil
	}
}

// Wait waits indefinitely until the count down is completed.
// Wait returns immediately if the count down has already been completed.
func (latch *CountDownLatch) Wait() {
	<-latch.countDownCompleteCh
}

// WaitTimeout waits until a given timeout for the count down to complete.
// If the count down is completed before the timeout, WaitTimeout returns true.
// Otherwise it returns false.
// WaitTimeout returns immediatley if the count down has already been completed.
func (latch *CountDownLatch) WaitTimeout(timeout time.Duration) bool {
	select {
	case <-latch.countDownCompleteCh:
		return true
	case <-time.After(timeout):
		return false
	}
}

