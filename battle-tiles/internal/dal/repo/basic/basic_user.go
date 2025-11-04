package basic

import (
	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/infra"
	repox "battle-tiles/pkg/plugin/gormx/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type BasicUserRepo interface {
	repox.CORMImpl[basicModel.BasicUser]
}

type basicUserRepo struct {
	repox.CORMImpl[basicModel.BasicUser]
	data *infra.Data
	log  *log.Helper
}

// NewBasicUserRepo .
func NewBasicUseRepo(data *infra.Data, logger log.Logger) BasicUserRepo {

	return &basicUserRepo{
		CORMImpl: repox.NewCORMImplRepo[basicModel.BasicUser](data),
		data:     data,
		log:      log.NewHelper(log.With(logger, "module", "repo/basicUser")),
	}
}
