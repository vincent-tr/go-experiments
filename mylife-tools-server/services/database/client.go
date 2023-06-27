package database

import (
	"context"
	"net/url"

	"mylife-tools-server/config"
	"mylife-tools-server/log"
	"mylife-tools-server/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

var logger = log.CreateLogger("mylife:server:database")

func init() {
	services.Register(&databaseService{})
}

type Collection = mongo.Collection

type databaseService struct {
	client   *mongo.Client
	database *mongo.Database
	session  mongo.Session
}

func (service *databaseService) Init(arg interface{}) error {
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

	session, err := client.StartSession()
	if err != nil {
		return err
	}

	service.client = client
	service.database = client.Database(dbName)
	service.session = session

	return nil
}

func (service *databaseService) Terminate() error {
	service.session.EndSession(context.TODO())
	return service.client.Disconnect(context.TODO())
}

func (service *databaseService) ServiceName() string {
	return "database"
}

func (service *databaseService) Dependencies() []string {
	return []string{}
}

func (service *databaseService) getCollection(name string) *Collection {
	return service.database.Collection(name)
}

func (service *databaseService) sessionId() bson.Raw {
	return service.session.ID()
}

func (service *databaseService) withTransaction(callback func() (interface{}, error)) (interface{}, error) {
	return service.session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		return callback()
	})
}

func getService() *databaseService {
	return services.GetService[*databaseService]("database")
}

// Public access

func GetCollection(name string) *Collection {
	return getService().getCollection(name)
}

func MakeId() string {
	return primitive.NewObjectID().Hex()
}

func EncodeId(value string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(value)
}

func DecodeId(value primitive.ObjectID) string {
	return value.Hex()
}

func SessionId() bson.Raw {
	return getService().sessionId()
}

func WithTransaction(callback func() (interface{}, error)) (interface{}, error) {
	return getService().withTransaction(callback)
}
