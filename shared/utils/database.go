package utils

import (
	c "context"
	"database/sql"

	q "github.com/MMortari/go-query-builder"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ListQuery[T any](
	ctx c.Context,
	db *sqlx.DB,
	span trace.Span,
	qb *q.QueryBuilder,
	capacity *int,
	withTotal bool,
) (*[]T, int, error) {
	query, args := qb.ToSelectSql()

	var resp []T

	if capacity != nil {
		resp = make([]T, 0, *capacity)
	} else {
		resp = make([]T, 0)
	}

	err := db.SelectContext(ctx, &resp, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, err
	}

	if !withTotal {
		return &resp, 0, nil
	}

	queryTotal, argsTotal := qb.ToSelectTotalSql()

	var total int

	err = db.GetContext(ctx, &total, queryTotal, argsTotal...)
	if err == sql.ErrNoRows {
		return nil, total, nil
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, total, err
	}

	return &resp, total, nil
}

func Insert(
	ctx c.Context,
	db *sqlx.DB,
	span trace.Span,
	query string,
	params any,
) error {
	_, err := db.NamedExecContext(ctx, query, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func Update(
	ctx c.Context,
	db *sqlx.DB,
	span trace.Span,
	query string,
	params any,
) error {
	_, err := db.NamedExecContext(ctx, query, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func Get[T any](
	ctx c.Context,
	db *sqlx.DB,
	span trace.Span,
	query string,
	args ...any,
) (*T, error) {
	var data T

	err := db.GetContext(ctx, &data, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &data, nil
}

func Delete(
	ctx c.Context,
	db *sqlx.DB,
	span trace.Span,
	query string,
	args ...any,
) error {
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
