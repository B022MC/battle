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
	shopGroupAdminService  *game.ShopGroupAdminService
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

	// 我的圈子/圈子列表
	r.shopGroupAdminService.RegisterRouter(root)
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
	shopGroupAdminService *game.ShopGroupAdminService,
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
		shopGroupAdminService:  shopGroupAdminService,
	}
}
