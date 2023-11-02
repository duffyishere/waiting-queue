package main

import "testing"

func TestAddEntryNumber(t *testing.T) {
	client, ctx := ConnRedis()
	AddEntryNumber(client, ctx, 1)
}
