package easyhttp

import (
	"errors"
	"sync"
)

// Safe errors
type safErrors struct {
	errs  []error
	mutex sync.Mutex
}

func (errs *safErrors) Append(err error) {
	if err == nil {
		return
	}
	errs.mutex.Lock()
	defer errs.mutex.Unlock()
	errs.errs = append(errs.errs, err)
}

func (errs *safErrors) Join() error {
	errs.mutex.Lock()
	defer errs.mutex.Unlock()
	return errors.Join(errs.errs...)
}

func (errs *safErrors) Error() string {
	return errs.Join().Error()
}
