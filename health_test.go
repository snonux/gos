package main

import "testing"

func TestHealthStatus(t *testing.T) {
	t.Parallel()

	h := newHealthStatus()
	h.set(warning, "fooService", "this is not good")
	h.set(critical, "barService", "this is not good either")
	h.set(warning, "bazService", "urgh!")
	h.set(unknown, "bazService", "don't know what happened here!")
	h.clear("fooService")

	result := h.String()
	expected := `UNKNOWN: don't know what happened here!
CRITICAL: this is not good either
`

	if result != expected {
		t.Error("expected", expected, "but got", result)
	}
	t.Log("got as expexted", result)
}
