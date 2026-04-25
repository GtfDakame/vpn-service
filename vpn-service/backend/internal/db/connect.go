package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"vpn-service/internal/config"
)

func Connect(cfg *config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}
	return conn, nil
}
