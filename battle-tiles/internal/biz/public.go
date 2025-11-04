package biz

import (
	"battle-tiles/internal/dal/repo"
	"github.com/go-kratos/kratos/v2/log"
)

// PublicUseCase is a Public usecase.
type PublicUseCase struct {
	repo repo.PublicRepo
	log  *log.Helper
}

// NewPublicUseCase new a Public usecase.
func NewPublicUseCase(repo repo.PublicRepo, logger log.Logger) *PublicUseCase {
	return &PublicUseCase{repo: repo, log: log.NewHelper(logger)}
}
