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
)
