package biz

import (
	"battle-tiles/internal/biz/basic"
	"battle-tiles/internal/biz/cloud"
	"battle-tiles/internal/biz/game"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewPublicUseCase,
	NewAsyNQUseCase,
	basic.NewBasicUserUseCase,
	basic.NewBasicLoginUseCase,
	basic.NewBasicMenuUseCase,

	cloud.NewPlatformUsecase,

	game.NewGameAccountUseCase,
	game.NewGameAccountGroupUseCase, // 新增：游戏账号圈子业务逻辑
	game.NewFundsUseCase,
	game.NewCtrlAccountUseCase,
	game.NewShopAdminUseCase,
	game.NewHouseSettingsUseCase,
	game.NewGameStatsUseCase,
	game.NewCtrlSessionUseCase,
	game.NewMemberRuleUseCase,
	game.NewBattleRecordUseCase,
	game.NewBattleSyncManager,
	game.NewShopGroupUseCase,
	game.NewMemberUseCase,
	game.NewBattleQueryUseCase,
	game.NewBalanceQueryUseCase,
	game.NewRoomCreditLimitUseCase, // 保留：用于玩家坐下时检查额度
	game.NewRoomCreditEventHandler, // 保留：用于玩家坐下时检查额度
)
