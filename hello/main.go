package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goroute/route"
)

func logger() route.MiddlewareFunc {
	return func(c route.Context, next route.HandlerFunc) error {
		fmt.Println("Request Path:", c.Request().URL.Path)
		return next(c)
	}
}

func main() {
	mux := route.NewServeMux()

	mux.Use(logger())

	mux.GET("/", func(c route.Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}
