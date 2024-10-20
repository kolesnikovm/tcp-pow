package wisdomclient

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/tcp-pow/internal/controller/tcp"
	pb "github.com/kolesnikovm/tcp-pow/pkg/proto/gen"
	"google.golang.org/protobuf/proto"
)

func (w *WisdomClient) GetQuote(ctx context.Context) (string, error) {
	const op = "wisdomclient.GetQuote"

	req := &pb.WrapperRequest{
		Request: &pb.WrapperRequest_QuoteRequest{
			QuoteRequest: &pb.QuoteRequest{},
		},
	}

	reqBuf, err := proto.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("%s: failed to marshall quote request: %w", op, err)
	}

	_, err = w.conn.Write(tcp.PackMessage(reqBuf))
	if err != nil {
		return "", fmt.Errorf("%s: failed to write quote request: %w", op, err)
	}

	respBuf, err := tcp.UnpackMessage(w.rd)
	if err != nil {
		return "", fmt.Errorf("%s: failed to read response: %w", op, err)
	}

	resp := &pb.WrapperResponse{}

	err = proto.Unmarshal(respBuf, resp)
	if err != nil {
		return "", fmt.Errorf("%s:failed to unmarshal response: %w", op, err)
	}

	switch resp.GetResponse().(type) {
	case *pb.WrapperResponse_Quote:
		return resp.GetQuote().GetText(), nil
	case *pb.WrapperResponse_Challenge:
		pow := w.pow.NewPowShield(int(resp.GetChallenge().GetDifficulty()))
		pow.SetChallenge(resp.GetChallenge().GetData())

		req := &pb.WrapperRequest{
			Request: &pb.WrapperRequest_Solution{
				Solution: &pb.Solution{
					Nonce: pow.GetSolution(),
				},
			},
		}

		reqBuf, err := proto.Marshal(req)
		if err != nil {
			return "", fmt.Errorf("%s: failed to marshall quote request: %w", op, err)
		}

		_, err = w.conn.Write(tcp.PackMessage(reqBuf))
		if err != nil {
			return "", fmt.Errorf("%s: failed to write quote request: %w", op, err)
		}

		_, err = tcp.UnpackMessage(w.rd)
		if err != nil {
			return "", fmt.Errorf("%s: failed to read response: %w", op, err)
		}

		return "", nil
	case *pb.WrapperResponse_Error:
		return "", fmt.Errorf("%s: server returned error with code %s", op, resp.GetError().GetCode())
	}

	return resp.GetQuote().GetText(), nil
}
