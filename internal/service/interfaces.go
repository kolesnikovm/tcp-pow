package service

import (
	"context"

	"github.com/kolesnikovm/tcp-pow/internal/domain"
)

type Wisdom interface {
	GetQuote(context.Context) (domain.Quote, error)
}
