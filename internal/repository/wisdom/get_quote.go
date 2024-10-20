package wisdomrepo

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/kolesnikovm/tcp-pow/internal/domain"
)

func (w *WisdomRepo) GetQuote(ctx context.Context) (domain.Quote, error) {
	const op = "wisdomrepo.GetQuote"

	quoteData := w.quotes[rand.Intn(len(w.quotes))]

	if len(quoteData) < 2 {
		return domain.Quote{}, fmt.Errorf("%s: failed to retrive quote", op)
	}

	return domain.Quote{
		Text: quoteData[1],
	}, nil
}
