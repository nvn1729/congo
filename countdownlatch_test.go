package congo

import (
	"testing"
	"time"
	"fmt"
)

func ExampleCountDownLatch() {
	latch := NewCountDownLatch(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer latch.CountDown() // count down after work is complete
			// do work
			// ...
			fmt.Println("Counting down")
			
		} ()
	}
	
	
	if (latch.WaitTimeout(5*time.Second)) {
		fmt.Println("Count down complete")
	} else {
		fmt.Println("Count down not complete")
	}
	// Output:
	// Counting down
	// Counting down
	// Counting down
	// Count down complete
}


func TestCountDownLatch_zero(t *testing.T) {
	latch := NewCountDownLatch(0)

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

func TestCountDownLatch_one(t *testing.T) {
	latch := NewCountDownLatch(1)

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

func TestCountDownLatch_oneasync(t *testing.T) {
	latch := NewCountDownLatch(1)

	go func() {
		time.Sleep(time.Second)
		latch.CountDown()
	} ()

	// WaitTimeout should return false because count down will take place after a second
	assertEqual(t, false, latch.WaitTimeout(100*time.Millisecond))

	// Wait will return after ~900 msec
	latch.Wait()

	// Check count is 0
	assertEqual(t, uint(0), latch.Count())
	
	// WaitTimeout now will return true
	assertEqual(t, true, latch.WaitTimeout(100*time.Millisecond))
}

func TestCountDownLatch_complete(t *testing.T) {
	latch1 := NewCountDownLatch(0)
	latch2 := NewCountDownLatch(1)
	latch3 := NewCountDownLatch(1e6)

	assertNotNil(t, latch1.Complete())

	assertNil(t, latch2.Complete())
	assertEqual(t, uint(0), latch2.Count())
	assertNotNil(t, latch2.Complete())

	latch4 := NewCountDownLatch(1)
	go func() {
		latch3.Wait()
		latch4.CountDown()
	}()
	assertNil(t, latch3.WeightedCountDown(3e5))
	assertEqual(t, false, latch4.WaitTimeout(500*time.Millisecond)) // test will pause half a second
	assertEqual(t, uint(7e5), latch3.Count())
	assertEqual(t, uint(1), latch4.Count())
	assertNil(t, latch3.Complete())
	assertEqual(t, true, latch4.WaitTimeout(time.Second))
	assertEqual(t, uint(0), latch3.Count())
	
}

func TestCountDownLatch_manyasync(t *testing.T) {
	count := 1e6
	latch1 := NewCountDownLatch(uint(count*(count+1)/2)) //sum of numbers 1 to count
	latch2 := NewCountDownLatch(uint(2*count + 1))
	latch3 := NewCountDownLatch(uint(3*count))
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
