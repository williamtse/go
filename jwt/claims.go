package jwt

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtV5 "github.com/golang-jwt/jwt/v5"
)

func GetValFromContext(ctx context.Context, key string) (string, error) {
	token, ok := jwt.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("获取登录用户ID失败1")
	}
	claims, ok := token.(jwtV5.MapClaims)
	if !ok {
		return "", fmt.Errorf("获取登录用户ID失败2")
	}
	userId, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("获取登录用户ID失败3")
	}
	return userId, nil
}
