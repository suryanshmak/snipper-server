package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.recoverPanic, app.logRequest, secureHeaders, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}))
	mux.Get("/", app.home)
	mux.Post("/snippet/create", app.createSnippet)
	mux.Post("/user/sign-up", app.signup)
	mux.Post("/user/login", app.login)
	mux.Post("/user/logout", app.logout)

	return mux
}
