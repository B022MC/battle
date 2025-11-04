// internal/dal/req/ctrl_account.go
package req

type CtrlListFilter struct {
	LoginMode *int32 // 1=account, 2=mobile
	Status    *int32 // 0/1
	Keyword   string // 按 identifier 模糊
}
type ListAllCtrlAccountsRequest struct {
	LoginMode string `json:"login_mode,omitempty" binding:"omitempty,oneof=account mobile"` // account|mobile
	Status    *int32 `json:"status,omitempty"`                                              // 0/1
	Keyword   string `json:"keyword,omitempty"`                                             // 模糊匹配 identifier
	Page      int    `json:"page" binding:"omitempty,min=1"`
	Size      int    `json:"size" binding:"omitempty,min=1,max=200"`
}
