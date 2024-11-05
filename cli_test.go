package main

import (
	"strings"
	"testing"
)

type FakeReader struct {
	output string
}

func (r FakeReader) Read(p []byte) (n int, err error) {
	p = []byte(r.output)
	return len(p), nil
}

func TestFInput(t *testing.T) {
	want := "red!"
	r := strings.NewReader(want)

	got := FInput(r)

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
