package easyhttp

import (
	"errors"
	"sync"
)

// Safe errors
type SafErrors struct {
	errs  []error
	mutex sync.Mutex
}

func (errs *SafErrors) Append(err error) {
	if err == nil {
		return
	}
	errs.mutex.Lock()
	defer errs.mutex.Unlock()
	errs.errs = append(errs.errs, err)
}

func (errs *SafErrors) Join() error {
	errs.mutex.Lock()
	defer errs.mutex.Unlock()
	return errors.Join(errs.errs...)
}

func (errs *SafErrors) Error() string {
	return errs.Join().Error()
}
