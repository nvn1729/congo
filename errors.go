package congo

import "errors"

// These are errors related to CountDownLatch.
var (
	// ErrCountDownLatchCompleted is returned when CountDown or WeightedCountDown is called after the count down is already complete
	ErrCountDownLatchCompleted = errors.New("Latch count down already complete")
)
	

