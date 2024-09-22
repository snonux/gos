package oi

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func EnsureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func EnsureParentDirExists(dir string) error {
	return EnsureDirExists(filepath.Dir(dir))
}

func ReadDirFilter(dir string, filter func(file os.DirEntry) bool) (chan string, error) {
	ch := make(chan string)

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
			if filter(file) {
				ch <- filepath.Join(dir, file.Name())
			}
		}
	}()

	return ch, nil
}

func ReadDirSlurp(dir string, filter func(file os.DirEntry) bool) ([]string, error) {
	var files []string

	ch, err := ReadDirFilter(dir, filter)
	if err != err {
		return files, err
	}

	for file := range ch {
		files = append(files, file)
	}

	return files, nil
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
