package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/unistack-org/micro/v3/logger"
	serviceapi "github.com/vielendanke/go-micro-test-service/service"
)

type message string

func setupDB(url string, errCh chan<- error, dbCh chan<- *sql.DB) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		errCh <- err
	}
	if err := db.Ping(); err != nil {
		errCh <- err
	}
	for i := 0; i < 2; i++ {
		dbCh <- db
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	dbCh := make(chan *sql.DB, 2)

	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.TraceLevel))

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		errCh <- fmt.Errorf("%s", <-sigChan)
	}()

	go setupDB("host=localhost dbname=messages_db sslmode=disable user=user password=userpassword", errCh, dbCh)

	go serviceapi.StartHTTPService(ctx, errCh, dbCh)
	go serviceapi.StartGRPCService(ctx, errCh, dbCh)
	go serviceapi.StartGithubService(ctx, errCh)

	fmt.Printf("Terminated: %v\n", <-errCh)
	// fmt.Println("Unregister service from consul")
	// consulregister.Unregister()
}
