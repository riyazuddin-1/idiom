package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var logger = log.Default()

type Postgres struct {
	pool *pgxpool.Pool
}

func Init(ctx context.Context, dsn string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.pool.Ping(ctx)
	if err != nil {
		logger.Printf("Failed to ping Postgres database: %v", err)
		return err
	}
	return nil
}

func (p *Postgres) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		logger.Printf("Failed to execute query: %v", err)
		return pgconn.CommandTag{}, err
	}
	return result, nil
}

func (p *Postgres) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return p.pool.QueryRow(ctx, query, args...)
}

func (p *Postgres) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		logger.Printf("Failed to query Postgres database: %v", err)
		return nil, err
	}
	return rows, nil
}

func (p *Postgres) Begin(ctx context.Context) (pgx.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		logger.Printf("Failed to begin Postgres transaction: %v", err)
		return nil, err
	}

	return tx, nil
}
