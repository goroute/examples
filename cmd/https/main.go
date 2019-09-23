package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/goroute/route"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	mux := route.NewServeMux()

	mux.GET("/", func(c route.Context) error {
		return c.String(http.StatusOK, "Hello TLS!")
	})

	log.Fatal(serveTLS(mux))
}

func serveTLS(mux *route.Mux) error {
	hostName := "example.com"
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(hostName, fmt.Sprintf("www.%s", hostName)),
	}
	srv := &http.Server{
		Handler:           mux,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		TLSConfig:         m.TLSConfig(),
	}
	return srv.ListenAndServeTLS("", "")
}
