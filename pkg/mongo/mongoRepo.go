package mongo

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "github.com/sgahlot/go-service-quickstart/pkg/common"
    "time"
)

const (
    DEFAULT_DB_NAME = "sample-db"
    DEFAULT_DB_HOST = "localhost"
    DEFAULT_DB_PORT = "27017"
    DEFAULT_DB_USER = "test-user"
    DEFAULT_DB_PASS = "test-pass"
    DEFAULT_DB_SRV = "false"
    DEFAULT_COLLECTION = "fruit"
    // dbUrl = "mongodb://test-user:test-pass@localhost:27017/sample-db"
)

var dbName = common.GetEnvOrDefault(common.DB_NAME_KEY, DEFAULT_DB_NAME)
var dbUser = common.GetEnvOrDefault(common.DB_USER_KEY, DEFAULT_DB_USER)
var dbPass = common.GetEnvOrDefault(common.DB_PASS_KEY, DEFAULT_DB_PASS)
var dbHost = common.GetEnvOrDefault(common.DB_HOST_KEY, DEFAULT_DB_HOST)
var dbPort = common.GetEnvOrDefault(common.DB_PORT_KEY, DEFAULT_DB_PORT)
var srv = common.GetEnvOrDefault(common.DB_SRV_KEY, DEFAULT_DB_SRV)

func GetContext() context.Context {
    return context.Background()
}

func DbContextWithCancel() (context.Context, context.CancelFunc) {
    return context.WithTimeout(GetContext(), 20*time.Second)
}

func getMongoDbConnectionString() string {
    sqlConnectionStr := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
        dbUser, dbPass, dbHost, dbPort, dbName)
    log.Printf("DB Connection Str: %s\n", sqlConnectionStr)

    return sqlConnectionStr
}

func GetMongoDbCollection(collection string) *mongo.Collection {
    ctx, cancel := DbContextWithCancel()
    defer cancel()

    options := options.Client().ApplyURI(getMongoDbConnectionString());

    client, err := mongo.Connect(ctx, options)
    common.CheckErrorWithPanic(err, "while connecting to MongoDB")

    return client.Database(dbName).Collection(collection)
}
