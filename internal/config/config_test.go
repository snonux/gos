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
	}
	t.Logf("got '%s' as expected", expected)

	expected = "default value"
	got = FromENV("GOS_JAJAJA", expected)
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)

	os.Unsetenv("JUJUJU_NOT_EXISTANT_ENV")
	if got = FromENV("JUJUJU_NOT_EXISTANT_ENV"); got != "" {
		t.Errorf("got '%s' but expected empty string", got)
	}
	t.Logf("got empty string as expected")

	expected = "casio g-shock"
	got = FromENV("GOS_WATCH", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)
}

func TestSecondENV(t *testing.T) {
	t.Parallel()

	os.Unsetenv("GOS_NONEXISTANT")
	os.Setenv("EDITOR", "hx")

	var (
		expected = "hx"
		got      = FromENV("GOS_NONEXISTANT", "EDITOR", "notepad.exe")
	)

	if expected != got {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
}dfdf

func TestIsAllUpperCase(t *testing.T) {
	if isAllUpperCase("foo_bar") {
		t.Errorf("lowercas letters in test case")
	}
	if isAllUpperCase("FOO123") {
		t.Errorf("numbers in string should not evaluate to is all upper")
	}
	if !isAllUpperCase("FOO_BAR") {
		t.Errorf("should be all upper")
	}
}
