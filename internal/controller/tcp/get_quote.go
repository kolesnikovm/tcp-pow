package tcp

import (
	"context"

	pb "github.com/kolesnikovm/tcp-pow/pkg/proto/gen"
)

func (s *Connection) getQuote(ctx context.Context, _ *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	quote, err := s.wisdomService.GetQuote(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.QuoteResponse{
		Text: quote.Text,
	}

	return resp, nil
}
