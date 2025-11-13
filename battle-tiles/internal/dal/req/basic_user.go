package req

import "battle-tiles/pkg/utils"

type AddBasicUserReq struct {
	UserName string `json:"username" from:"username" mapstructure:"username" title:"用户名称"`
	Password string `json:"password" from:"password" mapstructure:"password" title:"用户密码"`
}

type UpdateBasicUserReq struct {
	AddBasicUserReq `mapstructure:",squash"`
	Id              int32 `json:"id" from:"id" mapstructure:"id" title:"主键" binding:"required"` // 主键
}

type ListBasicUserReq struct {
	*utils.PageParam
}

type UsernamePasswordLoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"`   // 密码
}
type PhoneCodeLoginRequest struct {
	Phone string `json:"phone" binding:"required" example:"13800001111"` // 手机号
	Code  string `json:"code" binding:"required" example:"123456"`       // 验证码
}
type WeChatPhoneLoginRequest struct {
	Code          string `json:"code" binding:"required" example:"081xXx..."` // wx.login 拿到的临时 code
	EncryptedData string `json:"encrypted_data" binding:"required"`           // wx.getPhoneNumber 拿到的加密数据
	IV            string `json:"iv" binding:"required"`                       // 加密数据的 iv
}

type RegisterRequest struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"` // 前端走现在的 RSA 加密再传
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	WechatID string `json:"wechat_id"` //  微信号（人工填写）
	// 游戏账号绑定字段（可选）
	GameAccountMode string `json:"game_account_mode"` // "account" 或 "mobile"
	GameAccount     string `json:"game_account"`      // 游戏账号或手机号
	GamePassword    string `json:"game_password"`     // 游戏密码（MD5）
}
