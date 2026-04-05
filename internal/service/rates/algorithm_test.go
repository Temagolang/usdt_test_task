package rates

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func d(s string) decimal.Decimal { return decimal.RequireFromString(s) }

func entries(prices ...string) []OrderBookEntry {
	out := make([]OrderBookEntry, len(prices))
	for i, p := range prices {
		out[i] = OrderBookEntry{Price: d(p), Volume: d("1")}
	}
	return out
}

// --- TopN ---

func TestTopN(t *testing.T) {
	t.Parallel()

	book := entries("10.5", "20.3", "30.1", "40.9", "50.7")

	tests := []struct {
		name    string
		entries []OrderBookEntry
		n       int
		want    string
		wantErr string
	}{
		{name: "first", entries: book, n: 1, want: "10.5"},
		{name: "last", entries: book, n: 5, want: "50.7"},
		{name: "middle", entries: book, n: 3, want: "30.1"},

		{name: "empty_slice", entries: nil, n: 1, wantErr: "empty entries"},
		{name: "n_zero", entries: book, n: 0, wantErr: "n must be >= 1"},
		{name: "n_negative", entries: book, n: -1, wantErr: "n must be >= 1"},
		{name: "n_out_of_range", entries: book, n: 6, wantErr: "out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := TopN(tt.entries, tt.n)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.True(t, got.Equal(d(tt.want)),
				"got %s, want %s", got, tt.want)
		})
	}
}

// --- AvgNM ---

func TestAvgNM(t *testing.T) {
	t.Parallel()

	book := entries("10", "20", "30", "40", "50")

	tests := []struct {
		name    string
		entries []OrderBookEntry
		n, m    int
		want    string
		wantErr string
	}{
		{name: "single_element", entries: book, n: 3, m: 3, want: "30"},
		{name: "normal_range", entries: book, n: 2, m: 4, want: "30"},
		{name: "full_range", entries: book, n: 1, m: 5, want: "30"},
		{name: "first_two", entries: book, n: 1, m: 2, want: "15"},

		{name: "empty_slice", entries: nil, n: 1, m: 1, wantErr: "empty entries"},
		{name: "n_zero", entries: book, n: 0, m: 3, wantErr: "n must be >= 1"},
		{name: "m_zero", entries: book, n: 1, m: 0, wantErr: "m must be >= 1"},
		{name: "n_negative", entries: book, n: -1, m: 3, wantErr: "n must be >= 1"},
		{name: "n_greater_than_m", entries: book, n: 4, m: 2, wantErr: "n=4 > m=2"},
		{name: "m_out_of_range", entries: book, n: 1, m: 6, wantErr: "out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := AvgNM(tt.entries, tt.n, tt.m)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.True(t, got.Equal(d(tt.want)),
				"got %s, want %s", got, tt.want)
		})
	}
}

func TestAvgNM_DecimalPrecision(t *testing.T) {
	t.Parallel()

	// 10 + 20 + 30 = 60, 60/3 = 20 — exact
	// 10 + 20 = 30, 30/3... wait, let's use values that don't divide evenly.
	book := entries("10", "20", "33")

	got, err := AvgNM(book, 1, 3)
	require.NoError(t, err)

	// (10 + 20 + 33) / 3 = 63 / 3 = 21
	require.True(t, got.Equal(d("21")), "got %s, want 21", got)

	// Now a truly non-integer result: (10 + 20) / 3 won't happen since n..m is inclusive.
	// Use: 10 + 11 = 21, 21/2 = 10.5
	book2 := entries("10", "11")

	got2, err := AvgNM(book2, 1, 2)
	require.NoError(t, err)
	require.True(t, got2.Equal(d("10.5")), "got %s, want 10.5", got2)

	// Thirds: 1 + 2 + 3 = 6, 6/3 = 2
	// More interesting: 1 + 2 = 3, 3/2 = 1.5
	// Truly repeating: 10 + 20 + 30 = 60... all divide evenly.
	// Use 3 values that give repeating decimal: 1 + 1 + 1 = 3, but /3 = 1.
	// 1 + 2 + 4 = 7, 7/3 = 2.333...
	book3 := entries("1", "2", "4")

	got3, err := AvgNM(book3, 1, 3)
	require.NoError(t, err)

	// shopspring/decimal: 7/3 = 2.3333333333333333 (16 digits)
	want := d("7").Div(d("3"))
	require.True(t, got3.Equal(want),
		"got %s, want %s", got3, want)
}
