package main

import (
	"net/http"

	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/packages/config"
	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/packages/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// // using pat package
	// mux := pat.New()
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	// mux.Get("/favicon.ico", http.HandlerFunc(handlers.Repo.EmptyFunc))

	//using chi package
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/favicon.ico", handlers.Repo.EmptyFunc)

	return mux
}
