package config

import (
	"os"
	"testing"
)

func TestCamelToSnake(t *testing.T) {
	t.Parallel()

	var (
		expected = "GOS_FOO_BAR_BAZ"
		got      = camelToSnakeWithPrefix("GOS", "fooBarBaz")
	)

	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)
}

func TestFromEnv(t *testing.T) {
	t.Parallel()

	os.Setenv("GOS_TEST_FROM_ENV", "foobarbaz")

	var (
		expected = "foobarbaz"
		got      = fromEnv("testFromEnv")
	)

	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)

	expected = "default value"
	got = fromEnv("jajaja", expected)
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)

	if got = fromEnv("jujuju"); got != "" {
		t.Errorf("got '%s' but expected empty string", got)
		return
	}
	t.Logf("got empty string as expected")

	expected = "casio g-shock"
	got = fromEnv("watch", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
		return
	}
	t.Logf("got '%s' as expected", expected)
}
