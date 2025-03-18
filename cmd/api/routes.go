package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// custom handler
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// healthcheck handler
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	// movie handler
	router.HandlerFunc(http.MethodPost, "/api/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies/:id", app.showMovieHandler)

	return router
}
