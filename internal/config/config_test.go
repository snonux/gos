package config

import (
	"os"
	"slices"
	"testing"
)

func TestEnvToStr(t *testing.T) {
	t.Parallel()

	os.Unsetenv("NON_EXISTENT_ENV")
	os.Setenv("GOS_TEST_FROM_ENV", "foobarbaz")

	var (
		expected = "foobarbaz"
		got      = Env[ToString]("GOS_TEST_FROM_ENV")
	)

	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)

	expected = "default value"
	got = Env[ToString]("NON_EXISTENT_ENV", expected)
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)

	if got = Env[ToString]("NON_EXISTENT_ENV"); got != "" {
		t.Errorf("got '%s' but expected empty string", got)
	}
	t.Logf("got empty string as expected")

	expected = "casio g-shock"
	got = Env[ToString]("GOS_WATCH", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%s' but expected '%s'", got, expected)
	}
	t.Logf("got '%s' as expected", expected)
}

func TestEnvToStrSlice(t *testing.T) {
	t.Parallel()

	os.Setenv("GOS_TEST_SLICE_FROM_ENV", "foo,bar,baz")

	var (
		expected = []string{"foo", "bar", "baz"}
		got      = Env[ToStringSlice]("GOS_TEST_SLICE_FROM_ENV")
	)
	if !slices.Equal(got, expected) {
		t.Errorf("got '%v' but expected '%v'", got, expected)
	}
	t.Logf("got '%v' as expected", expected)

	expected = []string{"default value"}
	got = Env[ToStringSlice]("NON_EXISTENT_ENV_SLICE", "default value")
	if !slices.Equal(got, expected) {
		t.Errorf("got '%v' but expected '%v'", got, expected)
	}
	t.Logf("got '%v' as expected", expected)

	os.Unsetenv("NON_EXISTENT_ENV")
	if got = Env[ToStringSlice]("NON_EXISTENT_ENV"); len(got) > 0 {
		t.Errorf("got '%s' of len '%d' but expected empty slice", got, len(got))
	}
	t.Logf("got empty slice as expected")

	expected = []string{"casio", "g-shock"}
	got = Env[ToStringSlice]("NON_EXISTENT_ENV", "", "", "", "casio,g-shock", "")
	if !slices.Equal(got, expected) {
		t.Errorf("got '%v' but expected '%v'", got, expected)
	}
	t.Logf("got '%v' as expected", expected)
}

func TestEnvToInt(t *testing.T) {
	t.Parallel()

	os.Unsetenv("NON_EXISTENT_ENV")
	os.Setenv("GOS_TEST_INT_FROM_ENV", "1")

	var (
		expected = 1
		got      = Env[ToInteger](t, "GOS_TEST_INT_FROM_ENV")
	)

	if got != expected {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
	t.Logf("got '%d' as expected", expected)

	expected = 999
	got = Env[ToInteger]("NON_EXISTENT_ENV", expected)
	if got != expected {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
	t.Logf("got '%d' as expected", expected)

	if got = Env[ToInteger]("NON_EXISTENT_ENV"); got != 0 {
		t.Errorf("got '%d' but expected zero", got)
	}
	t.Logf("got zero as expected")

	expected = 1234
	got = Env[ToInteger]("GOS_WATCH", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
	t.Logf("got '%d' as expected", expected)
}

func TestEnvToBool(t *testing.T) {
	t.Parallel()

	os.Unsetenv("NON_EXISTENT_ENV")
	os.Setenv("GOS_TEST_BOOL_FROM_ENV", "true")

	var (
		expected = true
		got      = Env[ToBool](t, "GOS_TEST_BOOL_FROM_ENV")
	)

	if got != expected {
		t.Errorf("got '%t' but expected '%t'", got, expected)
	}
	t.Logf("got '%t' as expected", expected)

	expected = false
	got = Env[ToBool]("NON_EXISTENT_ENV", expected)
	if got != expected {
		t.Errorf("got '%t' but expected '%t'", got, expected)
	}
	t.Logf("got '%t' as expected", expected)

	if got = Env[ToBool]("NON_EXISTENT_ENV"); got {
		t.Errorf("got '%t' but expected false", got)
	}
	t.Logf("got 'false' as expected")

	expected = true
	got = Env[ToBool]("NON_EXISTENT_ENV", "", "", "", expected, "")
	if got != expected {
		t.Errorf("got '%t' but expected '%t'", got, expected)
	}
	t.Logf("got '%t' as expected", expected)
}

func TestSecondENV(t *testing.T) {
	t.Parallel()

	os.Unsetenv("GOS_NONEXISTANT")
	os.Setenv("EDITOR", "hx")

	var (
		expected = "hx"
		got      = Env[ToString]("GOS_NONEXISTANT", "EDITOR", "notepad.exe")
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
		got      = Env[ToString]("GOS_NONEXISTANT", func() string {
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
		got      = Env[ToInteger]("GOS_NONEXISTANT", func() int {
			return 666
		})
	)

	if expected != got {
		t.Errorf("got '%d' but expected '%d'", got, expected)
	}
}
