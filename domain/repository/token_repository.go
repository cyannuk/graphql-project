package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-project/core"
)

type TokenRepository DataSource

func (r *TokenRepository) GetTokenByID(ctx context.Context, id int64) (token string, err error) {
	connection, err := (*pgxpool.Pool)(r).Acquire(ctx)
	if err != nil {
		return
	}
	defer connection.Release()
	rows, err := connection.Query(ctx, `SELECT "token" FROM tokens WHERE "id" = $1`, id)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&token)
		return
	}
	err = rows.Err()
	if err == nil {
		err = core.ErrNotFound
	}
	return
}

func (r *TokenRepository) CreateToken(ctx context.Context, userId int64, token string) (err error) {
	connection, err := (*pgxpool.Pool)(r).Acquire(ctx)
	if err != nil {
		return
	}
	_, err = connection.Exec(ctx, `INSERT INTO tokens("id", "token") VALUES($1, $2) ON CONFLICT("id") DO UPDATE SET "token" = $2`, userId, token)
	connection.Release()
	return
}

func NewTokenRepository(dataSource *DataSource) *TokenRepository {
	return (*TokenRepository)(dataSource)
}
