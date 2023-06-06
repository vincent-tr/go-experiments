package sessions

import "sync"

type idGenerator struct {
	mutex sync.Mutex
	last  int
}

func newIdGenerator() idGenerator {
	return idGenerator{last: 0}
}

func (generator *idGenerator) Next() int {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	generator.last += 1
	return generator.last
}
