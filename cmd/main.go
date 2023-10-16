package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weathersvc/internal/http"
	"weathersvc/internal/open_meteo"
	"weathersvc/internal/storage"
	"weathersvc/internal/worker"
)

var version = "latest"

const (
	openMeteoAPIURL = "https://api.open-meteo.com/v1/forecast"
	postgresURL     = "postgres://admin:admin@localhost:5432/weathersvc"
	listenAddr      = "localhost:8080"
	workerInterval  = time.Minute * 15
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	logger := makeLogger("main")
	logger.Printf("Weathersvc ver: %s", version)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		logger.Println("Interrupt signal")
		cancel()
	}()
	api := open_meteo.NewOpenMeteoAPIClient(openMeteoAPIURL, time.Second*5)
	strg, err := storage.NewPostgresStorage(ctx, postgresURL)
	if err != nil {
		return fmt.Errorf("make pg storage %s", err)
	}
	wrkr := worker.NewWorker(makeLogger("worker"), workerInterval, api, strg)
	go wrkr.Start(ctx)
	server := http.NewServer(makeLogger("http_server"), strg)

	server.Start(ctx, listenAddr)

	return nil
}

func makeLogger(name string) *log.Logger {
	prefix := fmt.Sprintf("[%s] ", name)
	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
}
