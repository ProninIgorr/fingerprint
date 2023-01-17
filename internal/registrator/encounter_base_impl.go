package registrator

import (
	"sync"
)

func NewEncounter(initCap int) Encounter {
	return &encounter{hits: make(map[interface{}]Counter, initCap)}
}

type encounter struct {
	sync.RWMutex
	hits map[interface{}]Counter
}

func (r *encounter) GetCounterPairs() []CounterPair {
	r.RLock()
	defer r.RUnlock()
	pairs := make([]CounterPair, 0, len(r.hits))
	// Note: counter values may change during iteration (due to concurrency),
	// if not applicable, use Encounter, it has "read commited" isolation level
	for k, c := range r.hits {
		pairs = append(pairs, CounterPair{k, c.GetScore()}) //
	}
	return pairs
}

func (r *encounter) CheckIn(key interface{}) int {
	r.RLock()
	if c, ok := r.hits[key]; ok {
		r.RUnlock()
		return c.Add(1)
	}
	r.RUnlock()
	// ...
	r.Lock()
	defer r.Unlock()
	if _, ok := r.hits[key]; !ok {
		r.hits[key] = NewCounter(0)
	}
	return r.hits[key].Add(1)
}

func (r *encounter) GetScores() map[interface{}]int {
	r.RLock()
	defer r.RUnlock()
	result := make(map[interface{}]int, len(r.hits))
	for k, c := range r.hits {
		result[k] = c.GetScore()
	}
	return result
}

func (r *encounter) KeysCount() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.hits)
}

func (r *encounter) TotalCount() (totalCount int) {
	r.RLock()
	defer r.RUnlock()
	for _, count := range r.hits {
		totalCount += count.GetScore()
	}
	return
}
