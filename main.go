package main

import (
	"blog/router"
	"net/http"
)

func main() {
	server := http.Server{
		Addr: "localhost:8080",
	}

	router.Router()

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
