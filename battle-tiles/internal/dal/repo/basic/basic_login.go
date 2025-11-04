package basic

import (
	"battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/infra"
	"battle-tiles/pkg/plugin/gormx/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type BasicLoginRepo interface {
	FindByUsername(ctx context.Context, username string) (*basic.BasicUser, error)
	FindByPhone(ctx context.Context, phone string) (*basic.BasicUser, error)
	Create(ctx context.Context, user *basic.BasicUser) (int32, error)
	UpdateLastLoginAt(ctx context.Context, id int32) error
	FindByOpenID(ctx context.Context, openID string) (*basic.BasicUser, error)
	UpdateByPKWithMap(ctx context.Context, pk interface{}, updateData map[string]interface{}) (int64, error)
	SendSMSCode(ctx context.Context, phone, code string) error
	VerifySMSCode(ctx context.Context, phone, inputCode string) bool
}

type basicLoginRepo struct {
	repo.CORMImpl[basic.BasicUser]
	data *infra.Data
	log  *log.Helper
}

func NewBasicLoginRepo(data *infra.Data, logger log.Logger) BasicLoginRepo {
	return &basicLoginRepo{
		CORMImpl: repo.NewCORMImplRepo[basic.BasicUser](data),
		data:     data,
		log:      log.NewHelper(log.With(logger, "module", "repo/basicLogin")),
	}
}

func (r *basicLoginRepo) FindByUsername(ctx context.Context, username string) (*basic.BasicUser, error) {
	var user basic.BasicUser
	if err := r.data.GetDBWithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *basicLoginRepo) FindByPhone(ctx context.Context, phone string) (*basic.BasicUser, error) {
	var user basic.BasicUser
	if err := r.data.GetDBWithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *basicLoginRepo) Create(ctx context.Context, user *basic.BasicUser) (int32, error) {
	if err := r.data.GetDBWithContext(ctx).Create(user).Error; err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *basicLoginRepo) UpdateLastLoginAt(ctx context.Context, id int32) error {
	return r.data.GetDBWithContext(ctx).Model(&basic.BasicUser{}).Where("id = ?", id).
		Update("last_login_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
func (r *basicLoginRepo) FindByOpenID(ctx context.Context, openID string) (*basic.BasicUser, error) {
	var user *basic.BasicUser
	if err := r.data.GetDBWithContext(ctx).Where("open_id = ?", openID).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *basicLoginRepo) UpdateByPKWithMap(ctx context.Context, pk interface{}, updateData map[string]interface{}) (int64, error) {
	return r.data.GetDBWithContext(ctx).
		Model(&basic.BasicUser{}).
		Where("id = ?", pk).
		Updates(updateData).
		RowsAffected, nil
}

// 存储验证码（过期时间 5 分钟）
func (r *basicLoginRepo) SendSMSCode(ctx context.Context, phone, code string) error {
	key := fmt.Sprintf("login:code:%s", phone)
	return r.data.RDB.WithContext(ctx).Set(ctx, key, code, 5*time.Minute).Err()
}

// 校验验证码
func (r *basicLoginRepo) VerifySMSCode(ctx context.Context, phone, inputCode string) bool {
	failKey := fmt.Sprintf("login:fail:%s", phone)
	failCount, _ := r.data.RDB.WithContext(ctx).Incr(ctx, failKey).Result()
	if failCount >= 5 {
		return false // 禁止尝试
	}
	r.data.RDB.WithContext(ctx).Expire(ctx, failKey, 10*time.Minute)
	key := fmt.Sprintf("login:code:%s", phone)
	val, err := r.data.RDB.WithContext(ctx).Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return val == inputCode
}
