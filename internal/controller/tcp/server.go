package tcp

import (
	"bufio"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kolesnikovm/tcp-pow/internal/configs"
	"github.com/kolesnikovm/tcp-pow/internal/metrics"
	"github.com/kolesnikovm/tcp-pow/internal/pow"
	"github.com/kolesnikovm/tcp-pow/internal/service"
	"github.com/panjf2000/gnet/v2"
	"github.com/rs/zerolog/log"
)

type TCPServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool

	wisdomService service.Wisdom
	config        *configs.Config
	powFactory    *pow.PowShieldFatory
}

func NewTCPServer(config *configs.Config, wisdomService service.Wisdom, powFactory *pow.PowShieldFatory) *TCPServer {
	return &TCPServer{
		addr:          config.ListenAddress,
		multicore:     true,
		wisdomService: wisdomService,
		config:        config,
		powFactory:    powFactory,
	}
}

func (s *TCPServer) Run() error {
	return gnet.Run(s, "tcp://"+s.addr, gnet.WithMulticore(s.multicore))
}

func (s *TCPServer) Stop(ctx context.Context) error {
	return s.eng.Stop(ctx)
}

func (s *TCPServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	log.Info().Msgf("running server on tcp://%s with multi-core=%t", s.addr, s.multicore)

	s.eng = eng

	return
}

func (s *TCPServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	metrics.ClientConnections.Inc()

	clientConn := &Connection{
		id:            uuid.New(),
		conn:          c,
		config:        s.config,
		wisdomService: s.wisdomService,
		pow:           s.powFactory.NewPowShield(s.config.PowDifficulty),
	}
	c.SetContext(clientConn)

	return
}

func (s *TCPServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		log.Info().Msgf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}

	metrics.ClientConnections.Dec()

	return
}

func (s *TCPServer) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	metrics.RequestsTotal.Inc()

	clientConn := conn.Context().(*Connection)

	buf, err := UnpackMessage(bufio.NewReader(conn))
	if err != nil {
		log.Error().Err(err).Msg("failed to read from connection, closing")
		return gnet.Close
	}

	newBuf := make([]byte, len(buf))
	copy(newBuf, buf)

	ctx := context.TODO()

	go func() {
		start := time.Now()
		clientConn.handleRequest(ctx, newBuf)
		metrics.RequestsDuration.Observe(float64(time.Since(start).Nanoseconds()))
	}()

	return
}
