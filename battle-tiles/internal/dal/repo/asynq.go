package repo

import (
	"battle-tiles/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type AsyNQRepo interface {
}

type asyNQRepo struct {
	data *infra.Data
	log  *log.Helper
}

// NewAsyNQRepo .
func NewAsyNQRepo(data *infra.Data, logger log.Logger) AsyNQRepo {

	return &asyNQRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/asyNQ")),
	}
}
