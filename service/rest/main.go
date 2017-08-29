package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vincentserpoul/gorestarter/pkg/rest"
	"github.com/vincentserpoul/gorestarter/pkg/storage"
)

func main() {

	// Get the config
	conf := newConfig()

	// Get the MySQL conn pool
	sqlConnPool, errQ := storage.NewMySQLDBConnPool(conf.MySQLDBConf)
	if errQ != nil {
		log.Fatal(errQ)
	}

	// Initiate the logger
	logger := logrus.New()
	// logger.Formatter = &logrus.JSONFormatter{}

	srv := rest.New(conf.HTTPPort, sqlConnPool, logger)
	fmt.Printf("Listening on port :%d\n", conf.HTTPPort)

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	errS := srv.Shutdown(ctx)
	if errS != nil {
		log.Fatal(errS)
	}

	log.Println("Server gracefully stopped")
}
