package repo

import (
	"battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/repo/cloud"
	"battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/dal/repo/rbac"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewPublicRepo,
	NewAsyNQRepo,
	basic.NewBasicUseRepo,
	basic.NewBasicLoginRepo,
	basic.NewAuthRepo,
	basic.NewBaseMenuRepo,
	basic.NewBaseRoleMenuRelRepo,
	basic.NewBaseRoleMenuBtnRelRepo,

	cloud.NewBasePlatformRepo,

	game.NewGameAccountRepo,
	game.NewSessionRepo,
	game.NewCtrlAccountRepo,
	game.NewShopAdminRepo,
	game.NewWalletRepo,
	game.NewStatsRepo,
	game.NewWalletReadRepo,
	game.NewHouseSettingsRepo,
	game.NewGameAccountHouseRepo,
	game.NewCtrlAccountHouseRepo,
	game.NewMemberRuleRepo,
	game.NewFeeSettleRepo,
	game.NewUserApplicationRepo,
	game.NewBattleRecordRepo,
	game.NewShopGroupRepo,
	game.NewShopGroupMemberRepo,
	rbac.NewStore,
)
