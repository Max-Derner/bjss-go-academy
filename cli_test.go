package main

import (
	"strings"
	"testing"
)

func TestFInput(t *testing.T) {
	want := "I'm red!"
	r := strings.NewReader(want + "\n")

	got := FInput(r)

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
