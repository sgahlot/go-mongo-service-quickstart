package main

import (
	"errors"
	"fmt"
	"github.com/sgahlot/go-mongo-service-quickstart/pkg/common"
	"log"

	db "github.com/sgahlot/go-mongo-service-quickstart/pkg/db"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	SERVER_PORT = ":8080"
)

func mongoRoute(service common.Service) http.Handler {
	ctx := db.GetContext()

	endPoints := common.EndPoints{
		InsertFruit: common.InsertFruit(service),
		DeleteFruit: common.DeleteFruits(service),
		GetFruit:    common.GetFruit(service),
		GetFruits:   common.GetFruits(service),
	}

	router := common.CreateHandlers(ctx, endPoints)

	return router
}

func main() {
	router := mongoRoute(&db.FruitService{})

	errChan := make(chan error)
	go func() {
		log.Printf("Starting FruitShop server at port %s\n", SERVER_PORT)
		handler := router
		errChan <- http.ListenAndServe(SERVER_PORT, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- errors.New(fmt.Sprintf("Exit status %v", <-c))
	}()

	fmt.Println(<-errChan)
}
