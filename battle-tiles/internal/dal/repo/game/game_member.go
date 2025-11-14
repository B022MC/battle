package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type GameMemberRepo interface {
	// 鏍规嵁 game_id 鍜?group_id 鏌ヨ鎴愬憳
	GetByGameIDAndGroup(ctx context.Context, houseGID int32, gameID int32, groupID *int32) (*model.GameMember, error)

	// 鏌ヨ鎴愬憳鍦ㄦ墍鏈夊湀瀛愮殑璁板綍
	ListByGameID(ctx context.Context, houseGID int32, gameID int32) ([]*model.GameMember, error)

	// 鏌ヨ鍦堝瓙鐨勬墍鏈夋垚鍛?
	ListByGroup(ctx context.Context, houseGID int32, groupID int32, page, size int32) ([]*model.GameMember, int64, error)

	// 鏍规嵁 member_id 鏌ヨ鎴愬憳
	GetByID(ctx context.Context, memberID int32) (*model.GameMember, error)
}

type gameMemberRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewGameMemberRepo(data *infra.Data, logger log.Logger) GameMemberRepo {
	return &gameMemberRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_member")),
	}
}

func (r *gameMemberRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

// GetByGameIDAndGroup 鏍规嵁 game_id 鍜?group_id 鏌ヨ鎴愬憳
func (r *gameMemberRepo) GetByGameIDAndGroup(ctx context.Context, houseGID int32, gameID int32, groupID *int32) (*model.GameMember, error) {
	var member model.GameMember
	db := r.db(ctx).Where("house_gid = ? AND game_id = ?", houseGID, gameID)
	
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	} else {
		db = db.Where("group_id IS NULL")
	}
	
	if err := db.First(&member).Error; err != nil {
		return nil, err
	}
	
	return &member, nil
}

// ListByGameID 鏌ヨ鎴愬憳鍦ㄦ墍鏈夊湀瀛愮殑璁板綍
func (r *gameMemberRepo) ListByGameID(ctx context.Context, houseGID int32, gameID int32) ([]*model.GameMember, error) {
	var members []*model.GameMember
	if err := r.db(ctx).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Order("group_id").
		Find(&members).Error; err != nil {
		return nil, err
	}
	
	return members, nil
}

// ListByGroup 鏌ヨ鍦堝瓙鐨勬墍鏈夋垚鍛?
func (r *gameMemberRepo) ListByGroup(ctx context.Context, houseGID int32, groupID int32, page, size int32) ([]*model.GameMember, int64, error) {
	db := r.db(ctx).Model(&model.GameMember{}).
		Where("house_gid = ? AND group_id = ?", houseGID, groupID)
	
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	offset := (page - 1) * size
	
	var members []*model.GameMember
	if err := db.
		Order("balance DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&members).Error; err != nil {
		return nil, 0, err
	}
	
	return members, total, nil
}

// GetByID 鏍规嵁 member_id 鏌ヨ鎴愬憳
func (r *gameMemberRepo) GetByID(ctx context.Context, memberID int32) (*model.GameMember, error) {
	var member model.GameMember
	if err := r.db(ctx).Where("id = ?", memberID).First(&member).Error; err != nil {
		return nil, err
	}
	
	return &member, nil
}
