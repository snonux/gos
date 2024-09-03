package config

import (
	"os"
	"strconv"
	"strings"
)

type configTypes interface {
	~int | ~bool | ~string | []string
}

type enver[T configTypes] interface {
	fromStr(value string) (T, error) // Return T value from input string
	zero() T                         // Return T's zero value
}

func Str(keys ...any) string {
	return fromEnv[toStr](keys...)
}

func StrSlice(keys ...any) []string {
	return fromEnv[toStrSlice](keys...)
}

func Int(keys ...any) int {
	return fromEnv[toInt](keys...)
}

func Bool(keys ...any) bool {
	return fromEnv[toBool](keys...)
}

func fromEnv[U enver[T], T configTypes](keys ...any) T {
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

type toStr struct{}

func (toStr) fromStr(str string) (string, error) {
	return str, nil
}

func (toStr) zero() string {
	return ""
}

type toStrSlice struct{}

func (s toStrSlice) fromStr(str string) ([]string, error) {
	result := strings.Split(str, ",")
	if len(result) == 1 && result[0] == "" {
		return s.zero(), nil
	}
	return result, nil
}

func (toStrSlice) zero() []string {
	return []string{}
}

type toInt struct{}

func (toInt) fromStr(str string) (int, error) {
	return strconv.Atoi(str)
}

func (toInt) zero() int {
	return 0
}

type toBool struct{}

func (toBool) fromStr(str string) (bool, error) {
	return strconv.ParseBool(str)
}

func (toBool) zero() bool {
	return false
}
