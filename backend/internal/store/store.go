package store

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/keto-granola/keto-granola/internal/store/db"
	"github.com/keto-granola/keto-granola/internal/store/db/generated"
	"github.com/keto-granola/keto-granola/internal/utils"
)

const (
	dbMaxRetries = 5
	dbBaseDelay  = 100 * time.Millisecond
	pingTimeout  = 5 * time.Second
)

type Store struct {
	pool    *pgxpool.Pool
	Queries *generated.Queries
}

func New(ctx context.Context, dbURL string) (*Store, error) {
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("parse pool config %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create db pool %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err = pool.Ping(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err = db.RunMigrations(dbURL); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &Store{
		pool:    pool,
		Queries: generated.New(pool),
	}, nil
}

func (s *Store) PingDB(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) Close() {
	s.pool.Close()
}

func ExecWithResult[T any](ctx context.Context, query func() (T, error)) (T, error) {
	return utils.RetryWithExponentialBackoff(ctx, query, dbMaxRetries, dbBaseDelay, isRetryableDbError)
}

func ExecNoResult(ctx context.Context, command func() error) error {
	_, err := ExecWithResult(ctx, func() (*struct{}, error) { return nil, command() })
	return err
}

var transientPostgresErrorCodes = []string{
	"08", // Connection exceptions (network problems, can't reach database)
	"40", // Transaction rollback (like deadlocks or serialisation failures)
	"53", // Insufficient resources (out of memory, disk full)
	"55", // Object not in prerequisite state (like trying to use a prepared statement that doesn't exist)
	"57", // Operator intervention (admin killed the query, database shutting down)
}

func isRetryableDbError(err error) bool {
	if err == nil {
		return false
	}

	if isNotFoundErr(err) {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		errClass := pgErr.Code[:2]
		return slices.Contains(transientPostgresErrorCodes, errClass)
	}

	return false
}

func isNotFoundErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
