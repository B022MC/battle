package router

import (
	"battle-tiles/internal/service/game"

	"github.com/gin-gonic/gin"
)

type GameRouter struct {
	accountService         *game.AccountService
	sessionService         *game.SessionService
	fundsService           *game.FundsService
	ctrlAccountService     *game.CtrlAccountService
	shopAdminService       *game.ShopAdminService
	shopTableService       *game.ShopTableService
	gameShopMemberService  *game.GameShopMemberService
	gameStatsService       *game.GameStatsService
	walletQueryService     *game.WalletQueryService
	shopApplicationService *game.ShopApplicationService
	gameGroupService       *game.GameGroupService
	houseSettingsService   *game.HouseSettingsService
	battleRecordService    *game.BattleRecordService
	shopGroupService       *game.ShopGroupService
	memberService          *game.MemberService
	roomCreditLimitService *game.RoomCreditLimitService
}

func (r *GameRouter) InitRouter(root *gin.RouterGroup) {

	r.accountService.RegisterRouter(root)

	r.sessionService.RegisterRouter(root)

	r.ctrlAccountService.RegisterRouter(root)

	r.shopAdminService.RegisterRouter(root)

	r.shopTableService.RegisterRouter(root)

	r.fundsService.RegisterRouter(root)

	r.gameShopMemberService.RegisterRouter(root)

	r.gameStatsService.RegisterRouter(root)

	r.walletQueryService.RegisterRouter(root)

	r.shopApplicationService.RegisterRouter(root)

	r.gameGroupService.RegisterRouter(root)

	// 店铺设置
	r.houseSettingsService.RegisterRouter(root)

	// 战绩查询
	r.battleRecordService.RegisterRouter(root)

	// 圈子管理
	r.shopGroupService.RegisterRouter(root)

	// 成员管理
	r.memberService.RegisterRouter(root)

	// 房间额度设置
	r.roomCreditLimitService.RegisterRouter(root)
}

func NewGameRouter(
	accountService *game.AccountService,
	sessionService *game.SessionService,
	fundsService *game.FundsService,
	ctrlAccountService *game.CtrlAccountService,
	shopAdminService *game.ShopAdminService,
	shopTableService *game.ShopTableService,
	gameShopMemberService *game.GameShopMemberService,
	gameStatsService *game.GameStatsService,
	walletQueryService *game.WalletQueryService,
	shopApplicationService *game.ShopApplicationService,
	gameGroupService *game.GameGroupService,
	houseSettingsService *game.HouseSettingsService,
	battleRecordService *game.BattleRecordService,
	shopGroupService *game.ShopGroupService,
	memberService *game.MemberService,
	roomCreditLimitService *game.RoomCreditLimitService,
) *GameRouter {
	return &GameRouter{
		accountService:         accountService,
		sessionService:         sessionService,
		fundsService:           fundsService,
		ctrlAccountService:     ctrlAccountService,
		shopAdminService:       shopAdminService,
		shopTableService:       shopTableService,
		gameShopMemberService:  gameShopMemberService,
		gameStatsService:       gameStatsService,
		walletQueryService:     walletQueryService,
		shopApplicationService: shopApplicationService,
		gameGroupService:       gameGroupService,
		houseSettingsService:   houseSettingsService,
		battleRecordService:    battleRecordService,
		shopGroupService:       shopGroupService,
		memberService:          memberService,
		roomCreditLimitService: roomCreditLimitService,
	}
}
