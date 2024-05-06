package health

import "testing"

func TestHealthStatus(t *testing.T) {
	t.Parallel()

	h := NewStatus()
	h.Set(Warning, "fooService", "this is not good")
	h.Set(Critical, "barService", "this is not good either")
	h.Set(Warning, "bazService", "urgh!")
	h.Set(Unknown, "bazService", "don't know what happened here!")
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
