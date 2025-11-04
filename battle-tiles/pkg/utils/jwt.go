package utils

import (
	"battle-tiles/pkg/utils/request"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/sync/singleflight"
)

const (
	SigningKey                = "mc"
	HeaderSignTokenBufferTime = "1d" // HeaderSignTokenBufferTime 签名验证 缓冲时间
	HeaderSignTokenExpires    = "7d" // HeaderSignTokenExpires 签名有效期为 天
)

var singleGroup = &singleflight.Group{}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")

		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}

	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(SigningKey),
	}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	bf, _ := ParseDuration(HeaderSignTokenBufferTime)

	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失

		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{SigningKey},              // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(bf)),    // 过期时间 1天  配置文件
			Issuer:    SigningKey,                                // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//// 创建两个token
//func (j *JWT) CreateToken2(claims request.CustomClaims) (string, string, error) {
//	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.SigningKey)
//	if err != nil {
//		return "", "", err
//	}
//	ep, _ := ParseDuration(HeaderSignTokenExpires)
//	// rToken 不需要存储任何自定义数据 存 不从旧token取
//	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, request.CustomClaims{
//		BaseClaims: claims.BaseClaims,
//		BufferTime: int64(ep / time.Second),
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)), // 过期时间
//			Issuer:    SigningKey,                             // 签发人
//		},
//	}).SignedString(j.SigningKey)
//	if err != nil {
//		return "", "", err
//	}
//	return accessToken, refreshToken, nil
//}

// CreateToken2 创建两个token 单点登录的过期时间需要设置为单点的过期时间  (expireTimeStampVar ...int64。  为了让expireTimeStampVar变为可变参数，且不影响其他项目2025-02-27)
func (j *JWT) CreateToken2(claims request.CustomClaims, expireTimeStampVar ...int64) (string, string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.SigningKey)
	if err != nil {
		return "", "", err
	}

	// expireTimeStamp := int64(0)
	// if len(expireTimeStampVar) != 0 {
	// 	expireTimeStamp = expireTimeStampVar[0]
	// }

	var ep time.Duration
	ep, _ = ParseDuration(HeaderSignTokenExpires)

	// if expireTimeStamp == 0 {
	// } else {
	// 	expireTime := time.Unix(expireTimeStamp, 0)
	// 	ep = time.Now().Sub(expireTime)
	// }
	// rToken 不需要存储任何自定义数据 存 不从旧token取
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, request.CustomClaims{
		BaseClaims: claims.BaseClaims,
		BufferTime: int64(ep / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)), // 过期时间
			Issuer:    SigningKey,                             // 签发人
		},
	}).SignedString(j.SigningKey)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := singleGroup.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	if tokenString == "" {
		return nil, TokenNotValidYet
	}
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				if claims, ok := token.Claims.(*request.CustomClaims); ok {
					return claims, nil
				}
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
