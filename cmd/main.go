package main

import (
	"fmt"
	"github.com/go-kit/kit/log"

	mongoSvc "github.com/sgahlot/go-service-quickstart/pkg/mongo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	SERVER_PORT = ":9090"
)

func mongoRoute(logger log.Logger) http.Handler {
	ctx := mongoSvc.GetContext()

	// var svc mongoSvc.Service
	svc := &mongoSvc.FruitService{}
	endPoints := mongoSvc.EndPoints{
		InsertFruit: mongoSvc.InsertFruit(svc),
		DeleteFruit: mongoSvc.DeleteFruits(svc),
		GetFruit:    mongoSvc.GetFruit(svc),
		GetFruits:   mongoSvc.GetFruits(svc),
	}

	{
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	router := mongoSvc.CreateHandlers(ctx, endPoints)

	return router
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	router := mongoRoute(logger)

	logger.Log("msg", "HTTP", "addr", SERVER_PORT)

	errChan := make(chan error)
	go func() {
		// logger = log.With(logger, "status", "Starting fruit shop at port %s", SERVER_PORT)
		handler := router
		errChan <- http.ListenAndServe(SERVER_PORT, handler)
	}()

	// logger.Log("err", http.ListenAndServe(SERVER_PORT, nil))

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- logger.Log("Exit status", <-c)
	}()

	fmt.Println(<-errChan)
}
