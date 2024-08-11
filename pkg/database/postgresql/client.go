package postgresql

import (
	"context"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type ClientConfig struct {
	Host     string
	Login    string
	Password string
	DBName   string
}

type Client struct {
	conn   *pgx.Conn
	logger *slog.Logger
}

func NewClient(ctx context.Context, cfg *ClientConfig, logger *slog.Logger) (*Client, error) {
	conn, err := pgx.Connect(ctx, connectionString(cfg))
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		logger: logger,
	}, nil
}

func connectionString(cfg *ClientConfig) string {
	const baseURL = "postgresql://"
	const colonSep = ":"
	const slashSep = "/"
	const sep = "@"

	var builder strings.Builder

	if cfg.Login == "" && cfg.Password == "" {
		builder.Grow(len(baseURL) + len(cfg.Host) + len(slashSep) + len(cfg.DBName))
		builder.WriteString(baseURL)
		builder.WriteString(cfg.Host)
		builder.WriteString(slashSep)
		builder.WriteString(cfg.DBName)

		return builder.String()
	}

	builder.Grow(len(baseURL) +
		len(cfg.Login) +
		len(colonSep) +
		len(cfg.Password) +
		len(sep) +
		len(cfg.Host) +
		len(slashSep) +
		len(cfg.DBName))
	builder.WriteString(baseURL)
	builder.WriteString(cfg.Login)
	builder.WriteString(colonSep)
	builder.WriteString(cfg.Password)
	builder.WriteString(sep)
	builder.WriteString(cfg.Host)
	builder.WriteString(slashSep)
	builder.WriteString(cfg.DBName)

	return builder.String()
}

func (c *Client) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

func (c *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.conn.Query(ctx, sql, args...)
}

func (c *Client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.conn.Exec(ctx, sql, args...)
}

func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.conn.QueryRow(ctx, sql, args...)
}

func (c *Client) ExecTx(ctx context.Context, f func(ctx context.Context, tx Tx) error) error {
	const fName = "ExecTx"

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx) {
		if err == nil {
			return
		}

		if err := tx.Rollback(ctx); err != nil {
			c.logger.InfoContext(ctx, fName, slog.String("failed to rollback transaction", err.Error()))
		}
	}(tx)

	if err = f(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (c *Client) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.conn.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
