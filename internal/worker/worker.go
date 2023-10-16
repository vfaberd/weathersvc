package worker

import (
	"context"
	"fmt"
	"log"
	"time"
	"weathersvc/internal/core"
)

type Worker struct {
	logger   *log.Logger
	interval time.Duration

	fetcher core.WeatherFetcher
	storage core.WeatherStorage
}

func NewWorker(logger *log.Logger, interval time.Duration, fetcher core.WeatherFetcher, storage core.WeatherStorage) *Worker {
	return &Worker{
		logger:   logger,
		interval: interval,
		fetcher:  fetcher,
		storage:  storage,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	w.logger.Println("Worker started")
	w.loop(ctx)
	for {
		select {
		case <-ctx.Done():
			w.logger.Println("Worker stopped")
			return
		case <-ticker.C:
			w.loop(ctx)
		}
	}
}

func (w *Worker) loop(ctx context.Context) {
	var errCount int
	for _, loc := range core.KnownLocations {
		err := w.processLocation(ctx, loc)
		if err != nil {
			errCount++
			w.logger.Printf("[ERR] %s", err)
		}
	}
	w.logger.Printf("Loop complete. Successful: %d of %d", len(core.KnownLocations)-errCount, len(core.KnownLocations))
}

func (w *Worker) processLocation(ctx context.Context, loc core.Location) error {
	weather, err := w.fetcher.GetCurrent(ctx, loc)
	if err != nil {
		return fmt.Errorf("get current: %w", err)
	}
	err = w.storage.Save(ctx, weather)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}
