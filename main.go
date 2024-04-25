package main

import (
	diylog "github.com/Ho-Go-Music/GoServer/log"
	"github.com/Ho-Go-Music/GoServer/tools/httpfunc"
	"github.com/Ho-Go-Music/GoServer/tools/middleware"
	"net/http"
)

func main() {
	// server routing
	mux := http.NewServeMux()
	mux.HandleFunc("/login", httpfunc.LogInHandler)
	mux.HandleFunc("/", httpfunc.RootHandler)
	mux.HandleFunc("/logout", httpfunc.LogOutHandler)
	mux.HandleFunc("/test", httpfunc.TestHandler)
	// url path mapping fileSystem
	fs := http.FileServer(http.Dir("./public"))
	// preview file in browser
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	// browser automatically download file
	mux.HandleFunc("/download/", httpfunc.FileDownload)
	// middleware
	var Handler http.Handler
	//Handler = middleware.Identify(mux)
	Handler = middleware.CorsMiddleware(mux)

	err := http.ListenAndServe(":8080", Handler)
	if err != nil {
		diylog.Sugar.Errorln(err)
		return
	}
}
