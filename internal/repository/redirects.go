package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"techodom/pkg/cache/inmemory"
)

type RedirectRepository struct {
	pool  *pgxpool.Pool
	log   *zap.Logger
	cache inmemory.Cache
}

func NewRedirectRepository(pool *pgxpool.Pool, log *zap.Logger, cache inmemory.Cache) *RedirectRepository {
	return &RedirectRepository{pool: pool, log: log, cache: cache}
}

func (r *RedirectRepository) Find(code string) (string, error) {
	if v, ok := r.cache.Get(code); ok {
		return v, nil
	}
	// check is active?
	var link string
	err := r.pool.QueryRow(context.TODO(), "SELECT history_link FROM links WHERE active_link  = $1", code).Scan(&link)
	if errors.Is(err, pgx.ErrNoRows) {
		// check is history?
		err = r.pool.QueryRow(context.TODO(), "SELECT active_link FROM links WHERE history_link  = $1", code).Scan(&link)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("not found")
		} else if err != nil {
			r.log.Error("error while getting link", zap.Error(err))
			return "", err
		}
		r.cache.Add(code, link)
		return link, nil
	}
	r.cache.Add(code, code)
	return code, nil
}
