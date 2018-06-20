# Congo - Concurrency Primitives for Go

Package congo includes a CountDownLatch primitive, similar to what's available in the [`java.util.concurrent` package](https://docs.oracle.com/javase/9/docs/api/java/util/concurrent/CountDownLatch.html). A CountDownLatch is useful for signaling the completion of the one or more events occurring across multiple goroutines.

## Usage

In typical usage, a main goroutine creates the latch with a provided count and passes the latch to a number of goroutines. The goroutines invoke `CountDown` on the latch to signal completion and reduce the remaining count on the latch. The main goroutine invokes `Wait` or `WaitTimeout` on the latch to wait for all events to complete.

For example, the following code:

```
latch := congo.NewCountDownLatch(3)
for i := 0; i < 3; i++ {
	go func() {
		// do work
		// ...
		fmt.Println("Counting down")
		latch.CountDown() // done with work
	} ()
}
	
	
if (latch.WaitTimeout(5*time.Second)) {
	fmt.Println("Count down complete")
} else {
	fmt.Println("Count down not complete")
}
```

Produces output:

```
Counting down
Counting down
Counting down
Count down complete
```

## CountDownLatch vs. sync.WaitGroup

`CountDownLatch` is similar to `WaitGroup` in the standard Go `sync` package with a few notable differences:

* Caller can wait with a timeout (`WaitTimeout`) for the count down to complete.
* Couple of extra ways to `CountDown`:
  * `WeightedCountDown` reduces the remaining count by a specified number.
  * `Complete` reduces the remaining count to 0 and signals any waiting goroutines immediately.
* The starting count is set once at the time of creating the CountDownLatch. This avoids the potential for misuse of the `WaitGroup.Add` function, which should only be invoked in the main goroutine.


## Documentation

[![GoDoc](https://godoc.org/github.com/nvn1729/congo?status.svg)](https://godoc.org/github.com/nvn1729/congo)

