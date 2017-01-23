package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

// Server binds API endpoints to their handlers and initializes middleware.
func Server() http.Handler {
	mux := bone.New()

	// Status
	mux.Get("/health", http.HandlerFunc(healthCheck))

	mux.Post("/sessions", http.HandlerFunc(login))
	mux.Post("/users", http.HandlerFunc(registerUser))

	mux.Post("/access-checks", http.HandlerFunc(checkAccess))

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
