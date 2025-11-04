package resp

type MenuTree []Menu

// 与 DB 字段对齐，方便前端直用
type Menu struct {
	ID              int32   `json:"id"`
	ParentID        int32   `json:"parent_id"`
	MenuType        int32   `json:"menu_type"`
	Title           string  `json:"title"`
	Name            string  `json:"name"`
	Path            string  `json:"path"`
	Component       string  `json:"component"`
	Rank            *string `json:"rank,omitempty"`
	Redirect        string  `json:"redirect"`
	Icon            string  `json:"icon"`
	ExtraIcon       string  `json:"extra_icon"`
	EnterTransition string  `json:"enter_transition"`
	LeaveTransition string  `json:"leave_transition"`
	ActivePath      string  `json:"active_path"`
	Auths           string  `json:"auths"`
	FrameSrc        string  `json:"frame_src"`
	FrameLoading    bool    `json:"frame_loading"`
	KeepAlive       bool    `json:"keep_alive"`
	HiddenTag       bool    `json:"hidden_tag"`
	FixedTag        bool    `json:"fixed_tag"`
	ShowLink        bool    `json:"show_link"`
	ShowParent      bool    `json:"show_parent"`
	Children        []Menu  `json:"children,omitempty"`
}
