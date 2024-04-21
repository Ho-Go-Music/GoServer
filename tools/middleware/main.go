package middleware

import (
	"fmt"
	"net/http"
)

// CorsMiddleware is a middleware to set CORS headers
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set allowed origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
		// set allowed method
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// set allowed header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// make client carry cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 检查请求是否为 OPTIONS 预检请求
		// Check if the request is an OPTIONS preflight request
		if r.Method == "OPTIONS" {
			// Set additional headers for preflight request
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			w.WriteHeader(http.StatusOK)
			return
		}
		// 继续处理下一个中间件或请求处理函数
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware is a middleware to check if the request is authorized
func Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		fmt.Fprintf(w, "hello, %s", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
