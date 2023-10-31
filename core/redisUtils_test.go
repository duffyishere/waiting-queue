package main

import "testing"

func TestAddEntryNumber(t *testing.T) {
	client, ctx := connRedis()
	AddEntryNumber(client, ctx, 1)
}
