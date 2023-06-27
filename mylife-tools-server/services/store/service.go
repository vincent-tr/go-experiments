package store

import (
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/tasks"
)

const storeUpdateQueueId = "store-update"

var logger = log.CreateLogger("mylife:server:store")

func init() {
	services.Register(&storeService{})
}

type storeService struct {
	collections map[string]genericInternalCollection
}

func (service *storeService) Init(arg interface{}) error {

	if err := tasks.CreateQueue(storeUpdateQueueId); err != nil {
		return err
	}

	service.collections = make(map[string]genericInternalCollection)

	if arg != nil {
		builders := arg.([]interface{})

		for _, builder := range builders {
			if err := service.setupCollection(builder.(icollectionBuilder)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *storeService) Terminate() error {
	for _, col := range service.collections {
		col.terminate()
	}

	service.collections = nil

	return tasks.CloseQueue(storeUpdateQueueId)
}

func (service *storeService) ServiceName() string {
	return "store"
}

func (service *storeService) Dependencies() []string {
	// io because we use io queue
	return []string{"database", "io", "tasks"}
}

func (service *storeService) getCollection(name string) (genericInternalCollection, error) {
	value, exists := service.collections[name]
	if !exists {
		return nil, fmt.Errorf("Collection '%s' not found.", name)
	}

	return value, nil
}

func (service *storeService) setupCollection(builder icollectionBuilder) error {
	name := builder.name()
	_, exists := service.collections[name]

	if exists {
		return fmt.Errorf("Collection '%s' already exists.", name)
	}

	service.collections[name] = builder.build()

	return nil
}

func getService() *storeService {
	return services.GetService[*storeService]("store")
}

// Public access

func SetupCollection[TEntity Entity](builder *CollectionBuilder[TEntity]) error {
	return getService().setupCollection(builder)
}

func GetCollection[TEntity Entity](name string) (ICollection[TEntity], error) {
	value, err := getService().getCollection(name)
	if err != nil {
		return nil, err
	}

	col, ok := value.(ICollection[TEntity])
	if !ok {
		return nil, fmt.Errorf("Collection '%s' requested with bad entity type.", name)
	}

	return col, nil
}
