package req

import "battle-tiles/pkg/utils"

type AddBasicMenuReq struct {
	ParentID  int32  `gorm:"column:parent_id;type:integer;not null;comment:父级ID" json:"parent_id" binding:"omitempty" default:"-1" title:"父级路由"`  // 父级ID
	Path      string `gorm:"column:path;type:character varying(100);comment:路由path" json:"path" binding:"omitempty" title:"路由路径"`                 // 路由path
	Name      string `gorm:"column:name;type:character varying(100);comment:路由name" json:"name" binding:"required" title:"路由名称"`                  // 路由name
	Hidden    bool   `gorm:"column:hidden;type:boolean;not null;comment:是否隐藏" json:"hidden" binding:"omitempty" default:"false" title:"是否隐藏"`     // 是否隐藏
	Component string `gorm:"column:component;type:character varying(125);comment:对应前端文件路径" binding:"omitempty" json:"component" title:"对应前端文件路径"` // 对应前端文件路径
	Sort      int32  `gorm:"column:sort;type:integer;not null;default:1;comment:排序" json:"sort" binding:"required"  default:"10" title:"排序"`      // 排序
}

type UpdateBasicMenuReq struct {
	AddBasicMenuReq `mapstructure:",squash"`
	Id              int32 `json:"id" from:"id" mapstructure:"id" title:"主键" binding:"required"` // 主键
}

type ListBasicMenuReq struct {
	*utils.PageParam
}
type SaveMenuTree struct {
	MenuTree []MenuInfo `json:"menu_tree"`
}

// 用于批量保存树
type MenuInfo struct {
	ParentID        int32      `json:"parent_id" default:"-1"`
	MenuType        int32      `json:"menu_type"`
	Title           string     `json:"title"`
	Name            string     `json:"name"`
	Path            string     `json:"path"`
	Component       string     `json:"component"`
	Rank            *string    `json:"rank,omitempty"`
	Redirect        string     `json:"redirect"`
	Icon            string     `json:"icon"`
	ExtraIcon       string     `json:"extra_icon"`
	EnterTransition string     `json:"enter_transition"`
	LeaveTransition string     `json:"leave_transition"`
	ActivePath      string     `json:"active_path"`
	Auths           string     `json:"auths"`
	FrameSrc        string     `json:"frame_src"`
	FrameLoading    bool       `json:"frame_loading"`
	KeepAlive       bool       `json:"keep_alive"`
	HiddenTag       bool       `json:"hidden_tag"`
	FixedTag        bool       `json:"fixed_tag"`
	ShowLink        bool       `json:"show_link"`
	ShowParent      bool       `json:"show_parent"`
	Children        []MenuInfo `json:"children,omitempty"`
}
