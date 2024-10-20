package wisdomservice

import (
	"context"

	"github.com/kolesnikovm/tcp-pow/internal/domain"
)

func (w *WisdomService) GetQuote(ctx context.Context) (domain.Quote, error) {
	quote, err := w.wisdomRepo.GetQuote(ctx)
	if err != nil {
		return quote, err
	}

	return quote, nil
}
