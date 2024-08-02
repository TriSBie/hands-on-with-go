package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func GoRoutine() {
	f("direct")

	go f("goroutine")
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	time.Sleep(time.Second)
}

func Channels() {
	// Channel are typed by the values they convey
	message := make(chan string) // create channel string

	// Send a value into a channel using the channel <- syntax
	go func() {
		message <- "ping"
	}()

	msg := <-message
	fmt.Println(msg)
}

func SynchChannels() {
	done := make(chan bool, 1)
	go worker(done)
	<-done
}

func worker(done chan bool) {
	fmt.Println("Working")
	time.Sleep(time.Second)

	fmt.Println("Done")
	done <- true
}

func BufferedChannels() {
	message := make(chan string, 2)

	message <- "buffered"
	message <- "channel"

	// since channels only accepted 2 maximum receivers
	// If any exceed limit occurs -> deadlock

	// print out the message like a queuerf
	fmt.Println(<-message)
	fmt.Println(<-message)
}

// Specify ping channel only accepts strin value
func ping(pings chan<- string, message string) {
	pings <- message
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

// Select is used to wait on multiple channel operations
// Cast as switch case with channels operation
func Select() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "one"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

func ChannelDirections() {
	// using ping, pong  function
	pings := make(chan string, 1)
	pongs := make(chan string, 1)

	ping(pings, "Hello world")
	pong(pings, pongs)
	fmt.Println(<-pongs)
	<-pings
}

func Timeouts() {
	c1 := make(chan string, 1)
	go func() {
		// delay 2 seconds
		time.Sleep(2 * time.Second)
		c1 <- "Result 1"
	}()

	select {
	// since select will proceed with the values receive first
	case res := <-c1: // receives aftetr 2 seconds
		fmt.Println(res)
		// awaits a value to be sent aftter the timeout of 1s
	case <-time.After(1 * time.Second): // receives after 1 second
		fmt.Println("Time out1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "Result 2"
	}()

	select {
	case res2 := <-c2:
		fmt.Println(res2)
	case <-time.After(3 * time.Second):
		fmt.Println("time out 2")
	}
}

func NonBlockinChannel() {
	messages := make(chan string)
	signals := make(chan bool)
	msg := "hi"

	select {
	case msg := <-messages:
		fmt.Println("Receive mesasges", msg)
	default:
		fmt.Println("No messages received")
	}

	// Suggested code may be subject to a license. Learn more: ~LicenseLog:1240612970.
	msg = "hi"
	select {
	case messages <- msg:
		fmt.Println("Sent message", msg)
	default:
		fmt.Println("No message sent")
	}

	// Suggested code may be subject to a license. Learn more: ~LicenseLog:548428564.
	select {
	case msg := <-messages:
		fmt.Println("Received message", msg)
	case sig := <-signals:
		fmt.Println("Received signal", sig)
	default:
		fmt.Println("No activity")
	}
}

func ClosingChannels() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("Receive job", j)
			} else {
				fmt.Println("Received all jobs")
				done <- true
				return
			}
		}
	}()

	for i := 0; i < 3; i++ {
		jobs <- i
		fmt.Println("Sent job", i)
	}

	close(jobs) // close channel
	fmt.Println("Sent all jobs")

	<-done

	// reading from closed channel succeeds immediately return the zero value
	_, ok := <-jobs
	fmt.Println("Received more job", ok)
}

func RangeOverChannel() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	for elem := range queue {
		fmt.Println("Element", elem)
	}
}

func Timer() {
	/*
		- Timer : doing something once in the future
		- Ticker : doing something repeatedly at regular intervals
	*/
	timer1 := time.NewTimer(2 * time.Second)

	<-timer1.C // waits for the timer to fire
	fmt.Println("Timer 1 fired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()

	stop := timer2.Stop()
	if stop {
		fmt.Println("Timer 2 stopped")
	}

	time.Sleep(1 * time.Second)
}

func Tickers() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

// jobs is the emitter, results is the receiver
func workerV2(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func WorkerPool() {
	const numJobs = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go workerV2(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func workerV3(id int) {
	fmt.Printf("Worker %d starting \n", id)
	time.Sleep(time.Second)
	fmt.Printf("Workder %d done \n", id)
}

// Using waitgroup for waiting multiple goroutine finish
func WaitGroup() {
	var wg sync.WaitGroup

	// increase counter of wait group when several goroutines is lauch
	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			workerV3(i)
		}()
	}

	wg.Wait()
}

func runner1(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Runner 1")
}

func runner2(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Runner 2")
}

func WaitGroupReference() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go runner1(wg)
	go runner2(wg)

	wg.Wait() // wait until all goroutine finish
}

func RateLimiting() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}

	// close the requests channel
	close(requests)

	// define limiter using with timer Tick
	limiter := time.Tick(200 * time.Millisecond) // this is called by channel

	for req := range requests {
		<-limiter // get limiter value
		fmt.Println("Request", req, time.Now())
	}

	// define channel burstLimiter only accepts 3 buffered
	burstLimiter := make(chan time.Time, 3)

	// initializing burstLimiter with time now() as 3 times
	for i := 0; i < 3; i++ {
		burstLimiter <- time.Now()
	}

	// declare using with go routine with burstLimiter channel
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)

	for i := 0; i < 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)

	for req := range burstyRequests {
		<-burstLimiter
		fmt.Println("Request", req, time.Now())
	}
}

func AtomicsCounters() {
	var ops atomic.Uint64

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for c := 0; c < 1000; c++ {
				ops.Add(1)
			}
			wg.Done()
		}()
		wg.Wait()
	}
	// using atomic load for reading value while other goroutine updating it
	fmt.Println("ops: ", ops.Load())
}

/** Mutext **/
type Container struct {
	mu       sync.Mutex
	counters map[string]int
}

func (c *Container) inc(name string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	c.counters[name]++
}

// Access data across multiple go routines
// Want to synchronize access to shared state accross
// Ensuring that only one goroutine can access the critical section of code at a time.
func Mutexes() {
	c := Container{counters: map[string]int{
		"a": 0, "b": 0}}

	// since calling multiple go routines
	var wg sync.WaitGroup

	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			c.inc(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doIncrement("a", 1000)
	go doIncrement("b", 1000)

	wg.Wait()
	fmt.Println(c.counters)
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1

	// using infinite loop with for uncondition
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit: // this case should be listening when quit is receive error
			fmt.Println("Quit")
			return
		}
	}

}

func FibonacciWithChannel() {
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c) // wait for c to receive value ten times
		}
		quit <- 0
	}()
	fibonacci(c, quit)

}

/* Using with Stateful Goroutines*/

type readOp struct {
	key  int
	resp chan int // type as channels int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func StatefulGoRoutines() {
	var readOps uint64
	var writeOps uint64

	reads := make(chan readOp)
	writes := make(chan writeOp)

	// listener
	go func() {
		// state will persist across multiple go routines
		var state = make(map[int]int)
		for {
			select {
			// waiting for receives message from reads as could as possible
			case read := <-reads:
				read.resp <- state[read.key]
				// waiting for receives message from writes as could as possible
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	fmt.Println(rand.Intn(100))

	// emitter
	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read // write from read to reads
				<-read.resp   // read from read.resp channel
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Microsecond)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				writes <- write // send data write into writes channel
				<-write.resp    // read from write.resp channel
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Microsecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("ReadOps", readOpsFinal)

	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("WriteOps", writeOpsFinal)
}

func main() {
	// GoRoutine() - learning go routine
	// Channels() - learning channels
	// BufferedChannels() - using buffered channels with a limit accepts receiving message
	// SynchChannels() - using sync channels
	// Select() - handle multiple channels as switch case
	// Timeouts() - learning about timeouts
	// NonBlockinChannel() - non-blocking channels means
	// ClosingChannels()
	// RangeOverChannel() - iterating over channel using high-level with range
	// Timer() - cast as setTimeOut, doing once in the future
	// Tickers() - cast as setInterval, doing a specific task repeatedly
	// WorkerPool() - leanring how to communication between each channel using go
	// WaitGroup() - learning how to wait for multiple routines finish the execution
	// WaitGroupReference()
	// RateLimiting() - demonstration using rate-limiter using with channel time Tick

	/** Managing state **/
	// AtomicsCounters()
	// Mutexes()
	// FibonacciWithChannel()

	StatefulGoRoutines()
}
