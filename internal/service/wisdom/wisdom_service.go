package wisdomservice

import "github.com/kolesnikovm/tcp-pow/internal/repository"

type WisdomService struct {
	wisdomRepo repository.Wisdom
}

func NewWisdomService(wisdomRepo repository.Wisdom) *WisdomService {
	return &WisdomService{
		wisdomRepo: wisdomRepo,
	}
}
