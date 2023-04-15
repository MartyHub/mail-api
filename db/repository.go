package db

import (
	"context"

	db "github.com/MartyHub/mail-api/db/gen"
	"github.com/MartyHub/mail-api/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Close()
	Now() pgtype.Timestamp
	Ping(ctx context.Context) error
	Wrap(ctx context.Context, opts pgx.TxOptions, f func(querier db.Querier) error) error
}

func NewRepository(cfg Config, clock utils.Clock) (Repository, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		return nil, err
	}

	result := repository{
		clock: clock,
		pool:  pool,
	}

	return result, result.Ping(ctx)
}

type repository struct {
	clock utils.Clock
	pool  *pgxpool.Pool
}

func (repo repository) Close() {
	repo.pool.Close()
}

func (repo repository) Now() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  repo.clock.Now(),
		Valid: true,
	}
}

func (repo repository) Ping(ctx context.Context) error {
	return repo.pool.Ping(ctx)
}

func (repo repository) Wrap(ctx context.Context, opts pgx.TxOptions, f func(querier db.Querier) error) error {
	tx, err := repo.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)

			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = f(db.New(repo.pool).WithTx(tx))

	return err
}
