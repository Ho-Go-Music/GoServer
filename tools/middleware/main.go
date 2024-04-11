package middleware

import "net/http"

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set allowed origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// set allowed method
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// set allowed header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 检查请求是否为 OPTIONS 预检请求
		// Check if the request is an OPTIONS preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// 继续处理下一个中间件或请求处理函数
		next.ServeHTTP(w, r)
	})
}
