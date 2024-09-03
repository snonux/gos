package config

import (
	"os"
	"strconv"
	"strings"
)

type enverConstraint interface {
	~int | ~bool | ~string | []string
}

type enver[T enverConstraint] interface {
	// Return T value from input string
	fromStr(value string) T
	// Return T's zero value
	zero() T
}

func Env[U enver[T], T enverConstraint](keys ...any) T {
	var enver U

	for _, key := range keys {
		switch key := key.(type) {
		case string:
			if key == "" {
				continue
			}
			if !isAllUpperCase(key) {
				return enver.fromStr(key)
			}
			if value := os.Getenv(key); value != "" {
				return enver.fromStr(value)
			}
		case func() T:
			return key()
		}
	}

	return enver.zero()
}

type ToString struct{}

func (ToString) fromStr(str string) string {
	return str
}

func (ToString) zero() string {
	return ""
}

type ToStringSlice struct{}

func (s ToStringSlice) fromStr(str string) []string {
	result := strings.Split(str, ",")
	if len(result) == 1 && result[0] == "" {
		return s.zero()
	}
	return result
}

func (ToStringSlice) zero() []string {
	return []string{}
}

type ToInteger struct{}

// TODO: Return an error if can't convert to int
func (s ToInteger) fromStr(str string) int {
	if result, err := strconv.Atoi(str); err == nil {
		return result
	}
	return s.zero()
}

func (ToInteger) zero() int {
	return 0
}

type ToBool struct{}

// TODO: Return an error if can't convert to bool
func (s ToBool) fromStr(str string) bool {
	if result, err := strconv.ParseBool(str); err == nil {
		return result
	}
	return s.zero()
}

func (ToBool) zero() bool {
	return false
}
