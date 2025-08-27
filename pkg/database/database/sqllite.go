package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
)

type SQLiteConn struct {
	db *sql.DB
}

type SQLiteRow struct {
	row *sql.Row
}

type SQLiteCommandTag struct {
	res sql.Result
}

func (s *SQLiteConn) QueryRow(ctx context.Context, query string, args ...any) Row {
	return &SQLiteRow{row: s.db.QueryRowContext(ctx, query, args...)}
}

func (r *SQLiteRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (s *SQLiteConn) Exec(ctx context.Context, query string, args ...any) (CommandTag, error) {
	res, err := s.db.ExecContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return &SQLiteCommandTag{res: res}, nil
}

func (s *SQLiteConn) IsErrNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func (t *SQLiteCommandTag) RowsAffected() int64 {
	rows, _ := t.res.RowsAffected()
	return rows
}

func ConnSQLite(path string) DBClient {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	return &SQLiteConn{db: db}
}
