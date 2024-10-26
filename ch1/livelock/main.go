package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// livelocks are programs that are actively performing concurrent operations, but these operations do nothing
// to move the state of the program forward.

func main() {

	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(80 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		defer cadence.L.Unlock()
		cadence.Wait()
	}

	// This example demonstrates a very common reason livelock are written: two or more concurrent processes
	// attempting to prevent a deadlock without coordination.
	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		_, _ = fmt.Fprintf(out, " %v", dirName)
		// one goroutine increments dir atomically
		// then another goroutine increments it again atomically
		// end never reached the success condition
		// Livelocks are a subset of a larger set of problems called starvation.
		atomic.AddInt32(dir, 1)
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			_, _ = fmt.Fprint(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right, back, forward int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }
	tryBack := func(out *bytes.Buffer) bool { return tryDir("back", &back, out) }
	tryForward := func(out *bytes.Buffer) bool { return tryDir("forward", &forward, out) }

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		_, _ = fmt.Fprintf(&out, "%v is trying to scoot:", name)
		// If it has no limit of tries, it will get stuck in a livelock forever.
		for i := 0; i < 3; i++ {
			if tryLeft(&out) || tryRight(&out) || tryBack(&out) || tryForward(&out) {
				return
			}
		}
		_, _ = fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
