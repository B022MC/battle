package basic

import (
	basex "battle-tiles/pkg/plugin/gormx/base"
	"time"
)

const TableNameBasicUser = "basic_user"

type BasicUser struct {
	basex.Model[int32]

	// 基础登录信息
	Username string `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_basic_user_username;comment:用户名/员工工号" json:"username"` // 用户名
	Password string `gorm:"column:password;type:varchar(255);comment:密码" json:"password"`                                                  // 密码（哈希）
	Salt     string `gorm:"column:salt;type:varchar(50);comment:盐" json:"salt"`                                                            // 盐

	// 微信号（手填）
	WechatID string `gorm:"column:wechat_id;type:varchar(64);comment:微信号" json:"wechat_id"`

	// 展示信息
	Avatar   string `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像" json:"avatar"`      // 头像
	NickName string `gorm:"column:nick_name;type:varchar(50);not null;default:'';comment:昵称" json:"nick_name"` // 昵称

	Introduction string `gorm:"column:introduction;type:text;comment:个人介绍" json:"introduction"` // 个人介绍

	// 拼音信息
	PinyinCode  string `gorm:"column:pinyin_code;type:varchar(100);comment:拼音码" json:"pinyin_code"`  // 全拼
	FirstLetter string `gorm:"column:first_letter;type:varchar(50);comment:首字母" json:"first_letter"` // 首字母

	// 登录状态
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:timestamp with time zone;comment:最后登录时间" json:"last_login_at"` // 最后登录时间
}

func (BasicUser) TableName() string { return TableNameBasicUser }
