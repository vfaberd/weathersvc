package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
	"weathersvc/internal/core"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

var _ core.WeatherStorage = &PostgresStorage{}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse connStr: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new pgxpool: %w", err)
	}
	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return &PostgresStorage{
		db: dbpool,
	}, nil
}

func (s *PostgresStorage) Save(ctx context.Context, weather core.Weather) error {
	w := toInternal(weather)
	row := s.db.QueryRow(ctx, queryRecordExists, w.LocationName, w.Timestamp)
	var n int64
	err := row.Scan(&n)
	if errors.Is(err, pgx.ErrNoRows) {
		if err := s.add(ctx, w); err != nil {
			return fmt.Errorf("add: %s", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("query existing: %w", err)
	}
	if err := s.update(ctx, w); err != nil {
		return fmt.Errorf("update: %s", err)
	}
	return nil
}

func (s *PostgresStorage) GetLatest(ctx context.Context, location string) (core.Weather, error) {
	row := s.db.QueryRow(ctx, queryGetLatest, location)
	var w weather
	err := row.Scan(
		&w.LocationName,
		&w.Latitude,
		&w.Longitude,
		&w.Timestamp,
		&w.Temperature2M,
		&w.Relativehumidity2M,
		&w.PrecipitationProbability,
		&w.Visibility,
		&w.Windspeed10M,
	)
	if err != nil {
		return core.Weather{}, fmt.Errorf("get latest: %w", err)
	}
	res, err := w.toCore()
	if err != nil {
		return core.Weather{}, err
	}
	return res, nil
}

func (s *PostgresStorage) GetPeriod(ctx context.Context, location string, from, to time.Time) ([]core.Weather, error) {
	panic("not implemented")
}

func (s *PostgresStorage) add(ctx context.Context, item weather) error {
	_, err := s.db.Exec(ctx, queryAdd,
		item.LocationName,
		item.Latitude,
		item.Longitude,
		item.Timestamp,
		item.Temperature2M,
		item.Relativehumidity2M,
		item.PrecipitationProbability,
		item.Visibility,
		item.Windspeed10M,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) update(ctx context.Context, item weather) error {
	_, err := s.db.Exec(ctx, queryUpdateByLocationAndTime,
		item.Latitude,
		item.Longitude,
		item.Temperature2M,
		item.Relativehumidity2M,
		item.PrecipitationProbability,
		item.Visibility,
		item.Windspeed10M,
	)
	if err != nil {
		return err
	}
	return nil
}

type weather struct {
	LocationName             string  `db:"location_name"`
	Latitude                 float64 `db:"latitude"`
	Longitude                float64 `db:"longitude"`
	Timestamp                string  `db:"timestamp"`
	Temperature2M            float64 `db:"temperature_2m"`
	Relativehumidity2M       int     `db:"relativehumidity_2m"`
	PrecipitationProbability int     `db:"precipitation_probability"`
	Visibility               float64 `db:"visibility"`
	Windspeed10M             float64 `db:"windspeed_10m"`
}

func (w *weather) toCore() (core.Weather, error) {
	t, err := time.Parse(core.TimestampFormat, w.Timestamp)
	if err != nil {
		return core.Weather{}, fmt.Errorf("convert to core: %w", err)
	}
	loc := core.Location{
		Name:      w.LocationName,
		Latitude:  w.Latitude,
		Longitude: w.Longitude,
	}
	return core.Weather{
		Location:                 loc,
		Timestamp:                t,
		Temperature2M:            w.Temperature2M,
		Relativehumidity2M:       w.Relativehumidity2M,
		PrecipitationProbability: w.PrecipitationProbability,
		Visibility:               w.Visibility,
		Windspeed10M:             w.Windspeed10M,
	}, nil
}

func toInternal(w core.Weather) weather {
	return weather{
		LocationName:             w.Location.Name,
		Latitude:                 w.Location.Latitude,
		Longitude:                w.Location.Longitude,
		Timestamp:                w.Timestamp.Format(core.TimestampFormat),
		Temperature2M:            w.Temperature2M,
		Relativehumidity2M:       w.Relativehumidity2M,
		PrecipitationProbability: w.PrecipitationProbability,
		Visibility:               w.Visibility,
		Windspeed10M:             w.Windspeed10M,
	}
}
