package main

import (
	"fmt"
	"sync"
	"time"
)

// To run this program from the root directory checking DATA RACE:
// go run -race ./ch1/data-race/main.go

func main() {
	//raceConditions()
	//badResolutionRaceConditions()
	memoryAccessSync()
}

// Data race is when one concurrent operation attempts to read a variable while ate some undetermined time another
// concurrent operation is attempting to write to the same variable.
func raceConditions() {
	var data int
	go func() {
		data++
	}()
	if data == 0 {
		fmt.Printf("the value is %v\n", data)
	} else {
		fmt.Printf("the value %v is greater than 0\n", data)
	}
}

// In fact, it's still possible to get DATA RACE in this way too.
// Imagine if the function's goroutine takes more than 1 second to run.
// Even we have putting a sleep before reading the variable, it's still getting DATA RACE with we run the program with
// -race flag.
// What we need is a way to synchronize and make the goroutines to communicate each other.
func badResolutionRaceConditions() {
	var data int
	go func() { data++ }()
	time.Sleep(1 * time.Second) // this is terrible resolution
	if data == 0 {
		fmt.Printf("the value is %v\n", data)
	} else {
		fmt.Printf("the value %v is greater than 0\n", data)
	}
}

// memoryAccessSync demonstrates how to synchronize access to a variable.
// While it solves the DATA RACE issue, it haven't solved the 'race condition'.
// The order of operations in this function is still nondeterministic.
// It still don't know which will occur first in any given execution of this function.
// It can also create maintenance and performance issues.
// By synchronizing access to the memory in this manner, you are counting on all other developers to follow the same
// convention now and at the future.
func memoryAccessSync() {
	// Critical section is a section of code that is not thread safe or that needs exclusive access to a shared resource.
	// The following code is not idiomatic Go, but it very simply demonstrates memory access synchronization.
	var (
		memoryAccess sync.Mutex
		data         int
	)
	go func() {
		memoryAccess.Lock()
		defer memoryAccess.Unlock()
		data++ // critical section
	}()
	memoryAccess.Lock()
	if data == 0 { // critical section
		fmt.Printf("the value is %v\n", data)
	} else {
		fmt.Printf("the value %v is greater than 0\n", data) // critical section
	}
	memoryAccess.Unlock()
}
