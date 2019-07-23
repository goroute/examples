package main

import (
	"log"
	"net/http"

	"github.com/goroute/route"
)


func main() {
	mux := route.NewServeMux()

	mux.Use(func(c route.Context, next route.HandlerFunc) error {
		log.Println(c.Request().URL.Path)
		return next(c)
	})

	mux.GET("/", func(c route.Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}
