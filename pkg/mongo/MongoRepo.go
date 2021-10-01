package mongo

import (
	"context"
	"github.com/myeung18/service-binding-client/pkg/binding/convert"
	"github.com/sgahlot/go-mongo-service-quickstart/pkg/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	DEFAULT_DB_NAME    = "fruit"
	DEFAULT_COLLECTION = "fruit"
	DEFAULT_DB_URL     = "mongodb://test-user:test-pass@localhost:27017/fruit"
)

var (
	dbName        = common.GetEnvOrDefault(common.DB_NAME_KEY, DEFAULT_DB_NAME)
	mongoDbClient *mongo.Client
)

func init() {
	mongoDbClient = getMongoDbConnection()
}

func GetContext() context.Context {
	return context.Background()
}

func DbContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithTimeout(GetContext(), 20*time.Second)
}

func getMongoDbConnectionStringForNonBindingsRun() string {
	return common.GetEnvOrDefault(common.DB_URL_KEY, DEFAULT_DB_URL)
}

func getMongoDbConnectionStringForBindingsRun() string {
	sqlConnectionStr, err := convert.GetMongoDBConnectionString()
	common.CheckErrorWithPanic(err, "while trying to get MongoDB connection string from Bindings")

	return sqlConnectionStr
}

func getMongoDbConnectionString() string {
	bindingsDir := common.GetEnvOrDefault("SERVICE_BINDING_ROOT", "")
	var sqlConnectionStr string
	if bindingsDir == "" {
		sqlConnectionStr = getMongoDbConnectionStringForNonBindingsRun()
	} else {
		log.Printf("System property for Bindings dir [%s] found", bindingsDir)
		sqlConnectionStr = getMongoDbConnectionStringForBindingsRun()
	}

	log.Printf("DB Connection String = %s", sqlConnectionStr)
	return sqlConnectionStr
}

func getMongoDbConnection() *mongo.Client {
	ctx, cancel := DbContextWithCancel()
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(getMongoDbConnectionString())

	var err error

	client, err := mongo.Connect(ctx, clientOptions)
	common.CheckErrorWithPanic(err, "while connecting to MongoDB")

	err = client.Ping(context.TODO(), nil)
	common.CheckErrorWithPanic(err, "while trying to ping MongoDB")

	return client
}

func checkAndRefreshConnection() {
	err := mongoDbClient.Ping(context.TODO(), nil)
	if err != nil {
		// Try to get the connection again as we might be disconnected
		mongoDbClient = getMongoDbConnection()
	}
}

func GetMongoDbCollection(collection string) *mongo.Collection {

	log.Printf("Returning collection [%s] from DB [%s]", collection, dbName)

	checkAndRefreshConnection()

	return mongoDbClient.Database(dbName).Collection(collection)
}
