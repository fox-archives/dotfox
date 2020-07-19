package util

import "testing"

func TestContains(t *testing.T) {
	arr := []string{"alfa", "bravo"}

	// test string
	if !Contains(arr, "alfa") {
		t.Error("string 'alfa' not in array 'arr'. it is supposed to be in it")
	}
}
