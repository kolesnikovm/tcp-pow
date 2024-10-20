package test

import (
	"context"
	"testing"
	"time"

	"github.com/kolesnikovm/tcp-pow/internal/configs"
	"github.com/kolesnikovm/tcp-pow/internal/controller/tcp"
	"github.com/kolesnikovm/tcp-pow/internal/pow"
	repomocks "github.com/kolesnikovm/tcp-pow/internal/repository/mocks"
	wisdomservice "github.com/kolesnikovm/tcp-pow/internal/service/wisdom"
	wisdomclient "github.com/kolesnikovm/tcp-pow/pkg/client"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	t *testing.T

	config     *configs.Config
	client     *wisdomclient.WisdomClient
	wisdomRepo *repomocks.MockWisdom
}

func NewSuite(t *testing.T) (suite *Suite, cleanup func()) {
	config, err := configs.NewConfig("")
	require.NoError(t, err)

	wisdomRepo := repomocks.NewMockWisdom(t)

	wisdomService := wisdomservice.NewWisdomService(wisdomRepo)

	powFactory := pow.NewPowShieldFactory()

	tcpServer := tcp.NewTCPServer(config, wisdomService, powFactory)

	go func() {
		err := tcpServer.Run()
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	client, err := wisdomclient.NewWisdomClient(config.ServerAddress)
	require.NoError(t, err)

	suite = &Suite{
		t:          t,
		config:     config,
		client:     client,
		wisdomRepo: wisdomRepo,
	}

	cleanup = func() {
		err := tcpServer.Stop(context.Background())
		require.NoError(t, err)
	}

	return suite, cleanup
}
