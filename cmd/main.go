package main

import (
	"errors"
	"fmt"
	"log"

	mongoSvc "github.com/sgahlot/go-mongo-service-quickstart/pkg/mongo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	SERVER_PORT = ":8080"
)

func mongoRoute(service mongoSvc.Service) http.Handler {
	ctx := mongoSvc.GetContext()

	endPoints := mongoSvc.EndPoints{
		InsertFruit: mongoSvc.InsertFruit(service),
		DeleteFruit: mongoSvc.DeleteFruits(service),
		GetFruit:    mongoSvc.GetFruit(service),
		GetFruits:   mongoSvc.GetFruits(service),
	}

	router := mongoSvc.CreateHandlers(ctx, endPoints)

	return router
}

func main() {
	router := mongoRoute(&mongoSvc.FruitService{})

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
