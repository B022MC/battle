package service

import (
	cloudBiz "battle-tiles/internal/biz/cloud"
	"battle-tiles/internal/consts"
	"battle-tiles/pkg/utils/response"
	"context"

	"github.com/gin-gonic/gin"
)

// PlatformService 提供平台相关公共接口
type PlatformService struct{ uc *cloudBiz.PlatformUsecase }

func NewPlatformService(uc *cloudBiz.PlatformUsecase) *PlatformService {
	return &PlatformService{uc: uc}
}

func (s *PlatformService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/platforms")
	g.GET("/list", s.List)
	// Plaza related public constants for frontend consumption
	g.GET("/plaza/consts", s.PlazaConstants)
}

// List 列出全部平台
// @Summary      平台列表
// @Description  返回全部平台
// @Tags         Public
// @Produce      json
// @Success      200 {object} response.Body
// @Router       /platforms/list [get]
func (s *PlatformService) List(c *gin.Context) {
	items, err := s.uc.ListAll(context.Background())
	if err != nil {
		response.Fail(c, 500, err)
		return
	}
	response.Success(c, items)
}

type labeledInt struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Label string `json:"label"`
}

type labeledStr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Label string `json:"label"`
}

// PlazaConstants exposes selected plaza-related constants and their human-readable meanings.
// @Summary      Plaza constants
// @Description  Return plaza/game constants and labels for frontend
// @Tags         Public
// @Produce      json
// @Success      200 {object} response.Body
// @Router       /platforms/plaza/consts [get]
func (s *PlatformService) PlazaConstants(c *gin.Context) {
	out := map[string]any{}

	// Login modes
	out["modes"] = []labeledInt{
		{Key: "MODE_CHOOSE", Value: consts.MODE_CHOOSE, Label: "选择登录"},
		{Key: "MODE_ACCOUNT", Value: consts.MODE_ACCOUNT, Label: "帐号登录"},
		{Key: "MODE_REGISTER", Value: consts.MODE_REGISTER, Label: "注册帐号"},
		{Key: "MODE_SERVICE", Value: consts.MODE_SERVICE, Label: "服务界面"},
		{Key: "MODE_OPTION", Value: consts.MODE_OPTION, Label: "设置界面"},
		{Key: "MODE_WECHAT", Value: consts.MODE_WECHAT, Label: "微信登陆"},
		{Key: "MODE_STATION", Value: consts.MODE_STATION, Label: "工作站"},
		{Key: "MODE_PHONE", Value: consts.MODE_PHONE, Label: "手机登录"},
		{Key: "MODE_RESET", Value: consts.MODE_RESET, Label: "重置密码"},
	}

	// Scenes
	out["scenes"] = []labeledInt{
		{Key: "SCENE_LOGON", Value: consts.SCENE_LOGON, Label: "登录"},
		{Key: "SCENE_GAMELIST", Value: consts.SCENE_GAMELIST, Label: "游戏列表"},
		{Key: "SCENE_ROOMLIST", Value: consts.SCENE_ROOMLIST, Label: "房间列表"},
		{Key: "SCENE_USERINFO", Value: consts.SCENE_USERINFO, Label: "用户信息"},
		{Key: "SCENE_REDEEM_CODE", Value: consts.SCENE_REDEEM_CODE, Label: "兑换码"},
		{Key: "SCENE_BASEENSURE", Value: consts.SCENE_BASEENSURE, Label: "低保"},
		{Key: "SCENE_OPTION", Value: consts.SCENE_OPTION, Label: "设置"},
		{Key: "SCENE_SERVICE", Value: consts.SCENE_SERVICE, Label: "客服"},
		{Key: "SCENE_SYSTEM", Value: consts.SCENE_SYSTEM, Label: "系统"},
		{Key: "SCENE_SHOP", Value: consts.SCENE_SHOP, Label: "商城"},
		{Key: "SCENE_BATTLE_LIST", Value: consts.SCENE_BATTLE_LIST, Label: "约战列表"},
		{Key: "SCENE_BATTLE_CREATE", Value: consts.SCENE_BATTLE_CREATE, Label: "创建约战"},
		{Key: "SCENE_BATTLE_RECORD", Value: consts.SCENE_BATTLE_RECORD, Label: "约战记录"},
		{Key: "SCENE_BATTLE_FIND", Value: consts.SCENE_BATTLE_FIND, Label: "查找约战"},
		{Key: "SCENE_BATTLE_SCORE", Value: consts.SCENE_BATTLE_SCORE, Label: "约战结算"},
		{Key: "SCENE_BENEFIT", Value: consts.SCENE_BENEFIT, Label: "福利"},
		{Key: "SCENE_LUCKY_ROLL", Value: consts.SCENE_LUCKY_ROLL, Label: "幸运转盘"},
		{Key: "SCENE_RANK", Value: consts.SCENE_RANK, Label: "排行"},
		{Key: "SCENE_MORE_GAMES", Value: consts.SCENE_MORE_GAMES, Label: "更多游戏"},
		{Key: "SCENE_TEAHOUSE", Value: consts.SCENE_TEAHOUSE, Label: "茶馆"},
		{Key: "SCENE_MATCHLIST", Value: consts.SCENE_MATCHLIST, Label: "比赛列表"},
		{Key: "SCENE_MATCHWAIT", Value: consts.SCENE_MATCHWAIT, Label: "比赛等待"},
		{Key: "SCENE_HEALTH_DISPLAY", Value: consts.SCENE_HEALTH_DISPLAY, Label: "健康提示"},
		{Key: "SCENE_GAME", Value: consts.SCENE_GAME, Label: "游戏内"},
	}

	// User status
	out["user_status"] = []labeledInt{
		{Key: "US_NULL", Value: consts.US_NULL, Label: "无状态"},
		{Key: "US_FREE", Value: consts.US_FREE, Label: "站立"},
		{Key: "US_SIT", Value: consts.US_SIT, Label: "坐下"},
		{Key: "US_READY", Value: consts.US_READY, Label: "同意"},
		{Key: "US_LOOKON", Value: consts.US_LOOKON, Label: "旁观"},
		{Key: "US_PLAYING", Value: consts.US_PLAYING, Label: "游戏中"},
		{Key: "US_OFFLINE", Value: consts.US_OFFLINE, Label: "离线"},
	}

	// Member types (protocol-specific; adjust as needed)
	out["member_types"] = []labeledInt{
		{Key: "MEMBER_NORMAL", Value: 0, Label: "普通成员"},
		{Key: "MEMBER_OWNER", Value: 1, Label: "圈主/馆主"},
		{Key: "MEMBER_ADMIN", Value: 2, Label: "管理员"},
	}

	// Game genre
	out["game_genre"] = []labeledInt{
		{Key: "GAME_GENRE_GOLD", Value: consts.GAME_GENRE_GOLD, Label: "金币类型"},
		{Key: "GAME_GENRE_SCORE", Value: consts.GAME_GENRE_SCORE, Label: "点值类型"},
		{Key: "GAME_GENRE_MATCH", Value: consts.GAME_GENRE_MATCH, Label: "比赛类型"},
		{Key: "GAME_GENRE_EDUCATE", Value: consts.GAME_GENRE_EDUCATE, Label: "训练类型"},
	}

	// Table genre
	out["table_genre"] = []labeledInt{
		{Key: "TABLE_GENRE_NORMAL", Value: consts.TABLE_GENRE_NORMAL, Label: "普通房间"},
		{Key: "TABLE_GENRE_CREATE", Value: consts.TABLE_GENRE_CREATE, Label: "创建模式"},
		{Key: "TABLE_GENRE_DISTRIBUTE", Value: consts.TABLE_GENRE_DISTRIBUTE, Label: "分配模式"},
	}

	// Game kinds (KindID -> name)
	out["game_kinds"] = []labeledInt{
		{Key: "GameKindDingErHong", Value: consts.GameKindDingErHong, Label: consts.GetKindName(consts.GameKindDingErHong)},
		{Key: "GameKindHongErShi", Value: consts.GameKindHongErShi, Label: consts.GetKindName(consts.GameKindHongErShi)},
		{Key: "GameKindDuanGouQia3", Value: consts.GameKindDuanGouQia3, Label: consts.GetKindName(consts.GameKindDuanGouQia3)},
		{Key: "GameKindPaoDeKuai2", Value: consts.GameKindPaoDeKuai2, Label: consts.GetKindName(consts.GameKindPaoDeKuai2)},
		{Key: "GameKindDouShiSi", Value: consts.GameKindDouShiSi, Label: consts.GetKindName(consts.GameKindDouShiSi)},
		{Key: "GameKindDuanGouQia2", Value: consts.GameKindDuanGouQia2, Label: consts.GetKindName(consts.GameKindDuanGouQia2)},
		{Key: "GameKindHongZhong2", Value: consts.GameKindHongZhong2, Label: consts.GetKindName(consts.GameKindHongZhong2)},
	}

	// Message types
	out["system_message_types"] = []labeledInt{
		{Key: "SMT_CHAT", Value: consts.SMT_CHAT, Label: "聊天消息"},
		{Key: "SMT_EJECT", Value: consts.SMT_EJECT, Label: "弹出消息"},
		{Key: "SMT_GLOBAL", Value: consts.SMT_GLOBAL, Label: "全局消息"},
		{Key: "SMT_PROMPT", Value: consts.SMT_PROMPT, Label: "提示消息"},
		{Key: "SMT_TABLE_ROLL", Value: consts.SMT_TABLE_ROLL, Label: "滚动消息"},
		{Key: "SMT_CLOSE_ROOM", Value: consts.SMT_CLOSE_ROOM, Label: "关闭房间"},
		{Key: "SMT_CLOSE_GAME", Value: consts.SMT_CLOSE_GAME, Label: "关闭游戏"},
		{Key: "SMT_CLOSE_LINK", Value: consts.SMT_CLOSE_LINK, Label: "中断连接"},
		{Key: "SMT_SHOW_MOBILE", Value: consts.SMT_SHOW_MOBILE, Label: "手机显示"},
	}

	// URLs and versions
	out["versions"] = map[string]any{
		"LUA_VERSION":    consts.LUA_VERSION,
		"CLIENT_VERSION": consts.CLIENT_VERSION,
		"APP_VERSION":    consts.APP_VERSION,
	}
	out["urls"] = []labeledStr{
		{Key: "URL_LOGON_INFO", Value: consts.URL_LOGON_INFO, Label: "获取登录信息"},
		{Key: "URL_UPDATE_LUA", Value: consts.URL_UPDATE_LUA, Label: "热更新地址获取"},
		{Key: "URL_UPDATE_KERNEL", Value: consts.URL_UPDATE_KERNEL, Label: "安装包地址获取"},
		{Key: "URL_APP", Value: consts.URL_APP, Label: "默认安装包地址"},
		{Key: "URL_LUA", Value: consts.URL_LUA, Label: "默认LUA地址"},
	}

	response.Success(c, out)
}
