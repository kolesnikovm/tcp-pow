package test

import (
	"context"
	"testing"

	"github.com/kolesnikovm/tcp-pow/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetQuote(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	suite, cleanup := NewSuite(t)
	defer cleanup()

	quote := domain.Quote{
		Text: "test quote",
	}
	suite.wisdomRepo.EXPECT().GetQuote(mock.Anything).Return(quote, nil)

	quoteResp, err := suite.client.GetQuote(ctx)
	require.NoError(t, err)

	require.Equal(t, "", quoteResp)

	quoteResp, err = suite.client.GetQuote(ctx)
	require.NoError(t, err)

	require.Equal(t, quote.Text, quoteResp)
}
