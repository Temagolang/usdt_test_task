package grinex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestFetchDepth_HappyPath(t *testing.T) {
	t.Parallel()

	fixture, err := os.ReadFile("testdata/depth_response.json")
	require.NoError(t, err, "read fixture")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	c := New(srv.URL, "usdta7a5", nil)
	book, err := c.FetchDepth(context.Background())
	require.NoError(t, err)

	require.Len(t, book.Asks, 5)
	require.Len(t, book.Bids, 5)

	require.True(t, book.Asks[0].Price.Equal(decimal.RequireFromString("80.84")),
		"first ask price: got %s", book.Asks[0].Price)
	require.True(t, book.Bids[0].Price.Equal(decimal.RequireFromString("80.75")),
		"first bid price: got %s", book.Bids[0].Price)

	require.Equal(t, time.Unix(1775399602, 0), book.FetchedAt)
}

func TestFetchDepth_RequestPathAndSymbol(t *testing.T) {
	t.Parallel()

	var gotPath, gotSymbol string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotSymbol = r.URL.Query().Get("symbol")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"timestamp":0,"asks":[],"bids":[]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "usdta7a5", nil)
	_, err := c.FetchDepth(context.Background())
	require.NoError(t, err)

	require.Equal(t, "/api/v1/spot/depth", gotPath)
	require.Equal(t, "usdta7a5", gotSymbol)
}

func TestFetchDepth_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr string
	}{
		{
			name: "http_500",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: "HTTP 500",
		},
		{
			name: "malformed_json",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{invalid`))
			},
			wantErr: "grinex depth request",
		},
		{
			name: "invalid_ask_price",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"timestamp":1,"asks":[{"price":"not_a_number","volume":"1.0","amount":"1.0"}],"bids":[]}`))
			},
			wantErr: "invalid price",
		},
		{
			name: "invalid_bid_volume",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"timestamp":1,"asks":[],"bids":[{"price":"1.0","volume":"xxx","amount":"1.0"}]}`))
			},
			wantErr: "invalid volume",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			c := New(srv.URL, "usdta7a5", nil)
			_, err := c.FetchDepth(context.Background())
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestFetchDepth_EmptyBook(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"timestamp":1000,"asks":[],"bids":[]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "usdta7a5", nil)
	book, err := c.FetchDepth(context.Background())
	require.NoError(t, err)

	require.Empty(t, book.Asks)
	require.Empty(t, book.Bids)
}

func TestFetchDepth_VolumeParsed(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"timestamp":1,"asks":[{"price":"1.5","volume":"99.123","amount":"0"}],"bids":[]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "usdta7a5", nil)
	book, err := c.FetchDepth(context.Background())
	require.NoError(t, err)

	require.Len(t, book.Asks, 1)
	require.True(t, book.Asks[0].Volume.Equal(decimal.RequireFromString("99.123")),
		"volume: got %s", book.Asks[0].Volume)
}
