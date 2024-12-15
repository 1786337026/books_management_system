package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JWT 密钥（应该存放在环境变量中或配置文件里）
var SecretKey = []byte("your_secret_key")

// 用户结构体（JWT 中的数据部分）
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"` // 新增角色字段
	jwt.StandardClaims
}

// 生成 JWT Token
func GenerateToken(username, role string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(24 * time.Hour) // 设置为 24 小时有效期

	// 创建 Claims
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "book-management-system",
		},
	}

	// 使用 HMAC 签名算法生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用 SecretKey 签名生成 Token 字符串
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析和验证 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回用于验证的密钥
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 验证 token 是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
