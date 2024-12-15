package middleware

import (
	"context"
	"echoingsilence.cn/books-management-system/auth"
	"net/http"
	"strings"
)

// 验证请求中的 Token
func TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 Token
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, `{"error":"Token is missing"}`, http.StatusUnauthorized)
			return
		}

		// 去掉 "Bearer " 部分
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// 解析并验证 Token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// 将用户信息存入上下文中，供后续处理使用
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		r = r.WithContext(ctx)

		// 调用下一个处理函数
		next(w, r)
	}
}

// 验证是否是管理员
func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "admin" {
			http.Error(w, `{"error":"Permission denied"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// 验证是否是读者
func ReaderAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "reader" {
			http.Error(w, `{"error":"Permission denied"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// CORS 中间件
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置跨域头部
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // 允许所有域访问
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // 允许的 HTTP 方法
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // 允许的请求头部
		// 处理 OPTIONS 预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}
