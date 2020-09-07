package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
