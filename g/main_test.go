package main

import "testing"

func testSomething(t *testing.T) {
	result := 2 + 2
	expected := 4
	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}
