package biz

import (
	"battle-tiles/internal/dal/repo"
	"github.com/go-kratos/kratos/v2/log"
)

// AsyNQUseCase is a asyNQ usecase.
type AsyNQUseCase struct {
	repo repo.AsyNQRepo
	log  *log.Helper
}

// NewAsyNQUseCase new a asyNQ usecase.
func NewAsyNQUseCase(repo repo.AsyNQRepo, logger log.Logger) *AsyNQUseCase {
	return &AsyNQUseCase{repo: repo, log: log.NewHelper(logger)}
}
