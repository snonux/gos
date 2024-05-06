package health

import "testing"

func TestHealthStatus(t *testing.T) {
	t.Parallel()

	hs := NewStatus()
	hs.Set(Warning, "foo", "this is not good")
	hs.Set(Critical, "bar", "this is not good either")
	hs.Set(Warning, "baz", "urgh!")
	hs.Set(Unknown, "baz", "don't know what happened here!")
	hs.Clear("foo")

	result := hs.String()
	expected := `UNKNOWN: don't know what happened here! (handler baz)
CRITICAL: this is not good either (handler bar)
`

	if result != expected {
		t.Error("expected", expected, "but got", result)
	}
	t.Log("got as expexted", result)
}
