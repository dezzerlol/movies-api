package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *app) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.metrics)
	r.Use(app.recoverPanic)
	r.Use(app.enableCORS)
	r.Use(app.rateLimit)
	r.Use(app.authenticate)

	r.NotFound(app.err.notFoundResponse)
	r.MethodNotAllowed(app.err.notAllowedResponse)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/metrics", app.metricsRouter())

		r.Mount("/movies", app.moviesRouter())
		r.Mount("/users", app.usersRouter())
		r.Mount("/tokens", app.tokensRouter())
	})

	return r
}

// /movies
func (app *app) moviesRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.requirePermission("movies:read", app.listMoviesHandler))
	r.Post("/", app.requirePermission("movies:write", app.createMovieHandler))
	r.Get("/{id}", app.requirePermission("movies:read", app.showMovieHandler))
	r.Patch("/{id}", app.requirePermission("movies:write", app.updateMovieHandler))
	r.Delete("/{id}", app.requirePermission("movies:write", app.deleteMovieHandler))

	return app.requireActivatedUser(r)
}

// /users
func (app *app) usersRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", app.createUserHandler)
	r.Get("/", app.getUserHandler)
	r.Patch("/", app.updateUserHandler)
	r.Put("/activated", app.activateUserHandler)
	r.Put("/password", app.updateUserPasswordHandler)

	return r
}

// /tokens
func (app *app) tokensRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/authentication", app.createAuthTokenHandler)
	r.Post("/password-reset", app.createPasswordResetTokenHandler)
	r.Post("/activation", app.createActivationTokenHandler)

	return r
}

// /metrics
func (app *app) metricsRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthcheck", app.healthcheckHandler)
	r.Handle("/debug/vars", expvar.Handler())

	return r
}
