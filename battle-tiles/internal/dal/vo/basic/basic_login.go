package basic

type LoginResponse struct {
	User         *BaseUserInfo `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"` // 单位：秒
	Platform     string        `json:"cloud"`
}
type BaseUserInfo struct {
	ID           int32  `json:"id"`
	Username     string `json:"username"`
	Phone        string `json:"phone,omitempty"`
	Avatar       string `json:"avatar"`
	NickName     string `json:"nick_name"`
	RealName     string `json:"real_name,omitempty"`
	IdentityID   string `json:"identity_id,omitempty"`
	IdentityURL  string `json:"identity_url,omitempty"`
	Introduction string `json:"introduction,omitempty"`
}
