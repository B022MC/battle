package req

// VerifyAccountRequest 探活 82（不入库）
type VerifyAccountRequest struct {
	Mode    string `json:"mode"    binding:"required,oneof=account mobile"` // 登录方式
	Account string `json:"account" binding:"required"`                      // 账号或手机号
	PwdMD5  string `json:"pwd_md5" binding:"required,len=32"`               // 32位大写MD5
}

// BindMyAccountRequest 绑定“我的唯一游戏账号”
// 绑定“我的”账号（仅创建 game_account）
type BindMyAccountRequest struct {
	Mode     string `json:"mode"     binding:"required,oneof=account mobile"`
	Account  string `json:"account"  binding:"required"`
	PwdMD5   string `json:"pwd_md5"  binding:"required,len=32"`
	Nickname string `json:"nickname" binding:"omitempty"`
}
