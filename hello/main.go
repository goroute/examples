package main

import (
	"github.com/goroute/route"
	"log"
	"net/http"
)

func main() {
	mux := route.NewServeMux()

	mux.GET("/", func(c route.Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}
