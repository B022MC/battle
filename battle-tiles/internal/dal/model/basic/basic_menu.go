package basic

import basex "battle-tiles/pkg/plugin/gormx/base"

const TableNameBasicMenu = "basic_menu"

type BasicMenu struct {
	basex.Model[int32]

	ParentId        int32   `gorm:"column:parent_id;not null;default:-1;comment:父级ID 默认为-1" json:"parent_id"`
	MenuType        int32   `gorm:"column:menu_type;not null;comment:菜单类型" json:"menu_type"`
	Title           string  `gorm:"column:title;type:varchar;not null;comment:标题" json:"title"`
	Name            string  `gorm:"column:name;type:varchar;not null;comment:name" json:"name"`
	Path            string  `gorm:"column:path;type:varchar;not null;comment:path" json:"path"`
	Component       string  `gorm:"column:component;type:varchar;not null;comment:组件" json:"component"`
	Rank            *string `gorm:"column:rank;type:varchar;comment:排序" json:"rank"`
	Redirect        string  `gorm:"column:redirect;type:varchar;not null;comment:redirect" json:"redirect"`
	Icon            string  `gorm:"column:icon;type:varchar;not null;comment:图标" json:"icon"`
	ExtraIcon       string  `gorm:"column:extra_icon;type:varchar;not null;comment:额外图标" json:"extra_icon"`
	EnterTransition string  `gorm:"column:enter_transition;type:varchar;not null;comment:进入动画" json:"enter_transition"`
	LeaveTransition string  `gorm:"column:leave_transition;type:varchar;not null;comment:离开动画" json:"leave_transition"`
	ActivePath      string  `gorm:"column:active_path;type:varchar;not null;comment:激活路径" json:"active_path"`
	Auths           string  `gorm:"column:auths;type:varchar;not null;comment:权限" json:"auths"`
	FrameSrc        string  `gorm:"column:frame_src;type:varchar;not null;comment:内嵌 iframe 地址" json:"frame_src"`
	FrameLoading    bool    `gorm:"column:frame_loading;not null;default:false;comment:是否显示加载动画" json:"frame_loading"`
	KeepAlive       bool    `gorm:"column:keep_alive;not null;default:false;comment:是否 keep-alive" json:"keep_alive"`
	HiddenTag       bool    `gorm:"column:hidden_tag;not null;default:false;comment:是否隐藏" json:"hidden_tag"`
	FixedTag        bool    `gorm:"column:fixed_tag;not null;default:false;comment:是否固定" json:"fixed_tag"`
	ShowLink        bool    `gorm:"column:show_link;not null;default:true;comment:是否显示链接" json:"show_link"`
	ShowParent      bool    `gorm:"column:show_parent;not null;default:true;comment:是否显示父级" json:"show_parent"`

	// 仅用于入参/出参的树形结构，不入库
	Children []*BasicMenu `gorm:"-" json:"children,omitempty"`
}

func (*BasicMenu) TableName() string { return TableNameBasicMenu }
