package main

import (
	"net/http"

	"augustkang.com/keyfinder/pkg/handler"
)

func main() {
	http.HandleFunc("/", handler.RootHandler)
	http.ListenAndServe(":8080", nil)
}
