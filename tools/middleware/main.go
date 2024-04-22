package middleware

import (
	diylog "acaibird.com/log"
	sessions2 "acaibird.com/sessions"
	"encoding/gob"
	"net/http"
	"time"
)

// CorsMiddleware is a middleware to set CORS headers
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set allowed origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
			diylog.Sugar.Infoln("preflight request")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			w.WriteHeader(http.StatusOK)
			return
		}
		// 继续处理下一个中间件或请求处理函数
		next.ServeHTTP(w, r)
	})
}

// Identify AuthMiddleware is a middleware to check if the request is authorized
func Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gob.Register(time.Time{})
		session, err := sessions2.Store.Get(r, "session")
		if err != nil {
			diylog.Sugar.Errorln("get sessions error: ", err)
			http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
			return
		}
		// There is a situation where the user has not logged in yet
		// the /login routing request is directly handed over to the corresponding handler for processing.
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}
		// if session is new , user has not logged in
		if session.IsNew {
			//print("new session")
			//http.Error(w, "Unauthorized", http.StatusUnauthorized)
			//return
		}
		// get username from query parameters
		//username := strings.Join(r.URL.Query()["username"], "")
		// Check session integrity
		//if _, ok := session.Values[username]; !ok {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}
		next.ServeHTTP(w, r)
	})
}
