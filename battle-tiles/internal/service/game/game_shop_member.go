// internal/service/game/game_shop_member.go
package game

import (
	biz "battle-tiles/internal/biz/game"
	gameModel "battle-tiles/internal/dal/model/game"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	gameRepo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/dal/req"
	resp "battle-tiles/internal/dal/resp"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type GameShopMemberService struct {
	mgr              plaza.Manager
	rule             *biz.MemberRuleUseCase
	sAdm             gameRepo.GameShopAdminRepo
	users            basicRepo.BasicUserRepo
	apps             gameRepo.UserApplicationRepo
	gameAccount      gameRepo.GameAccountRepo
	gameMember       gameRepo.GameMemberRepo
	groupUC          *biz.ShopGroupUseCase
	groupRepo        gameRepo.ShopGroupRepo
	accountGroupRepo gameRepo.GameAccountGroupRepo // 添加圈子关系 repo
}

func NewGameShopMemberService(mgr plaza.Manager, rule *biz.MemberRuleUseCase, sAdm gameRepo.GameShopAdminRepo, users basicRepo.BasicUserRepo, apps gameRepo.UserApplicationRepo, gameAccount gameRepo.GameAccountRepo, gameMember gameRepo.GameMemberRepo, groupUC *biz.ShopGroupUseCase, groupRepo gameRepo.ShopGroupRepo, accountGroupRepo gameRepo.GameAccountGroupRepo) *GameShopMemberService {
	return &GameShopMemberService{mgr: mgr, rule: rule, sAdm: sAdm, users: users, apps: apps, gameAccount: gameAccount, gameMember: gameMember, groupUC: groupUC, groupRepo: groupRepo, accountGroupRepo: accountGroupRepo}
}

func (s *GameShopMemberService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())
	g.POST("/members/kick", middleware.RequirePerm("shop:member:kick"), s.Kick)
	g.POST("/members/list", middleware.RequirePerm("shop:member:view"), s.List)
	g.POST("/members/logout", middleware.RequirePerm("shop:member:logout"), s.Logout)
	g.POST("/diamond/query", middleware.RequirePerm("shop:member:view"), s.QueryDiamond)
	g.POST("/members/pull", middleware.RequirePerm("shop:member:view"), s.PullMembers)
	g.POST("/members/pull-to-group", middleware.RequirePerm("shop:member:update"), s.PullToGroup)
	g.POST("/members/remove-from-group", middleware.RequirePerm("shop:member:kick"), s.RemoveFromGroup)
	// 平台侧：按圈主返回“我圈子的成员”（基于已通过的入圈申请）
	g.POST("/members/list_platform", middleware.RequirePerm("shop:member:view"), s.ListPlatformMembers)
	// 平台侧：从圈中移除成员（标记该成员的入圈记录为移除）
	g.POST("/members/remove_platform", middleware.RequirePerm("shop:member:kick"), s.RemovePlatformMember)
	// 平台侧：直接将用户拉入圈子（创建已批准的入圈记录）
	g.POST("/members/add_platform", middleware.RequirePerm("shop:member:update"), s.AddToPlatformGroup)
	// 成员规则
	g.POST("/members/rules/vip", middleware.RequirePerm("shop:member:update"), s.SetVIP)
	g.POST("/members/rules/multi", middleware.RequirePerm("shop:member:update"), s.SetMulti)
	g.POST("/members/rules/temp_release", middleware.RequirePerm("shop:member:update"), s.SetTempRelease)
}

// AddToPlatformGroup 直接将用户拉入圈子（平台：创建已批准的入圈记录）
func (s *GameShopMemberService) AddToPlatformGroup(c *gin.Context) {
	var in struct {
		HouseGID     int32  `json:"house_gid" binding:"required"`
		GroupID      *int32 `json:"group_id"` // 非超管必填（=圈主ID）
		MemberUserID int32  `json:"member_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.MemberUserID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid params")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	caller := int32(claims.BaseClaims.UserID)
	// 计算圈主（admin）
	var admin int32
	if in.GroupID != nil && *in.GroupID > 0 {
		admin = *in.GroupID
	} else if caller != 1 { // 非超管默认自己的圈
		admin = caller
	} else {
		response.Fail(c, ecode.ParamsFailed, "group_id required for super admin")
		return
	}
	if s.apps == nil {
		response.Fail(c, ecode.Failed, "repo not ready")
		return
	}

	// 验证: 检查目标用户是否已经在某个圈子中
	existingJoin, err := s.apps.GetUserApprovedJoin(c.Request.Context(), in.HouseGID, in.MemberUserID)
	if err == nil && existingJoin != nil {
		response.Fail(c, ecode.Failed, "用户已经在圈子中")
		return
	}

	// 验证: 检查目标用户的角色
	if s.users == nil {
		response.Fail(c, ecode.Failed, "repo not ready")
		return
	}
	targetUser, err := s.users.SelectOneByPK(c.Request.Context(), in.MemberUserID)
	if err != nil {
		response.Fail(c, ecode.Failed, "用户不存在")
		return
	}
	// 不能拉入店铺管理员或超级管理员
	if targetUser.IsSuperAdmin() {
		response.Fail(c, ecode.Failed, "不能拉入超级管理员")
		return
	}
	if targetUser.IsStoreAdmin() {
		response.Fail(c, ecode.Failed, "不能拉入店铺管理员")
		return
	}

	if err := s.apps.AddApprovedJoin(c.Request.Context(), in.HouseGID, admin, in.MemberUserID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// RemovePlatformMember 将指定成员从指定圈移除（平台：更新入圈申请为 status=3）
func (s *GameShopMemberService) RemovePlatformMember(c *gin.Context) {
	var in struct {
		HouseGID     int32  `json:"house_gid" binding:"required"`
		GroupID      *int32 `json:"group_id"` // 非超管必填（=圈主ID）
		MemberUserID int32  `json:"member_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.MemberUserID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid params")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	caller := int32(claims.BaseClaims.UserID)
	// 计算圈主（admin）
	var admin int32
	if in.GroupID != nil && *in.GroupID > 0 {
		admin = *in.GroupID
	} else if caller != 1 { // 非超管默认自己的圈
		admin = caller
	} else {
		response.Fail(c, ecode.ParamsFailed, "group_id required for super admin")
		return
	}
	if s.apps == nil {
		response.Fail(c, ecode.Failed, "repo not ready")
		return
	}
	if _, err := s.apps.RemoveApprovedJoin(c.Request.Context(), in.HouseGID, admin, in.MemberUserID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// ListPlatformMembers 基于平台库的入圈申请（type=2, status=1）返回圈成员列表
// 超级管理员可指定 admin_user_id；店铺管理员默认使用当前用户为 admin_user_id
func (s *GameShopMemberService) ListPlatformMembers(c *gin.Context) {
	var in struct {
		HouseGID int32  `json:"house_gid" binding:"required"`
		GroupID  *int32 `json:"group_id"`      // 非超管必填：我的圈ID（= 我的 user_id）
		AdminUID *int32 `json:"admin_user_id"` // 可选：兼容旧参数
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	caller := int32(claims.BaseClaims.UserID)
	var admin int32
	if in.GroupID != nil && *in.GroupID > 0 {
		admin = *in.GroupID
	} else if in.AdminUID != nil && *in.AdminUID > 0 {
		admin = *in.AdminUID
	} else if caller != 1 { // 非超管默认查自己的圈
		admin = caller
	} else {
		admin = 0 // 超管且未指定 -> 全量
	}
	if s.apps == nil {
		response.Success(c, resp.ShopMemberListResponse{Items: []resp.ShopMemberListItem{}})
		return
	}
	var list []*gameModel.UserApplication
	if admin > 0 {
		list, err = s.apps.ListApprovedJoinsByAdmin(c.Request.Context(), in.HouseGID, admin)
	} else {
		list, err = s.apps.ListApprovedJoins(c.Request.Context(), in.HouseGID)
	}
	if err != nil || len(list) == 0 {
		response.Success(c, resp.ShopMemberListResponse{Items: []resp.ShopMemberListItem{}})
		return
	}
	// 平台圈ID规则：group_id = admin_user_id（若全量则为 0）
	gid := int(admin)
	// 批量昵称（包括成员和圈主）
	nameMap := map[int32]string{}
	groupNameMap := map[int32]string{} // 圈主ID -> 圈主昵称
	if s.users != nil {
		// 收集所有需要查询的用户ID（成员 + 圈主）
		userIDSet := make(map[int32]struct{})
		for _, a := range list {
			userIDSet[a.Applicant] = struct{}{}
			userIDSet[a.AdminUID] = struct{}{}
		}
		ids := make([]int32, 0, len(userIDSet))
		for uid := range userIDSet {
			ids = append(ids, uid)
		}
		rows, _ := s.users.SelectByPK(c.Request.Context(), ids)
		for _, u := range rows {
			if u != nil {
				nameMap[u.Id] = u.NickName
				groupNameMap[u.Id] = u.NickName
			}
		}
	}
	out := make([]resp.ShopMemberListItem, 0, len(list))
	for _, a := range list {
		item := resp.ShopMemberListItem{
			UserID:     uint32(a.Applicant),
			UserStatus: 0,
			GameID:     0,
			MemberID:   0,
			MemberType: 0, // 普通成员
			NickName:   nameMap[a.Applicant],
			GroupID:    gid,
		}
		// 添加圈子名称（圈主昵称）
		if groupName, ok := groupNameMap[a.AdminUID]; ok {
			item.GroupName = &groupName
		}
		out = append(out, item)
	}
	response.Success(c, resp.ShopMemberListResponse{Items: out})
}

// Kick
// @Summary      踢出成员
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.KickMemberRequest true "house_gid, member_id"
// @Success      200 {object} response.Body
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      409 {object} response.Body "no online session"
// @Router       /shops/members/kick [post]
func (s *GameShopMemberService) Kick(c *gin.Context) {
	var in req.KickMemberRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	actorUID := int(claims.BaseClaims.UserID)

	if err := s.mgr.KickMember(actorUID, in.HouseGID, in.MemberID); err != nil {
		if strings.Contains(err.Error(), "session not found") {
			response.Fail(c, ecode.Failed, "no online session")
			return
		}
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.SuccessWithOK(c)
}

// List
// @Summary      成员列表快照
// @Description  触发 GetGroupMembers 并返回最近一次的成员快照（若无则返回空列表）
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.ShopMemberListResponse} "data: { items: [] }"
// @Router       /shops/members/list [post]
func (s *GameShopMemberService) List(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			// 平台侧兜底：返回店铺管理员作为成员（member_type=管理员）
			if s.sAdm != nil {
				admins, err := s.sAdm.ListByHouse(c.Request.Context(), int32(in.HouseGID))
				if err == nil && len(admins) > 0 {
					nameMap := map[int32]string{}
					if s.users != nil {
						uidSet := make(map[int32]struct{}, len(admins))
						for _, a := range admins {
							uidSet[a.UserID] = struct{}{}
						}
						uids := make([]int32, 0, len(uidSet))
						for uid := range uidSet {
							uids = append(uids, uid)
						}
						if rows, err2 := s.users.SelectByPK(c.Request.Context(), uids); err2 == nil {
							for _, u := range rows {
								if u != nil {
									nameMap[u.Id] = u.NickName
								}
							}
						}
					}
					out := make([]resp.ShopMemberListItem, 0, len(admins))
					for _, a := range admins {
						out = append(out, resp.ShopMemberListItem{
							UserID:     uint32(a.UserID),
							UserStatus: 0,
							GameID:     0,
							MemberID:   0,
							MemberType: 2, // 管理员
							NickName:   nameMap[a.UserID],
							GroupID:    0,
						})
					}
					response.Success(c, resp.ShopMemberListResponse{Items: out})
					return
				}
			}
			response.Success(c, resp.ShopMemberListResponse{Items: []resp.ShopMemberListItem{}})
			return
		}
	}
	// 触发拉取
	sess.GetGroupMembers()
	// 返回快照（若无为 nil -> 空数组）
	mems := sess.ListMembers()

	// 构建 GameID 到游戏账号的映射
	// 注意：m.GameID 是游戏内玩家ID（game_player_id），不是 game_account.id
	gameIDToAccount := make(map[uint32]*gameModel.GameAccount)
	if s.gameAccount != nil && len(mems) > 0 {
		// 批量查询游戏账号（通过 GamePlayerID）
		for _, m := range mems {
			gamePlayerID := fmt.Sprintf("%d", m.GameID)
			if account, err := s.gameAccount.GetByGamePlayerID(c.Request.Context(), gamePlayerID); err == nil && account != nil {
				gameIDToAccount[m.GameID] = account
			}
		}
	}

	// 收集所有需要查询的平台用户ID
	userIDSet := make(map[int32]struct{})
	for _, account := range gameIDToAccount {
		if account.UserID != nil && *account.UserID > 0 {
			userIDSet[*account.UserID] = struct{}{}
		}
	}

	// 批量查询平台用户信息
	userIDToInfo := make(map[int32]*resp.UserInfo)
	if s.users != nil && len(userIDSet) > 0 {
		userIDs := make([]int32, 0, len(userIDSet))
		for uid := range userIDSet {
			userIDs = append(userIDs, uid)
		}

		if users, err := s.users.SelectByPK(c.Request.Context(), userIDs); err == nil {
			for _, u := range users {
				if u != nil {
					userIDToInfo[u.Id] = &resp.UserInfo{
						ID:           u.Id,
						Username:     u.Username,
						NickName:     u.NickName,
						Avatar:       u.Avatar,
						Role:         u.Role,
						Introduction: u.Introduction,
						CreatedAt:    u.CreatedAt.Format("2006-01-02 15:04:05"),
						UpdatedAt:    u.UpdatedAt.Format("2006-01-02 15:04:05"),
					}
				}
			}
		}
	}

	// 组装响应（包含平台用户信息）
	out := make([]resp.ShopMemberListItem, 0, len(mems))
	for _, m := range mems {
		// m.GameID 就是 game_player_id
		gamePlayerID := fmt.Sprintf("%d", m.GameID)

		item := resp.ShopMemberListItem{
			UserID:         m.UserID,
			UserStatus:     m.UserStatus,
			GameID:         m.GameID,
			MemberID:       m.MemberID,
			MemberType:     m.MemberType,
			NickName:       m.NickName,
			GroupID:        0,
			IsBindPlatform: false,
			GamePlayerID:   gamePlayerID, // 始终返回 game_player_id
		}

		// 添加游戏账号ID和平台用户信息
		if account, ok := gameIDToAccount[m.GameID]; ok {
			gameAccountID := uint32(account.Id)
			item.GameAccountID = &gameAccountID

			if account.UserID != nil && *account.UserID > 0 {
				item.IsBindPlatform = true
				if userInfo, exists := userIDToInfo[*account.UserID]; exists {
					item.PlatformUser = userInfo
				}
			}
		}

		// 查询当前圈子信息（使用游戏内玩家ID查询，不依赖 game_account）
		if s.gameMember != nil {
			// m.GameID 就是游戏内玩家ID，game_member.game_id 存储的也是游戏内玩家ID
			if member, err := s.gameMember.GetByGameID(c.Request.Context(), int32(in.HouseGID), int32(m.GameID)); err == nil && member != nil {
				if member.GroupID != nil {
					// 设置新字段
					item.CurrentGroupID = member.GroupID
					item.CurrentGroupName = member.GroupName
					// 同时设置旧字段（兼容前端）
					item.GroupID = int(*member.GroupID)
					item.GroupName = &member.GroupName
				}
			}
		}

		out = append(out, item)
	}
	response.Success(c, resp.ShopMemberListResponse{Items: out})
}

// QueryDiamond
// @Summary      查询当前中控账号财富（钻石）
// @Description  调用会话 GetDiamond；接口返回触发状态。
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.DiamondQueryResponse} "data: { triggered: bool }"
// @Router       /shops/diamond/query [post]
func (s *GameShopMemberService) QueryDiamond(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			response.Fail(c, ecode.Failed, "session not found or not online for this house")
			return
		}
	}
	sess.GetDiamond()
	response.Success(c, resp.DiamondQueryResponse{Triggered: true})
}

// Logout
// @Summary      强制用户下线（区别于踢出）
// @Description  使用 KickOffMember 强制下线（实现侧等价于从圈成员删除），与“踢人”动作在结果上相同。
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.KickMemberRequest true "house_gid, member_id"
// @Success      200 {object} response.Body
// @Router       /shops/members/logout [post]
func (s *GameShopMemberService) Logout(c *gin.Context) {
	var in req.KickMemberRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.mgr.KickMember(int(claims.BaseClaims.UserID), in.HouseGID, in.MemberID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// PullMembers
// @Summary      手动刷新成员列表（触发拉取）
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body
// @Router       /shops/members/pull [post]
func (s *GameShopMemberService) PullMembers(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			response.Fail(c, ecode.Failed, "session not found or not online for this house")
			return
		}
	}
	sess.GetGroupMembers()
	response.SuccessWithOK(c)
}

// 规则设置
func (s *GameShopMemberService) SetVIP(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
		MemberID int32 `json:"member_id" binding:"required"`
		VIP      bool  `json:"vip"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.rule.SetVIP(c.Request.Context(), uid, in.HouseGID, in.MemberID, in.VIP); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
func (s *GameShopMemberService) SetMulti(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
		MemberID int32 `json:"member_id" binding:"required"`
		Allow    bool  `json:"allow"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.rule.SetMultiGIDs(c.Request.Context(), uid, in.HouseGID, in.MemberID, in.Allow); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
func (s *GameShopMemberService) SetTempRelease(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
		MemberID int32 `json:"member_id" binding:"required"`
		Limit    int32 `json:"limit"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.rule.SetTempRelease(c.Request.Context(), uid, in.HouseGID, in.MemberID, in.Limit); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// PullToGroup 将游戏成员拉入圈子（直接通过游戏内玩家ID）
// POST /api/shops/members/pull-to-group
// 注意：此接口会自动创建 game_member 记录，不需要预先创建 game_account
func (s *GameShopMemberService) PullToGroup(c *gin.Context) {
	var in struct {
		HouseGID      int32    `json:"house_gid" binding:"required"`       // 店铺ID
		GroupID       int32    `json:"group_id" binding:"required"`        // 目标圈子ID
		GamePlayerIDs []string `json:"game_player_ids" binding:"required"` // game_player_id 列表（游戏内玩家ID）
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	if len(in.GamePlayerIDs) == 0 {
		response.Fail(c, ecode.ParamsFailed, "请选择要拉入的成员")
		return
	}

	// 查询圈子信息
	group, err := s.groupRepo.GetByID(c.Request.Context(), in.GroupID)
	if err != nil {
		response.Fail(c, ecode.Failed, fmt.Sprintf("查询圈子失败: %v", err))
		return
	}

	// 验证圈子属于指定店铺
	if group.HouseGID != in.HouseGID {
		response.Fail(c, ecode.Failed, "圈子不属于该店铺")
		return
	}

	// 从会话获取成员信息
	sess, ok := s.mgr.GetAnyByHouse(int(in.HouseGID))
	if !ok || sess == nil {
		response.Fail(c, ecode.Failed, "店铺会话不存在")
		return
	}

	// 触发拉取最新成员列表
	sess.GetGroupMembers()
	allMembers := sess.ListMembers()

	// 为每个游戏玩家创建或更新 game_member
	successCount := 0
	for _, gamePlayerID := range in.GamePlayerIDs {
		// 在会话中查找该成员
		var found bool
		var memberNickname string
		var gamePlayerIDInt int32

		for _, m := range allMembers {
			if fmt.Sprintf("%d", m.GameID) == gamePlayerID {
				found = true
				memberNickname = m.NickName
				gamePlayerIDInt = int32(m.GameID)
				break
			}
		}

		if !found {
			continue
		}

		// 1. 直接更新 game_account_group（使用 game_player_id）
		// 不再需要查询 game_account，直接用 game_player_id 操作
		if err := s.accountGroupRepo.UpdateGroupByGamePlayerAndHouse(
			c.Request.Context(),
			gamePlayerID, // 直接使用 game_player_id
			in.HouseGID,
			in.GroupID,
			group.GroupName,
		); err != nil {
			// game_account_group 更新失败，但继续更新 game_member
			// 这样至少前端能看到圈子信息
		}

		// 3. 更新 game_member（用于前端显示，必须更新）
		existingMember, err := s.gameMember.GetByGameID(c.Request.Context(), in.HouseGID, gamePlayerIDInt)
		if err != nil || existingMember == nil {
			// 不存在，创建新记录
			member := &gameModel.GameMember{
				HouseGID:  in.HouseGID,
				GameID:    gamePlayerIDInt, // 直接使用游戏内玩家ID
				GameName:  memberNickname,
				GroupID:   &in.GroupID,
				GroupName: group.GroupName,
			}
			if err := s.gameMember.Create(c.Request.Context(), member); err != nil {
				continue
			}
		} else {
			// 已存在，更新圈子信息
			if err := s.gameMember.UpdateGroup(c.Request.Context(), in.HouseGID, gamePlayerIDInt, in.GroupID, group.GroupName); err != nil {
				continue
			}
		}

		successCount++
	}

	if successCount == 0 {
		response.Fail(c, ecode.Failed, "没有成员被成功拉入圈子")
		return
	}

	response.Success(c, map[string]interface{}{
		"message": "拉圈成功",
		"count":   successCount,
	})
}

// RemoveFromGroup 将游戏成员从圈子中移除
// POST /api/shops/members/remove-from-group
func (s *GameShopMemberService) RemoveFromGroup(c *gin.Context) {
	var in struct {
		HouseGID      int32    `json:"house_gid" binding:"required"`       // 店铺ID
		GamePlayerIDs []string `json:"game_player_ids" binding:"required"` // game_player_id 列表（游戏内玩家ID）
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	if len(in.GamePlayerIDs) == 0 {
		response.Fail(c, ecode.ParamsFailed, "请选择要移除的成员")
		return
	}

	// 从会话获取成员信息
	sess, ok := s.mgr.GetAnyByHouse(int(in.HouseGID))
	if !ok || sess == nil {
		response.Fail(c, ecode.Failed, "店铺会话不存在")
		return
	}

	// 触发拉取最新成员列表
	sess.GetGroupMembers()
	allMembers := sess.ListMembers()

	// 移除成员圈子信息并踢出店铺
	successCount := 0
	kickedCount := 0

	for _, gamePlayerID := range in.GamePlayerIDs {
		// 在会话中查找该成员
		var found bool
		var gamePlayerIDInt int32
		var memberID int

		for _, m := range allMembers {
			if fmt.Sprintf("%d", m.GameID) == gamePlayerID {
				found = true
				gamePlayerIDInt = int32(m.GameID)
				memberID = int(m.MemberID) // 保存 member_id 用于踢出
				break
			}
		}

		if !found {
			continue
		}

		// 1. 查询 game_account
		account, err := s.gameAccount.GetByGamePlayerID(c.Request.Context(), gamePlayerID)
		if err != nil || account == nil {
			// 游戏账号不存在，跳过
			continue
		}

		// 2. 从 game_account_group 中移除（战绩同步和计费依赖这个表）
		// 注意：这里需要删除或标记为 inactive，具体取决于业务需求
		// 暂时跳过 game_account_group 的处理，因为可能需要保留历史记录
		// TODO: 需要实现 game_account_group 的处理逻辑

		// 3. 更新 game_member（清空圈子信息）
		existingMember, err := s.gameMember.GetByGameID(c.Request.Context(), in.HouseGID, gamePlayerIDInt)
		if err != nil || existingMember == nil {
			// 不存在记录，跳过
			continue
		}

		// 清空圈子信息（设置为 0 和空字符串）
		if err := s.gameMember.UpdateGroup(c.Request.Context(), in.HouseGID, gamePlayerIDInt, 0, ""); err != nil {
			continue
		}

		successCount++

		// 4. 踢出店铺：调用游戏服务器命令（类似 passing-dragonfly 的 "踢" 命令）
		if memberID > 0 {
			// 发送踢人命令到游戏服务器
			sess.KickOffMember(int(in.HouseGID), memberID)
			kickedCount++
		}
	}

	if successCount == 0 {
		response.Fail(c, ecode.Failed, "没有成员被成功移除")
		return
	}

	response.Success(c, map[string]interface{}{
		"message":      fmt.Sprintf("成功踢出 %d 名成员（已从圈子和店铺中移除）", successCount),
		"count":        successCount,
		"kicked_count": kickedCount,
	})
}
