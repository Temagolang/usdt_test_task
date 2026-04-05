package rates

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	repo "github.com/example/grinex-rates-service/internal/repo/rates"
)

// --- test doubles ---

type stubFetcher struct {
	book *OrderBook
	err  error
}

func (f *stubFetcher) FetchDepth(_ context.Context) (*OrderBook, error) {
	return f.book, f.err
}

type spyRepo struct {
	saved  *repo.SaveRateParams
	err    error
	called bool
}

func (r *spyRepo) SaveRate(_ context.Context, params repo.SaveRateParams) error {
	r.called = true
	r.saved = &params

	return r.err
}

// --- helpers ---

func testBook() *OrderBook {
	return &OrderBook{
		Asks: []OrderBookEntry{
			{Price: d("80.84"), Volume: d("100")},
			{Price: d("80.85"), Volume: d("200")},
			{Price: d("80.86"), Volume: d("300")},
		},
		Bids: []OrderBookEntry{
			{Price: d("80.75"), Volume: d("100")},
			{Price: d("80.74"), Volume: d("200")},
			{Price: d("80.73"), Volume: d("300")},
		},
		FetchedAt: time.Date(2026, 4, 5, 18, 50, 0, 0, time.UTC),
	}
}

// unsupportedRequest is a test-only Request to reach the default branch.
type unsupportedRequest struct{}

func (unsupportedRequest) algorithm() {}

// --- tests ---

func TestGetRates_TopN(t *testing.T) {
	t.Parallel()

	book := testBook()
	fetcher := &stubFetcher{book: book}
	repository := &spyRepo{}
	svc := NewService(fetcher, repository)

	rate, err := svc.GetRates(context.Background(), TopNRequest{N: 2})
	require.NoError(t, err)

	require.True(t, rate.Ask.Equal(d("80.85")), "ask: got %s", rate.Ask)
	require.True(t, rate.Bid.Equal(d("80.74")), "bid: got %s", rate.Bid)
	require.Equal(t, book.FetchedAt, rate.FetchedAt)

	require.True(t, repository.called, "repo.SaveRate must be called")
	require.True(t, repository.saved.Ask.Equal(d("80.85")), "saved ask: got %s", repository.saved.Ask)
	require.True(t, repository.saved.Bid.Equal(d("80.74")), "saved bid: got %s", repository.saved.Bid)
	require.Equal(t, book.FetchedAt, repository.saved.FetchedAt)
}

func TestGetRates_AvgNM(t *testing.T) {
	t.Parallel()

	fetcher := &stubFetcher{book: testBook()}
	repository := &spyRepo{}
	svc := NewService(fetcher, repository)

	rate, err := svc.GetRates(context.Background(), AvgNMRequest{N: 1, M: 3})
	require.NoError(t, err)

	// asks: (80.84 + 80.85 + 80.86) / 3 = 80.85
	require.True(t, rate.Ask.Equal(d("80.85")), "ask: got %s", rate.Ask)
	// bids: (80.75 + 80.74 + 80.73) / 3 = 80.74
	require.True(t, rate.Bid.Equal(d("80.74")), "bid: got %s", rate.Bid)
}

func TestGetRates_FetcherError(t *testing.T) {
	t.Parallel()

	fetcher := &stubFetcher{err: errors.New("connection refused")}
	repository := &spyRepo{}
	svc := NewService(fetcher, repository)

	_, err := svc.GetRates(context.Background(), TopNRequest{N: 1})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fetch depth")
	require.False(t, repository.called, "repo must not be called on fetch error")
}

func TestGetRates_RepoError(t *testing.T) {
	t.Parallel()

	fetcher := &stubFetcher{book: testBook()}
	repository := &spyRepo{err: errors.New("db is down")}
	svc := NewService(fetcher, repository)

	_, err := svc.GetRates(context.Background(), TopNRequest{N: 1})
	require.Error(t, err)
	require.Contains(t, err.Error(), "save rate")
}

func TestGetRates_EmptyBook(t *testing.T) {
	t.Parallel()

	fetcher := &stubFetcher{book: &OrderBook{
		FetchedAt: time.Now(),
	}}
	repository := &spyRepo{}
	svc := NewService(fetcher, repository)

	_, err := svc.GetRates(context.Background(), TopNRequest{N: 1})
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty entries")
	require.False(t, repository.called, "repo must not be called on algorithm error")
}

func TestGetRates_UnsupportedRequest(t *testing.T) {
	t.Parallel()

	fetcher := &stubFetcher{book: testBook()}
	repository := &spyRepo{}
	svc := NewService(fetcher, repository)

	_, err := svc.GetRates(context.Background(), unsupportedRequest{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported algorithm")
}
