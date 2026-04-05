package rates

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// TopN returns the price at position n (1-based) in the entries slice.
func TopN(entries []OrderBookEntry, n int) (decimal.Decimal, error) {
	if n <= 0 {
		return decimal.Decimal{}, fmt.Errorf("topN: n must be >= 1, got %d", n)
	}
	if len(entries) == 0 {
		return decimal.Decimal{}, fmt.Errorf("topN: empty entries")
	}
	if n > len(entries) {
		return decimal.Decimal{}, fmt.Errorf("topN: n=%d out of range (len=%d)", n, len(entries))
	}
	return entries[n-1].Price, nil
}

// AvgNM returns the average price of entries in range [n, m] (1-based, inclusive).
func AvgNM(entries []OrderBookEntry, n, m int) (decimal.Decimal, error) {
	if n <= 0 {
		return decimal.Decimal{}, fmt.Errorf("avgNM: n must be >= 1, got %d", n)
	}
	if m <= 0 {
		return decimal.Decimal{}, fmt.Errorf("avgNM: m must be >= 1, got %d", m)
	}
	if n > m {
		return decimal.Decimal{}, fmt.Errorf("avgNM: n=%d > m=%d", n, m)
	}
	if len(entries) == 0 {
		return decimal.Decimal{}, fmt.Errorf("avgNM: empty entries")
	}
	if m > len(entries) {
		return decimal.Decimal{}, fmt.Errorf("avgNM: m=%d out of range (len=%d)", m, len(entries))
	}

	sum := decimal.Zero
	for _, e := range entries[n-1 : m] {
		sum = sum.Add(e.Price)
	}
	count := decimal.NewFromInt(int64(m - n + 1))
	return sum.Div(count), nil
}
