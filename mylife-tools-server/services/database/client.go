package database

import (
	"context"
	"net/url"

	"mylife-tools-server/config"
	"mylife-tools-server/log"
	"mylife-tools-server/services"

	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

var logger = log.CreateLogger("mylife:server:mongo")

func init() {
	services.Register(&databaseService{})
}

type Collection = mongo.Collection

type databaseService struct {
	client   *mongo.Client
	database *mongo.Database
}

func (service *databaseService) Init() error {
	mongoUrl := config.GetString("mongo")

	parsedUrl, err := url.Parse(mongoUrl)
	if err != nil {
		return err
	}
	dbName := parsedUrl.Path[1:]

	logger.WithFields(log.Fields{"mongoUrl": mongoUrl, "dbName": dbName}).Info("Config")

	opts := options.Client().ApplyURI(mongoUrl).SetDirect(true)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	service.client = client
	service.database = client.Database(dbName)

	return nil
}

func (service *databaseService) Terminate() error {
	return service.client.Disconnect(context.TODO())
}

func (service *databaseService) ServiceName() string {
	return "database"
}

func (service *databaseService) Dependencies() []string {
	return []string{}
}

func (service *databaseService) GetCollection(name string) *Collection {
	return service.database.Collection(name)
}

func getService() *databaseService {
	return services.GetService[*databaseService]("database")
}

// Public access

func GetCollection(name string) *Collection {
	return getService().GetCollection(name)
}
