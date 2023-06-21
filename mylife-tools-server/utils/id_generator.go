package utils

import "sync"

type IdGenerator struct {
	mutex sync.Mutex
	last  uint64
}

func NewIdGenerator() IdGenerator {
	return IdGenerator{last: 0}
}

func (generator *IdGenerator) Next() uint64 {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	generator.last += 1
	return generator.last
}
