package server

import "testing"

func TestBytesToInt64(t *testing.T) {
	var expect int64 = 1

	bytes := make([]byte, 4)

	bytes[0] = 1
	bytes[1] = 0
	bytes[2] = 0
	bytes[3] = 0

	result := bytesToInt64(bytes)

	if result != expect {
		t.Errorf("Expected %v, got %v.", expect, result)
	}
}
