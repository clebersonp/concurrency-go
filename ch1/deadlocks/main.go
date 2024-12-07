package main

import (
	"fmt"
	"sync"
	"time"
)

// A deadlocked program is one in which all concurrent processes are waiting on one another.
// In this state, the program will never recover without outside intervention.

func main() {

	type value struct {
		mu    sync.Mutex
		value int
	}
	// There is a deadlock.
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock() // v1 = a lock; v1 = b lock
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)
		v2.mu.Lock() // v2 = b try to lock; v2 = a try to lock. But they already locked before.
		// fatal error: all goroutines are asleep - deadlock!
		defer v2.mu.Unlock()
		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	// these two goroutines will run concurrently and will cause a deadlock
	// the first goroutine is waiting on the second goroutine to unlock the lock
	go printSum(&a, &b)
	// the second goroutine is waiting on the first goroutine to unlock the lock
	go printSum(&b, &a)
	wg.Wait()
}
