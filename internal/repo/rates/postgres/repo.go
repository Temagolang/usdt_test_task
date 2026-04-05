// Package postgres implements the rates repository using PostgreSQL.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	repo "github.com/example/grinex-rates-service/internal/repo/rates"
	db "github.com/example/grinex-rates-service/internal/repo/rates/postgres/sqlcgen"
)

// Repo implements repo.Repository using sqlc-generated queries.
type Repo struct {
	q *db.Queries
}

// New creates a Repo backed by the given database connection.
func New(conn db.DBTX) *Repo {
	return &Repo{q: db.New(conn)}
}

// SaveRate persists a rate record.
func (r *Repo) SaveRate(ctx context.Context, params repo.SaveRateParams) error {
	err := r.q.InsertRate(ctx, db.InsertRateParams{
		Ask:       numericFrom(params.Ask),
		Bid:       numericFrom(params.Bid),
		FetchedAt: timestamptzFrom(params.FetchedAt),
	})
	if err != nil {
		return fmt.Errorf("insert rate: %w", err)
	}

	return nil
}

func numericFrom(d decimal.Decimal) pgtype.Numeric {
	// shopspring/decimal: coefficient is *big.Int, exponent is int32.
	// pgtype.Numeric expects the same layout.
	return pgtype.Numeric{
		Int:   d.Coefficient(),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

func timestamptzFrom(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
