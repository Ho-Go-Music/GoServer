package main

import (
	diylog "acaibird.com/log"
	"acaibird.com/tools/httpfunc"
	"acaibird.com/tools/middleware"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", httpfunc.LogInHandler)
	mux.HandleFunc("/", httpfunc.RootHandler)
	mux.HandleFunc("/logout", httpfunc.LogOutHandler)
	mux.HandleFunc("/re", httpfunc.RootHandler)
	err := http.ListenAndServe(":8080", middleware.CorsMiddleware(mux))
	if err != nil {
		diylog.Sugar.Errorln(err)
		return
	}
}
