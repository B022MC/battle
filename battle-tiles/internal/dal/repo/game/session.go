package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"
	"time"

	"gorm.io/gorm"
)

type SessionRepo interface {
	Insert(ctx context.Context, s *model.GameSession) error
	CloseActive(ctx context.Context, gameAccountID int32) error
	CloseActiveBy(ctx context.Context, gameAccountID int32, userID int32, houseGID int32) error
	List(ctx context.Context, userID int32, houseGID *int, state *string, limit, offset int) ([]*model.GameSession, error)
	EnsureOnlineByHouse(ctx context.Context, houseGID int32, ctrlAccountID int32) error
	SetOfflineByHouse(ctx context.Context, houseGID int32) error
	// UpsertOnlineByHouse: 若存在任意该店铺记录，则更新其最新一条为 online；否则插入一条
	UpsertOnlineByHouse(ctx context.Context, ctrlAccountID int32, userID int32, houseGID int32) error
	// UpsertErrorByHouse: 若存在任意该店铺记录，则更新其最新一条为 error；否则插入一条
	UpsertErrorByHouse(ctx context.Context, ctrlAccountID int32, userID int32, houseGID int32, errorMsg string) error
}

type sessionRepo struct{ data *infra.Data }

func NewSessionRepo(data *infra.Data) SessionRepo { return &sessionRepo{data: data} }

func (r *sessionRepo) Insert(ctx context.Context, s *model.GameSession) error {
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = time.Now()
	}
	return r.data.GetDBWithContext(ctx).Create(s).Error
}

func (r *sessionRepo) CloseActive(ctx context.Context, gameAccountID int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Exec(`
    UPDATE game_session
       SET end_at = now(), updated_at = now()
     WHERE id = (
       SELECT id FROM game_session
        WHERE game_ctrl_account_id = ?  AND end_at IS NULL
        ORDER BY created_at DESC
        LIMIT 1
     )`, gameAccountID).Error
}

// CloseActiveBy 关闭指定账号+用户+店铺下未结束的在线记录（幂等）
func (r *sessionRepo) CloseActiveBy(ctx context.Context, gameAccountID int32, userID int32, houseGID int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Exec(`
    UPDATE game_session
       SET end_at = now(), updated_at = now()
     WHERE game_ctrl_account_id = ? AND user_id = ? AND house_gid = ? AND end_at IS NULL
    `, gameAccountID, userID, houseGID).Error
}

func (r *sessionRepo) List(ctx context.Context, userID int32, houseGID *int, state *string, limit, offset int) ([]*model.GameSession, error) {
	db := r.data.GetDBWithContext(ctx).Where("user_id = ? ", userID)
	if houseGID != nil {
		db = db.Where("house_gid = ?", *houseGID)
	}
	if state != nil {
		db = db.Where("state = ?", *state)
	}
	var out []*model.GameSession
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&out).Error
	return out, err
}

// EnsureOnlineByHouse 确保指定店铺存在且仅存在一条在线记录；无则插入，有则置为在线
func (r *sessionRepo) EnsureOnlineByHouse(ctx context.Context, houseGID int32, ctrlAccountID int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		// 仅将未结束(end_at IS NULL)的记录置为 online，并清空 end_at（容错）
		if err := tx.Exec(`
            UPDATE game_session
               SET state = 'online', end_at = NULL, updated_at = now()
             WHERE house_gid = ? AND end_at IS NULL
        `, houseGID).Error; err != nil {
			return err
		}
		// 若无任何未结束记录，则插入一条新的在线记录（user_id=0 占位）
		if err := tx.Exec(`
            INSERT INTO game_session (game_ctrl_account_id, user_id, house_gid, state, device_ip, error_msg, created_at, updated_at)
            SELECT ?, 0, ?, 'online', '', '', now(), now()
            WHERE NOT EXISTS (SELECT 1 FROM game_session WHERE house_gid = ? AND end_at IS NULL)
        `, ctrlAccountID, houseGID, houseGID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *sessionRepo) SetOfflineByHouse(ctx context.Context, houseGID int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Exec(`
        UPDATE game_session
           SET end_at = now(), updated_at = now(), state = 'offline'
         WHERE house_gid = ? AND end_at IS NULL
    `, houseGID).Error
}

// UpsertOnlineByHouse 将该店铺最新一条记录更新为 online；若不存在任何记录则插入新记录
func (r *sessionRepo) UpsertOnlineByHouse(ctx context.Context, ctrlAccountID int32, userID int32, houseGID int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Exec(`
        WITH updated AS (
            UPDATE game_session
               SET game_ctrl_account_id = ?, user_id = ?, state = 'online', end_at = NULL, error_msg = '', updated_at = now()
             WHERE id = (
                 SELECT id FROM game_session WHERE house_gid = ? ORDER BY created_at DESC LIMIT 1
             )
            RETURNING id
        )
        INSERT INTO game_session (game_ctrl_account_id, user_id, house_gid, state, device_ip, error_msg, created_at, updated_at)
        SELECT ?, ?, ?, 'online', '', '', now(), now()
         WHERE NOT EXISTS (SELECT 1 FROM updated);
    `, ctrlAccountID, userID, houseGID, ctrlAccountID, userID, houseGID).Error
}

// UpsertErrorByHouse 将该店铺最新一条记录更新为 error；若不存在任何记录则插入新记录
func (r *sessionRepo) UpsertErrorByHouse(ctx context.Context, ctrlAccountID int32, userID int32, houseGID int32, errorMsg string) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Exec(`
        WITH updated AS (
            UPDATE game_session
               SET game_ctrl_account_id = ?, user_id = ?, state = 'error', error_msg = ?, end_at = NULL, updated_at = now()
             WHERE id = (
                 SELECT id FROM game_session WHERE house_gid = ? ORDER BY created_at DESC LIMIT 1
             )
            RETURNING id
        )
        INSERT INTO game_session (game_ctrl_account_id, user_id, house_gid, state, device_ip, error_msg, created_at, updated_at)
        SELECT ?, ?, ?, 'error', '', ?, now(), now()
         WHERE NOT EXISTS (SELECT 1 FROM updated);
    `, ctrlAccountID, userID, errorMsg, houseGID, ctrlAccountID, userID, houseGID, errorMsg).Error
}
