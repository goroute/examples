package main

import (
	"flag"
	"github.com/goroute/compress"
	"log"
	"net/http"
	"time"

	"github.com/goroute/route"
)

var (
	addr = flag.String("addr", ":9000", "Server serve address")
)

func main() {
	flag.Parse()

	mux := route.NewServeMux(
		route.WithRenderer(route.NewDefaultTemplateRenderer("views/*html")),
	)

	mux.Use(compress.New())

	// Serve static files under /static path from assets folder.
	mux.Static("/static", "assets")

	// Render index.html.
	mux.GET("/", func(c route.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{"name": "Go Route!"})
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         *addr,
		Handler:      mux,
	}
	log.Fatal(srv.ListenAndServe())
}
