package main

// This file is mandatory as otherwise the gabeat.test binary is not generated correctly.

import (
	"errors"
	"flag"
	"fmt"
	"testing"
)

var systemTest *bool

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {

	if *systemTest {
		main()
	}
}

func TestExitStatusError(t *testing.T) {
	testExitStatus(t, 1, errors.New("Test error for testing"))
}

func TestExitStatusSuccess(t *testing.T) {
	testExitStatus(t, 0, nil)
}

func testExitStatus(t *testing.T, expected int, err error) {
	status := getExitStatus(err)
	if status != expected {
		message, _ := fmt.Printf("Expected %d, got ", expected)
		t.Error(message, status)
	}
}
