package latch

import (
	"testing"
	"time"
	//"fmt"
)

func TestZero(t *testing.T) {
	latch := New(0)

	// check 0 count
	assertEqual(t, uint(0), latch.Count())

	// following should result in errors
	assertNotNil(t, latch.CountDown())
	assertNotNil(t, latch.WeightedCountDown(0))

	// Wait should not block
	latch.Wait()
	
	// WaitTimeout should not block and should return true
	assertEqual(t, true, latch.WaitTimeout(time.Second))
}

func TestOne(t *testing.T) {
	latch := New(1)

	// check count of 1
	assertEqual(t, uint(1), latch.Count())

	// test should pause a half second, return v as false, no err
	assertEqual(t, false, latch.WaitTimeout(500*time.Millisecond))
	
	// count down
	assertNil(t, latch.CountDown())

	// repeat tests from TestCountDownZero
	// check 0 count
	assertEqual(t, uint(0), latch.Count())

	// following should result in errors
	assertNotNil(t, latch.CountDown())
	assertNotNil(t, latch.WeightedCountDown(0))

	// Wait should not block
	latch.Wait()
	
	// WaitTimeout should not block and should return true
	assertEqual(t, true, latch.WaitTimeout(time.Second))
}

func TestOneAsync(t *testing.T) {
	latch := New(1)

	go func(d time.Duration, latch *CountDownLatch) {
		time.Sleep(d)
		latch.CountDown()
	} (time.Second, latch)

	// WaitTimeout should return false because count down will take place after a second
	assertEqual(t, false, latch.WaitTimeout(100*time.Millisecond))

	// Wait will return after ~900 msec
	latch.Wait()

	// Check count is 0
	assertEqual(t, uint(0), latch.Count())
	
	// WaitTimeout now will return true
	assertEqual(t, true, latch.WaitTimeout(100*time.Millisecond))
}

func TestManyAsync(t *testing.T) {
	count := 1e6
	latch1 := New(uint(count*(count+1)/2)) //sum of numbers 1 to count
	latch2 := New(uint(2*count + 1))
	latch3 := New(uint(3*count))
	for i := 1; i <= int(count); i++ {
		go func (weight uint) {
			assertNil(t, latch1.WeightedCountDown(weight))
			latch2.Wait()
			assertNil(t, latch3.WeightedCountDown(3))
		} (uint(i))
	}
	
	for ; !latch1.WaitTimeout(50*time.Millisecond) ; {
	//	fmt.Println("latch 1 remaining:", latch1.Count())
	}
	
	assertEqual(t, uint(0), latch1.Count())
	assertEqual(t, uint(2*count + 1), latch2.Count())
	assertEqual(t, uint(3*count), latch3.Count())

	for i := 0; i < int(count); i++ {
		go func () {
			assertNil(t, latch2.WeightedCountDown(2))
			latch3.Wait()
		} ()
	}

	for ; latch2.Count() > 1 ; {
		latch2.WaitTimeout(50*time.Millisecond)
	//	fmt.Println("latch 2 remaining:", latch2.Count())
	}

	assertEqual(t, false, latch2.WaitTimeout(time.Second))
	assertEqual(t, uint(1), latch2.Count())
	assertEqual(t, uint(3*count), latch3.Count())

	assertNil(t, latch2.CountDown())
	assertEqual(t, true, latch2.WaitTimeout(time.Second))
	assertEqual(t, uint(0), latch2.Count())

	for ; !latch3.WaitTimeout(50*time.Millisecond) ; {
	//	fmt.Println("latch 3 remaining:", latch3.Count())
	}

	assertEqual(t, uint(0), latch3.Count())
}

// assertion helpers

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatal("Not equal:", "expected:", expected, ", actual:", actual)
	}
}

func assertNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Fatal("Value not nil, actual:", actual)
	}
}

func assertNotNil(t *testing.T, actual interface{}) {
	if actual == nil {
		t.Fatal("Value is nil")
	}
}
