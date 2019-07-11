package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/goroute/route"
)

var (
	addr = flag.String("addr", ":9000", "Server serve address")
)

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c route.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	flag.Parse()

	// Register views.
	renderer := &templateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	mux := route.NewServeMux(route.WithRenderer(renderer))

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
