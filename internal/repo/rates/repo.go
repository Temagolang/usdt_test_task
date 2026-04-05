// Package rates defines the repository contract for persisting exchange rates.
package rates

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// SaveRateParams holds the data needed to persist a single rate record.
// Defined in repo package (not service) to avoid import cycles.
type SaveRateParams struct {
	Ask       decimal.Decimal
	Bid       decimal.Decimal
	FetchedAt time.Time
}

// Repository persists exchange rate data.
type Repository interface {
	SaveRate(ctx context.Context, params SaveRateParams) error
}
