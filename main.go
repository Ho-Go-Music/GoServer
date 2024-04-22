package main

import (
	diylog "acaibird.com/log"
	"acaibird.com/tools/httpfunc"
	"acaibird.com/tools/middleware"
	"net/http"
)

func main() {
	// server routing
	mux := http.NewServeMux()
	mux.HandleFunc("/login", httpfunc.LogInHandler)
	mux.HandleFunc("/", httpfunc.RootHandler)
	mux.HandleFunc("/logout", httpfunc.LogOutHandler)
	mux.HandleFunc("/test", httpfunc.TestHandler)
	// middleware
	var Handler http.Handler
	Handler = middleware.Identify(mux)
	Handler = middleware.CorsMiddleware(Handler)

	err := http.ListenAndServe(":8080", Handler)
	if err != nil {
		diylog.Sugar.Errorln(err)
		return
	}
}
