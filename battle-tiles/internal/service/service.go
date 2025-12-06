package service

import (
	"battle-tiles/internal/service/basic"
	"battle-tiles/internal/service/game"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewPublicService,
	NewOpsService,
	NewPlatformService,
	NewAsyNQService,

	basic.NewBasicUserService,
	basic.NewBasicRoleService,
	basic.NewBasicLoginService,
	basic.NewBasicMenuService,
	basic.NewBasicPermissionService,

	game.NewSessionService,
	game.NewAccountService,
	game.NewFundsService,
	game.NewCtrlAccountService,
	game.NewShopAdminService,
	game.NewHouseSettingsService,
	game.NewShopTableService,
	game.NewGameShopMemberService,
	game.NewGameStatsService,
	game.NewWalletQueryService,
	game.NewShopApplicationService,
	game.NewGameGroupService,
	game.NewBattleRecordService,
	game.NewShopGroupService,
	game.NewMemberService,
	game.NewRoomCreditLimitService,
	NewSessionMonitor,
)
