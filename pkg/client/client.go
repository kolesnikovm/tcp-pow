package wisdomclient

import (
	"bufio"
	"fmt"
	"net"

	"github.com/kolesnikovm/tcp-pow/internal/pow"
)

type WisdomClient struct {
	conn net.Conn
	rd   *bufio.Reader
	pow  *pow.PowShieldFatory
}

func NewWisdomClient(addr string) (*WisdomClient, error) {
	const op = "wisdomclient.NewWisdomClient"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to server '%s': %w", op, addr, err)
	}

	rd := bufio.NewReader(conn)

	powShield := pow.NewPowShieldFactory()

	return &WisdomClient{
		conn: conn,
		rd:   rd,
		pow:  powShield,
	}, nil
}
