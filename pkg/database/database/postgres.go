package database

import (
	"context"
	"fmt"

	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type PgxConn struct {
	conn *pgx.Conn
}

func NewPgxConn(conn *pgx.Conn) *PgxConn {
	return &PgxConn{conn: conn}
}

func (p *PgxConn) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return p.conn.QueryRow(ctx, sql, args...)
}

func (p *PgxConn) Exec(ctx context.Context, query string, args ...any) (CommandTag, error) {
	tag, err := p.conn.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PgxCommandTag{tag: tag}, nil
}

type PgxCommandTag struct {
	tag pgconn.CommandTag
}

func (p *PgxCommandTag) RowsAffected() int64 {
	return p.tag.RowsAffected()
}

func ConnectPg(cfg *config.Config) (DBClient, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseDb,
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal().Err(err).Str("connStr", connStr).Msg("Failed to connect to database")
		return nil, err
	}

	log.Info().
		Str("db", cfg.DatabaseDb).
		Str("host", cfg.DatabaseHost).
		Str("user", cfg.DatabaseUser).
		Str("port", cfg.DatabasePort).
		Msg("Successfully connected to the database")

	return NewPgxConn(conn), nil
}
