package main

import (
	"fmt"
	"testing"
)

var requestId = "testRequestId2"

func TestAddLine(t *testing.T) {
	addQueue(requestId)
}

func TestGetWaitingNumBy(t *testing.T) {
	fmt.Printf("%s's waiting number: %d\n", requestId, getWaitingNumBy(requestId))
}

func TestGetLastEnterNum(t *testing.T) {
	fmt.Println("Last enter number: ", getLastEnterNum())
}

func TestIncreaseLastEnterNumBy(t *testing.T) {
	fmt.Println("Increased last enter number: ", increaseLastEnterNumBy(10))
}

func TestIncreaseWaitingNum(t *testing.T) {
	fmt.Println("Increased waiting line number: ", increaseWaitingNum())
}
