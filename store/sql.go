package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/deliergky/monitorink/data"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrorInsertingRecord = errors.New("error_inserting_record")

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	DB       string
	Password string
}

type PostgresStore struct {
	dbpool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, config PostgresConfig) (*PostgresStore, error) {
	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		config.Host,
		config.Port,
		config.User,
		config.DB,
		config.Password,
	)
	dbpool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{
		dbpool: dbpool,
	}, nil
}

func rollback(ctx context.Context, tx pgx.Tx, err error) error {
	if err := tx.Rollback(ctx); err != nil {
		return err
	}
	return err
}

func (s *PostgresStore) Close() {
	s.dbpool.Close()
}

func (s *PostgresStore) Persist(ctx context.Context, d data.ResponseData) error {
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	cmd, err := tx.Exec(ctx, "INSERT INTO response VALUES($1, $2, $3, $4, $5)", d.URL, d.StatusCode, d.ResponseTime, d.Matched, d.CreatedAt)

	if err != nil {
		return rollback(ctx, tx, err)
	}

	if cmd.RowsAffected() == int64(0) {
		return rollback(ctx, tx, ErrorInsertingRecord)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return rollback(ctx, tx, err)
	}
	return nil

}

func (s *PostgresStore) Retrieve(ctx context.Context) ([]data.ResponseData, error) {
	defer s.dbpool.Close()
	rows, err := s.dbpool.Query(ctx, "select * from response order by created_at desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	responses := make([]data.ResponseData, 0)
	for rows.Next() {
		var d data.ResponseData
		err := rows.Scan(&d.URL, &d.StatusCode, &d.ResponseTime, &d.Matched, &d.CreatedAt)
		if err != nil {
			return nil, err
		}
		responses = append(responses, d)
	}

	return responses, nil

}
