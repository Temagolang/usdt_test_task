package grinex

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"

	"github.com/example/grinex-rates-service/internal/service/rates"
)

const depthPath = "/api/v1/spot/depth"

// Client fetches order book data from Grinex exchange.
type Client struct {
	http   *resty.Client
	symbol string
}

// New creates a Grinex client for the given base URL and trading symbol.
func New(baseURL, symbol string) *Client {
	return &Client{
		http:   resty.New().SetBaseURL(baseURL),
		symbol: symbol,
	}
}

// FetchDepth retrieves the current order book from Grinex depth API.
func (c *Client) FetchDepth(ctx context.Context) (*rates.OrderBook, error) {
	var raw depthResponse

	resp, err := c.http.R().
		SetContext(ctx).
		SetQueryParam("symbol", c.symbol).
		SetResult(&raw).
		Get(depthPath)
	if err != nil {
		return nil, fmt.Errorf("grinex depth request: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("grinex depth: HTTP %d", resp.StatusCode())
	}

	return mapDepthResponse(&raw)
}

func mapDepthResponse(raw *depthResponse) (*rates.OrderBook, error) {
	asks, err := mapEntries(raw.Asks)
	if err != nil {
		return nil, fmt.Errorf("parse asks: %w", err)
	}

	bids, err := mapEntries(raw.Bids)
	if err != nil {
		return nil, fmt.Errorf("parse bids: %w", err)
	}

	return &rates.OrderBook{
		Asks:      asks,
		Bids:      bids,
		FetchedAt: time.Unix(raw.Timestamp, 0),
	}, nil
}

func mapEntries(raw []depthEntry) ([]rates.OrderBookEntry, error) {
	out := make([]rates.OrderBookEntry, len(raw))
	for i, e := range raw {
		price, err := decimal.NewFromString(e.Price)
		if err != nil {
			return nil, fmt.Errorf("entry %d: invalid price %q: %w", i, e.Price, err)
		}
		volume, err := decimal.NewFromString(e.Volume)
		if err != nil {
			return nil, fmt.Errorf("entry %d: invalid volume %q: %w", i, e.Volume, err)
		}
		out[i] = rates.OrderBookEntry{Price: price, Volume: volume}
	}
	return out, nil
}
