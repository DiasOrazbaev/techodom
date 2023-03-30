package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	timeout = 5 * time.Second
)

func NewConnectionPool(host, user, password, dbname, port string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
