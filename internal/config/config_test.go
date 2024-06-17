package config

import (
	"os"
	"testing"
)

func TestFromENV(t *testing.T) {
	t.Parallel()

	os.Setenv("GOS_TEST_FROM_ENV", "foobarbaz")

	var (
		expected = "foobarbaz"
		got      = FromENV("GOS_TEST_FROM_ENV")
	)

	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)

	expected = "default value"
	got = FromENV("GOS_JAJAJA", expected)
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)

	if got = FromENV("jujuju"); got != "" {
		t.Errorf("got '%s' but expected empty string", got)
		return
	}
	t.Logf("got empty string as expected")

	expected = "casio g-shock"
	got = FromENV("GOS_WATCH", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)
}
