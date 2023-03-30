package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"techodom/internal/entity"
	"techodom/internal/utils"
)

type Admin struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewAdmin(pool *pgxpool.Pool, logger *zap.Logger) *Admin {
	return &Admin{pool: pool, logger: logger}
}

func (a *Admin) FindByID(id string) (*entity.Redirect, error) {
	query := `SELECT active_link, history_link FROM links WHERE id = $1`
	row := a.pool.QueryRow(context.TODO(), query, id)
	var link entity.Redirect
	err := row.Scan(&link.ActiveLink, &link.HistoryLink)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, utils.ErrNotFound
	} else if err != nil {
		a.logger.Error("error while scanning row", zap.Error(err))
		return nil, err
	}
	return &link, nil
}

func (a *Admin) GetLinks(page int, perPage int) ([]*entity.Redirect, error) {
	offset := (page - 1) * perPage
	query := `SELECT history_link, active_link FROM links where id > $1 LIMIT $2
`
	rows, err := a.pool.Query(context.TODO(), query, offset, perPage)
	if err != nil {
		a.logger.Error("error while getting links", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var links []*entity.Redirect
	for rows.Next() {
		var r entity.Redirect
		if err := rows.Scan(&r.HistoryLink, &r.ActiveLink); err != nil {
			a.logger.Error("error while scanning links", zap.Error(err))
			return nil, err
		}
		links = append(links, &r)
	}

	return links, nil
}

func (a *Admin) All(page int, perPage int) ([]*entity.Redirect, error) {
	return a.GetLinks(page, perPage)
}

func (a *Admin) Create(old, new string) error {
	query := `INSERT INTO links (history_link, active_link) VALUES ($1, $2)`
	_, err := a.pool.Exec(context.TODO(), query, old, new)
	return err
}

func (a *Admin) Update(old, new string, id int) error {
	query := `UPDATE links SET history_link = $1, active_link = $2 WHERE id = $3`
	_, err := a.pool.Exec(context.TODO(), query, old, new, id)
	return err
}

func (a *Admin) Delete(id int) error {
	query := `DELETE FROM links WHERE id = $1`
	_, err := a.pool.Exec(context.TODO(), query, id)
	return err
}
