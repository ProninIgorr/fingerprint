package registrator

import (
	"sync"
)

func NewCounter(initValue int) Counter {
	return &counter{hits: initValue}
}

type counter struct {
	sync.RWMutex
	hits int
}

func (e *counter) Add(num int) int {
	e.Lock()
	defer e.Unlock()
	e.hits += num
	return e.hits
}

func (e *counter) GetScore() int {
	e.RLock()
	defer e.RUnlock()
	return e.hits
}
