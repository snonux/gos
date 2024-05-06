package health

import "testing"

func TestHealthStatus(t *testing.T) {
	t.Parallel()

	h := NewStatus()
	h.Set(warning, "fooService", "this is not good")
	h.Set(critical, "barService", "this is not good either")
	h.Set(warning, "bazService", "urgh!")
	h.Set(unknown, "bazService", "don't know what happened here!")
	h.Clear("fooService")

	result := h.String()
	expected := `UNKNOWN: don't know what happened here!
CRITICAL: this is not good either
`

	if result != expected {
		t.Error("expected", expected, "but got", result)
	}
	t.Log("got as expexted", result)
}
