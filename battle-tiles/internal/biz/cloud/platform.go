package cloud

import (
	cloudModel "battle-tiles/internal/dal/model/cloud"
	"battle-tiles/internal/dal/repo/cloud"

	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type PlatformUsecase struct {
	repo cloud.BasePlatformRepo
	log  log.Logger
}

func NewPlatformUsecase(repo cloud.BasePlatformRepo, logger log.Logger) *PlatformUsecase {
	return &PlatformUsecase{repo: repo, log: logger}
}

// ListAll 返回全部平台
func (uc *PlatformUsecase) ListAll(ctx context.Context) ([]*cloudModel.BasePlatform, error) {
	return uc.repo.ListPlatform()
}
