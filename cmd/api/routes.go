package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// custom handler
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowResponse)

	// healthcheck handler
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	// movie handler CRUD
	router.HandlerFunc(http.MethodPost, "/api/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/api/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/movies/:id", app.deleteMovieHandler)
	// movie handler endpoint param
	router.HandlerFunc(http.MethodGet, "/api/v1/movies", app.listMoviesHandler)

	return app.recoverPanic(app.rateLimit(router))
}
