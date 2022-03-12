package main

import (
	"net/http"
)

func main() {

	srv := &http.Server{
		Addr:    ":80",
		Handler: routes(),
	}

	srv.ListenAndServe()
}
