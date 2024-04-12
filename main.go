package main

import (
	diylog "acaibird.com/log"
)

func main() {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/login", httpfunc.LogInHandler)
	//http.ListenAndServe(":8080", middleware.CorsMiddleware(mux))
	diylog.Sugar.Info("Hello, World!")
}
