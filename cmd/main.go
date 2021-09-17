package main

import (
    // "database/sql"
    "fmt"
    "github.com/go-kit/kit/log"
    // "github.com/gorilla/mux"
    "github.com/sgahlot/go-service-quickstart/pkg/common"
    mongoSvc "github.com/sgahlot/go-service-quickstart/pkg/mongo"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

const (
    GET = "GET"
    DELETE = "DELETE"
    SERVER_PORT = ":9090"
    MONGO_DB = "mongodb"
    POSTGRES_DB = "postgres"
    DEFAULT_DB_TYPE = MONGO_DB
)

// func getDbConnection() *sql.DB {
//     dbType := common.GetEnvOrDefault(common.DB_TYPE_KEY, DEFAULT_DB_TYPE)
//     if dbType == MONGO_DB {
//         return nil
//     }
//     return service.GetPostgresDbConnection()
// }
//
// func postgresRoute(logger log.Logger) http.Handler {
//     db := getDbConnection()
//     svc := service.GetAccountService(db, logger, "postgres")
//
//     router := mux.NewRouter()
//     http.Handle("/", router)
//     http.Handle("/account", endpoint.GetCreateAccountHandler(svc))
//     http.Handle("/account/update", endpoint.UpdateCustomerHandler(svc))
//
//     router.Handle("/account/all", endpoint.AllCustomersHandler(svc)).Methods(GET)
//     router.Handle("/account/{customerId}", endpoint.CustomerIdHandler(svc)).Methods(GET)
//     router.Handle("/account/{custoemrId}", endpoint.DeleteCustomerHandler(svc)).Methods(DELETE)
//
//     return router
// }

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

    dbType := common.GetEnvOrDefault(common.DB_TYPE_KEY, DEFAULT_DB_TYPE)

    var router http.Handler
    if dbType == POSTGRES_DB {
         //router = postgresRoute(logger)
    } else {
        router = mongoRoute(logger)
    }

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
