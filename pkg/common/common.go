package common

import (
	"log"
	"os"
)

const (
	MONGO_DB_BINDING = "mongodb"
	DB_NAME_KEY      = "DB_NAME"
	DB_URL_KEY       = "DB_URL"
	RESPOSNE_SUCCESS = "Success"

	POST   = "POST"
	PUT    = "PUT"
	GET    = "GET"
	DELETE = "DELETE"
)

func CheckErrorWithPanic(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
		// panic(err.Error())
	}
}

func GetEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
