package consts

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	DefaultIntMinusOneValue = -1
	DefaultBatchSize        = 1000
	DefaultPk               = "id"
	RootRole                = "root" // 超管
	UID                     = "uid"
)
const (
	//STATION_ID = 1000
	CLIENT_SRC = "client.src."
	PLAZA_SRC  = "plaza.src."
	GAME_SRC   = ""
	GAME_RES   = ""

	LUA_VERSION = 1.5
	LOGONSERVER = ""
	LOGONPORT   = ""

	MODE_CHOOSE   = 1 //选择登录
	MODE_ACCOUNT  = 2 //帐号登录
	MODE_REGISTER = 3 //注册帐号
	MODE_SERVICE  = 4 //服务界面
	MODE_OPTION   = 5 //设置界面
	MODE_WECHAT   = 6 //微信登陆
	MODE_STATION  = 7
	MODE_PHONE    = 8 //手机登录
	MODE_RESET    = 9 //重置密码

	SCENE_LOGON          = 1
	SCENE_GAMELIST       = 2
	SCENE_ROOMLIST       = 3
	SCENE_USERINFO       = 6
	SCENE_REDEEM_CODE    = 7
	SCENE_BASEENSURE     = 8
	SCENE_OPTION         = 11
	SCENE_SERVICE        = 12
	SCENE_SYSTEM         = 14
	SCENE_SHOP           = 15
	SCENE_BATTLE_LIST    = 16
	SCENE_BATTLE_CREATE  = 17
	SCENE_BATTLE_RECORD  = 18
	SCENE_BATTLE_FIND    = 19
	SCENE_BATTLE_SCORE   = 20
	SCENE_BENEFIT        = 21
	SCENE_LUCKY_ROLL     = 22
	SCENE_RANK           = 23
	SCENE_MORE_GAMES     = 24
	SCENE_TEAHOUSE       = 25
	SCENE_MATCHLIST      = 26
	SCENE_MATCHWAIT      = 27
	SCENE_HEALTH_DISPLAY = 28

	SCENE_GAME = 30

	SCENE_EX_END = 50

	MAIN_SOCKET_INFO = 0

	SUB_SOCKET_CONNECT = 1
	SUB_SOCKET_ERROR   = 2
	SUB_SOCKET_CLOSE   = 3

	US_NULL    = 0x00 //没有状态
	US_FREE    = 0x01 //站立状态
	US_SIT     = 0x02 //坐下状态
	US_READY   = 0x03 //同意状态
	US_LOOKON  = 0x04 //旁观状态
	US_PLAYING = 0x05 //游戏状态
	US_OFFLINE = 0x06 //断线状态

	FACE_X = 48
	FACE_Y = 48

	INVALID_TABLE = 65535
	INVALID_CHAIR = 65535
	INVALID_ITEM  = 65535

	GENDER_FEMALE  = 1 //女性性别
	GENDER_MANKIND = 2 //男性性别
	GENDER_SECRET  = 0 //保密

	GAME_GENRE_GOLD    = 0x0001 //金币类型
	GAME_GENRE_SCORE   = 0x0002 //点值类型
	GAME_GENRE_MATCH   = 0x0004 //比赛类型
	GAME_GENRE_EDUCATE = 0x0008 //训练类型

	TABLE_GENRE_NORMAL     = 0x0000 //普通房间
	TABLE_GENRE_CREATE     = 0x1000 //创建模式
	TABLE_GENRE_DISTRIBUTE = 0x2000 //分配模式

	LEN_GAME_LIST_ITEM     = 142
	LEN_GAME_SERVER_ITEM   = 236
	LEN_CREATE_OPTION_ITEM = 93 + 40 + 16

	LEN_MD5        = 33 //加密密码
	LEN_ACCOUNTS   = 32 //帐号长度
	LEN_NICKNAME   = 32 //昵称长度
	LEN_PASSWORD   = 33 //密码长度
	LEN_SERVER     = 32 //房名长度
	LEN_PROCESS    = 32
	LEN_DOMAIN     = 63
	LEN_GROUP_NAME = 32

	LEN_MOBILE_PHONE = 16 //移动电话
	LEN_COMPELLATION = 16 //真实名字
	LEN_MACHINE_ID   = 33 //序列长度

	LEN_USER_CHAT = 128

	MDM_GP_LOGON = 1   //广场登录
	MDM_MB_LOGON = 100 //广场登录

	SUB_MB_LOGON_ACCOUNTS      = 2 //帐号登录
	SUB_MB_REGISTER_ACCOUNTS   = 3 //注册帐号
	SUB_MB_LOGON_OTHERPLATFORM = 4 //其他登陆

	SUB_MB_LOGON_GAMEID_LUA        = 10 //I D 登录
	SUB_MB_LOGON_ACCOUNTS_LUA      = 11 //帐号登录
	SUB_MB_LOGON_OTHERPLATFORM_LUA = 12 //其他登录
	SUB_MB_REGISTER_ACCOUNTS_LUA   = 13 //注册帐号

	SUB_MB_LOGON_AUTHCODE    = 21 //验证码登陆
	SUB_MB_LOGON_MOBILEPHONE = 22 //手机登陆

	SUB_MB_LOGON_SUCCESS = 100 //登录成功
	SUB_MB_LOGON_FAILURE = 101 //登录失败
	SUB_MB_LOGON_FINISH  = 102 //登录完成
	SUB_MB_UPDATE_NOTIFY = 200 //升级提示

	MDM_MB_SERVER_LIST = 101 //列表信息

	SUB_MB_LIST_KIND     = 100 //种类列表
	SUB_MB_LIST_SERVER   = 101 //房间列表
	SUB_MB_LIST_MATCH    = 102 //比赛列表
	SUB_MB_CREATE_OPTION = 103 //开桌选项
	SUB_MB_GAME_OPTION   = 104 //游戏配置
	SUB_MB_LIST_LOGON    = 105 //登录列表
	SUB_MB_LIST_AGENT    = 106 //代理列表

	SUB_MB_SERVER_AGENT  = 107 //房间代理
	SUB_MB_LIST_FINISH   = 200 //列表完成
	SUB_MB_SERVER_FINISH = 201 //房间完成

	SUB_MB_GET_LIST         = 1 //获取列表
	SUB_MB_GET_SERVER       = 2 //获取房间
	SUB_MB_GET_OPTION       = 3 //获取配置
	SUB_MB_GET_OPTION_LUA   = 4 //获取配置
	SUB_MB_GET_SERVER_AGENT = 5 //房间代理

	MDM_GP_USER_SERVICE = 3 //用户服务

	SUB_GP_USER_FACE_INFO   = 210 //头像信息
	SUB_GP_SYSTEM_FACE_INFO = 211 //系统头像
	SUB_GP_CUSTOM_FACE_INFO = 212 //自定义头像

	SUB_GP_MODIFY_ACCOUNTS   = 1 //修改账号
	SUB_GP_MODIFY_LOGON_PASS = 2 //修改登录密码
	SUB_GP_MODIFY_INDIVIDUAL = 5 //修改资料

	SUB_GP_BASEENSURE_QUERY = 20 //查询低保
	SUB_GP_BASEENSURE_TAKE  = 21 //领取低保
	SUB_GP_GET_VALIDATECDOE = 22 //获取验证码

	//手机绑定
	SUB_GP_BIND_MOBILEPHONE   = 30 //绑定手机
	SUB_GP_UNBIND_MOBILEPHONE = 31 //解绑手机
	SUB_GP_RESET_LOGONPASSWD  = 32 //重置密码
	SUB_GP_ACQUIRE_AUTHCODE   = 35 //获取验证码

	SUB_GP_BASEENSURE_INFO   = 300 //低保信息
	SUB_GP_BASEENSURE_RESULT = 301 //领取结果
	SUB_GP_BASEENSURE_FAILED = 302 //操作失败
	SUB_GP_VALIDATECDOE_INFO = 303 //验证码信息

	SUB_GP_QUERY_WEALTH     = 13
	SUB_GP_QUERY_WEALTH_LUA = 14 //查询财富
	SUB_GP_USER_WEALTH      = 201

	MDM_GP_BANK_OPERATE = 4 //银行操作

	SUB_GP_STORAGE = 1 //银行存储
	SUB_GP_DRAWOUT = 2 //银行取出
	SUB_GP_UPDATE  = 5 //更新金币

	SUB_GP_OPERATE_SUCCESS = 100 //操作成功
	SUB_GP_OPERATE_FAILURE = 101 //操作失败

	//手机绑定
	SUB_GP_BINDMP_SUCCESS   = 350 //绑定成功
	SUB_GP_UNBINDMP_SUCCESS = 351 //解绑成功
	SUB_GP_ACQUIREAC_RESULT = 355 //获取结果

	MDM_GP_REMOTE_SERVICE = 6   //远程服务
	SUB_GP_QUERY_TABLE    = 100 //查询桌子
	SUB_GP_CREATE_TABLE   = 101 //创建桌子
	SUB_GP_QUERY_RECORD   = 102 //查询记录

	SUB_GP_TABLE_INFO       = 200 //桌子信息
	SUB_GP_TABLE_LIST       = 201 //桌子列表
	SUB_GP_OPERATE_FAILED   = 202 //操作失败
	SUB_GP_BATTLE_RECORD    = 203 //约战记录
	SUB_GP_TABLE_PARAM_LIST = 204 //桌子参数

	SUB_GP_TABLE_INFO_LUA = 205 //桌子信息
	//                            //LUA新增
	SUB_GP_QUERY_TABLE_LUA    = 110 //查询桌子
	SUB_GP_CREATE_TABLE_LUA   = 111 //创建桌子
	SUB_GP_QUERY_RECORD_LUA   = 112 //查询记录
	SUB_GA_GET_TABLE_USERLIST = 113 //获取列表
	//////////////////////////////////////////////////////////////////-
	//登录信息
	//////////////////////////////////////////////////////////////////-
	MDM_GR_LOGON            = 1 //登录信息
	SUB_GR_LOGON_MOBILE     = 2 //手机登录
	SUB_GR_LOGON_MOBILE_LUA = 4 //手机登录
	SUB_GR_LOGON_MOBILE_NEW = 5 //锁用户登陆

	SUB_GR_LOGON_SUCCESS = 100 //登录成功
	SUB_GR_LOGON_FAILURE = 101 //登录失败
	SUB_GR_LOGON_FINISH  = 102 //登录完成

	SUB_GR_UPDATE_NOTIFY = 200 //升级提示

	//////////////////////////////////////////////////////////////////-
	//配置信息
	//////////////////////////////////////////////////////////////////-
	MDM_GR_CONFIG        = 2   //配置信息
	SUB_GR_CONFIG_COLUMN = 100 //列表配置
	SUB_GR_CONFIG_SERVER = 101 //房间配置
	SUB_GR_CONFIG_FINISH = 102 //配置完成

	//////////////////////////////////////////////////////////////////-
	//用户信息
	//////////////////////////////////////////////////////////////////-
	MDM_GR_USER = 3 //用户信息

	SUB_GR_USER_RULE               = 1  //用户规则
	SUB_GR_USER_LOOKON             = 2  //旁观请求
	SUB_GR_USER_SITDOWN            = 3  //坐下请求
	SUB_GR_USER_STANDUP            = 4  //起立请求
	SUB_GR_USER_APPLY              = 5  //报名请求
	SUB_GR_USER_FEE_QUERY          = 6  //费用查询
	SUB_GR_USER_REPULSE_SIT        = 7  //拒绝玩家
	SUB_GR_USER_INFO_REQ           = 8  //请求用户信息
	SUB_GR_USER_CHAIR_REQ          = 9  //请求更换位置
	SUB_GR_USER_CHAIR_INFO_REQUEST = 10 //请求椅子用户信息
	SUB_GR_USER_DISMISS_TABLE      = 11 //解散桌子

	SUB_GR_USER_ENTER           = 100 //用户进入
	SUB_GR_USER_SCORE           = 101 //用户分数
	SUB_GR_USER_STATUS          = 102 //用户状态
	SUB_GR_USER_SEGMENT         = 103 //用户段位
	SUB_GR_REQUEST_FAILURE      = 104 //请求失败
	SUB_GR_USER_WEALTH_EX       = 105 //财富信息(有同名 ex)
	SUB_GR_USER_WAIT_DISTRIBUTE = 106
	SUB_GR_USER_DISMISS_RESULT  = 107 //解散结果
	SUB_GR_USER_MATCH_SHARE     = 108 //比赛分享

	SUB_GR_USER_CHAT         = 200 //聊天消息
	SUB_GR_USER_WEALTH       = 201 //财富信息
	SUB_GR_USER_CONVERSATION = 202 //会话消息
	SUB_GR_USER_BUGLE        = 203 //喇叭消息
	SUB_GR_WHSPER_REPLY      = 204 //自动回复

	MDM_GF_FRAME = 100 //框架命令

	SUB_GF_GAME_OPTION   = 1 //游戏配置
	SUB_GF_USER_READY    = 2 //用户准备
	SUB_GF_LOOKON_CONFIG = 3 //旁观配置

	SUB_GF_LAUNCH_DISMISS = 10 //发起解散
	SUB_GF_BALLOT_DISMISS = 11 //投票解散

	SUB_GF_DISMISS_NOTIFY  = 160 //解散提醒
	SUB_GF_DISMISS_BALLOT  = 161 //解散投票
	SUB_GF_DISMISS_SUCCESS = 162 //解散成功

	SUB_GF_GAME_STATUS   = 100 //游戏状态
	SUB_GF_GAME_SCENE    = 101 //游戏场景
	SUB_GF_LOOKON_STATUS = 102 //旁观状态

	SUB_GF_TABLE_PARAM   = 150 //桌子参数
	SUB_GF_TABLE_BATTLE  = 151 //桌子战况
	SUB_GF_TABLE_GRADE   = 152 //桌子战绩
	SUB_GF_TABLE_PARAMEX = 153 //桌子参数
	SUB_GF_TABLE_RENEW   = 155 //续桌通知
	SUB_GF_USER_VOICE    = 9   //用户语音

	SUB_GF_USER_CHAT = 3

	SUB_GF_SYSTEM_MESSAGE = 200 //系统消息
	SUB_GF_ACTION_MESSAGE = 201 //动作消息

	//////////////////////////////////////////////////////////////////-
	//状态信息
	//////////////////////////////////////////////////////////////////-
	MDM_GR_STATUS = 4 //状态信息

	SUB_GR_TABLE_INFO   = 100 //桌子信息
	SUB_GR_TABLE_STATUS = 101 //桌子状态

	MDM_GF_GAME = 200 //游戏命令

	REQUEST_FAILURE_NORMAL   = 0 //常规原因
	REQUEST_FAILURE_NOGOLD   = 1 //金币不足
	REQUEST_FAILURE_NOSCORE  = 2 //积分不足
	REQUEST_FAILURE_PASSWORD = 3 //密码错误

	REQUEST_FAILURE_CMD_NONE     = 0 //原始命令
	REQUEST_FAILURE_CMD_CONCLUDE = 1 //关闭命令

	MDM_CM_SYSTEM = 1000 //系统命令

	SUB_CM_SYSTEM_MESSAGE = 1 //系统消息
	SUB_CM_ACTION_MESSAGE = 2 //动作消息

	SMT_CHAT       = 0x0001 //聊天消息
	SMT_EJECT      = 0x0002 //弹出消息
	SMT_GLOBAL     = 0x0004 //全局消息
	SMT_PROMPT     = 0x0008 //提示消息
	SMT_TABLE_ROLL = 0x0010 //滚动消息

	SMT_CLOSE_ROOM = 0x0100 //关闭房间
	SMT_CLOSE_GAME = 0x0200 //关闭游戏
	SMT_CLOSE_LINK = 0x0400 //中断连接

	SMT_SHOW_MOBILE = 0x1000 //手机显示

	//////////////////////////////////////////////////////////////////-
	//////////////////////////////////////////////////////////////////-

	//携带信息
	DTP_GP_UI_USER_NOTE      = 1 //用户说明
	DTP_GP_UI_COMPELLATION   = 2 //真实名字
	DTP_GP_UI_SEAT_PHONE     = 3 //固定电话
	DTP_GP_UI_MOBILE_PHONE   = 4 //移动电话
	DTP_GP_UI_QQ             = 5 //Q Q 号码
	DTP_GP_UI_EMAIL          = 6 //电子邮件
	DTP_GP_UI_DWELLING_PLACE = 7 //联系地址
	DTP_GP_UI_NICKNAME       = 8 //用户昵称

	DTP_GR_TABLE_PASSWORD = 1 //桌子密码

	DTP_GR_NICK_NAME   = 10 //用户昵称
	DTP_GR_GROUP_NAME  = 11 //社团名字
	DTP_GR_UNDER_WRITE = 12 //个性签名

	DTP_GR_USER_NOTE   = 20 //用户备注
	DTP_GR_CUSTOM_FACE = 21 //自定头像

	RFC_NULL                 = 0 //无错误
	RFC_PASSWORD_INCORRECT   = 1 //密码错误
	RFC_SCORE_UNENOUGH       = 3 //游戏币不足
	RFC_CREATE_TABLE_FAILURE = 4 //创建失败
	RFC_ENTER_TABLE_FAILURE  = 5 //进入失败

	QUERY_KIND_NUMBER = 0 //编号类型
	QUERY_KIND_USERID = 1 //标识类型
	QUERY_KIND_GROUP  = 2 //标识类型

	SETTLE_KIND_TIME  = 0 //时间结算
	SETTLE_KIND_COUNT = 1 //局数结算
	SETTLE_KIND_NONE  = 2

	//财富掩码
	WEALTH_MASK_INGOT    = 0x01 //钻石掩码
	WEALTH_MASK_MEDAL    = 0x02 //奖牌掩码
	WEALTH_MASK_SCORE    = 0x04 //金币掩码
	WEALTH_MASK_ROOMCARD = 0x08 //房卡掩码

	//货币类型
	CURRENCY_KIND_INGOT    = 0 //货币类型
	CURRENCY_KIND_ROOMCARD = 1 //货币类型

	//支付类型
	PAY_TYPE_OWNER = 0x01 //房主支付
	PAY_TYPE_USER  = 0x02 //玩家支付

	//配置掩码
	OPTION_MASK_TIME     = 0x01 //时间配置
	OPTION_MASK_COUNT    = 0x02 //局数配置
	OPTION_MASK_INGOT    = 0x04 //钻石配置
	OPTION_MASK_ROOMCARD = 0x08 //房卡配置

	//配置类型
	OPTION_TYPE_NONE     = 0x00 //无效配置
	OPTION_TYPE_SINGLE   = 0x01 //单选配置
	OPTION_TYPE_MULTIPLE = 0x02 //多选配置
	OPTION_TYPE_INPUT    = 0x04 //输入配置

	FO_FORBID_RECHARGE = 0x00000001 //用户权限

	//解散状态
	DISMISS_STATE_START = 1 //发起解散
	DISMISS_STATE_OTHER = 2 //解散房间
	DISMISS_STATE_WAIT  = 3 //等待解散
	DISMISS_STATE_OVER  = 4 //解散结果
)

const (
	MDM_GA_MESSAGE_SERVICE = 5
	SUB_GA_ENTER_MESSAGE   = 1
)

const (
	STATION_ID = 2000 //1000//

	//获取地址
	URL_LOGON_INFO = "http://service.foxuc.com/GetAppService.ashx?action=GetLVersion"
	//热更新地址获取
	URL_UPDATE_LUA = "http://service.foxuc.com/GetAppService.ashx?action=GetLuaBasic"
	//安装包地址获取
	URL_UPDATE_KERNEL = "http://service.foxuc.com/GetAppService.ashx?action=GetLuaKernel"

	//默认安装包地址
	URL_APP = "http://download.foxuc.com/Loader/SCGame/App_EShop/Android/Plaza_Lua_SC.apk"
	//默认LUA地址
	URL_LUA = "http://download.foxuc.com/Loader/SCGame/App_EShop/Android/"

	//获取登录信息
	DomainGetLogon = "androidsc.foxuc.com"
	PortGetLogon   = 8200
	WebAddress     = "phone.foxuc.com"

	//基础版本
	BASE_LOGIC_VERSION = "3" //强制更新

	BASE_RES_VERSION = "1" //提示更新

	MARKET_ID = 1

	//设备版本号
	DEVICE_TYPE = 0x10

	//设备来源
	APP_SOURCE = 0x100107d0
)

const (

	//逻辑服务
	MDM_GA_LOGIC_SERVICE = 2 //逻辑服务

	//请求命令
	SUB_GA_LOGON_SERVER  = 1  //登录服务
	SUB_GA_SEARCH_GROUP  = 2  //搜索群组
	SUB_GA_CREATE_GROUP  = 3  //创建群组
	SUB_GA_UPDATE_GROUP  = 4  //更新群组
	SUB_GA_DELETE_GROUP  = 5  //删除群组
	SUB_GA_UPDATE_MEMBER = 6  //更新成员
	SUB_GA_DELETE_MEMBER = 7  //删除成员
	SUB_GA_APPLY_REQUEST = 8  //申请请求
	SUB_GA_APPLY_RESPOND = 9  //申请响应
	SUB_GA_APPLY_DELETE  = 10 //申请删除
	SUB_GA_APPLY_REPEAL  = 11 //申请撤销
	SUB_GA_SETTLE_BATTLE = 12 //约战结算

	SUB_GA_LEAVE_SERVER = 19 //离开服务

	SUB_GA_APPEND_CONFIG = 30 //20                                    //添加玩法
	SUB_GA_MODIFY_CONFIG = 31 //21									//修改玩法
	SUB_GA_DELETE_CONFIG = 32 //22

	UPMEMBER_KIND_TYPE  = 1
	UPMEMBER_KIND_RIGHT = 2 //--馆员权限

	MDM_GA_GROUP_SERVICE = 3 //					--群组服务
	SUB_GA_ENTER_GROUP   = 1

	MDM_GA_BATTLE_SERVICE = 1 //约战命令

	SUB_GA_QUERY_TABLE     = 110
	SUB_GA_DISMISS_TABLE   = 112
	SUB_GA_TABLE_ITEM      = 200 //	--桌子信息
	SUB_GA_TABLE_LIST      = 201 //--桌子列表
	SUB_GA_OPERATE_FAILED  = 202 //--操作失败
	SUB_GA_BATTLE_RECORD   = 203 //	--约战记录
	SUB_GA_TABLE_PARAMLIST = 204 //--桌子参数
	SUB_GA_TABLE_USERLIST  = 205 //--用户列表
	SUB_GA_DISMISS_RESULT  = 206 //--解散结果
	SUB_GA_USER_SITDOWN    = 300 //--用户坐下
	SUB_GA_USER_STANDUP    = 301 //--用户起立
	SUB_GA_TABLE_PARAM     = 302 //--桌子参数
	SUB_GA_TABLE_DISMISS   = 303 //	--桌子解散
	SUB_GA_TABLE_RENEW     = 304 //-续桌通知

	SUB_GA_GROUP_OPTION     = 100 //					--群组配置
	SUB_GA_APPLY_MESSAGE    = 101 //					--申请消息
	SUB_GA_LOGON_FAILURE    = 102 //					--登录失败
	SUB_GA_SEARCH_RESULT    = 103 //					--搜索结果
	SUB_GA_WEALTH_UPDATE    = 104 //					--财富更新
	SUB_GA_APPLY_DEL_RESULT = 105 //				--删除结果
	SUB_GA_OPERATE_SUCCESS  = 200 //				--操作成功
	SUB_GA_OPERATE_FAILURE  = 201 //				--操作失败
	SUB_GA_SYSTEM_MESSAGE   = 300 //				--系统消息

	SUB_GA_GROUP_ITEM     = 100 //					--群组信息
	SUB_GA_GROUP_PROPERTY = 101 //				--群组属性
	SUB_GA_GROUP_MEMBER   = 102 //				--群组成员
	SUB_GA_GROUP_UPDATE   = 103 //				--群组更新
	SUB_GA_GROUP_DELETE   = 104 //				--群组移除
	SUB_GA_MEMBER_INSERT  = 105 //					--添加成员
	SUB_GA_MEMBER_DELETE  = 106 //					--删除成员
	SUB_GA_MEMBER_UPDATE  = 107 //				--成员更新
	SUB_GA_BATTLE_CONFIG  = 120 //				--约战玩法
	SUB_GA_CONFIG_APPEND  = 121 //				--玩法添加
	SUB_GA_CONFIG_MODIFY  = 122 //				--玩法修改
	SUB_GA_CONFIG_DELETE  = 123 //				--玩法删除
	SUB_GA_ENTER_SUCCESS  = 200 //					--进入成功
	SUB_GA_ENTER_FAILURE  = 201 //					--进入失败

)

var (
	//大厅版本号
	CLIENT_VERSION = Version(1, 1, 0, 0)
	//APP版本号
	APP_VERSION = Version(1, 1, 0, 2)
)

var (
	APPLY_MSG_TYPE_REQUEST = 1
	APPLY_MSG_TYPE_RESPOND = 2

	APPLY_STATUS_NONE   = 0 //--审核状态
	APPLY_STATUS_AGREE  = 1 //--同意状态
	APPLY_STATUS_REFUSE = 2 //	--拒绝状态
	APPLY_STATUS_REPEAL = 3 //--撤销状态

)

func Version(p, m, s, b int) uint32 {
	var v uint32
	v = uint32((p & 0xff) << 24)
	v += uint32((m & 0xff) << 16)
	v += uint32((s & 0xff) << 8)
	v += uint32((b & 0xff) << 0)
	return v

}
func MachineID() string {
	//return "9989F9878AE99BAE8D746E73FE9A715C"
	return strings.ReplaceAll(uuid.Must(uuid.NewV4(), nil).String(), "-", "")
	//return fmt.Sprintf("%02d-%02d-%02d-%02d-%02d", rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100))
}

// 登录模式
//type GameLoginMode string
//
//const (
//	GameLoginModeAccount GameLoginMode = "account"
//	GameLoginModeMobile  GameLoginMode = "mobile"
//)

type GameLoginMode int32

const (
	GameLoginModeAccount GameLoginMode = 1
	GameLoginModeMobile  GameLoginMode = 2
)
