package config

import (
	"os"
	"testing"
)

func TestEnvToStr(t *testing.T) {
	t.Parallel()

	os.Setenv("GOS_TEST_FROM_ENV", "foobarbaz")

	var (
		expected = "foobarbaz"
		got      = EnvToStr("GOS_TEST_FROM_ENV")
	)

	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)

	expected = "default value"
	got = EnvToStr("GOS_JAJAJA", expected)
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)

	os.Unsetenv("JUJUJU_NOT_EXISTANT_ENV")
	if got = EnvToStr("JUJUJU_NOT_EXISTANT_ENV"); got != "" {
		t.Errorf("got '%s' but expected empty string", got)
	}
	t.Logf("got empty string as expected")

	expected = "casio g-shock"
	got = EnvToStr("GOS_WATCH", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)
}

func TestEnvToInt(t *testing.T) {
	t.Parallel()

	os.Setenv("GOS_TEST_INT_FROM_ENV", "1")

	var (
		expected = 1
		got      = EnvToInt(t, "GOS_TEST_INT_FROM_ENV")
	)

	if got != expected {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
	t.Logf("got '%d' as expected", expected)

	expected = 999
	got = EnvToInt("GOS_JAJAJA", expected)
	if got != expected {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
	t.Logf("got '%d' as expected", expected)

	os.Unsetenv("JUJUJU_NOT_EXISTANT_ENV")
	if got = EnvToInt("JUJUJU_NOT_EXISTANT_ENV"); got != 0 {
		t.Errorf("got '%d' but expected zero", got)
	}
	t.Logf("got zero as expected")

	expected = 1234
	got = EnvToInt("GOS_WATCH", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
	t.Logf("got '%d' as expected", expected)
}

func TestSecondENV(t *testing.T) {
	t.Parallel()

	os.Unsetenv("GOS_NONEXISTANT")
	os.Setenv("EDITOR", "hx")

	var (
		expected = "hx"
		got      = EnvToStr("GOS_NONEXISTANT", "EDITOR", "notepad.exe")
	)

	if expected != got {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
}

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

func TestDefaultStrCB(t *testing.T) {
	t.Parallel()
	os.Unsetenv("GOS_NONEXISTANT")

	var (
		expected = "hello"
		got      = EnvToStr("GOS_NONEXISTANT", func() string {
			return "hello"
		})
	)

	if expected != got {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
}

func TestDefaultIntCB(t *testing.T) {
	t.Parallel()
	os.Unsetenv("GOS_NONEXISTANT")

	var (
		expected = 666
		got      = EnvToInt("GOS_NONEXISTANT", func() int {
			return 666
		})
	)

	if expected != got {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
}
