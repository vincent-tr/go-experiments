package tasks

import (
	"errors"
	"fmt"
	"mylife-tools-server/log"
	"sync"
	"time"
)

type taskQueueStatus int

const (
	Running = iota
	Closing
	Closed
)

const queueMaxBuffer = 1024

type taskQueue struct {
	id     string
	tasks  chan *task // queuing a nil task indicates exit
	done   chan struct{}
	status taskQueueStatus
	mutex  sync.Mutex
}

func newTaskQueue(id string) *taskQueue {
	queue := &taskQueue{
		id:     id,
		tasks:  make(chan *task, queueMaxBuffer),
		done:   make(chan struct{}, 1),
		status: Running,
	}

	go queue.workerEntry()

	return queue
}

func (queue *taskQueue) close() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.status == Closing {
		logger.WithField("queueId", queue.id).Error("Queue is already closing")
		return
	} else if queue.status == Closed {
		logger.WithField("queueId", queue.id).Error("Queue is already closed")
		return
	}

	queue.status = Closing
	queue.tasks <- nil
	<-queue.done
	queue.status = Closed
}

func (queue *taskQueue) submit(name string, taskImpl Task) error {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.status == Closing {
		return errors.New(fmt.Sprintf("Cannot add tasks while closing on queue '%s'", queue.id))
	} else if queue.status == Closed {
		return errors.New(fmt.Sprintf("Cannot add tasks on closed queue '%s'", queue.id))
	}

	queue.tasks <- newTask(name, taskImpl)

	return nil
}

func (queue *taskQueue) workerEntry() {
	defer func() {
		queue.done <- struct{}{}
	}()

	for {
		task := <-queue.tasks

		if task == nil {
			break
		}

		task.run(queue)
	}

}

type task struct {
	name string
	impl Task
}

func newTask(name string, impl Task) *task {
	return &task{
		name: name,
		impl: impl,
	}
}

func (t *task) run(queue *taskQueue) {
	logger.WithFields(log.Fields{"queueId": queue.id, "taskName": t.name}).Debug("Task begin")

	tmr := newTimer()
	t.impl()
	elapsed := tmr.elapsed()

	logger.WithFields(log.Fields{"queueId": queue.id, "taskName": t.name, "elapsedMs": elapsed}).Debug("Task end")
}

type timer struct {
	begin time.Time
}

func newTimer() timer {
	return timer{begin: time.Now()}
}

func (t *timer) elapsed() float64 {
	return time.Since(t.begin).Seconds() * 1000
}
