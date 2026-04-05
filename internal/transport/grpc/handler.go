package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	ratesv1 "github.com/example/grinex-rates-service/gen/rates/v1"
	"github.com/example/grinex-rates-service/internal/service/rates"
)

// RatesHandler implements ratesv1.RatesServiceServer.
type RatesHandler struct {
	ratesv1.UnimplementedRatesServiceServer
	svc *rates.Service
}

// NewRatesHandler creates a handler and registers it on the gRPC server.
func NewRatesHandler(reg grpc.ServiceRegistrar, svc *rates.Service) *RatesHandler {
	h := &RatesHandler{svc: svc}
	ratesv1.RegisterRatesServiceServer(reg, h)

	return h
}

// GetRates maps proto request to domain, calls service, maps result back to proto.
func (h *RatesHandler) GetRates(ctx context.Context, req *ratesv1.GetRatesRequest) (*ratesv1.GetRatesResponse, error) {
	domainReq, err := mapRequest(req)
	if err != nil {
		return nil, err
	}

	rate, err := h.svc.GetRates(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get rates: %v", err)
	}

	return &ratesv1.GetRatesResponse{
		Ask:       rate.Ask.String(),
		Bid:       rate.Bid.String(),
		Timestamp: timestamppb.New(rate.FetchedAt),
	}, nil
}

func mapRequest(req *ratesv1.GetRatesRequest) (rates.Request, error) {
	switch alg := req.GetAlgorithm().(type) {
	case *ratesv1.GetRatesRequest_TopN:
		n := int(alg.TopN.GetN())
		if n <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "top_n.n must be >= 1, got %d", n)
		}

		return rates.TopNRequest{N: n}, nil

	case *ratesv1.GetRatesRequest_AvgNm:
		n := int(alg.AvgNm.GetN())
		m := int(alg.AvgNm.GetM())

		if n <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "avg_nm.n must be >= 1, got %d", n)
		}
		if m <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "avg_nm.m must be >= 1, got %d", m)
		}
		if n > m {
			return nil, status.Errorf(codes.InvalidArgument, "avg_nm.n must be <= m, got n=%d m=%d", n, m)
		}

		return rates.AvgNMRequest{N: n, M: m}, nil

	default:
		return nil, status.Errorf(codes.InvalidArgument, "algorithm is required")
	}
}
