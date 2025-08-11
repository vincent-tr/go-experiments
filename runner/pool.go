package runner

import (
	"runtime"
	"sync"
)

type Task func()

type TaskPool struct {
	wg     sync.WaitGroup
	tasks  chan Task
	closed bool
}

func NewTaskPool() *TaskPool {
	workers := runtime.NumCPU() // Use number of CPU cores as default

	p := &TaskPool{
		tasks: make(chan Task, 10000000),
	}

	p.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer p.wg.Done()
			for task := range p.tasks {
				task()
			}
		}(i)
	}

	return p
}

func (p *TaskPool) Submit(task Task) {
	if p.closed {
		panic("submit on closed pool")
	}
	p.tasks <- task
}

func (p *TaskPool) Close() {
	if !p.closed {
		p.closed = true
		close(p.tasks)
		p.wg.Wait()
	}
}
