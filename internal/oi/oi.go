package oi

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/exp/rand"
)

var ErrNotFound = errors.New("no file/entry found")

func EnsureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func EnsureParentDirExists(dir string) error {
	return EnsureDirExists(filepath.Dir(dir))
}

// Rename to ReadDirCh
func ReadDirFilter[T any](dir string, cb func(file os.DirEntry) (T, bool)) (chan T, error) {
	ch := make(chan T)

	if err := EnsureDirExists(dir); err != nil {
		return ch, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return ch, err
	}

	go func() {
		defer close(ch)
		for _, file := range files {
			if val, ok := cb(file); ok {
				ch <- val
			}
		}
	}()

	return ch, nil
}
func TraverseDir(dir string, cb func(file os.DirEntry) error) error {
	if err := EnsureDirExists(dir); err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var errs []error
	for _, file := range files {
		if err := cb(file); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Rename to ReadDir
func ReadDirSlurp[T any](dir string, cb func(file os.DirEntry) (T, bool)) ([]T, error) {
	var results []T

	ch, err := ReadDirFilter(dir, cb)
	if err != err {
		return results, err
	}

	for file := range ch {
		results = append(results, file)
	}

	return results, nil
}

// Rename to ReadDirRandom
func ReadDirRandomEntry[T any](dir string, cb func(file os.DirEntry) (T, bool)) (T, error) {
	results, err := ReadDirSlurp(dir, cb)

	if err != nil {
		var zero T
		return zero, err
	}
	if len(results) == 0 {
		var zero T
		return zero, ErrNotFound
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	return results[rand.Intn(len(results))], nil
}

func IsRegular(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func CopyFile(srcPath, dstPath string) error {
	if !IsRegular(srcPath) {
		return fmt.Errorf("%s is not a regular file", srcPath)
	}

	source, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer source.Close()

	if err := EnsureParentDirExists(dstPath); err != nil {
		return err
	}

	destination, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func Rename(srcPath, dstPath string) error {
	if err := EnsureParentDirExists(dstPath); err != nil {
		return err
	}
	return os.Rename(srcPath, dstPath)
}
