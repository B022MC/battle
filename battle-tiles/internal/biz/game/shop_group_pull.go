package game

import (
	"context"
	"fmt"

	model "battle-tiles/internal/dal/model/game"
)

// PullMembersToGroup 将游戏成员拉入指定圈子
// 支持：
// 1. 无圈子成员 -> 拉入圈子
// 2. 其他圈子成员 -> 转移到新圈子
//
// memberIDs: game_member.id 列表
func (uc *ShopGroupUseCase) PullMembersToGroup(ctx context.Context, groupID int32, adminUserID int32, memberIDs []int32) error {
	// 1. 验证圈子权限
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("圈子不存在: %w", err)
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	// 2. 批量处理每个成员
	for _, memberID := range memberIDs {
		// 2.1 查询成员信息
		member, err := uc.gameMemberRepo.GetByID(ctx, memberID)
		if err != nil {
			uc.log.Warnf("查询成员 %d 失败: %v", memberID, err)
			continue
		}

		// 2.2 验证成员是否在同一店铺
		if member.HouseGID != group.HouseGID {
			uc.log.Warnf("成员 %d 不在店铺 %d 中，跳过", memberID, group.HouseGID)
			continue
		}

		// 2.3 检查是否已在目标圈子
		if member.GroupID != nil && *member.GroupID == groupID {
			uc.log.Infof("成员 %d 已在圈子 %d 中，跳过", memberID, groupID)
			continue
		}

		// 2.4 更新 game_account_group（主表）
		// 使用 game_player_id 而非 game_account_id
		gamePlayerID := fmt.Sprintf("%d", member.GameID)
		if err := uc.accountGroupRepo.UpdateGroupByGamePlayerAndHouse(
			ctx, gamePlayerID, member.HouseGID, groupID, group.GroupName,
		); err != nil {
			uc.log.Errorf("更新成员 %d 的 game_account_group 失败: %v", memberID, err)
			continue
		}

		// 2.5 更新 game_member（从表）
		if err := uc.gameMemberRepo.UpdateGroup(
			ctx, member.HouseGID, member.GameID, groupID, group.GroupName,
		); err != nil {
			uc.log.Errorf("更新成员 %d 的 game_member 失败: %v", memberID, err)
			continue
		}

		uc.log.Infof("成功将成员 %d (game_id=%d) 拉入圈子 %d", memberID, member.GameID, groupID)
	}

	return nil
}

// PullGameAccountsToGroup 通过游戏账号ID将成员拉入圈子
// 现在使用 game_account 的 game_player_id 来操作
func (uc *ShopGroupUseCase) PullGameAccountsToGroup(ctx context.Context, groupID int32, adminUserID int32, gameAccountIDs []int32) error {
	// 1. 验证圈子权限
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("圈子不存在: %w", err)
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	// 2. 批量处理每个游戏账号
	for _, gameAccountID := range gameAccountIDs {
		// 2.1 查询游戏账号，获取 game_player_id
		gameAccount, err := uc.accountRepo.GetByID(ctx, gameAccountID)
		if err != nil || gameAccount == nil || gameAccount.GamePlayerID == "" {
			uc.log.Warnf("游戏账号 %d 不存在或缺少 game_player_id，跳过", gameAccountID)
			continue
		}

		// 2.2 更新 game_account_group（主表）使用 game_player_id
		if err := uc.accountGroupRepo.UpdateGroupByGamePlayerAndHouse(
			ctx, gameAccount.GamePlayerID, group.HouseGID, groupID, group.GroupName,
		); err != nil {
			uc.log.Errorf("更新游戏账号 %d 的 game_account_group 失败: %v", gameAccountID, err)
			continue
		}

		// 2.3 更新或创建 game_member（从表）- 使用 game_player_id 转成 int32
		// game_member.game_id 存储的是游戏内玩家ID（int32）
		var gamePlayerIDInt int32
		fmt.Sscanf(gameAccount.GamePlayerID, "%d", &gamePlayerIDInt)

		if err := uc.gameMemberRepo.UpdateGroup(
			ctx, group.HouseGID, gamePlayerIDInt, groupID, group.GroupName,
		); err != nil {
			// 如果是因为不存在而失败，尝试创建
			newMember := &model.GameMember{
				HouseGID:  group.HouseGID,
				GameID:    gamePlayerIDInt,
				GameName:  gameAccount.Nickname,
				GroupID:   &groupID,
				GroupName: group.GroupName,
				Balance:   0,
				Credit:    0,
				Forbid:    false,
			}
			if createErr := uc.gameMemberRepo.Create(ctx, newMember); createErr != nil {
				uc.log.Errorf("创建成员记录失败 (game_id=%d): %v", gamePlayerIDInt, createErr)
				continue
			}
		}

		uc.log.Infof("成功将游戏账号 %d (game_player_id=%s) 拉入圈子 %d", gameAccountID, gameAccount.GamePlayerID, groupID)
	}

	return nil
}

// ListHouseMembers 查询店铺下的所有游戏成员（可按是否有圈子筛选）
func (uc *ShopGroupUseCase) ListHouseMembers(ctx context.Context, houseGID int32, onlyWithoutGroup bool, page, size int32) ([]*model.GameMember, int64, error) {
	if onlyWithoutGroup {
		members, err := uc.gameMemberRepo.ListByHouseWithoutGroup(ctx, houseGID)
		return members, int64(len(members)), err
	}
	return uc.gameMemberRepo.ListByHouse(ctx, houseGID, page, size)
}
