// Package jwt JWT Token 生成和验证工具
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Claims JWT 载荷结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
//
// 参数:
//   - userID: 用户ID
//   - username: 用户名
//   - email: 邮箱
//   - secret: JWT 密钥
//   - expireHours: 过期时间（小时）
//
// 返回:
//   - string: JWT Token 字符串
//   - error: 错误信息
func GenerateToken(userID uint, username, email, secret string, expireHours int) (string, error) {
	if secret == "" {
		return "", gerror.New("JWT 密钥不能为空")
	}

	if expireHours <= 0 {
		expireHours = 24 // 默认 24 小时
	}

	// 创建 Claims
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "allfi",
			Subject:   username,
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", gerror.Wrap(err, "签名 Token 失败")
	}

	return tokenString, nil
}

// ParseToken 解析和验证 JWT Token
//
// 参数:
//   - tokenString: JWT Token 字符串
//   - secret: JWT 密钥
//
// 返回:
//   - *Claims: 解析后的载荷
//   - error: 错误信息
func ParseToken(tokenString, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, gerror.New("Token 不能为空")
	}

	if secret == "" {
		return nil, gerror.New("JWT 密钥不能为空")
	}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gerror.Newf("无效的签名算法: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, gerror.Wrap(err, "解析 Token 失败")
	}

	// 验证 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, gerror.New("无效的 Token")
}

// RefreshToken 刷新 Token（如果 Token 即将过期）
//
// 参数:
//   - tokenString: 旧的 JWT Token
//   - secret: JWT 密钥
//   - refreshThreshold: 刷新阈值（小时）- 如果 Token 剩余有效期少于此值，则刷新
//   - expireHours: 新 Token 的过期时间（小时）
//
// 返回:
//   - string: 新的 JWT Token（如果需要刷新）或原 Token
//   - bool: 是否刷新了 Token
//   - error: 错误信息
func RefreshToken(tokenString, secret string, refreshThreshold, expireHours int) (string, bool, error) {
	// 解析旧 Token
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return "", false, err
	}

	// 检查是否需要刷新
	expiresAt := claims.ExpiresAt.Time
	now := time.Now()
	remainingTime := expiresAt.Sub(now)

	// 如果剩余时间少于阈值，则刷新
	if remainingTime < time.Hour*time.Duration(refreshThreshold) {
		newToken, err := GenerateToken(
			claims.UserID,
			claims.Username,
			claims.Email,
			secret,
			expireHours,
		)
		if err != nil {
			return "", false, err
		}
		return newToken, true, nil
	}

	// 不需要刷新，返回原 Token
	return tokenString, false, nil
}
