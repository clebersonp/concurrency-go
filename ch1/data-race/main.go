package main

import (
	"fmt"
	"time"
)

// To run this program from the root directory checking DATA RACE:
// go run -race ./ch1/data-race/main.go

func main() {
	raceConditions()
	badResolutionRaceConditions()
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
	}
}
