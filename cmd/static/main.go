package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/goroute/compress"

	"github.com/goroute/route"
)

var (
	addr = flag.String("addr", ":9000", "Server serve address")
)

var (
	// Base layout templates.
	baseTemplates = []string{
		"views/base.layout.html",
		"views/footer.partial.html",
		"views/navbar.partial.html",
	}

	// Page templates.
	pageTemplates = []string{
		"views/home.page.html",
		"views/about.page.html",
		"views/users.page.html",
	}
)

// newTemplatesRenderer returns templates renderer.
func newTemplatesRenderer() route.Renderer {
	groups := map[string]*template.Template{}
	for _, v := range pageTemplates {
		addTemplate(groups, v)
	}
	return &templateRenderer{groups: groups}
}

// addTemplate adds template to grouped templates map.
func addTemplate(groups map[string]*template.Template, name string) {
	templates := []string{name}
	templates = append(templates, baseTemplates...)
	groups[name] = template.Must(template.ParseFiles(templates...))
}

type templateRenderer struct {
	groups map[string]*template.Template
}

// Render renders template with given name and data.
func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c route.Context) error {
	if tmpl, ok := t.groups[name]; ok {
		return tmpl.Execute(w, data)
	}
	return fmt.Errorf("could not find template: %s", name)
}

func main() {
	flag.Parse()

	mux := route.NewServeMux(
		route.WithRenderer(newTemplatesRenderer()),
	)

	mux.Use(compress.New())
	mux.Use(func(c route.Context, next route.HandlerFunc) error {
		err := next(c)
		if err != nil {
			log.Println("ERR", err)
		}
		return err
	})

	// Serve static files under /static path from assets folder.
	mux.Static("/static", "assets")

	// Render home page.
	mux.GET("/", home)
	mux.GET("/about", about)
	mux.GET("/users", users)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         *addr,
		Handler:      mux,
	}
	log.Fatal(srv.ListenAndServe())
}

func home(c route.Context) error {
	return c.Render(http.StatusOK, "views/home.page.html", struct{}{})
}

func about(c route.Context) error {
	model := struct {
		Now time.Time
	}{
		Now: time.Now(),
	}
	return c.Render(http.StatusOK, "views/about.page.html", model)
}

type User struct {
	ID          int
	FirstName   string
	LastName    string
	AccountName string
}

type usersModel struct {
	Users []User
}

func users(c route.Context) error {
	model := usersModel{
		Users: []User{
			{
				ID:          1,
				FirstName:   "Mark",
				LastName:    "Otto",
				AccountName: "@mdo",
			},
			{
				ID:          2,
				FirstName:   "Jacob",
				LastName:    "Thornton",
				AccountName: "@fat",
			},
			{
				ID:          3,
				FirstName:   "Larry",
				LastName:    "The Bird",
				AccountName: "@twitter",
			},
		},
	}
	return c.Render(http.StatusOK, "views/users.page.html", model)
}
