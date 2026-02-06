package jwt

import (
	"errors"
	"fmt"
	"go-todo-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint
	Username  string
	IsRefresh bool
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired    = errors.New("token已过期")
	ErrTokenNoValidYet = errors.New("token未生效")
	ErrTokenMalformed  = errors.New("token格式错误")
	ErrTokenInvalid    = errors.New("无效的token")
	ErrTokenNotFound   = errors.New("token不存在")
)

func GenerateToken(userID uint, userName string, isRefresh bool) (string, error) {
	jwtConfig := config.GlobalConfig.JWT
	var expiration time.Duration
	if !isRefresh {
		expiration = time.Duration(jwtConfig.AccessExpire) * time.Second
	} else {
		expiration = time.Duration(jwtConfig.RefreshExpire) * time.Second
	}
	Claims := Claims{
		UserID:    userID,
		Username:  userName,
		IsRefresh: isRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   "todo-api-access",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNoValidYet
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

func IstokenExpired(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return errors.Is(err, ErrTokenExpired)
}

func GettokenExpiration(tokenString string) (time.Time, error) {
	Claims, err := ParseToken(tokenString)
	if err != nil {
		return time.Time{}, ErrTokenExpired
	}
	return Claims.ExpiresAt.Time, nil
}
