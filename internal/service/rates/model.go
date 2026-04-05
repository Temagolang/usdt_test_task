// Package rates contains domain types and business logic for USDT rate calculation.
package rates

import (
	"time"

	"github.com/shopspring/decimal"
)

// Request is a sealed interface for algorithm selection.
// Implemented only by TopNRequest and AvgNMRequest.
type Request interface {
	algorithm()
}

// TopNRequest selects the value at position N (1-based) in the order book.
type TopNRequest struct {
	N int
}

func (TopNRequest) algorithm() {}

// AvgNMRequest calculates the average of entries in range [N, M] (1-based, inclusive).
type AvgNMRequest struct {
	N int
	M int
}

func (AvgNMRequest) algorithm() {}

// OrderBookEntry is a single price/volume level from the exchange depth response.
type OrderBookEntry struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
}

// OrderBook holds the ask and bid sides of the exchange order book
// along with the timestamp when it was captured by the exchange.
type OrderBook struct {
	Asks      []OrderBookEntry
	Bids      []OrderBookEntry
	FetchedAt time.Time
}

// Rate is the computed result: ask and bid prices at a point in time.
type Rate struct {
	Ask       decimal.Decimal
	Bid       decimal.Decimal
	FetchedAt time.Time
}
