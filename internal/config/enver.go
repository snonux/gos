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
	fromStr(value string) (T, error)
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
				if val, err := enver.fromStr(key); err == nil {
					return val
				}
			} else if strVal := os.Getenv(key); strVal != "" {
				if val, err := enver.fromStr(strVal); err == nil {
					return val
				}
			}
		case T:
			return key
		case func() T:
			return key()
		}
	}

	return enver.zero()
}

type Str struct{}

func (Str) fromStr(str string) (string, error) {
	return str, nil
}

func (Str) zero() string {
	return ""
}

type StrSlice struct{}

func (s StrSlice) fromStr(str string) ([]string, error) {
	result := strings.Split(str, ",")
	if len(result) == 1 && result[0] == "" {
		return s.zero(), nil
	}
	return result, nil
}

func (StrSlice) zero() []string {
	return []string{}
}

type Int struct{}

func (Int) fromStr(str string) (int, error) {
	return strconv.Atoi(str)
}

func (Int) zero() int {
	return 0
}

type Bool struct{}

func (Bool) fromStr(str string) (bool, error) {
	return strconv.ParseBool(str)
}

func (Bool) zero() bool {
	return false
}
