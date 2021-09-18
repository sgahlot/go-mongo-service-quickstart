package common

import (
	"log"
	"os"
)

const (
	DB_NAME_KEY      = "DB_NAME"
	DB_HOST_KEY      = "DB_HOST"
	DB_PORT_KEY      = "DB_PORT"
	DB_USER_KEY      = "DB_USER"
	DB_PASS_KEY      = "DB_PASS"
	DB_SRV_KEY       = "DB_SRV"
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
