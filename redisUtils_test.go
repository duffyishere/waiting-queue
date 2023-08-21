package main

import (
	"fmt"
	"testing"
)

func TestAddLine(t *testing.T) {
	addLine("testRequestId2")
}

func TestGetRequestN(t *testing.T) {
	result := getLineRange(0, 1)
	fmt.Println(result)
}

func TestIncreaseCount(t *testing.T) {
	result := increaseCountBy(1)
	fmt.Println(result)
}
