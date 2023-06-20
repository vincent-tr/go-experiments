package tasks

import (
	"errors"
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
)

var logger = log.CreateLogger("mylife:server:tasks")

type Task func()

func init() {
	services.Register(&taskService{})
}

type taskService struct {
	queues map[string]*taskQueue
}

func (service *taskService) Init() error {
	service.queues = make(map[string]*taskQueue)

	return nil
}

func (service *taskService) Terminate() error {
	for id, _ := range service.queues {
		service.closeQueue(id)
	}

	return nil
}

func (service *taskService) ServiceName() string {
	return "tasks"
}

func (service *taskService) Dependencies() []string {
	return []string{}
}

func (service *taskService) createQueue(id string) error {
	if _, exists := service.queues[id]; exists {
		return errors.New(fmt.Sprintf("Cannot create queue '%s': already exists", id))
	}

	queue := newTaskQueue(id)
	service.queues[id] = queue

	logger.WithField("queueId", queue.id).Debug("Queue created")

	return nil
}

func (service *taskService) closeQueue(id string) error {
	queue, exists := service.queues[id]
	if !exists {
		return errors.New(fmt.Sprintf("Cannot close queue '%s': does not exists", id))
	}

	queue.close()
	delete(service.queues, id)

	logger.WithField("queueId", queue.id).Debug("Queue closed")

	return nil
}

func (service *taskService) getQueue(id string) (*taskQueue, error) {
	queue, exists := service.queues[id]
	if !exists {
		return nil, errors.New(fmt.Sprintf("Cannot get queue '%s': does not exists", id))
	}

	return queue, nil
}

func getService() *taskService {
	return services.GetService[*taskService]("tasks")
}

// Public access

func CreateQueue(id string) error {
	return getService().createQueue(id)
}

func CloseQueue(id string) error {
	return getService().closeQueue(id)
}

func Submit(id string, taskName string, taskImpl Task) error {
	queue, err := getService().getQueue(id)
	if err != nil {
		return err
	}

	return queue.submit(taskName, taskImpl)
}
