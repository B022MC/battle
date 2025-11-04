package basic

import (
	basicModel "battle-tiles/internal/dal/model/basic"
	"time"
)

type BasicMenuVo struct {
	*basicModel.BasicMenu `mapstructure:",squash"`
	// eg: Name string `json:"name" from:"name" mapstructure:"name" title:"名称" binding:"required"`
}
type BasicMenuDoc struct {
	ID              int32          `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	ParentId        int32          `json:"parent_id"`
	MenuType        int32          `json:"menu_type"`
	Title           string         `json:"title"`
	Name            string         `json:"name"`
	Path            string         `json:"path"`
	Component       string         `json:"component"`
	Rank            *string        `json:"rank,omitempty"`
	Redirect        string         `json:"redirect"`
	Icon            string         `json:"icon"`
	ExtraIcon       string         `json:"extra_icon"`
	EnterTransition string         `json:"enter_transition"`
	LeaveTransition string         `json:"leave_transition"`
	ActivePath      string         `json:"active_path"`
	Auths           string         `json:"auths"`
	FrameSrc        string         `json:"frame_src"`
	FrameLoading    bool           `json:"frame_loading"`
	KeepAlive       bool           `json:"keep_alive"`
	HiddenTag       bool           `json:"hidden_tag"`
	FixedTag        bool           `json:"fixed_tag"`
	ShowLink        bool           `json:"show_link"`
	ShowParent      bool           `json:"show_parent"`
	Children        []BasicMenuDoc `json:"children,omitempty"`
}
