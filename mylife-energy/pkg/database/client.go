package database

import (
	"context"
	"net/url"

	"mylife-energy/pkg/config"
	"mylife-energy/pkg/log"
	"mylife-energy/pkg/services"

	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

var logger = log.CreateLogger("mongo:client")

func init() {
	services.Register(&DatabaseService{})
}

type Collection = mongo.Collection

type DatabaseService struct {
	client   *mongo.Client
	database *mongo.Database
}

func (service *DatabaseService) Init() error {
	mongoUrl := config.GetString("mongo")

	parsedUrl, err := url.Parse(mongoUrl)
	if err != nil {
		return err
	}
	dbName := parsedUrl.Path[1:]

	logger.WithFields(log.Fields{"mongoUrl": mongoUrl, "dbName": dbName}).Info("Config")

	opts := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	service.client = client
	service.database = client.Database(dbName)

	return nil
}

func (service *DatabaseService) Terminate() error {
	return service.client.Disconnect(context.TODO())
}

func (service *DatabaseService) ServiceName() string {
	return "database"
}

func (service *DatabaseService) Dependencies() []string {
	return []string{}
}

func (service *DatabaseService) GetCollection(name string) *Collection {
	return service.database.Collection(name)
}

// Shortcuts

func GetCollection(name string) *Collection {
	return getService().GetCollection(name)
}

func getService() *DatabaseService {
	return services.GetService[*DatabaseService]("database")
}
