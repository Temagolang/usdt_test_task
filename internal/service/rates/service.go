package rates

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	repo "github.com/example/grinex-rates-service/internal/repo/rates"
)

// DepthFetcher retrieves order book data from the exchange.
type DepthFetcher interface {
	FetchDepth(ctx context.Context) (*OrderBook, error)
}

// Service orchestrates rate calculation: fetch → compute → persist.
type Service struct {
	fetcher DepthFetcher
	repo    repo.Repository
}

// NewService creates a rates service.
func NewService(fetcher DepthFetcher, repo repo.Repository) *Service {
	return &Service{
		fetcher: fetcher,
		repo:    repo,
	}
}

// GetRates fetches the order book, applies the requested algorithm,
// persists the result, and returns the computed rate.
func (s *Service) GetRates(ctx context.Context, req Request) (*Rate, error) {
	book, err := s.fetcher.FetchDepth(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch depth: %w", err)
	}

	ask, bid, err := s.calculate(book, req)
	if err != nil {
		return nil, fmt.Errorf("calculate rate: %w", err)
	}

	rate := &Rate{
		Ask:       ask,
		Bid:       bid,
		FetchedAt: book.FetchedAt,
	}

	err = s.repo.SaveRate(ctx, repo.SaveRateParams{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		FetchedAt: rate.FetchedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("save rate: %w", err)
	}

	return rate, nil
}

func (s *Service) calculate(book *OrderBook, req Request) (ask, bid decimal.Decimal, err error) {
	switch r := req.(type) {
	case TopNRequest:
		ask, err = TopN(book.Asks, r.N)
		if err != nil {
			return ask, bid, fmt.Errorf("ask topN: %w", err)
		}

		bid, err = TopN(book.Bids, r.N)
		if err != nil {
			return ask, bid, fmt.Errorf("bid topN: %w", err)
		}

	case AvgNMRequest:
		ask, err = AvgNM(book.Asks, r.N, r.M)
		if err != nil {
			return ask, bid, fmt.Errorf("ask avgNM: %w", err)
		}

		bid, err = AvgNM(book.Bids, r.N, r.M)
		if err != nil {
			return ask, bid, fmt.Errorf("bid avgNM: %w", err)
		}

	default:
		return ask, bid, fmt.Errorf("unsupported algorithm: %T", req)
	}

	return ask, bid, nil
}
