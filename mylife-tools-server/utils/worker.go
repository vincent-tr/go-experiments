package utils

type Worker struct {
	exit chan struct{}
	done chan struct{}
}

func NewWorker(callback func(exit chan struct{})) *Worker {
	worker := &Worker{
		exit: make(chan struct{}, 1),
		done: make(chan struct{}, 1),
	}

	go func() {
		defer func() {
			worker.done <- struct{}{}
		}()

		callback(worker.exit)
	}()

	return worker
}

func (worker *Worker) Terminate() {
	worker.exit <- struct{}{}
	<-worker.done
}
