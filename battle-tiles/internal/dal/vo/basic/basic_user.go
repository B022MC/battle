package basic

import (
	basicModel "battle-tiles/internal/dal/model/basic"
	"time"
)

type BasicUserVo struct {
	*basicModel.BasicUser
}
type RegisterReq struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

// BasicUserDoc 仅用于 Swagger 文档，避免 basex.Model[int32] 泛型带来的解析问题。
type BasicUserDoc struct {
	ID           int32      `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Username     string     `json:"username"`      // 用户名/员工工号
	Password     string     `json:"password"`      // 密码（哈希）
	Salt         string     `json:"salt"`          // 盐
	WechatID     string     `json:"wechat_id"`     // 微信号
	Avatar       string     `json:"avatar"`        // 头像
	NickName     string     `json:"nick_name"`     // 昵称
	Introduction string     `json:"introduction"`  // 个人介绍
	PinyinCode   string     `json:"pinyin_code"`   // 全拼
	FirstLetter  string     `json:"first_letter"`  // 首字母
	LastLoginAt  *time.Time `json:"last_login_at"` // 最后登录时间
}
