package database

import "context"

type DBClient interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (CommandTag, error)
}

type CommandTag interface {
	RowsAffected() int64
}

type Row interface {
	Scan(dest ...any) error
}
