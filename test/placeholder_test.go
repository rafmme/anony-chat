package test

import "testing"

func TestPlaceholder(t *testing.T) {
	want := 4
	got := 2 + 2

	if got != want {
		t.Errorf("Got %v instead of %v", got, want)
	}
}
