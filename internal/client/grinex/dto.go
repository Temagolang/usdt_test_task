// Package grinex provides an HTTP client for the Grinex exchange API.
package grinex

// depthResponse is the raw JSON shape returned by GET /api/v1/spot/depth?symbol=...
//
// All numeric values arrive as strings from the exchange.
// Conversion to decimal.Decimal happens in the client layer during mapping to domain types.
type depthResponse struct {
	Timestamp int64        `json:"timestamp"`
	Asks      []depthEntry `json:"asks"`
	Bids      []depthEntry `json:"bids"`
}

// depthEntry is a single level in the order book.
//
// Fields from Grinex:
//   - price:  quote price as string (e.g. "80.84")
//   - volume: base asset volume as string (e.g. "24752.6039")
//   - amount: quote amount (price * volume), not used in calculations
type depthEntry struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
}
