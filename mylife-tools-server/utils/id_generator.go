package utils

import "sync"

type IdGenerator struct {
	mutex sync.Mutex
	last  int
}

func NewIdGenerator() IdGenerator {
	return IdGenerator{last: 0}
}

func (generator *IdGenerator) Next() int {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	generator.last += 1
	return generator.last
}
