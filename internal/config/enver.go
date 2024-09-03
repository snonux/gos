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

func Str(keys ...any) string {
	return fromEnv[ToStr](keys...)
}

func StrSlice(keys ...any) []string {
	return fromEnv[ToStrSlice](keys...)
}

func Int(keys ...any) int {
	return fromEnv[ToInt](keys...)
}

func Bool(keys ...any) bool {
	return fromEnv[ToBool](keys...)
}

func fromEnv[U enver[T], T enverConstraint](keys ...any) T {
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

type ToStr struct{}

func (ToStr) fromStr(str string) (string, error) {
	return str, nil
}

func (ToStr) zero() string {
	return ""
}

type ToStrSlice struct{}

func (s ToStrSlice) fromStr(str string) ([]string, error) {
	result := strings.Split(str, ",")
	if len(result) == 1 && result[0] == "" {
		return s.zero(), nil
	}
	return result, nil
}

func (ToStrSlice) zero() []string {
	return []string{}
}

type ToInt struct{}

func (ToInt) fromStr(str string) (int, error) {
	return strconv.Atoi(str)
}

func (ToInt) zero() int {
	return 0
}

type ToBool struct{}

func (ToBool) fromStr(str string) (bool, error) {
	return strconv.ParseBool(str)
}

func (ToBool) zero() bool {
	return false
}
