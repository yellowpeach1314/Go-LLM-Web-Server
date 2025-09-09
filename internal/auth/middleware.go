package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserStorage 用户存储接口
type UserStorage interface {
	GetUserByID(id int) (interface{}, error)
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtService *JWTService, userStorage UserStorage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从Authorization header获取token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "缺少Authorization header")
				return
			}

			// 检查Bearer前缀
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "无效的Authorization header格式")
				return
			}

			tokenString := tokenParts[1]

			// 验证token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("token验证失败: %v", err))
				return
			}

			// 从数据库获取用户信息（确保用户仍然存在且活跃）
			user, err := userStorage.GetUserByID(claims.UserID)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "用户不存在或已被禁用")
				return
			}

			// 将用户信息添加到请求上下文
			ctx := context.WithValue(r.Context(), "user", user)
			ctx = context.WithValue(ctx, "user_id", claims.UserID)

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware 可选认证中间件（允许匿名访问，但如果有token则验证）
func OptionalAuthMiddleware(jwtService *JWTService, userStorage UserStorage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			// 如果没有Authorization header，直接继续
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// 如果有Authorization header，尝试验证
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				tokenString := tokenParts[1]

				if claims, err := jwtService.ValidateToken(tokenString); err == nil {
					if user, err := userStorage.GetUserByID(claims.UserID); err == nil {
						// 将用户信息添加到请求上下文
						ctx := context.WithValue(r.Context(), "user", user)
						ctx = context.WithValue(ctx, "user_id", claims.UserID)
						r = r.WithContext(ctx)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext 从请求上下文获取用户信息
func GetUserFromContext(r *http.Request) (interface{}, bool) {
	user := r.Context().Value("user")
	return user, user != nil
}

// GetUserIDFromContext 从请求上下文获取用户ID
func GetUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value("user_id").(int)
	return userID, ok
}

// respondWithError 返回错误响应
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":  message,
		"status": "error",
		"code":   fmt.Sprintf("%d", statusCode),
	})
}
