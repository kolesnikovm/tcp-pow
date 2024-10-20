package tcp

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kolesnikovm/tcp-pow/internal/configs"
	"github.com/kolesnikovm/tcp-pow/internal/pow"
	"github.com/kolesnikovm/tcp-pow/internal/service"
	pb "github.com/kolesnikovm/tcp-pow/pkg/proto/gen"
	"github.com/panjf2000/gnet/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	mu sync.Mutex

	id uuid.UUID

	config        *configs.Config
	conn          gnet.Conn
	wisdomService service.Wisdom
	pow           *pow.PowShield

	isVerified bool
	requests   int64
}

func (c *Connection) handleRequest(ctx context.Context, reqBuf []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := &pb.WrapperRequest{}

	c.requests++

	err := proto.Unmarshal(reqBuf, req)
	if err != nil {
		c.handleError(pb.Error_ERR_BAD_RQUEST)
		return
	}

	var wrapperResp *pb.WrapperResponse

	switch {
	case c.requests <= c.config.MaxRequests && c.isVerified && req.GetQuoteRequest() != nil:
		wrapperResp, err = c.handleQuoteRequest(ctx, req)
	case !c.isVerified && req.GetSolution() != nil:
		wrapperResp, err = c.handleSolution(ctx, req)
	case !c.isVerified && req.GetSolution() == nil || c.requests > c.config.MaxRequests:
		wrapperResp, err = c.handleChallenge()
	default:
		log.Debug().
			Str("id", c.id.String()).
			Bool("isVerified", c.isVerified).
			Int64("requests", c.requests).
			Type("request", req.GetRequest()).
			Msgf("unexpected request")

		c.handleError(pb.Error_ERR_BAD_RQUEST)
		return
	}

	if err != nil {
		c.handleError(pb.Error_ERR_INTERNAL)
		return
	}

	respBuf, err := proto.Marshal(wrapperResp)
	if err != nil {
		c.handleError(pb.Error_ERR_INTERNAL)
		return
	}

	c.conn.AsyncWrite(PackMessage(respBuf), nil)
}

func (c *Connection) handleQuoteRequest(ctx context.Context, req *pb.WrapperRequest) (*pb.WrapperResponse, error) {
	resp, err := c.getQuote(ctx, req.GetQuoteRequest())
	if err != nil {
		return nil, err
	}

	return &pb.WrapperResponse{
		Response: &pb.WrapperResponse_Quote{
			Quote: resp,
		},
	}, nil
}

func (c *Connection) handleSolution(ctx context.Context, req *pb.WrapperRequest) (*pb.WrapperResponse, error) {
	isValid := c.pow.VerifySolution(req.GetSolution().GetNonce())
	if !isValid {
		log.Info().Msg("invalid solution")
		c.handleError(pb.Error_ERR_INVALID_SOLUTION)
		return nil, nil
	}

	c.isVerified = true
	c.requests = 1

	return c.handleQuoteRequest(ctx, req)
}

func (c *Connection) handleChallenge() (*pb.WrapperResponse, error) {
	c.isVerified = false

	challenge := uuid.New()
	c.pow.SetChallenge(challenge[:])

	return &pb.WrapperResponse{
		Response: &pb.WrapperResponse_Challenge{
			Challenge: &pb.Challenge{
				Data:       challenge[:],
				Difficulty: int32(c.pow.GetDifficulty()),
			},
		},
	}, nil
}

func (c *Connection) handleError(errCode pb.Error_Code) {
	resp := &pb.Error{
		Code: errCode,
	}

	buf, err := proto.Marshal(resp)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshall error response")
		return
	}

	c.conn.AsyncWrite(PackMessage(buf), nil)
}

func PackMessage(msg []byte) []byte {
	length := len(msg)
	lengthPrefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthPrefix, uint32(length))

	return append(lengthPrefix, msg...)
}

func UnpackMessage(rd *bufio.Reader) ([]byte, error) {
	const op = "tcp.UnpackMessage"

	lengthPrefix := make([]byte, 4)

	_, err := rd.Read(lengthPrefix)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read response length: %w", op, err)
	}

	buf := make([]byte, binary.LittleEndian.Uint32(lengthPrefix))

	_, err = rd.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read response: %w", op, err)
	}

	return buf, nil
}
