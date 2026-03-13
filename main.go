package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	svr := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	svr.ListenAndServe()
}
