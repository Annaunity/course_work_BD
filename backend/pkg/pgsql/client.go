package pgsql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"sync"
	"time"
)

var once sync.Once

func NewClient(ctx context.Context, host, port, username, password, database string) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)

	var client *pgxpool.Pool
	var err error

	once.Do(func() {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		client, err = pgxpool.Connect(ctx, connString)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client. error: %w", err)
	}

	err = client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres. error: %w", err)
	}

	return client, nil
}
