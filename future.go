package commands

import (
	"errors"
	"sync"
)

type promise struct {
	result CommandResult
	err error
}

// future
type Future struct {
	promiseChan chan promise
	closed bool
	mu sync.Mutex
}

func (f *Future) WaitAndGet() (CommandResult, error) {
	p, ok := <- f.promiseChan
	if !ok {
		return nil, errors.New("Get command handle result failed, the chan is closed.")
	}
	if p.err != nil {
		return nil, p.err
	}
	return p.result, nil
}

func (f *Future) closeChan() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return
	}
	close(f.promiseChan)
	f.closed = true
}

func (f *Future) setPromise(result CommandResult, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return
	}
	if err != nil {
		f.promiseChan <- promise{nil, err}
		return
	}
	f.promiseChan <- promise{result, nil}
}