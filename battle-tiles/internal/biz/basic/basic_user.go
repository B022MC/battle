package basic

import (
	basicModel "battle-tiles/internal/dal/model/basic"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/req"
	basicVo "battle-tiles/internal/dal/vo/basic"
	"battle-tiles/pkg/utils"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

var (
	ErrBasicUserExisted = errors.New("名称重复")
)

// BasicUserUseCase is a BasicUser usecase.
type BasicUserUseCase struct {
	repo basicRepo.BasicUserRepo
	log  *log.Helper
}

// NewBasicUserUseCase new a BasicUser usecase.
func NewBasicUserUseCase(
	repo basicRepo.BasicUserRepo,
	logger log.Logger,
) *BasicUserUseCase {
	return &BasicUserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/basic_user")),
	}
}

func (uc *BasicUserUseCase) beforeCheck(ctx context.Context, req *req.AddBasicUserReq) error {
	old, err := uc.repo.SelectOne(ctx, &basicModel.BasicUser{
		Username: req.UserName,
	})
	if err != nil {
		return err
	}
	if old != nil {
		return errors.Wrapf(ErrBasicUserExisted, "请检查[%s]", req.UserName)
	}
	return nil
}

func (uc *BasicUserUseCase) Create(ctx context.Context, req *req.AddBasicUserReq) (*basicModel.BasicUser, error) {
	if err := uc.beforeCheck(ctx, req); err != nil {
		return nil, err
	}
	basicUserModel := basicModel.BasicUser{}
	err := uc.repo.Insert(ctx, &basicUserModel)
	return &basicUserModel, err
}

func (uc *BasicUserUseCase) Delete(ctx context.Context, id interface{}) error {
	_, err := uc.repo.DeleteByPK(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *BasicUserUseCase) Get(ctx context.Context, id interface{}) (*basicVo.BasicUserVo, error) {
	one, err := uc.repo.SelectOneByPK(ctx, id)
	if err != nil {
		return nil, err
	}
	return &basicVo.BasicUserVo{
		BasicUser: one,
	}, nil
}

func (uc *BasicUserUseCase) ListByOption(ctx context.Context, req *req.ListBasicUserReq) ([]*utils.Option, int64, error) {
	query := uc.repo.DB(ctx).Model(&basicModel.BasicUser{})
	query = query.Select("id as value, name as label")
	var res []*utils.Option
	total, err := uc.repo.ListPageByOption(ctx, query, req.PageParam, &res)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		total = int64(len(res))
	}
	return res, total, nil
}

func (uc *BasicUserUseCase) List(ctx context.Context, req *req.ListBasicUserReq) ([]*basicModel.BasicUser, int64, error) {
	query := uc.repo.DB(ctx).Model(&basicModel.BasicUser{})

	return uc.repo.ListPage(ctx, query, req.PageParam)
}

func (uc *BasicUserUseCase) Update(ctx context.Context, id interface{}, field map[string]interface{}) (*basicVo.BasicUserVo, error) {
	row, err := uc.repo.UpdateByPKWithMap(ctx, id, field)
	if err != nil {
		return nil, err
	}
	if row == 0 {
		return nil, ErrBasicUserExisted
	}
	return uc.Get(ctx, id)
}
