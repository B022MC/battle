package request

import (
	"github.com/golang-jwt/jwt/v4"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UserID   int32
	Platform string
	Username string
	NickName string

	Roles []int32  `json:"roles,omitempty"`
	Perms []string `json:"perms,omitempty"`
}

func (b BaseClaims) IsSuperAdmin() bool {
	// 这里随意：比如角色ID=1认为是超管；或在 Perms 里有 "root:*"
	for _, r := range b.Roles {
		if r == 1 {
			return true
		}
	}
	return false
}
func (b BaseClaims) GetPerms() []string { return b.Perms }
