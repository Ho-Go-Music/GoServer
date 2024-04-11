package main

import (
	"acaibird.com/tools/httpfunc"
	"acaibird.com/tools/middleware"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", httpfunc.LogInHandler)
	//// 创建 CORS 中间件
	//corsMiddleware := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"http://localhost:5173"}, // 允许的源
	//	AllowedMethods:   []string{"GET", "POST"},           // 允许的请求方法
	//	AllowedHeaders:   []string{"Content-Type"},          // 允许的请求头
	//	AllowCredentials: true,                              // 允许携带凭证（如Cookie）
	//})
	//handlerWithCors := corsMiddleware.Handler(mux)
	http.ListenAndServe(":8080", middleware.CorsMiddleware(mux))
}
