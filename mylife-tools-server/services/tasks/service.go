package tasks

import (
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

func (service *taskService) Init(arg interface{}) error {
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
		return fmt.Errorf("Cannot create queue '%s': already exists", id)
	}

	queue := newTaskQueue(id)
	service.queues[id] = queue

	logger.WithField("queueId", queue.id).Debug("Queue created")

	return nil
}

func (service *taskService) closeQueue(id string) error {
	queue, exists := service.queues[id]
	if !exists {
		return fmt.Errorf("Cannot close queue '%s': does not exists", id)
	}

	queue.close()
	delete(service.queues, id)

	logger.WithField("queueId", queue.id).Debug("Queue closed")

	return nil
}

func (service *taskService) getQueue(id string) (*taskQueue, error) {
	queue, exists := service.queues[id]
	if !exists {
		return nil, fmt.Errorf("Cannot get queue '%s': does not exists", id)
	}

	return queue, nil
}

func getService() *taskService {
	return services.GetService[*taskService]("tasks")
}

// Public access

func CreateQueue(queueId string) error {
	return getService().createQueue(queueId)
}

func CloseQueue(queueId string) error {
	return getService().closeQueue(queueId)
}

func Submit(queueId string, taskName string, taskImpl Task) error {
	queue, err := getService().getQueue(queueId)
	if err != nil {
		return err
	}

	return queue.submit(taskName, taskImpl)
}
