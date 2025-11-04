package cloud

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"battle-tiles/internal/consts"
	cloudModel "battle-tiles/internal/dal/model/cloud"
	"battle-tiles/internal/infra"
	pdb "battle-tiles/pkg/plugin/dbx"
)

type BasePlatformRepo interface {
	ListPlatform() ([]*cloudModel.BasePlatform, error)
	UpsertPlatform(platform cloudModel.BasePlatform) (*cloudModel.BasePlatform, error)
	GetPlatformInfo(ctx context.Context) (*cloudModel.BasePlatform, error)
}

type basePlatformRepo struct {
	data *infra.Data
	log  *log.Helper
}

func (rp *basePlatformRepo) ListPlatform() ([]*cloudModel.BasePlatform, error) {
	db, ok := rp.data.DBMap[consts.CloudPlatformDB]
	if !ok {
		return nil, errors.New("cloud cloud db not found")
	}

	var platformList []*cloudModel.BasePlatform
	if err := db.Find(&platformList).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return platformList, nil
}

func (rp *basePlatformRepo) GetPlatformInfo(ctx context.Context) (*cloudModel.BasePlatform, error) {
	dbName := pdb.GetDBKeyFromCtx(ctx)
	db, ok := rp.data.DBMap[consts.CloudPlatformDB]
	if !ok {
		return nil, errors.New("cloud cloud db not found")
	}

	var platformInfo cloudModel.BasePlatform
	if err := db.Table(cloudModel.TableNameBasePlatform).Where("db_name = ?", dbName).First(&platformInfo).Error; err != nil {
		return nil, err
	}

	return &platformInfo, nil
}

func (rp *basePlatformRepo) UpsertPlatform(platform cloudModel.BasePlatform) (*cloudModel.BasePlatform, error) {
	db, ok := rp.data.DBMap[consts.CloudPlatformDB]
	if !ok {
		return nil, errors.New("cloud cloud db not found")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var existing cloudModel.BasePlatform
		err := tx.Where("cloud = ?", platform.Platform).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(&platform).Error
		}

		return tx.Table(cloudModel.TableNameBasePlatform).Where("cloud = ?", platform.Platform).Updates(map[string]interface{}{
			"name": platform.Name,
		}).Error
	})

	if err != nil {
		return nil, err
	}

	return &platform, nil
}

func NewBasePlatformRepo(data *infra.Data, logger log.Logger) (BasePlatformRepo, error) {
	return &basePlatformRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/cloud")),
	}, nil
}
