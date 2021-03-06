# Congo - Concurrency Primitives for Go

Package congo includes a CountDownLatch primitive, similar to what's available in the [`java.util.concurrent` package](https://docs.oracle.com/javase/9/docs/api/java/util/concurrent/CountDownLatch.html). A CountDownLatch is useful for signaling the completion of the one or more events occurring across multiple goroutines.

## Usage

In typical usage, a main goroutine creates the latch with a provided count and passes the latch to a number of goroutines. The goroutines invoke `CountDown` on the latch to signal completion and reduce the remaining count on the latch. The main goroutine invokes `Wait` or `WaitTimeout` on the latch to wait for all events to complete.

For example, the following code:

```go
latch := congo.NewCountDownLatch(3)
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
```

Produces output:

```plain
Counting down
Counting down
Counting down
Count down complete
```

## CountDownLatch vs. sync.WaitGroup

`CountDownLatch` is similar to `WaitGroup` in the standard Go `sync` package with a few notable differences:

* Caller can wait with a timeout (`WaitTimeout`) for the count down to complete.
* Caller can retrieve the current count using `Count` to track the progress of latch count down completion.
* Couple of extra ways to `CountDown`:
  * `WeightedCountDown` reduces the remaining count by a specified number.
  * `Complete` reduces the remaining count to 0 and signals any waiting goroutines immediately.
* The starting count is set once at the time of creating the CountDownLatch. This avoids the potential for misuse of the `WaitGroup.Add` function, which should only be invoked in the main goroutine.

## Installation

To install congo, use `go get`:

```go
go get github.com/nvn1729/congo
```

Then import the `github.com/nvn1729/congo` into your code:

```go
import "github.com/nvn1729/congo"
```

## Documentation

[![GoDoc](https://godoc.org/github.com/nvn1729/congo?status.svg)](https://godoc.org/github.com/nvn1729/congo)

[Blog Post](https://naveensunkavally.com/2018/06/21/writing-robust-concurrency-tests-in-go-using-a-countdownlatch/)

